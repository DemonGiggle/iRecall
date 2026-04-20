package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gigol/irecall/config"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
)

type RuntimeState struct {
	Engine   *core.Engine
	Settings *core.Settings
	Profile  *core.UserProfile
	Paths    AppPaths
}

func OpenRuntime(root string) (*RuntimeState, error) {
	paths, err := resolvePaths(root)
	if err != nil {
		return nil, err
	}
	if err := ensurePaths(paths); err != nil {
		return nil, err
	}

	store, err := db.Open(paths.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open app database: %w", err)
	}

	defaults := core.DefaultSettings()
	engine := core.New(store, defaults)

	settings, err := engine.LoadSettings(context.Background())
	if err != nil || settings == nil {
		settings = defaults
	}
	settings.RootDir = paths.RootDir
	engine.UpdateSettings(settings)

	profile, err := engine.LoadUserProfile(context.Background())
	if err != nil {
		_ = engine.Close()
		return nil, fmt.Errorf("load app user profile: %w", err)
	}
	if err := engine.BootstrapQuoteIdentity(context.Background()); err != nil {
		_ = engine.Close()
		return nil, fmt.Errorf("bootstrap app quote identity: %w", err)
	}

	return &RuntimeState{
		Engine:   engine,
		Settings: settings,
		Profile:  profile,
		Paths:    paths,
	}, nil
}

func SwitchRuntime(current *RuntimeState, nextSettings *core.Settings) (*RuntimeState, error) {
	if current == nil || current.Engine == nil {
		return nil, fmt.Errorf("runtime is not initialized")
	}
	if nextSettings == nil {
		return nil, fmt.Errorf("settings are required")
	}

	nextSettings.RootDir = strings.TrimSpace(nextSettings.RootDir)
	currentRoot := strings.TrimSpace(current.Paths.RootDir)
	nextRoot, err := normalizeRootDir(nextSettings.RootDir)
	if err != nil {
		return nil, err
	}
	nextSettings.RootDir = nextRoot

	if nextRoot == currentRoot {
		if err := current.Engine.SaveSettings(context.Background(), nextSettings); err != nil {
			return nil, err
		}
		if err := config.SavePreferredRootPath(nextRoot); err != nil {
			return nil, fmt.Errorf("persist preferred root: %w", err)
		}
		current.Settings = nextSettings
		current.Paths.RootDir = nextRoot
		return current, nil
	}

	targetPaths, err := resolvePaths(nextRoot)
	if err != nil {
		return nil, err
	}
	targetHasData, err := runtimeHasData(targetPaths)
	if err != nil {
		return nil, err
	}

	sourceRoot := currentRoot
	if err := current.Engine.Close(); err != nil {
		return nil, err
	}

	restoreOnError := func(cause error) (*RuntimeState, error) {
		restored, restoreErr := OpenRuntime(sourceRoot)
		if restoreErr == nil {
			return restored, cause
		}
		return nil, fmt.Errorf("%w (restore original runtime: %v)", cause, restoreErr)
	}

	if !targetHasData {
		if err := copyRuntimeData(current.Paths, targetPaths); err != nil {
			return restoreOnError(fmt.Errorf("copy runtime data: %w", err))
		}
	}

	nextRuntime, err := OpenRuntime(nextRoot)
	if err != nil {
		return restoreOnError(err)
	}
	if err := nextRuntime.Engine.SaveSettings(context.Background(), nextSettings); err != nil {
		_ = nextRuntime.Engine.Close()
		return restoreOnError(err)
	}
	nextRuntime.Settings = nextSettings

	if err := config.SavePreferredRootPath(nextRoot); err != nil {
		_ = nextRuntime.Engine.Close()
		return restoreOnError(fmt.Errorf("persist preferred root: %w", err))
	}

	return nextRuntime, nil
}

func normalizeRootDir(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "", nil
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve absolute data root: %w", err)
	}
	return abs, nil
}

func runtimeHasData(paths AppPaths) (bool, error) {
	for _, candidate := range []string{paths.DBPath, paths.DataDir, paths.ConfigDir, paths.StateDir} {
		info, err := os.Stat(candidate)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return false, err
		}
		if !info.IsDir() {
			return true, nil
		}
		entries, err := os.ReadDir(candidate)
		if err != nil {
			return false, err
		}
		for _, entry := range entries {
			if entry.Name() == config.PreferredRootFileName {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

func copyRuntimeData(source, target AppPaths) error {
	for _, item := range []struct {
		src string
		dst string
	}{
		{src: source.DataDir, dst: target.DataDir},
		{src: source.ConfigDir, dst: target.ConfigDir},
		{src: source.StateDir, dst: target.StateDir},
	} {
		if err := copyDir(item.src, item.dst); err != nil {
			return err
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dst, 0o700)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", src)
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return err
	}
	skipTarget, err := isSameOrNestedPath(src, dst)
	if err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == src {
			return nil
		}
		if skipTarget {
			sameOrNested, err := isSameOrNestedPath(path, dst)
			if err != nil {
				return err
			}
			if sameOrNested {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		if d.IsDir() {
			dirInfo, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(targetPath, dirInfo.Mode().Perm())
		}
		return copyFile(path, targetPath)
	})
}

func copyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func isSameOrNestedPath(parent, child string) (bool, error) {
	parent, err := filepath.Abs(parent)
	if err != nil {
		return false, err
	}
	child, err = filepath.Abs(child)
	if err != nil {
		return false, err
	}
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false, err
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))), nil
}
