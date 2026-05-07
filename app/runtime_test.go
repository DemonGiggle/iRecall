package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gigol/irecall/config"
)

func TestSwitchRuntimeRejectsMissingInputs(t *testing.T) {
	if _, err := SwitchRuntime(nil, nil); err == nil || !strings.Contains(err.Error(), "runtime is not initialized") {
		t.Fatalf("SwitchRuntime(nil, nil) error = %v, want runtime initialization failure", err)
	}

	if _, err := SwitchRuntime(&RuntimeState{}, nil); err == nil || !strings.Contains(err.Error(), "runtime is not initialized") {
		t.Fatalf("SwitchRuntime(empty runtime, nil) error = %v, want runtime initialization failure", err)
	}

	current := openRuntimeForTest(t, filepath.Join(t.TempDir(), "current"))
	defer func() { _ = current.Engine.Close() }()

	if _, err := SwitchRuntime(current, nil); err == nil || !strings.Contains(err.Error(), "settings are required") {
		t.Fatalf("SwitchRuntime(valid runtime, nil) error = %v, want missing settings failure", err)
	}
}

func TestSwitchRuntimeSameRootUpdatesSettingsAndPreferredRoot(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))

	root := filepath.Join(t.TempDir(), "same-root")
	current := openRuntimeForTest(t, root)
	defer func() { _ = current.Engine.Close() }()

	next := *current.Settings
	next.RootDir = "  " + root + "  "
	next.Theme = "forest"

	switched, err := SwitchRuntime(current, &next)
	if err != nil {
		t.Fatalf("SwitchRuntime() error = %v", err)
	}
	if switched != current {
		t.Fatal("SwitchRuntime() returned a different runtime for same-root update")
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("Abs(root) error = %v", err)
	}
	if current.Paths.RootDir != absRoot {
		t.Fatalf("current.Paths.RootDir = %q, want %q", current.Paths.RootDir, absRoot)
	}
	if current.Settings.Theme != "forest" {
		t.Fatalf("current.Settings.Theme = %q, want forest", current.Settings.Theme)
	}

	saved, err := current.Engine.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("LoadSettings() error = %v", err)
	}
	if saved.Theme != "forest" {
		t.Fatalf("saved.Theme = %q, want forest", saved.Theme)
	}

	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != absRoot {
		t.Fatalf("preferred root = %q, want %q", preferredRoot, absRoot)
	}
}

func TestSwitchRuntimeUsesExistingTargetWithoutCopyingSourceData(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg-config"))

	sourceRoot := filepath.Join(t.TempDir(), "source")
	targetRoot := filepath.Join(t.TempDir(), "target")
	seedQuoteForRuntimeTest(t, sourceRoot, "source-only quote")
	seedQuoteForRuntimeTest(t, targetRoot, "target-only quote")

	current := openRuntimeForTest(t, sourceRoot)

	next := *current.Settings
	next.RootDir = targetRoot
	next.Theme = "midnight"

	switched, err := SwitchRuntime(current, &next)
	if err != nil {
		t.Fatalf("SwitchRuntime() error = %v", err)
	}
	defer func() { _ = switched.Engine.Close() }()

	quotes, err := switched.Engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].Content != "target-only quote" {
		t.Fatalf("quotes after switching to existing target = %+v, want target data only", quotes)
	}

	absTarget, err := filepath.Abs(targetRoot)
	if err != nil {
		t.Fatalf("Abs(targetRoot) error = %v", err)
	}
	if switched.Paths.RootDir != absTarget {
		t.Fatalf("switched.Paths.RootDir = %q, want %q", switched.Paths.RootDir, absTarget)
	}

	preferredRoot, err := config.LoadPreferredRootPath()
	if err != nil {
		t.Fatalf("LoadPreferredRootPath() error = %v", err)
	}
	if preferredRoot != absTarget {
		t.Fatalf("preferred root = %q, want %q", preferredRoot, absTarget)
	}
}

func TestSwitchRuntimeRestoresOriginalRuntimeOnCopyFailure(t *testing.T) {
	sourceRoot := filepath.Join(t.TempDir(), "source")
	seedQuoteForRuntimeTest(t, sourceRoot, "restore me")

	current := openRuntimeForTest(t, sourceRoot)
	badFile := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(badFile, []byte("x"), 0o600); err != nil {
		t.Fatalf("WriteFile(badFile) error = %v", err)
	}
	current.Paths.DataDir = badFile

	next := *current.Settings
	next.RootDir = filepath.Join(t.TempDir(), "target")

	restored, err := SwitchRuntime(current, &next)
	if err == nil || !strings.Contains(err.Error(), "copy runtime data") {
		t.Fatalf("SwitchRuntime() error = %v, want copy runtime data failure", err)
	}
	if restored == nil {
		t.Fatal("SwitchRuntime() restored runtime = nil, want reopened original runtime")
	}
	defer func() { _ = restored.Engine.Close() }()

	absSource, err := filepath.Abs(sourceRoot)
	if err != nil {
		t.Fatalf("Abs(sourceRoot) error = %v", err)
	}
	if restored.Paths.RootDir != absSource {
		t.Fatalf("restored.Paths.RootDir = %q, want %q", restored.Paths.RootDir, absSource)
	}

	quotes, err := restored.Engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes(restored) error = %v", err)
	}
	if len(quotes) != 1 || quotes[0].Content != "restore me" {
		t.Fatalf("restored quotes = %+v, want original data preserved", quotes)
	}
}

