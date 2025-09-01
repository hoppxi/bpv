package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
	var audioFiles []string
	var mu sync.Mutex
	var totalFiles int

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			totalFiles++
		}
		return nil
	})

	fw.sendProgress(0, totalFiles, "Scanning files...")
	
	startTime := time.Now()
	processedFiles := 0

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if fw.audioExtensions[ext] {
			mu.Lock()
			audioFiles = append(audioFiles, path)
			mu.Unlock()
		}

		processedFiles++
		if processedFiles%100 == 0 || processedFiles == totalFiles {
			elapsed := time.Since(startTime)
			fw.sendProgress(processedFiles, totalFiles, 
				fmt.Sprintf("Scanned %d/%d files (%v elapsed)", processedFiles, totalFiles, elapsed.Round(time.Second)))
		}

		return nil
	})

	fw.sendProgress(totalFiles, totalFiles, "Scan completed")
	return audioFiles, err
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
			// Don't block if channel is full
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