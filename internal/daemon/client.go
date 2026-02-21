package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hoppxi/bpv/internal/cache"
	"github.com/hoppxi/bpv/internal/store"
)

type Client struct {
	conn    net.Conn
	mu      sync.Mutex
	scanner *bufio.Scanner
}

func Connect() (*Client, error) {
	sockPath := SocketPath()

	conn, err := net.DialTimeout("unix", sockPath, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("daemon not running: %v)", err)
	}

	sc := bufio.NewScanner(conn)
	sc.Buffer(make([]byte, 0), 64*1024*1024)

	return &Client{
		conn:    conn,
		scanner: sc,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) send(req Request) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	if _, err := c.conn.Write(data); err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}
		return nil, fmt.Errorf("connection closed")
	}

	var resp Response
	if err := json.Unmarshal(c.scanner.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &resp, nil
}

func (c *Client) Ping() error {
	resp, err := c.send(Request{Action: "ping"})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("ping failed: %s", resp.Error)
	}
	return nil
}

func (c *Client) GetLibrary(dir string) (*cache.CachedLibrary, error) {
	resp, err := c.send(Request{Action: "library", Dir: dir})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("library error: %s", resp.Error)
	}
	return resp.Library, nil
}

func (c *Client) Scan(dir string) (*cache.CachedLibrary, error) {
	resp, err := c.send(Request{Action: "scan", Dir: dir})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("scan error: %s", resp.Error)
	}
	return resp.Library, nil
}

func (c *Client) GetCoverArt(filePath string) (string, string, error) {
	resp, err := c.send(Request{Action: "cover-art", FilePath: filePath})
	if err != nil {
		return "", "", err
	}
	if !resp.OK {
		return "", "", fmt.Errorf("cover art error: %s", resp.Error)
	}
	return resp.CoverArt, resp.CoverMime, nil
}

func (c *Client) GetFavorites() ([]string, error) {
	resp, err := c.send(Request{Action: "get-favorites"})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("favorites error: %s", resp.Error)
	}
	return resp.Favorites, nil
}

func (c *Client) AddFavorite(filePath string) error {
	resp, err := c.send(Request{Action: "add-favorite", FilePath: filePath})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("add favorite error: %s", resp.Error)
	}
	return nil
}

func (c *Client) RemoveFavorite(filePath string) error {
	resp, err := c.send(Request{Action: "remove-favorite", FilePath: filePath})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("remove favorite error: %s", resp.Error)
	}
	return nil
}

func (c *Client) IsFavorite(filePath string) (bool, error) {
	resp, err := c.send(Request{Action: "is-favorite", FilePath: filePath})
	if err != nil {
		return false, err
	}
	if !resp.OK {
		return false, fmt.Errorf("is favorite error: %s", resp.Error)
	}
	return resp.IsFav, nil
}

func (c *Client) GetSettings() (*store.Settings, error) {
	resp, err := c.send(Request{Action: "get-settings"})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("settings error: %s", resp.Error)
	}
	return resp.Settings, nil
}

func (c *Client) SaveSettings(settings *store.Settings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	resp, err := c.send(Request{Action: "save-settings", Value: string(data)})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("save settings error: %s", resp.Error)
	}
	return nil
}

func (c *Client) GetStats() (map[string]int, error) {
	resp, err := c.send(Request{Action: "get-stats"})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("stats error: %s", resp.Error)
	}
	return resp.Stats, nil
}

func (c *Client) RecordPlay(filePath string) error {
	resp, err := c.send(Request{Action: "record-play", FilePath: filePath})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("record play error: %s", resp.Error)
	}
	return nil
}

func (c *Client) GetQueue() (*store.QueueState, error) {
	resp, err := c.send(Request{Action: "get-queue"})
	if err != nil {
		return nil, err
	}
	if !resp.OK {
		return nil, fmt.Errorf("get queue error: %s", resp.Error)
	}
	return resp.Queue, nil
}

func (c *Client) SaveQueue(q *store.QueueState) error {
	data, err := json.Marshal(q)
	if err != nil {
		return err
	}
	resp, err := c.send(Request{Action: "save-queue", Value: string(data)})
	if err != nil {
		return err
	}
	if !resp.OK {
		return fmt.Errorf("save queue error: %s", resp.Error)
	}
	return nil
}

func IsRunning() bool {
	sockPath := SocketPath()
	conn, err := net.DialTimeout("unix", sockPath, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
