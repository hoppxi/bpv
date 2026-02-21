package tui

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/metadata"
)

type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatAll:
		return "All"
	case RepeatOne:
		return "One"
	default:
		return "Off"
	}
}

func (r RepeatMode) Icon() string {
	switch r {
	case RepeatAll:
		return "🔁"
	case RepeatOne:
		return "🔂"
	default:
		return "➡"
	}
}

const standardSampleRate beep.SampleRate = 44100

type Player struct {
	mu sync.Mutex

	queue        []metadata.AudioFile
	queueIndex   int
	shuffle      bool
	shuffleOrder []int
	repeat       RepeatMode

	playing      bool
	paused       bool
	streamer     beep.StreamSeekCloser
	resampled    *beep.Resampler
	ctrl         *beep.Ctrl
	volume       *effects.Volume
	volLevel     float64
	format       beep.Format
	duration     time.Duration
	currentTrack *metadata.AudioFile
	trackEnded   bool

	speakerInit bool

	favorites map[string]bool
}

func NewPlayer() *Player {
	return &Player{
		repeat:    RepeatOff,
		volLevel:  0,
		favorites: make(map[string]bool),
	}
}

func (p *Player) SetQueue(tracks []metadata.AudioFile, startIndex int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.queue = make([]metadata.AudioFile, len(tracks))
	copy(p.queue, tracks)
	p.queueIndex = startIndex

	if p.shuffle {
		p.buildShuffleOrder()
	}
}

func (p *Player) Queue() []metadata.AudioFile {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]metadata.AudioFile, len(p.queue))
	copy(out, p.queue)
	return out
}

func (p *Player) MoveQueueItem(from, to int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if from < 0 || from >= len(p.queue) || to < 0 || to >= len(p.queue) {
		return
	}
	item := p.queue[from]
	p.queue = append(p.queue[:from], p.queue[from+1:]...)
	newQueue := make([]metadata.AudioFile, 0, len(p.queue)+1)
	newQueue = append(newQueue, p.queue[:to]...)
	newQueue = append(newQueue, item)
	newQueue = append(newQueue, p.queue[to:]...)
	p.queue = newQueue

	// Adjust current index.
	if p.queueIndex == from {
		p.queueIndex = to
	} else if from < p.queueIndex && to >= p.queueIndex {
		p.queueIndex--
	} else if from > p.queueIndex && to <= p.queueIndex {
		p.queueIndex++
	}
}

func (p *Player) RemoveFromQueue(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if index < 0 || index >= len(p.queue) {
		return
	}
	p.queue = append(p.queue[:index], p.queue[index+1:]...)
	if p.queueIndex > index {
		p.queueIndex--
	}
	if p.queueIndex >= len(p.queue) {
		p.queueIndex = len(p.queue) - 1
	}
}

// ─── Playback ───────────────────────────────────────────────────────────────

func (p *Player) PlayCurrent() error {
	p.mu.Lock()
	if len(p.queue) == 0 {
		p.mu.Unlock()
		return fmt.Errorf("queue is empty")
	}

	idx := p.resolveIndex()
	if idx < 0 || idx >= len(p.queue) {
		p.mu.Unlock()
		return fmt.Errorf("invalid queue index")
	}

	track := p.queue[idx]
	p.mu.Unlock()

	return p.playFile(track)
}

func (p *Player) playFile(track metadata.AudioFile) error {
	p.stopInternal()

	f, err := os.Open(track.FilePath)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", track.FilePath, err)
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	ext := strings.ToLower(filepath.Ext(track.FilePath))
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	case ".m4a", ".aac":
		streamer, format, err = DecodeAAC(f)
	default:
		f.Close()
		logger.Log.Warn("Unsupported format %s: skipping %s", ext, track.FileName)
		return p.Next()
	}

	if err != nil {
		f.Close()
		logger.Log.Error("decode error for %s: %v", track.FileName, err)
		return p.Next()
	}

	// Initialize speaker ONCE at the standard sample rate.
	p.mu.Lock()
	if !p.speakerInit {
		p.mu.Unlock()
		err = speaker.Init(standardSampleRate, standardSampleRate.N(time.Second/10))
		if err != nil {
			streamer.Close()
			return fmt.Errorf("speaker init error: %w", err)
		}
		p.mu.Lock()
		p.speakerInit = true
	}
	p.mu.Unlock()

	// Resample to the standard rate if the file has a different sample rate.
	var finalStreamer beep.Streamer
	var resampled *beep.Resampler
	if format.SampleRate != standardSampleRate {
		resampled = beep.Resample(4, format.SampleRate, standardSampleRate, streamer)
		finalStreamer = resampled
	} else {
		finalStreamer = streamer
	}

	ctrl := &beep.Ctrl{Streamer: finalStreamer, Paused: false}
	vol := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   p.volLevel,
		Silent:   false,
	}

	totalSamples := streamer.Len()
	dur := format.SampleRate.D(totalSamples)

	p.mu.Lock()
	p.streamer = streamer
	p.resampled = resampled
	p.ctrl = ctrl
	p.volume = vol
	p.format = format
	p.playing = true
	p.paused = false
	p.duration = dur
	p.trackEnded = false
	trackCopy := track
	p.currentTrack = &trackCopy
	p.mu.Unlock()

	speaker.Play(beep.Seq(vol, beep.Callback(func() {
		p.mu.Lock()
		p.playing = false
		p.paused = false
		p.trackEnded = true
		p.mu.Unlock()
	})))

	return nil
}

