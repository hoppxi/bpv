package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/scanner"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Port      int    `json:"port"`
	MusicDir  string `json:"music_dir"`
	Version   string `json:"version"`
}

type LibraryResponse struct {
	Status      string               `json:"status"`
	TotalFiles  int                  `json:"total_files"`
	AudioFiles  int                  `json:"audio_files"`
	Artists     map[string]int       `json:"artists"`
	Albums      map[string]int       `json:"albums"`
	Genres      map[string]int       `json:"genres"`
	Files       []metadata.AudioFile `json:"files"`
	ScanTime    string               `json:"scan_time"`
	Errors      []string             `json:"errors,omitempty"`
}

type MetadataResponse struct {
	Status   string             `json:"status"`
	File     metadata.AudioFile `json:"file"`
}

type ScanRequest struct {
	FullScan    bool `json:"full_scan"`
	ExtractCovers bool `json:"extract_covers"`
}

// handleRoot serves a simple placeholder page if React app isn't built
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
	<title>BPV Music Player</title>
	<style>
		body { 
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
			margin: 0; 
			padding: 20px; 
			background: linear-gradient(135deg, #000a35ff 0%%, #230044ff 100%%);
			color: white;
			min-height: 100vh;
		}
		.container { 
			max-width: 600px; 
			margin: 0 auto; 
			background: rgba(255, 255, 255, 0.1);
			padding: 30px;
			border-radius: 15px;
			backdrop-filter: blur(10px);
			border: 1px solid rgba(255, 255, 255, 0.2);
		}
		h1 { 
			text-align: center; 
			margin-bottom: 30px;
			font-size: 2.5em;
			text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
		}
		.status { 
			background: rgba(255, 255, 255, 0.2); 
			padding: 20px; 
			border-radius: 10px; 
			margin: 20px 0; 
		}
		.endpoints { 
			margin: 30px 0; 
		}
		.endpoint { 
			background: rgba(255, 255, 255, 0.15); 
			padding: 15px; 
			margin: 10px 0; 
			border-radius: 8px; 
			border-left: 4px solid #ff6b6b;
		}
		.endpoint:hover {
			background: rgba(255, 255, 255, 0.25);
			transform: translateX(5px);
			transition: all 0.3s ease;
		}
		a { 
			color: #ff6b6b; 
			text-decoration: none; 
			font-weight: bold;
		}
		a:hover { 
			text-decoration: underline; 
		}
		.code {
			background: rgba(0, 0, 0, 0.3);
			padding: 2px 6px;
			border-radius: 4px;
			font-family: 'Monaco', 'Menlo', monospace;
			font-size: 0.9em;
		}
		.footer {
			text-align: center;
			margin-top: 40px;
			opacity: 0.8;
			font-size: 0.9em;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>BPV Music Player</h1>
		
		<div class="status">
			<p><strong>Status:</strong> Server is running</p>
			<p><strong>Music Directory:</strong> <code>%s</code></p>
			<p><strong>Port:</strong> <code>%d</code></p>
			<p><strong>API Base URL:</strong> <code>http://%s/api/</code></p>
		</div>

		<div class="endpoints">
			<h3>Available API Endpoints:</h3>
			
			<div class="endpoint">
				<strong><a href="/api/health" target="_blank">/api/health</a></strong>
				<p>Server health check and status information</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/library" target="_blank">/api/library</a></strong>
				<p>Get complete music library with metadata</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/scan" target="_blank">/api/scan</a></strong>
				<p>Initiate a library scan (POST request)</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/artists" target="_blank">/api/artists</a></strong>
				<p>List all artists in the library</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/albums" target="_blank">/api/albums</a></strong>
				<p>List all albums in the library</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/genres" target="_blank">/api/genres</a></strong>
				<p>List all genres in the library</p>
			</div>

			<div class="endpoint">
				<strong><a href="/api/search?q=rock" target="_blank">/api/search?q=rock</a></strong>
				<p>Search through music library</p>
			</div>
		</div>

		<div class="status">
			<h3>Next Steps:</h3>
			<p>To use the React web interface:</p>
			<ol>
				<li>Navigate to the <span class="code">web/</span> directory</li>
				<li>Run <span class="code">npm install</span> to install dependencies</li>
				<li>Run <span class="code">npm run build</span> to build the React app</li>
				<li>Refresh this page to see the music player interface</li>
			</ol>
		</div>

		<div class="footer">
			<p>BPV Music Player. Built with Go + React. <a href="https://github.com/hoppxi/bpv" target="_blank">GitHub</a></p>
		</div>
	</div>
</body>
</html>
	`, s.musicDir, s.port, r.Host)
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
		Version:  "0.1.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLibrary(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if err := recover(); err != nil {
           logger.Log.Error("Panic in handleLibrary: %v", err)
            http.Error(w, "Internal server error during library processing", http.StatusInternalServerError)
        }
    }()

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if s.library == nil {
        fileWalker := scanner.NewFileWalker(nil)
        audioFilePaths, err := fileWalker.WalkDirectory(s.musicDir)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error scanning directory: %v", err), http.StatusInternalServerError)
            return
        }

        extractor := metadata.NewExtractor()
        var files []metadata.AudioFile
        artists := make(map[string]int)
        albums := make(map[string]int)
        genres := make(map[string]int)
        var errors []string

       logger.Log.Debug("Scanning %d audio files...", len(audioFilePaths))

        for i, filePath := range audioFilePaths {
            audioFile, err := extractor.ExtractFromFile(filePath)
            if err != nil {
                errors = append(errors, fmt.Sprintf("Failed to extract metadata from %s: %v", 
                    filepath.Base(filePath), err))
                continue
            }

            if audioFile == nil {
                errors = append(errors, fmt.Sprintf("Extracted nil metadata for %s", filePath))
                continue
            }

            files = append(files, *audioFile)
            if audioFile.Artist != "" && audioFile.Artist != "Unknown Artist" {
                artists[audioFile.Artist]++
            }
            if audioFile.Album != "" && audioFile.Album != "Unknown Album" {
                albums[audioFile.Album]++
            }
            if audioFile.Genre != "" && audioFile.Genre != "Unknown Genre" {
                genres[audioFile.Genre]++
            }

            if (i+1)%100 == 0 || i+1 == len(audioFilePaths) {
               logger.Log.Debug("Processed %d/%d files", i+1, len(audioFilePaths))
            }
        }

       logger.Log.Debug("Library scan completed: %d files, %d errors", len(files), len(errors))

        s.library = &scanner.ScanResult{
            TotalFiles: len(audioFilePaths),
            AudioFiles: len(files),
            Artists:    artists,
            Albums:     albums,
            Genres:     genres,
            Files:      files,
            Duration:   time.Duration(0),
            Errors:     errors,
        }
        s.lastScan = time.Now()
    }

    if s.library == nil {
        http.Error(w, "Library not initialized", http.StatusInternalServerError)
        return
    }

    response := LibraryResponse{
        Status:     "ok",
        TotalFiles: s.library.TotalFiles,
        AudioFiles: s.library.AudioFiles,
        Artists:    s.library.Artists,
        Albums:     s.library.Albums,
        Genres:     s.library.Genres,
        Files:      s.library.Files,
        ScanTime:   s.lastScan.Format(time.RFC3339),
        Errors:     s.library.Errors,
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

	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	extractor := metadata.NewExtractor()
	extractor.SetExtractCoverArt(req.ExtractCovers)

	go func() {
		result, err := s.scanner.ScanLibrary(s.musicDir)
		if err != nil {
			logger.Log.Error("Scan failed: %v\n", err)
			return
		}

		s.library = result
		s.lastScan = time.Now()
		logger.Log.Success("Scan completed: %d audio files found\n", result.AudioFiles)
	}()

	response := map[string]any{
		"status":  "started",
		"message": "Library scan initiated",
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

	progressChan := s.scanner.GetProgressChannel()

	for {
		select {
		case progress := <-progressChan:
			data, _ := json.Marshal(progress)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		
		case <-r.Context().Done():
			return
		
		case <-time.After(30 * time.Second):
	
			fmt.Fprintf(w, ": keep-alive\n\n")
			w.(http.Flusher).Flush()
		}
	}
}

func (s *Server) handleArtists(w http.ResponseWriter, r *http.Request) {
	if s.library == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	artists := make([]map[string]any, 0, len(s.library.Artists))
	for artist, count := range s.library.Artists {
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

	var artistSongs []metadata.AudioFile
	for _, file := range s.library.Files {
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
	if s.library == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	albums := make([]map[string]any, 0, len(s.library.Albums))
	for album, count := range s.library.Albums {
		albums = append(albums, map[string]any{
			"name":  album,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"albums":  albums,
		"total":   len(albums),
	})
}

func (s *Server) handleAlbum(w http.ResponseWriter, r *http.Request) {
	albumName := strings.TrimPrefix(r.URL.Path, "/api/album/")
	if albumName == "" {
		http.Error(w, "Album name required", http.StatusBadRequest)
		return
	}

	var albumSongs []metadata.AudioFile
	for _, file := range s.library.Files {
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
	if s.library == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	genres := make([]map[string]any, 0, len(s.library.Genres))
	for genre, count := range s.library.Genres {
		genres = append(genres, map[string]any{
			"name":  genre,
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"genres":  genres,
		"total":   len(genres),
	})
}

func (s *Server) handleGenre(w http.ResponseWriter, r *http.Request) {
	genreName := strings.TrimPrefix(r.URL.Path, "/api/genre/")
	if genreName == "" {
		http.Error(w, "Genre name required", http.StatusBadRequest)
		return
	}

	var genreSongs []metadata.AudioFile
	for _, file := range s.library.Files {
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

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if s.library == nil {
		http.Error(w, "Library not scanned yet", http.StatusNotFound)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	query = strings.ToLower(query)
	var results []metadata.AudioFile

	for _, file := range s.library.Files {
		if strings.Contains(strings.ToLower(file.Title), query) ||
			strings.Contains(strings.ToLower(file.Artist), query) ||
			strings.Contains(strings.ToLower(file.Album), query) ||
			strings.Contains(strings.ToLower(file.Genre), query) {
			results = append(results, file)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "ok",
		"query":    query,
		"results":  results,
		"count":    len(results),
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
		"status":       "ok",
		"music_dir":    s.musicDir,
		"file_counts":  fileCounts,
		"total_files":  len(files),
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

	s.library = nil

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "Library scan initiated",
	})
}