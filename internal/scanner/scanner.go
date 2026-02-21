package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/metadata"
)

type ScanProgress struct {
	Current   int       `json:"current"`
	Total     int       `json:"total"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ScanResult struct {
	TotalFiles int                  `json:"total_files"`
	AudioFiles int                  `json:"audio_files"`
	Artists    map[string]int       `json:"artists"`
	Albums     map[string]int       `json:"albums"`
	Genres     map[string]int       `json:"genres"`
	Composers  map[string]int       `json:"composers"`
	Files      []metadata.AudioFile `json:"files"`
	Duration   time.Duration        `json:"duration"`
	Errors     []string             `json:"errors"`
}

type Scanner struct {
	fileWalker        *FileWalker
	metadataExtractor *metadata.Extractor
	progressChan      chan ScanProgress
}

func NewScanner() *Scanner {
	progressChan := make(chan ScanProgress, 100)
	return &Scanner{
		fileWalker:        NewFileWalker(progressChan),
		metadataExtractor: metadata.NewExtractor(),
		progressChan:      progressChan,
	}
}

func (s *Scanner) ScanLibrary(rootPath string) (*ScanResult, error) {
	startTime := time.Now()

	audioFilePaths, err := s.fileWalker.WalkDirectory(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	s.sendProgress(len(audioFilePaths), len(audioFilePaths), "Extracting metadata...")

	result := &ScanResult{
		Artists:   make(map[string]int),
		Albums:    make(map[string]int),
		Genres:    make(map[string]int),
		Composers: make(map[string]int),
		Errors:    []string{},
	}

	extractor := metadata.NewExtractor()
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, filePath := range audioFilePaths {
		wg.Add(1)

		go func(path string, index int) {
			defer wg.Done()

			audioFile, err := extractor.ExtractFromFile(path)
			if err != nil {
				mu.Lock()
				result.Errors = append(result.Errors,
					fmt.Sprintf("Failed to extract metadata from %s: %v", filepath.Base(path), err))
				mu.Unlock()
				return
			}

			mu.Lock()
			result.Files = append(result.Files, *audioFile)

			if audioFile.Artist != "" && audioFile.Artist != "Unknown Artist" {
				result.Artists[audioFile.Artist]++
			}
			if audioFile.Album != "" && audioFile.Album != "Unknown Album" {
				result.Albums[audioFile.Album]++
			}
			if audioFile.Genre != "" && audioFile.Genre != "Unknown Genre" {
				result.Genres[audioFile.Genre]++
			}
			if audioFile.Composer != "" && audioFile.Composer != "Unknown Composer" {
				result.Composers[audioFile.Composer]++
			}
			mu.Unlock()

			if index%10 == 0 || index == len(audioFilePaths)-1 {
				s.sendProgress(index+1, len(audioFilePaths),
					fmt.Sprintf("Processed %d/%d files", index+1, len(audioFilePaths)))
			}
		}(filePath, i)
	}

	wg.Wait()

	result.TotalFiles = len(audioFilePaths)
	result.AudioFiles = len(result.Files)
	result.Duration = time.Since(startTime)

	s.sendProgress(len(audioFilePaths), len(audioFilePaths),
		fmt.Sprintf("Scan completed: %d audio files found in %v", result.AudioFiles, result.Duration.Round(time.Millisecond)))

	return result, nil
}

func (s *Scanner) GetProgressChannel() <-chan ScanProgress {
	return s.progressChan
}

func (s *Scanner) sendProgress(current, total int, message string) {
	select {
	case s.progressChan <- ScanProgress{
		Current:   current,
		Total:     total,
		Message:   message,
		Timestamp: time.Now(),
	}:
	default:
	}
}

func (s *Scanner) QuickScan(rootPath string) ([]metadata.AudioFile, error) {
	audioFilePaths, err := s.fileWalker.WalkDirectory(rootPath)
	if err != nil {
		return nil, err
	}

	var files []metadata.AudioFile
	for _, path := range audioFilePaths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		files = append(files, metadata.AudioFile{
			FilePath: path,
			FileName: filepath.Base(path),
			FileSize: info.Size(),
			Modified: info.ModTime(),
			FileType: strings.TrimPrefix(filepath.Ext(path), "."),
			Title:    filepath.Base(path),
			Artist:   "Unknown Artist",
			Album:    "Unknown Album",
			Composer: "Unknown Composer",
		})
	}

	return files, nil
}
