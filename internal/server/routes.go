package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/store"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Port     int    `json:"port"`
	MusicDir string `json:"music_dir"`
	Version  string `json:"version"`
}

type LibraryResponse struct {
	Status     string               `json:"status"`
	MusicDir   string               `json:"music_dir"`
	TotalFiles int                  `json:"total_files"`
	AudioFiles int                  `json:"audio_files"`
	Artists    map[string]int       `json:"artists"`
	Albums     map[string]int       `json:"albums"`
	Genres     map[string]int       `json:"genres"`
	Composers  map[string]int       `json:"composers"`
	Files      []metadata.AudioFile `json:"files"`
	ScanTime   string               `json:"scan_time"`
	Errors     []string             `json:"errors,omitempty"`
}

type MetadataResponse struct {
	Status string             `json:"status"`
	File   metadata.AudioFile `json:"file"`
}

type ScanRequest struct {
	FullScan      bool `json:"full_scan"`
	ExtractCovers bool `json:"extract_covers"`
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/text")
	fmt.Fprintf(w, `could not find web build. set BPV_DIST_DIR env var and try again`)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:   "ok",
		Port:     s.port,
		MusicDir: s.musicDir,
		Version:  "0.3.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLibrary(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Error("Panic in handleLibrary: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.lib == nil {
		if s.client != nil {
			lib, err := s.client.GetLibrary(s.musicDir)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to load library: %v", err), http.StatusInternalServerError)
				return
			}
			s.lib = lib
		} else {
			http.Error(w, "Library not available", http.StatusServiceUnavailable)
			return
		}
	}

	files := s.getFiles()

	response := LibraryResponse{
		Status:     "ok",
		MusicDir:   s.musicDir,
		TotalFiles: s.lib.FileCount,
		AudioFiles: s.lib.FileCount,
		Artists:    s.lib.Artists,
		Albums:     s.lib.Albums,
		Genres:     s.lib.Genres,
		Composers:  s.lib.Composers,
		Files:      files,
		ScanTime:   s.lib.ScanTime.Format(time.RFC3339),
		Errors:     s.lib.Errors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/api/metadata/")
	if filePath == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(s.musicDir, filePath)

	extractor := metadata.NewExtractor()
	audioFile, err := extractor.ExtractFromFile(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract metadata: %v", err), http.StatusInternalServerError)
		return
	}

	response := MetadataResponse{
		Status: "ok",
		File:   *audioFile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		if s.client != nil {
			lib, err := s.client.Scan(s.musicDir)
			if err != nil {
				logger.Log.Error("Scan failed: %v", err)
				return
			}
			s.lib = lib
			logger.Log.Success("Scan completed: %d audio files found", lib.FileCount)
		}
	}()

	response := map[string]any{
		"status":  "started",
		"message": "Library scan initiated via daemon",
		"path":    s.musicDir,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleScanProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, _ := json.Marshal(map[string]any{
		"current": 100,
		"total":   100,
		"message": "Scan complete (via daemon cache)",
	})
	fmt.Fprintf(w, "data: %s\n\n", data)
	w.(http.Flusher).Flush()
}

func (s *Server) handleArtists(w http.ResponseWriter, r *http.Request) {
	if s.lib == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	artists := make([]map[string]any, 0, len(s.lib.Artists))
	for artist, count := range s.lib.Artists {
		artists = append(artists, map[string]any{
			"name":  artist,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"artists": artists,
		"total":   len(artists),
	})
}

func (s *Server) handleArtist(w http.ResponseWriter, r *http.Request) {
	artistName := strings.TrimPrefix(r.URL.Path, "/api/artist/")
	if artistName == "" {
		http.Error(w, "Artist name required", http.StatusBadRequest)
		return
	}

	files := s.getFiles()
	var artistSongs []metadata.AudioFile
	for _, file := range files {
		if file.Artist == artistName {
			artistSongs = append(artistSongs, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"artist": artistName,
		"songs":  artistSongs,
		"count":  len(artistSongs),
	})
}

func (s *Server) handleAlbums(w http.ResponseWriter, r *http.Request) {
	if s.lib == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	albums := make([]map[string]any, 0, len(s.lib.Albums))
	for album, count := range s.lib.Albums {
		albums = append(albums, map[string]any{
			"name":  album,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"albums": albums,
		"total":  len(albums),
	})
}

func (s *Server) handleAlbum(w http.ResponseWriter, r *http.Request) {
	albumName := strings.TrimPrefix(r.URL.Path, "/api/album/")
	if albumName == "" {
		http.Error(w, "Album name required", http.StatusBadRequest)
		return
	}

	files := s.getFiles()
	var albumSongs []metadata.AudioFile
	for _, file := range files {
		if file.Album == albumName {
			albumSongs = append(albumSongs, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"album":  albumName,
		"songs":  albumSongs,
		"count":  len(albumSongs),
	})
}

func (s *Server) handleGenres(w http.ResponseWriter, r *http.Request) {
	if s.lib == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	genres := make([]map[string]any, 0, len(s.lib.Genres))
	for genre, count := range s.lib.Genres {
		genres = append(genres, map[string]any{
			"name":  genre,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"genres": genres,
		"total":  len(genres),
	})
}

func (s *Server) handleGenre(w http.ResponseWriter, r *http.Request) {
	genreName := strings.TrimPrefix(r.URL.Path, "/api/genre/")
	if genreName == "" {
		http.Error(w, "Genre name required", http.StatusBadRequest)
		return
	}

	files := s.getFiles()
	var genreSongs []metadata.AudioFile
	for _, file := range files {
		if file.Genre == genreName {
			genreSongs = append(genreSongs, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"genre":  genreName,
		"songs":  genreSongs,
		"count":  len(genreSongs),
	})
}

func (s *Server) handleComposers(w http.ResponseWriter, r *http.Request) {
	if s.lib == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	composers := make([]map[string]any, 0, len(s.lib.Composers))
	for composer, count := range s.lib.Composers {
		composers = append(composers, map[string]any{
			"name":  composer,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"composers": composers,
		"total":     len(composers),
	})
}

func (s *Server) handleComposer(w http.ResponseWriter, r *http.Request) {
	composerName := strings.TrimPrefix(r.URL.Path, "/api/composer/")
	if composerName == "" {
		http.Error(w, "Composer name required", http.StatusBadRequest)
		return
	}

	files := s.getFiles()
	var composerSongs []metadata.AudioFile
	for _, file := range files {
		if file.Composer == composerName {
			composerSongs = append(composerSongs, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "ok",
		"composer": composerName,
		"songs":    composerSongs,
		"count":    len(composerSongs),
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if s.lib == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	query = strings.ToLower(query)
	files := s.getFiles()
	var results []metadata.AudioFile

	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Title), query) ||
			strings.Contains(strings.ToLower(file.Artist), query) ||
			strings.Contains(strings.ToLower(file.Album), query) ||
			strings.Contains(strings.ToLower(file.Genre), query) {
			results = append(results, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

func (s *Server) handleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(s.musicDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read music directory: %v", err), http.StatusInternalServerError)
		return
	}

	fileCounts := make(map[string]int)
	for _, file := range files {
		if file.IsDir() {
			fileCounts["directories"]++
		} else {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			fileCounts[ext]++
		}
	}

	response := map[string]any{
		"status":        "ok",
		"music_dir":     s.musicDir,
		"file_counts":   fileCounts,
		"total_files":   len(files),
		"server_uptime": time.Since(s.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleScanSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		if s.client != nil {
			lib, err := s.client.Scan(s.musicDir)
			if err != nil {
				logger.Log.Error("Simple scan failed: %v", err)
				return
			}
			s.lib = lib
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "Library scan initiated via daemon",
	})
}

func (s *Server) handleBaseFilePath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]any{
		"status":    "ok",
		"base_path": s.musicDir,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if s.client != nil {
			favs, err := s.client.GetFavorites()
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{
					"status":    "ok",
					"favorites": []string{},
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]any{
				"status":    "ok",
				"favorites": favs,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"status":    "ok",
				"favorites": []string{},
			})
		}

	case http.MethodPost:
		var req struct {
			FilePath string `json:"file_path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FilePath == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if s.client != nil {
			s.client.AddFavorite(req.FilePath)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	case http.MethodDelete:
		var req struct {
			FilePath string `json:"file_path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FilePath == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if s.client != nil {
			s.client.RemoveFavorite(req.FilePath)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRecordPlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FilePath == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if s.client != nil {
		s.client.RecordPlay(req.FilePath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var stats map[string]int
	if s.client != nil {
		var err error
		stats, err = s.client.GetStats()
		if err != nil {
			stats = make(map[string]int)
		}
	} else {
		stats = make(map[string]int)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"plays":  stats,
	})
}

func (s *Server) handleQueue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if s.client != nil {
			q, err := s.client.GetQueue()
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{
					"status": "ok",
					"queue":  store.QueueState{},
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"queue":  q,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"queue":  store.QueueState{},
			})
		}

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		var q store.QueueState
		if err := json.Unmarshal(body, &q); err != nil {
			http.Error(w, "Invalid queue JSON", http.StatusBadRequest)
			return
		}
		if s.client != nil {
			s.client.SaveQueue(&q)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSettingsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if s.client != nil {
			settings, err := s.client.GetSettings()
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{
					"status":   "ok",
					"settings": store.Settings{},
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]any{
				"status":   "ok",
				"settings": settings,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{
				"status":   "ok",
				"settings": store.Settings{},
			})
		}

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		var settings store.Settings
		if err := json.Unmarshal(body, &settings); err != nil {
			http.Error(w, "Invalid settings JSON", http.StatusBadRequest)
			return
		}
		if s.client != nil {
			s.client.SaveSettings(&settings)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
