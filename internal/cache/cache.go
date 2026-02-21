package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/xdg"
)

type CachedLibrary struct {
	Dir       string               `json:"dir"`
	ScanTime  time.Time            `json:"scan_time"`
	FileCount int                  `json:"file_count"`
	Files     []metadata.AudioFile `json:"files"`
	Artists   map[string]int       `json:"artists"`
	Albums    map[string]int       `json:"albums"`
	Genres    map[string]int       `json:"genres"`
	Composers map[string]int       `json:"composers"`
	Errors    []string             `json:"errors,omitempty"`
}

type Cache struct {
	dir string
	mu  sync.RWMutex

	hot map[string]*CachedLibrary
}

func NewCache() (*Cache, error) {
	dir := xdg.CacheDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Cache{
		dir: dir,
		hot: make(map[string]*CachedLibrary),
	}, nil
}

func NewCacheAt(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Cache{
		dir: dir,
		hot: make(map[string]*CachedLibrary),
	}, nil
}

func hashDir(dir string) string {
	h := sha256.Sum256([]byte(dir))
	return fmt.Sprintf("%x", h[:8])
}

func (c *Cache) cacheFile(dir string) string {
	return filepath.Join(c.dir, hashDir(dir)+".json")
}

func (c *Cache) Load(dir string) *CachedLibrary {
	c.mu.RLock()
	if lib, ok := c.hot[dir]; ok {
		c.mu.RUnlock()
		return lib
	}
	c.mu.RUnlock()

	path := c.cacheFile(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var lib CachedLibrary
	if err := json.Unmarshal(data, &lib); err != nil {
		return nil
	}

	c.mu.Lock()
	c.hot[dir] = &lib
	c.mu.Unlock()

	return &lib
}

func (c *Cache) Save(lib *CachedLibrary) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.hot[lib.Dir] = lib

	data, err := json.Marshal(lib)
	if err != nil {
		return err
	}

	path := c.cacheFile(lib.Dir)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (c *Cache) Invalidate(dir string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.hot, dir)
	os.Remove(c.cacheFile(dir))
}

func (c *Cache) IsStale(dir string) bool {
	lib := c.Load(dir)
	if lib == nil {
		return true
	}

	info, err := os.Stat(dir)
	if err != nil {
		return true
	}

	if info.ModTime().After(lib.ScanTime) {
		return true
	}

	return false
}
