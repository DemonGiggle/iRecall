package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const appName = "irecall"

var rootOverride string

// SetRootPath overrides the default XDG directories and stores all app files
// beneath the provided root path.
func SetRootPath(root string) {
	rootOverride = root
}

// RootPath returns the active override root path, if any.
func RootPath() string {
	return rootOverride
}

// DataDir returns the XDG data directory for iRecall.
// Falls back to ~/.local/share/irecall.
func DataDir() string {
	if rootOverride != "" {
		return filepath.Join(rootOverride, "data")
	}
	if runtime.GOOS == "windows" {
		if root := windowsRoot(); root != "" {
			return filepath.Join(root, "data")
		}
	}
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", appName)
}

// ConfigDir returns the XDG config directory for iRecall.
// Falls back to ~/.config/irecall.
func ConfigDir() string {
	if rootOverride != "" {
		return filepath.Join(rootOverride, "config")
	}
	if runtime.GOOS == "windows" {
		if root := windowsRoot(); root != "" {
			return filepath.Join(root, "config")
		}
	}
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", appName)
}

// StateDir returns the XDG state directory for iRecall.
// Falls back to ~/.local/state/irecall.
func StateDir() string {
	if rootOverride != "" {
		return filepath.Join(rootOverride, "state")
	}
	if runtime.GOOS == "windows" {
		if root := windowsRoot(); root != "" {
			return filepath.Join(root, "state")
		}
	}
	if d := os.Getenv("XDG_STATE_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", appName)
}

// DBPath returns the full path to the SQLite database file.
func DBPath() string {
	return filepath.Join(DataDir(), "irecall.db")
}

// LogPath returns the full path to the log file.
func LogPath() string {
	return filepath.Join(StateDir(), "irecall.log")
}

// EnsureDirs creates all required application directories.
func EnsureDirs() error {
	for _, dir := range []string{DataDir(), ConfigDir(), StateDir()} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func windowsRoot() string {
	if d := os.Getenv("LOCALAPPDATA"); d != "" {
		return filepath.Join(d, appName)
	}
	if d, err := os.UserCacheDir(); err == nil && d != "" {
		return filepath.Join(d, appName)
	}
	return ""
}