func (p *Player) stopInternal() {
	speaker.Clear()
	p.mu.Lock()
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
	p.resampled = nil
	p.ctrl = nil
	p.volume = nil
	p.playing = false
	p.paused = false
	p.trackEnded = false
	p.mu.Unlock()
}

func (p *Player) Stop() {
	p.stopInternal()
}

func (p *Player) TogglePause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ctrl == nil {
		return
	}
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	p.paused = p.ctrl.Paused
	speaker.Unlock()
}

// ─── Track Navigation ───────────────────────────────────────────────────────

func (p *Player) Next() error {
	p.mu.Lock()
	if len(p.queue) == 0 {
		p.mu.Unlock()
		return nil
	}

	p.queueIndex++
	if p.queueIndex >= len(p.queue) {
		if p.repeat == RepeatAll {
			p.queueIndex = 0
		} else {
			p.queueIndex = len(p.queue) - 1
			p.mu.Unlock()
			p.Stop()
			return nil
		}
	}
	p.mu.Unlock()
	return p.PlayCurrent()
}

func (p *Player) Previous() error {
	p.mu.Lock()
	pos := p.positionUnsafe()
	if pos > 3*time.Second {
		p.mu.Unlock()
		return p.PlayCurrent()
	}

	if len(p.queue) == 0 {
		p.mu.Unlock()
		return nil
	}

	p.queueIndex--
	if p.queueIndex < 0 {
		if p.repeat == RepeatAll {
			p.queueIndex = len(p.queue) - 1
		} else {
			p.queueIndex = 0
		}
	}
	p.mu.Unlock()
	return p.PlayCurrent()
}

// CheckTrackEnd checks if the current track has ended and handles auto-advance.
func (p *Player) CheckTrackEnd() bool {
	p.mu.Lock()
	if !p.trackEnded {
		p.mu.Unlock()
		return false
	}
	p.trackEnded = false

	if p.repeat == RepeatOne {
		p.mu.Unlock()
		_ = p.PlayCurrent()
		return true
	}

	p.queueIndex++
	if p.queueIndex >= len(p.queue) {
		if p.repeat == RepeatAll {
			p.queueIndex = 0
		} else {
			p.queueIndex = len(p.queue) - 1
			p.mu.Unlock()
			return false
		}
	}
	p.mu.Unlock()
	_ = p.PlayCurrent()
	return true
}

// ─── Volume ─────────────────────────────────────────────────────────────────

func (p *Player) VolumeUp() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.volume == nil {
		p.volLevel += 0.5
		if p.volLevel > 5 {
			p.volLevel = 5
		}
		return
	}
	speaker.Lock()
	p.volume.Volume += 0.5
	if p.volume.Volume > 5 {
		p.volume.Volume = 5
	}
	p.volLevel = p.volume.Volume
	speaker.Unlock()
}

func (p *Player) VolumeDown() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.volume == nil {
		p.volLevel -= 0.5
		if p.volLevel < -5 {
			p.volLevel = -5
		}
		return
	}
	speaker.Lock()
	p.volume.Volume -= 0.5
	if p.volume.Volume < -5 {
		p.volume.Volume = -5
	}
	p.volLevel = p.volume.Volume
	speaker.Unlock()
}

func (p *Player) ToggleMute() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.volume == nil {
		return
	}
	speaker.Lock()
	p.volume.Silent = !p.volume.Silent
	speaker.Unlock()
}

// ─── Seeking ────────────────────────────────────────────────────────────────

func (p *Player) SeekForward(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.streamer == nil {
		return
	}
	speaker.Lock()
	pos := p.streamer.Position()
	newPos := pos + p.format.SampleRate.N(d)
	if newPos >= p.streamer.Len() {
		newPos = p.streamer.Len() - 1
	}
	_ = p.streamer.Seek(newPos)
	speaker.Unlock()
}

func (p *Player) SeekBackward(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.streamer == nil {
		return
	}
	speaker.Lock()
	pos := p.streamer.Position()
	newPos := pos - p.format.SampleRate.N(d)
	if newPos < 0 {
		newPos = 0
	}
	_ = p.streamer.Seek(newPos)
	speaker.Unlock()
}

