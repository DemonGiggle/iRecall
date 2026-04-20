package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const appName = "irecall"
const preferredRootFileName = "root-path"
const PreferredRootFileName = preferredRootFileName

var rootOverride string

// SetRootPath overrides the default XDG directories and stores all app files
// beneath the provided root path.
func SetRootPath(root string) {
	rootOverride = strings.TrimSpace(root)
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
	return DefaultDataDir()
}

// ConfigDir returns the XDG config directory for iRecall.
// Falls back to ~/.config/irecall.
func ConfigDir() string {
	if rootOverride != "" {
		return filepath.Join(rootOverride, "config")
	}
	return DefaultConfigDir()
}

// StateDir returns the XDG state directory for iRecall.
// Falls back to ~/.local/state/irecall.
func StateDir() string {
	if rootOverride != "" {
		return filepath.Join(rootOverride, "state")
	}
	return DefaultStateDir()
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

// DefaultDataDir returns the platform default data directory without applying
// any root override.
func DefaultDataDir() string {
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

// DefaultConfigDir returns the platform default config directory without
// applying any root override.
func DefaultConfigDir() string {
	if runtime.GOOS == "windows" {
		if root := windowsRoot(); root != "" {
			return filepath.Join(root, "config")
		}
	}
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	// Fall back to the data dir so config files live with other app data by default
	return DefaultDataDir()
}

// DefaultStateDir returns the platform default state directory without
// applying any root override.
func DefaultStateDir() string {
	if runtime.GOOS == "windows" {
		if root := windowsRoot(); root != "" {
			return filepath.Join(root, "state")
		}
	}
	if d := os.Getenv("XDG_STATE_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	// Fall back to the data dir so state files live with other app data by default
	return DefaultDataDir()
}

func preferredRootPathFile() string {
	return filepath.Join(DefaultConfigDir(), preferredRootFileName)
}

// LoadPreferredRootPath returns the persisted root override used on startup.
func LoadPreferredRootPath() (string, error) {
	data, err := os.ReadFile(preferredRootPathFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SavePreferredRootPath persists the selected root override used on startup.
// An empty root clears the persisted override and returns to default XDG paths.
func SavePreferredRootPath(root string) error {
	root = strings.TrimSpace(root)
	path := preferredRootPathFile()
	if root == "" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(root+"\n"), 0o600)
}

// ApplyPreferredRootPath loads the persisted root override into the active
// process if no explicit override is already set.
func ApplyPreferredRootPath() (string, error) {
	if RootPath() != "" {
		return RootPath(), nil
	}
	root, err := LoadPreferredRootPath()
	if err != nil {
		return "", err
	}
	if root != "" {
		SetRootPath(root)
	}
	return root, nil
}
