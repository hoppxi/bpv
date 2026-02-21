package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/logger"
)

type FileWalker struct {
	audioExtensions map[string]bool
	progressChan    chan<- ScanProgress
}

func NewFileWalker(progressChan chan<- ScanProgress) *FileWalker {
	return &FileWalker{
		audioExtensions: map[string]bool{
			".mp3":  true,
			".flac": true,
			".wav":  true,
			".ogg":  true,
			".m4a":  true,
			".aac":  true,
			".wma":  true,
			".opus": true,
			".alac": true,
			".aiff": true,
			".dsf":  true,
		},
		progressChan: progressChan,
	}
}

func (fw *FileWalker) WalkDirectory(rootPath string) ([]string, error) {
	logger.Log.Debug("Starting directory walk: %s", rootPath)

	var audioFiles []string
	var mu sync.Mutex
	var totalFiles int

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				logger.Log.Error("Permission denied: %s", path)
				return nil
			}
			logger.Log.Error("Error accessing %s: %v", path, err)
			return err
		}
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error counting files: %v", err)
	}

	logger.Log.Debug("Total files to scan: %d", totalFiles)
	fw.sendProgress(0, totalFiles, "Scanning files...")

	startTime := time.Now()
	processedFiles := 0

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				logger.Log.Error("Permission denied: %s", path)
				return nil
			}
			logger.Log.Error("Error accessing %s: %v", path, err)
			return err
		}

		if d.IsDir() {
			if d.Name() != "." && d.Name() != ".." && strings.HasPrefix(d.Name(), ".") {
				logger.Log.Debug("Skipping hidden directory: %s", path)
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if fw.audioExtensions[ext] {
			mu.Lock()
			audioFiles = append(audioFiles, path)
			mu.Unlock()
			logger.Log.Debug("Found audio file: %s", path)
		}

		processedFiles++
		if processedFiles%100 == 0 || processedFiles == totalFiles {
			elapsed := time.Since(startTime)
			logger.Log.Debug("Scan progress: %d/%d files", processedFiles, totalFiles)
			fw.sendProgress(processedFiles, totalFiles,
				fmt.Sprintf("Scanned %d/%d files (%v elapsed)", processedFiles, totalFiles, elapsed.Round(time.Second)))
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	logger.Log.Debug("Directory walk completed: found %d audio files", len(audioFiles))
	fw.sendProgress(totalFiles, totalFiles, "Scan completed")
	return audioFiles, nil
}

func (fw *FileWalker) sendProgress(current, total int, message string) {
	if fw.progressChan != nil {
		select {
		case fw.progressChan <- ScanProgress{
			Current:   current,
			Total:     total,
			Message:   message,
			Timestamp: time.Now(),
		}:
		default:
		}
	}
}

func (fw *FileWalker) IsAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return fw.audioExtensions[ext]
}

func (fw *FileWalker) GetAudioExtensions() []string {
	extensions := make([]string, 0, len(fw.audioExtensions))
	for ext := range fw.audioExtensions {
		extensions = append(extensions, ext)
	}
	return extensions
}
