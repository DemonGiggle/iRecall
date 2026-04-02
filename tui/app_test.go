package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
	"github.com/gigol/irecall/tui/pages"
)

func TestAppTabNavigationAndQuotesReload(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	if app.page != pageRecall {
		t.Fatalf("initial page = %v, want %v", app.page, pageRecall)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageQuotes {
		t.Fatalf("page after tab = %v, want %v", app.page, pageQuotes)
	}

	app = updateAppWithKey(t, app, tea.KeyShiftTab)
	if app.page != pageRecall {
		t.Fatalf("page after shift+tab from quotes = %v, want %v", app.page, pageRecall)
	}

	app = updateAppWithKey(t, app, tea.KeyShiftTab)
	if app.page != pageSettings {
		t.Fatalf("page after shift+tab from recall = %v, want %v", app.page, pageSettings)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageRecall {
		t.Fatalf("page after tab from settings = %v, want %v", app.page, pageRecall)
	}
}

func TestAppStartsWithProfilePromptWhenNameMissing(t *testing.T) {
	t.Parallel()

	settings := core.DefaultSettings()
	app := NewApp(nil, settings, &core.UserProfile{UserID: "user-1", DisplayName: ""}, 120, 40)

	if app.overlay != overlayUserProfilePrompt {
		t.Fatalf("overlay = %v, want %v", app.overlay, overlayUserProfilePrompt)
	}
}

func newTestApp(t *testing.T) App {
	t.Helper()

	path := filepath.Join(t.TempDir(), "irecall.db")
	store, err := db.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	settings := core.DefaultSettings()
	engine := core.New(store, settings)
	profile := &core.UserProfile{UserID: "user-1", DisplayName: "Alice"}
	engine.UpdateUserProfile(profile)

	if _, err := store.InsertQuote("Test quote for page reload.", db.QuoteIdentity{
		GlobalID:     "quote-1",
		AuthorUserID: profile.UserID,
		AuthorName:   profile.DisplayName,
		SourceUserID: profile.UserID,
		SourceName:   profile.DisplayName,
		Version:      1,
	}); err != nil {
		t.Fatalf("insert quote: %v", err)
	}

	return NewApp(engine, settings, profile, 120, 40)
}

func updateAppWithKey(t *testing.T, app App, key tea.KeyType) App {
	t.Helper()

	model, cmd := app.Update(tea.KeyMsg{Type: key})
	next, ok := model.(App)
	if !ok {
		t.Fatalf("model type = %T, want tui.App", model)
	}

	if cmd == nil {
		return next
	}

	msg := cmd()
	if msg == nil {
		return next
	}
	if loaded, ok := msg.(pages.QuotesLoadedMsg); ok {
		if loaded.Err != nil {
			t.Fatalf("quotes reload returned error: %v", loaded.Err)
		}
		if len(loaded.Quotes) != 1 {
			t.Fatalf("quotes reload count = %d, want 1", len(loaded.Quotes))
		}
	}

	model, _ = next.Update(msg)
	next, ok = model.(App)
	if !ok {
		t.Fatalf("model type after command = %T, want tui.App", model)
	}

	return next
}
