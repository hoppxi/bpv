package metadata

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/hoppxi/bpv/internal/logger"
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
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	audioFile := &AudioFile{
		FilePath:    filePath,
		FileName:    filepath.Base(filePath),
		FileSize:    info.Size(),
		FileType:    strings.TrimPrefix(filepath.Ext(filePath), "."),
		Modified:    info.ModTime(),
		RawMetadata: make(map[string]any),
		Title:       strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)),
		Artist:      "Unknown Artist",
		Album:       "Unknown Album",
		Genre:       "Unknown Genre",
	}

	if info.Size() == 0 {
		audioFile.Error = "File is empty"
		return audioFile, nil
	}

	file.Seek(0, io.SeekStart)
	metadata, err := tag.ReadFrom(file)
	if err != nil && err != io.EOF {
		if logger.Log.IsVerbose() {
			logger.Log.Warn("Failed to extract metadata from %s: %v", filePath, err)
		}
		audioFile.Error = fmt.Sprintf("Metadata extraction error: %v", err)
		return audioFile, nil
	}

	e.populateBasicMetadata(audioFile, metadata)

	if e.extractCoverArt && metadata != nil {
		e.extractCoverArtData(audioFile, metadata)
	}

	// Reset file pointer again for technical metadata extraction
	file.Seek(0, io.SeekStart)
	e.populateTechnicalMetadata(audioFile, file, metadata)

	return audioFile, nil
}

