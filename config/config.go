package config

import (
	"os"
	"path/filepath"
)

const appName = "irecall"

// DataDir returns the XDG data directory for iRecall.
// Falls back to ~/.local/share/irecall.
func DataDir() string {
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", appName)
}

// ConfigDir returns the XDG config directory for iRecall.
// Falls back to ~/.config/irecall.
func ConfigDir() string {
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", appName)
}

// StateDir returns the XDG state directory for iRecall.
// Falls back to ~/.local/state/irecall.
func StateDir() string {
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