func TestRuntimeHasDataIgnoresPreferredRootMarker(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	paths, err := resolvePaths(root)
	if err != nil {
		t.Fatalf("resolvePaths() error = %v", err)
	}
	if err := ensurePaths(paths); err != nil {
		t.Fatalf("ensurePaths() error = %v", err)
	}

	markerPath := filepath.Join(paths.ConfigDir, config.PreferredRootFileName)
	if err := os.WriteFile(markerPath, []byte("/tmp/other-root\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(markerPath) error = %v", err)
	}

	hasData, err := runtimeHasData(paths)
	if err != nil {
		t.Fatalf("runtimeHasData(marker only) error = %v", err)
	}
	if hasData {
		t.Fatal("runtimeHasData(marker only) = true, want false")
	}

	realFile := filepath.Join(paths.ConfigDir, "settings.json")
	if err := os.WriteFile(realFile, []byte("{}"), 0o600); err != nil {
		t.Fatalf("WriteFile(realFile) error = %v", err)
	}
	hasData, err = runtimeHasData(paths)
	if err != nil {
		t.Fatalf("runtimeHasData(real file) error = %v", err)
	}
	if !hasData {
		t.Fatal("runtimeHasData(real file) = false, want true")
	}
}

func TestCopyDirNestedTargetSkipsRecursiveCopy(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	if err := os.MkdirAll(src, 0o700); err != nil {
		t.Fatalf("MkdirAll(src) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "quote.txt"), []byte("nested copy"), 0o600); err != nil {
		t.Fatalf("WriteFile(src file) error = %v", err)
	}

	dst := filepath.Join(src, "nested", "copy")
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("copyDir() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, "quote.txt"))
	if err != nil {
		t.Fatalf("ReadFile(copied file) error = %v", err)
	}
	if string(data) != "nested copy" {
		t.Fatalf("copied file content = %q, want nested copy", string(data))
	}

	if _, err := os.Stat(filepath.Join(dst, "nested")); !os.IsNotExist(err) {
		t.Fatalf("nested recursive path error = %v, want not exist", err)
	}
}

func TestCopyDirRejectsFileSourceAndCopyFileCreatesParents(t *testing.T) {
	root := t.TempDir()
	srcFile := filepath.Join(root, "source.txt")
	if err := os.WriteFile(srcFile, []byte("copied"), 0o600); err != nil {
		t.Fatalf("WriteFile(srcFile) error = %v", err)
	}

	if err := copyDir(srcFile, filepath.Join(root, "dst")); err == nil || !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("copyDir(file source) error = %v, want directory failure", err)
	}

	dstFile := filepath.Join(root, "deep", "path", "copied.txt")
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}
	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("ReadFile(dstFile) error = %v", err)
	}
	if string(data) != "copied" {
		t.Fatalf("copied file content = %q, want copied", string(data))
	}
}

func TestIsSameOrNestedPath(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "parent")
	child := filepath.Join(parent, "child")
	sibling := filepath.Join(root, "sibling")

	same, err := isSameOrNestedPath(parent, parent)
	if err != nil {
		t.Fatalf("isSameOrNestedPath(parent, parent) error = %v", err)
	}
	if !same {
		t.Fatal("isSameOrNestedPath(parent, parent) = false, want true")
	}

	nested, err := isSameOrNestedPath(parent, child)
	if err != nil {
		t.Fatalf("isSameOrNestedPath(parent, child) error = %v", err)
	}
	if !nested {
		t.Fatal("isSameOrNestedPath(parent, child) = false, want true")
	}

	notNested, err := isSameOrNestedPath(parent, sibling)
	if err != nil {
		t.Fatalf("isSameOrNestedPath(parent, sibling) error = %v", err)
	}
	if notNested {
		t.Fatal("isSameOrNestedPath(parent, sibling) = true, want false")
	}
}

func openRuntimeForTest(t *testing.T, root string) *RuntimeState {
	t.Helper()

	runtimeState, err := OpenRuntime(root)
	if err != nil {
		t.Fatalf("OpenRuntime(%q) error = %v", root, err)
	}
	return runtimeState
}

func seedQuoteForRuntimeTest(t *testing.T, root, content string) {
	t.Helper()

	app, err := NewApp(root)
	if err != nil {
		t.Fatalf("NewApp(%q) error = %v", root, err)
	}
	if _, err := app.AddQuote(content); err != nil {
		_ = app.engine.Close()
		t.Fatalf("AddQuote(%q) error = %v", content, err)
	}
	if err := app.engine.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
