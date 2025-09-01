package metadata

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/hoppxi/bpv/pkg/logger"
)

type AudioFile struct {
	FilePath     string         `json:"file_path"`
	FileName     string         `json:"file_name"`
	FileSize     int64          `json:"file_size"`
	FileType     string         `json:"file_type"`
	Modified     time.Time      `json:"modified"`
	Title        string         `json:"title"`
	Artist       string         `json:"artist"`
	Album        string         `json:"album"`
	AlbumArtist  string         `json:"album_artist"`
	Composer     string         `json:"composer"`
	Genre        string         `json:"genre"`
	Year         int            `json:"year"`
	Track        int            `json:"track"`
	TotalTracks  int            `json:"total_tracks"`
	Disc         int            `json:"disc"`
	TotalDiscs   int            `json:"total_discs"`
	Duration     time.Duration  `json:"duration"`
	Bitrate      int            `json:"bitrate"`
	SampleRate   int            `json:"sample_rate"`
	Channels     int            `json:"channels"`
	Comment      string         `json:"comment"`
	Lyrics       string         `json:"lyrics"`
	BPM          int            `json:"bpm"`
	CoverArt     string         `json:"cover_art,omitempty"` // Base64 encoded
	CoverArtMime string         `json:"cover_art_mime,omitempty"`
	RawMetadata  map[string]any `json:"raw_metadata,omitempty"`
	Error        string         `json:"error,omitempty"`
}

type Extractor struct {
	extractCoverArt bool
	maxCoverSize    int
}

func NewExtractor() *Extractor {
	return &Extractor{
		extractCoverArt: true,
		maxCoverSize:    300, // max dimension for cover art
	}
}

func (e *Extractor) ExtractFromFile(filePath string) (*AudioFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	metadata, err := tag.ReadFrom(file)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read metadata: %v", err)
	}

	audioFile := &AudioFile{
		FilePath:    filePath,
		FileName:    filepath.Base(filePath),
		FileSize:    info.Size(),
		FileType:    strings.TrimPrefix(filepath.Ext(filePath), "."),
		Modified:    info.ModTime(),
		RawMetadata: make(map[string]any),
	}

	e.populateBasicMetadata(audioFile, metadata)

	if e.extractCoverArt {
		e.extractCoverArtData(audioFile, metadata)
	}

	e.populateTechnicalMetadata(audioFile, file)

	return audioFile, nil
}

func (e *Extractor) populateBasicMetadata(audioFile *AudioFile, metadata tag.Metadata) {
	// Title
	if title := metadata.Title(); title != "" {
		audioFile.Title = title
	} else {
		audioFile.Title = strings.TrimSuffix(audioFile.FileName, filepath.Ext(audioFile.FileName))
	}

	if artist := metadata.Artist(); artist != "" {
		audioFile.Artist = artist
	}

	if album := metadata.Album(); album != "" {
		audioFile.Album = album
	}
	
	if albumArtist := metadata.AlbumArtist(); albumArtist != "" {
		audioFile.AlbumArtist = albumArtist
	}

	if composer := metadata.Composer(); composer != "" {
		audioFile.Composer = composer
	}

	if genre := metadata.Genre(); genre != "" {
		audioFile.Genre = genre
	}

	if year := metadata.Year(); year > 0 {
		audioFile.Year = year
	}

	track, totalTracks := metadata.Track()
	audioFile.Track = track
	audioFile.TotalTracks = totalTracks

	disc, totalDiscs := metadata.Disc()
	audioFile.Disc = disc
	audioFile.TotalDiscs = totalDiscs

	if comment := metadata.Comment(); comment != "" {
		audioFile.Comment = comment
	}

	// Lyrics (if available in raw format)
	raw := metadata.Raw()
	if raw != nil {
		if lyrics, ok := raw["lyrics"]; ok {
			if lyricStr, ok := lyrics.(string); ok {
				audioFile.Lyrics = lyricStr
			}
		}
		
		// Store raw metadata for advanced use
		maps.Copy(audioFile.RawMetadata, raw)
	}
}

func (e *Extractor) extractCoverArtData(audioFile *AudioFile, metadata tag.Metadata) {
	picture := metadata.Picture()
	if picture == nil {
		return
	}

	img, _, err := image.Decode(bytes.NewReader(picture.Data))
	if err != nil {
		return
	}

	img = e.resizeImage(img)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return
	}

	audioFile.CoverArt = base64.StdEncoding.EncodeToString(buf.Bytes())
	audioFile.CoverArtMime = "image/jpeg"
}

func (e *Extractor) resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width <= e.maxCoverSize && height <= e.maxCoverSize {
		return img
	}

	var newWidth, newHeight int
	if width > height {
		newWidth = e.maxCoverSize
		newHeight = int(float64(height) * float64(e.maxCoverSize) / float64(width))
	} else {
		newHeight = e.maxCoverSize
		newWidth = int(float64(width) * float64(e.maxCoverSize) / float64(height))
	}

	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * width / newWidth
			srcY := y * height / newHeight
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

func (e *Extractor) populateTechnicalMetadata(audioFile *AudioFile, file *os.File) {
	type formatMeta struct {
		Bitrate    int // kbps
		SampleRate int // Hz
		Channels   int
	}

	formatDefaults := map[string]formatMeta{
		"flac": {Bitrate: 1411, SampleRate: 44100, Channels: 2},
		"wav":  {Bitrate: 1411, SampleRate: 44100, Channels: 2},
		"aiff": {Bitrate: 1411, SampleRate: 44100, Channels: 2},
		"mp3":  {Bitrate: 320, SampleRate: 44100, Channels: 2},
		"m4a":  {Bitrate: 256, SampleRate: 44100, Channels: 2},
		"aac":  {Bitrate: 256, SampleRate: 44100, Channels: 2},
	}

	lowerType := strings.ToLower(audioFile.FileType)
	if meta, ok := formatDefaults[lowerType]; ok {
		audioFile.Bitrate = meta.Bitrate
		audioFile.SampleRate = meta.SampleRate
		audioFile.Channels = meta.Channels
	} else {
		audioFile.Bitrate = 128
		audioFile.SampleRate = 44100
		audioFile.Channels = 2
	}

	if fileInfo, err := file.Stat(); err == nil {
		audioFile.FileSize = fileInfo.Size()

		if audioFile.Bitrate > 0 {
			bits := audioFile.FileSize * 8
			bps := int64(audioFile.Bitrate * 1000)
			durationSec := float64(bits) / float64(bps)
			audioFile.Duration = time.Duration(durationSec * float64(time.Second))
		}
	} else {
		audioFile.Duration = 0
		logger.Log.Error("Failed to stat file: %v", err)
	}
}

func (e *Extractor) SetExtractCoverArt(extract bool) {
	e.extractCoverArt = extract
}

func (e *Extractor) SetMaxCoverSize(size int) {
	e.maxCoverSize = size
}