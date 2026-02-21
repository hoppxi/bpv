package library

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/scanner"
)

type Library struct {
	MusicDir   string
	ScanResult *scanner.ScanResult
	ScanTime   time.Time

	mu        sync.RWMutex
	favorites []string
	playStats map[string]int
}

func ScanDirectory(musicDir string) (*Library, error) {
	sc := scanner.NewScanner()
	result, err := sc.ScanLibrary(musicDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan library: %w", err)
	}

	return &Library{
		MusicDir:   musicDir,
		ScanResult: result,
		ScanTime:   time.Now(),
		favorites:  []string{},
		playStats:  make(map[string]int),
	}, nil
}

func (l *Library) Files() []metadata.AudioFile {
	if l.ScanResult == nil {
		return nil
	}
	return l.ScanResult.Files
}

func (l *Library) Artists() map[string]int {
	if l.ScanResult == nil {
		return nil
	}
	return l.ScanResult.Artists
}

func (l *Library) Albums() map[string]int {
	if l.ScanResult == nil {
		return nil
	}
	return l.ScanResult.Albums
}

func (l *Library) Genres() map[string]int {
	if l.ScanResult == nil {
		return nil
	}
	return l.ScanResult.Genres
}

func (l *Library) Composers() map[string]int {
	if l.ScanResult == nil {
		return nil
	}
	return l.ScanResult.Composers
}

func (l *Library) ArtistTracks(artist string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range l.Files() {
		if f.Artist == artist {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (l *Library) AlbumTracks(album string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range l.Files() {
		if f.Album == album {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (l *Library) GenreTracks(genre string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range l.Files() {
		if f.Genre == genre {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (l *Library) Search(query string) []metadata.AudioFile {
	q := strings.ToLower(query)
	var results []metadata.AudioFile
	for _, f := range l.Files() {
		if strings.Contains(strings.ToLower(f.Title), q) ||
			strings.Contains(strings.ToLower(f.Artist), q) ||
			strings.Contains(strings.ToLower(f.Album), q) ||
			strings.Contains(strings.ToLower(f.Genre), q) {
			results = append(results, f)
		}
	}
	return results
}

func (l *Library) TotalTracks() int {
	if l.ScanResult == nil {
		return 0
	}
	return l.ScanResult.AudioFiles
}

func (l *Library) GetFavorites() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]string, len(l.favorites))
	copy(out, l.favorites)
	return out
}

func (l *Library) AddFavorite(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if slices.Contains(l.favorites, path) {
		return
	}
	l.favorites = append(l.favorites, path)
}

func (l *Library) RemoveFavorite(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, f := range l.favorites {
		if f == path {
			l.favorites = append(l.favorites[:i], l.favorites[i+1:]...)
			return
		}
	}
}

func (l *Library) RecordPlay(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.playStats[path]++
}

func (l *Library) GetPlayStats() map[string]int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]int, len(l.playStats))
	maps.Copy(out, l.playStats)
	return out
}
