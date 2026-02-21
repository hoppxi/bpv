package xdg

import (
	"fmt"
	"os"
	"path/filepath"
)

func DataDir() string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "bpv")
}

func StateDir() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "bpv")
}

func RuntimeDir() string {
	base := os.Getenv("XDG_RUNTIME_DIR")
	if base == "" {
		base = filepath.Join(os.TempDir(), "bpv-"+uidStr())
		os.MkdirAll(base, 0700)
		return base
	}
	dir := filepath.Join(base, "bpv")
	os.MkdirAll(dir, 0700)
	return dir
}

func CacheDir() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "bpv")
}

func SocketPath() string {
	return filepath.Join(RuntimeDir(), "bpv.sock")
}

func LockPath() string {
	return filepath.Join(RuntimeDir(), "bpv.lock")
}

func LogPath() string {
	return filepath.Join(StateDir(), "bpvd.log")
}

func PidPath() string {
	return filepath.Join(RuntimeDir(), "bpvd.pid")
}

func uidStr() string {
	return fmt.Sprintf("%d", os.Getuid())
}
