package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/cache"
	"github.com/hoppxi/bpv/internal/logger"
	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/scanner"
	"github.com/hoppxi/bpv/internal/store"
	"github.com/hoppxi/bpv/internal/xdg"
)

type Request struct {
	Action   string `json:"action"`
	Dir      string `json:"dir,omitempty"`
	FilePath string `json:"file_path,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
}

type Response struct {
	OK        bool                 `json:"ok"`
	Error     string               `json:"error,omitempty"`
	Library   *cache.CachedLibrary `json:"library,omitempty"`
	Favorites []string             `json:"favorites,omitempty"`
	Settings  *store.Settings      `json:"settings,omitempty"`
	Stats     map[string]int       `json:"stats,omitempty"`
	CoverArt  string               `json:"cover_art,omitempty"`
	CoverMime string               `json:"cover_mime,omitempty"`
	IsFav     bool                 `json:"is_fav,omitempty"`
	Queue     *store.QueueState    `json:"queue,omitempty"`
}

type Daemon struct {
	store    *store.Store
	cache    *cache.Cache
	listener net.Listener
	mu       sync.Mutex
	scanning map[string]bool
}

func SocketPath() string {
	return xdg.SocketPath()
}

func NewDaemon() (*Daemon, error) {
	st, err := store.NewStore()
	if err != nil {
		return nil, logger.Log.Error("failed to create store: %w", err)
	}

	ch, err := cache.NewCache()
	if err != nil {
		return nil, logger.Log.Error("failed to create cache: %w", err)
	}

	return &Daemon{
		store:    st,
		cache:    ch,
		scanning: make(map[string]bool),
	}, nil
}

func (d *Daemon) Start() error {
	sockPath := SocketPath()

	os.Remove(sockPath)
	os.MkdirAll(filepath.Dir(sockPath), 0755)

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", sockPath, err)
	}
	d.listener = listener
	os.Chmod(sockPath, 0700)

	logger.Log.Info("Daemon started")
	logger.Log.Info("Socket: %s", sockPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return nil
		}
		go d.handleConnection(conn)
	}
}

func (d *Daemon) Stop() {
	if d.listener != nil {
		d.listener.Close()
		os.Remove(SocketPath())
	}
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0), 64*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			d.sendResponse(conn, Response{OK: false, Error: "invalid request: " + err.Error()})
			continue
		}

		resp := d.handleRequest(req)
		d.sendResponse(conn, resp)
	}
}

func (d *Daemon) sendResponse(conn net.Conn, resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	data = append(data, '\n')
	conn.Write(data)
}

func (d *Daemon) handleRequest(req Request) Response {
	switch req.Action {
	case "ping":
		return Response{OK: true}
	case "library":
		return d.handleLibrary(req.Dir)
	case "scan":
		return d.handleScan(req.Dir)
	case "cover-art":
		return d.handleCoverArt(req.FilePath)
	case "get-favorites":
		return d.handleGetFavorites()
	case "add-favorite":
		return d.handleAddFavorite(req.FilePath)
	case "remove-favorite":
		return d.handleRemoveFavorite(req.FilePath)
	case "is-favorite":
		return d.handleIsFavorite(req.FilePath)
	case "get-settings":
		return d.handleGetSettings()
	case "save-settings":
		return d.handleSaveSettings(req.Value)
	case "get-stats":
		return d.handleGetStats()
	case "record-play":
		return d.handleRecordPlay(req.FilePath)
	case "get-queue":
		return d.handleGetQueue()
	case "save-queue":
		return d.handleSaveQueue(req.Value)
	default:
		return Response{OK: false, Error: "unknown action: " + req.Action}
	}
}

func (d *Daemon) handleLibrary(dir string) Response {
	if dir == "" {
		return Response{OK: false, Error: "dir is required"}
	}

	lib := d.cache.Load(dir)
	if lib != nil {
		return Response{OK: true, Library: lib}
	}

	return d.scanAndCache(dir)
}

func (d *Daemon) handleScan(dir string) Response {
	if dir == "" {
		return Response{OK: false, Error: "dir is required"}
	}
	d.cache.Invalidate(dir)
	return d.scanAndCache(dir)
}

func (d *Daemon) scanAndCache(dir string) Response {
	d.mu.Lock()
	if d.scanning[dir] {
		d.mu.Unlock()
		return Response{OK: false, Error: "scan already in progress for " + dir}
	}
	d.scanning[dir] = true
	d.mu.Unlock()

	defer func() {
		d.mu.Lock()
		delete(d.scanning, dir)
		d.mu.Unlock()
	}()

	sc := scanner.NewScanner()
	result, err := sc.ScanLibrary(dir)
	if err != nil {
		return Response{OK: false, Error: "scan failed: " + err.Error()}
	}

	lib := &cache.CachedLibrary{
		Dir:       dir,
		ScanTime:  time.Now(),
		FileCount: result.AudioFiles,
		Files:     result.Files,
		Artists:   result.Artists,
		Albums:    result.Albums,
		Genres:    result.Genres,
		Composers: result.Composers,
		Errors:    result.Errors,
	}

	if err := d.cache.Save(lib); err != nil {
		logger.Log.Error("Failed to save cache: %v", err)
	}

	settings, _ := d.store.GetSettings()
	settings.LastDir = dir
	d.store.SaveSettings(settings)

	return Response{OK: true, Library: lib}
}

func (d *Daemon) handleCoverArt(filePath string) Response {
	if filePath == "" {
		return Response{OK: false, Error: "file_path is required"}
	}

	extractor := metadata.NewExtractor()
	audioFile, err := extractor.ExtractFromFile(filePath)
	if err != nil {
		return Response{OK: false, Error: "failed to extract cover art: " + err.Error()}
	}

	return Response{
		OK:        true,
		CoverArt:  audioFile.CoverArt,
		CoverMime: audioFile.CoverArtMime,
	}
}

func (d *Daemon) handleGetFavorites() Response {
	favs, err := d.store.GetFavorites()
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	if favs == nil {
		favs = []string{}
	}
	return Response{OK: true, Favorites: favs}
}

func (d *Daemon) handleAddFavorite(filePath string) Response {
	if err := d.store.AddFavorite(filePath); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true}
}

func (d *Daemon) handleRemoveFavorite(filePath string) Response {
	if err := d.store.RemoveFavorite(filePath); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true}
}

func (d *Daemon) handleIsFavorite(filePath string) Response {
	isFav, err := d.store.IsFavorite(filePath)
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, IsFav: isFav}
}

func (d *Daemon) handleGetSettings() Response {
	settings, err := d.store.GetSettings()
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Settings: settings}
}

func (d *Daemon) handleSaveSettings(value string) Response {
	var settings store.Settings
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return Response{OK: false, Error: "invalid settings JSON: " + err.Error()}
	}
	if err := d.store.SaveSettings(&settings); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true}
}

func (d *Daemon) handleGetStats() Response {
	stats, err := d.store.GetPlayStats()
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Stats: stats}
}

func (d *Daemon) handleRecordPlay(filePath string) Response {
	if err := d.store.RecordPlay(filePath); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true}
}

func (d *Daemon) handleGetQueue() Response {
	q, err := d.store.GetQueue()
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Queue: q}
}

func (d *Daemon) handleSaveQueue(value string) Response {
	var q store.QueueState
	if err := json.Unmarshal([]byte(value), &q); err != nil {
		return Response{OK: false, Error: "invalid queue JSON: " + err.Error()}
	}
	if err := d.store.SaveQueue(&q); err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true}
}
