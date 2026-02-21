package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/hoppxi/bpv/internal/xdg"
)

type Store struct {
	dir string
	mu  sync.RWMutex
}

type Settings struct {
	Volume         float64 `json:"volume"`
	Shuffle        bool    `json:"shuffle"`
	Repeat         int     `json:"repeat"` // 0=off, 1=all, 2=one
	LastDir        string  `json:"last_dir"`
	LastPort       int     `json:"last_port"`
	VisualizerType string  `json:"visualizer_type,omitempty"`
	ShowVisualizer *bool   `json:"show_visualizer,omitempty"`
	AutoPlay       *bool   `json:"auto_play,omitempty"`
	Crossfade      *bool   `json:"crossfade,omitempty"`
	Gapless        *bool   `json:"gapless,omitempty"`
	EqBass         float64 `json:"eq_bass,omitempty"`
	EqMid          float64 `json:"eq_mid,omitempty"`
	EqTreble       float64 `json:"eq_treble,omitempty"`
	EqEnabled      *bool   `json:"eq_enabled,omitempty"`
}

type QueueState struct {
	FilePaths    []string `json:"file_paths"`
	CurrentIndex int      `json:"current_index"`
	Shuffle      bool     `json:"shuffle"`
	Repeat       int      `json:"repeat"`
}

func NewStore() (*Store, error) {
	dir := xdg.DataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func NewStoreAt(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Dir() string {
	return s.dir
}

func (s *Store) favPath() string {
	return filepath.Join(s.dir, "favorites.json")
}

func (s *Store) GetFavorites() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readStringSlice(s.favPath())
}

func (s *Store) SetFavorites(favs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writeJSON(s.favPath(), favs)
}

func (s *Store) AddFavorite(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	favs, _ := s.readStringSlice(s.favPath())
	if slices.Contains(favs, filePath) {
		return nil
	}

	favs = append(favs, filePath)
	return s.writeJSON(s.favPath(), favs)
}

func (s *Store) RemoveFavorite(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	favs, _ := s.readStringSlice(s.favPath())
	for i, f := range favs {
		if f == filePath {
			favs = append(favs[:i], favs[i+1:]...)
			return s.writeJSON(s.favPath(), favs)
		}
	}
	return nil
}

func (s *Store) IsFavorite(filePath string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	favs, _ := s.readStringSlice(s.favPath())
	if slices.Contains(favs, filePath) {
		return true, nil
	}
	return false, nil
}

func (s *Store) settingsPath() string {
	return filepath.Join(s.dir, "settings.json")
}

func (s *Store) GetSettings() (*Settings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	settings := &Settings{
		Volume:   0,
		Repeat:   0,
		LastPort: 8080,
	}

	data, err := os.ReadFile(s.settingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return settings, err
	}

	if err := json.Unmarshal(data, settings); err != nil {
		return settings, err
	}
	return settings, nil
}

func (s *Store) SaveSettings(settings *Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writeJSON(s.settingsPath(), settings)
}

func (s *Store) queuePath() string {
	return filepath.Join(s.dir, "queue.json")
}

func (s *Store) GetQueue() (*QueueState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.queuePath())
	if err != nil {
		if os.IsNotExist(err) {
			return &QueueState{}, nil
		}
		return nil, err
	}

	var q QueueState
	if err := json.Unmarshal(data, &q); err != nil {
		return &QueueState{}, nil
	}
	return &q, nil
}

func (s *Store) SaveQueue(q *QueueState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writeJSON(s.queuePath(), q)
}

func (s *Store) statsPath() string {
	return filepath.Join(s.dir, "playstats.json")
}

func (s *Store) GetPlayStats() (map[string]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.statsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]int), nil
		}
		return nil, err
	}

	var stats map[string]int
	if err := json.Unmarshal(data, &stats); err != nil {
		return make(map[string]int), nil
	}
	return stats, nil
}

func (s *Store) RecordPlay(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, _ := os.ReadFile(s.statsPath())
	stats := make(map[string]int)
	if data != nil {
		json.Unmarshal(data, &stats)
	}
	stats[filePath]++
	return s.writeJSON(s.statsPath(), stats)
}

func (s *Store) readStringSlice(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var result []string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Store) writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
