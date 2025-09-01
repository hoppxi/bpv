package server

import (
	"encoding/base64"
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

type Server struct {
	port     int
	musicDir string
	server   *http.Server
	scanner  *scanner.Scanner
	library  *scanner.ScanResult
	lastScan time.Time
	startTime time.Time
}

// NewServer creates a new BPV server instance
func NewServer(port int, musicDir string) *Server {
	s := &Server{
		port:     port,
		musicDir: musicDir,
		scanner:  scanner.NewScanner(),
		startTime: time.Now(),
	}
	return s
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	
	s.setupRoutes(mux)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      Logger(Recovery(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	logger.Log.Info("Server starting on http://localhost:%d", s.port)
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes(mux *http.ServeMux) {
	// Serve music files with CORS headers
	musicHandler := http.StripPrefix("/files/", http.FileServer(http.Dir(s.musicDir)))
	mux.Handle("/files/", s.enableCORS(musicHandler))
	
	// Serve React app (if built)
	webDir := "./web/dist"
	if s.serveWebApp(mux, webDir) {
		logger.Log.Info("Serving React app from: %s", webDir)
	} else {
		logger.Log.Error("React app not found at: %s, serving placeholder", webDir)
		mux.HandleFunc("/", s.handleRoot)
	}
	
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/library", s.handleLibrary)
	mux.HandleFunc("/api/metadata/", s.handleMetadata)
	mux.HandleFunc("/api/scan", s.handleScan)
	mux.HandleFunc("/api/scan-simple", s.handleScanSimple)
	mux.HandleFunc("/api/scan/progress", s.handleScanProgress)
	mux.HandleFunc("/api/artists", s.handleArtists)
	mux.HandleFunc("/api/artist/", s.handleArtist)
	mux.HandleFunc("/api/albums", s.handleAlbums)
	mux.HandleFunc("/api/album/", s.handleAlbum)
	mux.HandleFunc("/api/genres", s.handleGenres)
	mux.HandleFunc("/api/genre/", s.handleGenre)
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/debug", s.handleDebug)

	mux.HandleFunc("/api/cover/", s.handleCoverArt)
}

func (s *Server) serveWebApp(mux *http.ServeMux, webDir string) bool {
	absPath, err := filepath.Abs(webDir)
	if err != nil {
		return false
	}
	
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return false
	}
	
	fs := http.FileServer(http.Dir(absPath))
	mux.Handle("/", fs)
	return true
}

func (s *Server) handleCoverArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/api/cover/")
	if filePath == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(s.musicDir, filePath)
	extractor := metadata.NewExtractor()
	audioFile, err := extractor.ExtractFromFile(fullPath)
	if err != nil || audioFile.CoverArt == "" {
		http.Error(w, "No cover art found", http.StatusNotFound)
		return
	}

	coverData, err := base64.StdEncoding.DecodeString(audioFile.CoverArt)
	if err != nil {
		http.Error(w, "Invalid cover art", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", audioFile.CoverArtMime)
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.Write(coverData)
}

func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}