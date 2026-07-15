package pages

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
)

func TestQuoteSharePageExportsAndSavesPayload(t *testing.T) {
	t.Parallel()

	engine := newShareTestEngine(t)
	first, err := engine.AddQuote(context.Background(), "first shared quote")
	if err != nil {
		t.Fatalf("AddQuote(first) error = %v", err)
	}
	second, err := engine.AddQuote(context.Background(), "second shared quote")
	if err != nil {
		t.Fatalf("AddQuote(second) error = %v", err)
	}

	page := NewQuoteSharePage(engine, 120, 40)
	page.Reset([]core.Quote{*first, *second})

	msg := page.Init()()
	model, _ := page.Update(msg)
	page = model

	if page.payload == "" {
		t.Fatal("payload = empty, want exported manifest")
	}
	if !containsAll(
		page.payload,
		"\"schema_version\": 3",
		"\"source_backend\": \"local\"",
		"\"source_entity_type\": \"quote\"",
		"first shared quote",
		"second shared quote",
	) {
		t.Fatalf("payload missing expected content:\n%s", page.payload)
	}

	path := filepath.Join(t.TempDir(), "exports", "quotes.irecall")
	page.pathInput.SetValue(path)
	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	page = model
	if cmd == nil {
		t.Fatal("save command = nil, want command")
	}

	model, _ = page.Update(cmd())
	page = model

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open saved bundle: %v", err)
	}
	if len(zr.File) == 0 || zr.File[0].Name != "manifest.json" {
		t.Fatalf("bundle entries = %#v, want manifest.json", zr.File)
	}
	if !strings.Contains(page.statusMsg, path) {
		t.Fatalf("status = %q, want path %q", page.statusMsg, path)
	}
}

func newShareTestEngine(t *testing.T) *core.Engine {
	t.Helper()

	path := filepath.Join(t.TempDir(), "irecall.db")
	store, err := db.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	engine := core.New(store, core.DefaultSettings())
	profile := &core.UserProfile{
		UserID:      "user-1",
		DisplayName: "Alice",
	}
	if err := engine.SaveUserProfile(context.Background(), profile); err != nil {
		t.Fatalf("SaveUserProfile() error = %v", err)
	}
	engine.UpdateUserProfile(profile)
	return engine
}
