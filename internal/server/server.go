package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hoppxi/bpv/internal/cache"
	"github.com/hoppxi/bpv/internal/daemon"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/xdg"
)

type Server struct {
	port      int
	musicDir  string
	server    *http.Server
	client    *daemon.Client
	lib       *cache.CachedLibrary
	startTime time.Time
}

func NewServer(port int, musicDir string) *Server {
	return &Server{
		port:      port,
		musicDir:  musicDir,
		startTime: time.Now(),
	}
}

func (s *Server) Start() error {
	client, err := daemon.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to daemon: %w", err)
	}
	s.client = client

	lib, err := client.GetLibrary(s.musicDir)
	if err != nil {
		return fmt.Errorf("failed to load library: %w", err)
	}
	s.lib = lib

	mux := http.NewServeMux()
	s.setupRoutes(mux)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      Logger(Recovery(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Log.Debug("Server started on http://localhost:%d", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.client != nil {
		s.client.Close()
	}
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func resolveWebDir() string {
	if dir := os.Getenv("BPV_WEB_DIR"); dir != "" {
		return dir
	}

	xdgDir := filepath.Join(xdg.DataDir(), "dist")
	if _, err := os.Stat(xdgDir); err == nil {
		return xdgDir
	}

	return "./web/dist"
}

func (s *Server) setupRoutes(mux *http.ServeMux) {
	musicHandler := http.StripPrefix("/files/", http.FileServer(http.Dir(s.musicDir)))
	mux.Handle("/files/", s.enableCORS(musicHandler))

	webDir := resolveWebDir()
	if s.serveWebApp(mux, webDir) {
		logger.Log.Info("Serving Vue app from: %s", webDir)
	} else {
		logger.Log.Error("Vue app not found at: %s, serving placeholder", webDir)
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
	mux.HandleFunc("/api/composers", s.handleComposers)
	mux.HandleFunc("/api/composer/", s.handleComposer)
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/base-path", s.handleBaseFilePath)
	mux.HandleFunc("/api/debug", s.handleDebug)
	mux.HandleFunc("/api/cover/", s.handleCoverArt)
	mux.HandleFunc("/api/favorites", s.handleFavorites)
	mux.HandleFunc("/api/stats/play", s.handleRecordPlay)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/queue", s.handleQueue)
	mux.HandleFunc("/api/settings", s.handleSettingsAPI)
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

	var coverArt, coverMime string
	if s.client != nil {
		var err error
		coverArt, coverMime, err = s.client.GetCoverArt(fullPath)
		if err != nil || coverArt == "" {
			logger.Log.ErrorP("Server", "%s", err)
			http.Error(w, "No cover art found", http.StatusNotFound)
			return
		}
	} else {
		extractor := metadata.NewExtractor()
		audioFile, err := extractor.ExtractFromFile(fullPath)
		if err != nil || audioFile.CoverArt == "" {
			http.Error(w, "No cover art found", http.StatusNotFound)
			return
		}
		coverArt = audioFile.CoverArt
		coverMime = audioFile.CoverArtMime
	}

	coverData, err := base64.StdEncoding.DecodeString(coverArt)
	if err != nil {
		http.Error(w, "Invalid cover art", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", coverMime)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(coverData)
}

func (s *Server) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getFiles() []metadata.AudioFile {
	if s.lib == nil {
		return nil
	}
	return s.lib.Files
}
