package tui

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
	"github.com/gigol/irecall/tui/pages"
	"github.com/gigol/irecall/tui/styles"
)

func TestAppTabNavigationAndQuotesReload(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	if app.page != pageRecall {
		t.Fatalf("initial page = %v, want %v", app.page, pageRecall)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageHistory {
		t.Fatalf("page after tab = %v, want %v", app.page, pageHistory)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageQuotes {
		t.Fatalf("page after second tab = %v, want %v", app.page, pageQuotes)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageSettings {
		t.Fatalf("page after third tab = %v, want %v", app.page, pageSettings)
	}

	app = updateAppWithKey(t, app, tea.KeyShiftTab)
	if app.page != pageQuotes {
		t.Fatalf("page after shift+tab from settings = %v, want %v", app.page, pageQuotes)
	}

	app = updateAppWithKey(t, app, tea.KeyShiftTab)
	if app.page != pageHistory {
		t.Fatalf("page after shift+tab from quotes = %v, want %v", app.page, pageHistory)
	}

	app = updateAppWithKey(t, app, tea.KeyShiftTab)
	if app.page != pageRecall {
		t.Fatalf("page after shift+tab from history = %v, want %v", app.page, pageRecall)
	}

	app = updateAppWithKey(t, app, tea.KeyTab)
	if app.page != pageHistory {
		t.Fatalf("page after tab from recall = %v, want %v", app.page, pageHistory)
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

func TestAppHeaderShowsUserGreeting(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	view := app.View()
	if !containsAllText(view, "Hi! Alice", "Recall", "Quotes", "History", "Settings") {
		t.Fatalf("header missing expected greeting:\n%s", view)
	}
}

func TestAppAppliesThemeFromSettings(t *testing.T) {
	settings := core.DefaultSettings()
	settings.Theme = "ocean"
	app := NewApp(nil, settings, &core.UserProfile{UserID: "user-1", DisplayName: "Alice"}, 120, 40)

	if app.settings.Theme != "ocean" {
		t.Fatalf("settings theme = %q, want ocean", app.settings.Theme)
	}
	if styles.CurrentThemeName() != "ocean" {
		t.Fatalf("CurrentThemeName() = %q, want ocean", styles.CurrentThemeName())
	}
}

func TestAppOpensAndClosesQuoteShareOverlay(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)
	quotes, err := app.engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("quote count = %d, want 1", len(quotes))
	}

	model, cmd := app.Update(pages.OpenQuoteShareMsg{Quotes: quotes})
	app, _ = model.(App)
	if app.overlay != overlayQuoteShare {
		t.Fatalf("overlay after open = %v, want %v", app.overlay, overlayQuoteShare)
	}
	if cmd == nil {
		t.Fatal("share init command = nil, want command")
	}

	msg := cmd()
	model, _ = app.Update(msg)
	app, _ = model.(App)
	if !containsAllText(app.View(), "Share Quotes", "Export Payload", "\"schema_version\": 2") {
		t.Fatalf("share overlay view missing expected content:\n%s", app.View())
	}

	model, _ = app.Update(pages.CloseQuoteShareMsg{})
	app, _ = model.(App)
	if app.overlay != overlayNone {
		t.Fatalf("overlay after close = %v, want %v", app.overlay, overlayNone)
	}
}

func TestAppReloadsQuotesAfterImportOverlayCloses(t *testing.T) {
	t.Parallel()

	app := newTestApp(t)

	model, cmd := app.Update(pages.OpenQuoteImportMsg{})
	app, _ = model.(App)
	if app.overlay != overlayQuoteImport {
		t.Fatalf("overlay after import open = %v, want %v", app.overlay, overlayQuoteImport)
	}
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = app.Update(msg)
			app, _ = model.(App)
		}
	}

	model, cmd = app.Update(pages.CloseQuoteImportMsg{Reload: true})
	app, _ = model.(App)
	if app.overlay != overlayNone {
		t.Fatalf("overlay after import close = %v, want %v", app.overlay, overlayNone)
	}
	if cmd == nil {
		t.Fatal("quotes reload command after import close = nil, want command")
	}
	msg := cmd()
	loaded, ok := msg.(pages.QuotesLoadedMsg)
	if !ok {
		t.Fatalf("reload msg type = %T, want pages.QuotesLoadedMsg", msg)
	}
	if loaded.Err != nil {
		t.Fatalf("quotes reload error = %v", loaded.Err)
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
		GlobalID:         "quote-1",
		AuthorUserID:     profile.UserID,
		AuthorName:       profile.DisplayName,
		SourceUserID:     profile.UserID,
		SourceName:       profile.DisplayName,
		SourceBackend:    "local",
		SourceNamespace:  "local:" + profile.UserID,
		SourceEntityType: "quote",
		SourceEntityID:   "quote-1",
		SourceLabel:      "Local quote",
		Version:          1,
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
	if loaded, ok := msg.(pages.HistoryLoadedMsg); ok {
		if loaded.Err != nil {
			t.Fatalf("history reload returned error: %v", loaded.Err)
		}
	}

	model, _ = next.Update(msg)
	next, ok = model.(App)
	if !ok {
		t.Fatalf("model type after command = %T, want tui.App", model)
	}

	return next
}

func containsAllText(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
