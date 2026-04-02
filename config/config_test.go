package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootPathOverrideChangesAllDirs(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_STATE_HOME", "")

	original := RootPath()
	SetRootPath(filepath.Join(t.TempDir(), "instance-a"))
	t.Cleanup(func() { SetRootPath(original) })

	if got, want := DataDir(), filepath.Join(RootPath(), "data"); got != want {
		t.Fatalf("DataDir() = %q, want %q", got, want)
	}
	if got, want := ConfigDir(), filepath.Join(RootPath(), "config"); got != want {
		t.Fatalf("ConfigDir() = %q, want %q", got, want)
	}
	if got, want := StateDir(), filepath.Join(RootPath(), "state"); got != want {
		t.Fatalf("StateDir() = %q, want %q", got, want)
	}
	if got, want := DBPath(), filepath.Join(RootPath(), "data", "irecall.db"); got != want {
		t.Fatalf("DBPath() = %q, want %q", got, want)
	}
	if got, want := LogPath(), filepath.Join(RootPath(), "state", "irecall.log"); got != want {
		t.Fatalf("LogPath() = %q, want %q", got, want)
	}
}

func TestEnsureDirsCreatesOverrideTree(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_STATE_HOME", "")

	original := RootPath()
	root := filepath.Join(t.TempDir(), "instance-b")
	SetRootPath(root)
	t.Cleanup(func() { SetRootPath(original) })

	if err := EnsureDirs(); err != nil {
		t.Fatalf("EnsureDirs() error = %v", err)
	}

	for _, dir := range []string{DataDir(), ConfigDir(), StateDir()} {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("stat %q: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", dir)
		}
	}
}