func (e *Extractor) populateBasicMetadata(audioFile *AudioFile, metadata tag.Metadata) {
	if metadata == nil {
		audioFile.Title = strings.TrimSuffix(audioFile.FileName, filepath.Ext(audioFile.FileName))
		audioFile.Artist = "Unknown Artist"
		audioFile.Album = "Unknown Album"
		audioFile.Genre = "Unknown Genre"
		return
	}

	fileName := strings.TrimSuffix(audioFile.FileName, filepath.Ext(audioFile.FileName))

	if title := metadata.Title(); title != "" {
		audioFile.Title = title
	} else {
		if parts := strings.SplitN(fileName, " - ", 2); len(parts) == 2 {
			audioFile.Title = strings.TrimSpace(parts[1])
		} else {
			audioFile.Title = fileName
		}
	}

	if artist := metadata.Artist(); artist != "" {
		audioFile.Artist = artist
	} else {
		if parts := strings.SplitN(fileName, " - ", 2); len(parts) == 2 {
			audioFile.Artist = strings.TrimSpace(parts[0])
		} else {
			audioFile.Artist = "Unknown Artist"
		}
	}

	if album := metadata.Album(); album != "" {
		audioFile.Album = album
	} else {
		parentDir := filepath.Base(filepath.Dir(audioFile.FilePath))
		if parentDir != "." && parentDir != ".." && parentDir != "" {
			audioFile.Album = parentDir
		} else {
			audioFile.Album = "Unknown Album"
		}
	}

	if albumArtist := metadata.AlbumArtist(); albumArtist != "" {
		audioFile.AlbumArtist = albumArtist
	} else {
		audioFile.AlbumArtist = audioFile.Artist
	}

	if composer := metadata.Composer(); composer != "" {
		audioFile.Composer = composer
	}

	if genre := metadata.Genre(); genre != "" {
		audioFile.Genre = genre
	} else {
		audioFile.Genre = "Unknown Genre"
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

	if lyrics := metadata.Lyrics(); lyrics != "" {
		audioFile.Lyrics = lyrics
	}
}

func (e *Extractor) extractCoverArtData(audioFile *AudioFile, metadata tag.Metadata) {
	if metadata == nil {
		return
	}

	picture := metadata.Picture()
	if picture == nil {
		return
	}

	if len(picture.Data) == 0 {
		return
	}

	audioFile.CoverArtMime = picture.MIMEType

	// For non-JPEG images, we'll convert to JPEG
	if picture.MIMEType != "image/jpeg" {
		img, _, err := image.Decode(bytes.NewReader(picture.Data))
		if err != nil {
			if logger.Log.IsVerbose() {
				logger.Log.Warn("Failed to decode cover art: %v", err)
			}
			return
		}

		img = e.resizeImage(img)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
		if err != nil {
			if logger.Log.IsVerbose() {
				logger.Log.Warn("Failed to encode cover art as JPEG: %v", err)
			}
			return
		}

		audioFile.CoverArt = base64.StdEncoding.EncodeToString(buf.Bytes())
		audioFile.CoverArtMime = "image/jpeg"
	} else {
		// For JPEG, we can use the original data but may need to resize
		img, _, err := image.Decode(bytes.NewReader(picture.Data))
		if err != nil {
			if logger.Log.IsVerbose() {
				logger.Log.Warn("Failed to decode JPEG cover art: %v", err)
			}
			return
		}

		img = e.resizeImage(img)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
		if err != nil {
			if logger.Log.IsVerbose() {
				logger.Log.Warn("Failed to re-encode JPEG cover art: %v", err)
			}
			return
		}

		audioFile.CoverArt = base64.StdEncoding.EncodeToString(buf.Bytes())
	}
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

func (e *Extractor) populateTechnicalMetadata(audioFile *AudioFile, _ *os.File, metadata tag.Metadata) {
	if metadata != nil {
		raw := metadata.Raw()

		if format, ok := raw["format"]; ok {
			if formatMap, ok := format.(map[string]interface{}); ok {
				if bitrate, ok := formatMap["bitrate"]; ok {
					if b, ok := bitrate.(int); ok {
						audioFile.Bitrate = b / 1000 // Convert to kbps
					}
				}
				if sampleRate, ok := formatMap["sampleRate"]; ok {
					if sr, ok := sampleRate.(int); ok {
						audioFile.SampleRate = sr
					}
				}
				if channels, ok := formatMap["channels"]; ok {
					if ch, ok := channels.(int); ok {
						audioFile.Channels = ch
					}
				}
			}
		}

		if duration, ok := raw["duration"]; ok {
			if d, ok := duration.(time.Duration); ok {
				audioFile.Duration = d
			} else if d, ok := duration.(float64); ok {
				audioFile.Duration = time.Duration(d * float64(time.Second))
			}
		}
	}

	if audioFile.Bitrate == 0 || audioFile.SampleRate == 0 || audioFile.Channels == 0 {
		formatDefaults := map[string]struct {
			Bitrate    int
			SampleRate int
			Channels   int
		}{
			"flac": {Bitrate: 1411, SampleRate: 44100, Channels: 2},
			"wav":  {Bitrate: 1411, SampleRate: 44100, Channels: 2},
			"aiff": {Bitrate: 1411, SampleRate: 44100, Channels: 2},
			"mp3":  {Bitrate: 320, SampleRate: 44100, Channels: 2},
			"m4a":  {Bitrate: 256, SampleRate: 44100, Channels: 2},
			"aac":  {Bitrate: 256, SampleRate: 44100, Channels: 2},
			"ogg":  {Bitrate: 192, SampleRate: 44100, Channels: 2},
			"opus": {Bitrate: 128, SampleRate: 48000, Channels: 2},
		}

		lowerType := strings.ToLower(audioFile.FileType)
		if meta, ok := formatDefaults[lowerType]; ok {
			if audioFile.Bitrate == 0 {
				audioFile.Bitrate = meta.Bitrate
			}
			if audioFile.SampleRate == 0 {
				audioFile.SampleRate = meta.SampleRate
			}
			if audioFile.Channels == 0 {
				audioFile.Channels = meta.Channels
			}
		} else {
			if audioFile.Bitrate == 0 {
				audioFile.Bitrate = 128
			}
			if audioFile.SampleRate == 0 {
				audioFile.SampleRate = 44100
			}
			if audioFile.Channels == 0 {
				audioFile.Channels = 2
			}
		}
	}

	if audioFile.Duration == 0 && audioFile.FileSize > 0 && audioFile.Bitrate > 0 {
		bits := audioFile.FileSize * 8
		bps := int64(audioFile.Bitrate * 1000)
		durationSec := float64(bits) / float64(bps)
		audioFile.Duration = time.Duration(durationSec * float64(time.Second))
	}
}

func (e *Extractor) SetExtractCoverArt(extract bool) {
	e.extractCoverArt = extract
}

func (e *Extractor) SetMaxCoverSize(size int) {
	if size > 0 {
		e.maxCoverSize = size
	}
}

func (e *Extractor) ExtractFromFiles(filePaths []string) ([]*AudioFile, []error) {
	var results []*AudioFile
	var errors []error

	for _, filePath := range filePaths {
		audioFile, err := e.ExtractFromFile(filePath)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		results = append(results, audioFile)
	}

	return results, errors
}