// ─── Shuffle / Repeat ───────────────────────────────────────────────────────

func (p *Player) ToggleShuffle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.shuffle = !p.shuffle
	if p.shuffle {
		p.buildShuffleOrder()
	}
}

func (p *Player) CycleRepeat() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.repeat = (p.repeat + 1) % 3
}

// ─── Favorites ──────────────────────────────────────────────────────────────

func (p *Player) ToggleFavorite(filePath string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.favorites[filePath] {
		delete(p.favorites, filePath)
	} else {
		p.favorites[filePath] = true
	}
}

func (p *Player) IsFavorite(filePath string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.favorites[filePath]
}

func (p *Player) FavoriteCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.favorites)
}

func (p *Player) FavoritePaths() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, 0, len(p.favorites))
	for k := range p.favorites {
		out = append(out, k)
	}
	return out
}

// SetFavoritePaths bulk-sets the favorites from daemon-loaded data.
func (p *Player) SetFavoritePaths(paths []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.favorites = make(map[string]bool, len(paths))
	for _, path := range paths {
		p.favorites[path] = true
	}
}

// ─── State Getters ──────────────────────────────────────────────────────────

func (p *Player) positionUnsafe() time.Duration {
	if p.streamer == nil {
		return 0
	}
	speaker.Lock()
	pos := p.streamer.Position()
	speaker.Unlock()
	return p.format.SampleRate.D(pos)
}

func (p *Player) Position() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.positionUnsafe()
}

func (p *Player) Duration() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.duration
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing && !p.paused
}

func (p *Player) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

func (p *Player) HasTrack() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentTrack != nil && (p.playing || p.paused)
}

func (p *Player) CurrentTrack() *metadata.AudioFile {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentTrack
}

func (p *Player) GetVolume() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.volLevel
}

func (p *Player) IsMuted() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.volume == nil {
		return false
	}
	return p.volume.Silent
}

func (p *Player) Shuffle() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.shuffle
}

func (p *Player) Repeat() RepeatMode {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.repeat
}

func (p *Player) QueueLen() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.queue)
}

func (p *Player) QueueIndex() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.resolveIndex()
}

// CurrentQueueItemIndex returns the index into the underlying queue slice of the
// currently playing item. Returns -1 if the queue is empty.
func (p *Player) CurrentQueueItemIndex() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.queue) == 0 {
		return -1
	}
	idx := p.resolveIndex()
	if idx < 0 || idx >= len(p.queue) {
		return -1
	}
	return idx
}

// SetShuffle enables/disables shuffle deterministically.
// When disabling shuffle, the queueIndex is converted back to the underlying
// queue item index so playback stays on the same track.
func (p *Player) SetShuffle(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shuffle == enabled {
		return
	}

	if !enabled {
		idx := p.resolveIndex()
		p.shuffle = false
		p.shuffleOrder = nil
		if idx >= 0 && idx < len(p.queue) {
			p.queueIndex = idx
		}
		return
	}

	p.shuffle = true
	if len(p.queue) > 0 {
		p.buildShuffleOrder()
	}
}

func (p *Player) SetRepeat(mode RepeatMode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.repeat = mode
}

func (p *Player) Progress() float64 {
	dur := p.Duration()
	if dur == 0 {
		return 0
	}
	pos := p.Position()
	prog := float64(pos) / float64(dur)
	if prog > 1 {
		prog = 1
	}
	if prog < 0 {
		prog = 0
	}
	return prog
}

func (p *Player) VolumePercent() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	pct := int((p.volLevel + 5) * 10)
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return pct
}

func (p *Player) Close() {
	p.Stop()
}

// ─── Internal ───────────────────────────────────────────────────────────────

func (p *Player) resolveIndex() int {
	if p.shuffle && len(p.shuffleOrder) > 0 {
		if p.queueIndex < len(p.shuffleOrder) {
			return p.shuffleOrder[p.queueIndex]
		}
	}
	return p.queueIndex
}

func (p *Player) buildShuffleOrder() {
	p.shuffleOrder = make([]int, len(p.queue))
	for i := range p.shuffleOrder {
		p.shuffleOrder[i] = i
	}
	current := p.queueIndex
	if current < len(p.shuffleOrder) {
		p.shuffleOrder[0], p.shuffleOrder[current] = p.shuffleOrder[current], p.shuffleOrder[0]
	}
	for i := len(p.shuffleOrder) - 1; i > 1; i-- {
		j := 1 + rand.Intn(i)
		p.shuffleOrder[i], p.shuffleOrder[j] = p.shuffleOrder[j], p.shuffleOrder[i]
	}
	p.queueIndex = 0
}
