package pages

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/core/db"
)

func TestHistoryPageListAndDetailFlow(t *testing.T) {
	t.Parallel()

	engine, quote := newHistoryTestEngine(t)
	page := NewHistoryPage(engine, 120, 40)

	entry, err := engine.SaveRecallHistory(context.Background(),
		"How should I remember this?",
		"Keep the response in the history entry.",
		[]core.Quote{quote},
	)
	if err != nil {
		t.Fatalf("SaveRecallHistory() error = %v", err)
	}

	model, _ := page.Update(HistoryLoadedMsg{
		Entries: []core.RecallHistorySummary{{
			ID:        entry.ID,
			Question:  entry.Question,
			Response:  entry.Response,
			CreatedAt: entry.CreatedAt,
		}},
	})
	page = model

	if !strings.Contains(page.View(), "How should I remember this?") {
		t.Fatalf("history list view missing question:\n%s", page.View())
	}

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if !page.detail || !page.detailLoading {
		t.Fatalf("detail state after enter = detail:%v loading:%v, want true/true", page.detail, page.detailLoading)
	}
	if cmd == nil {
		t.Fatal("detail load command = nil")
	}

	msg := cmd()
	loaded, ok := msg.(HistoryDetailLoadedMsg)
	if !ok {
		t.Fatalf("detail load msg type = %T, want HistoryDetailLoadedMsg", msg)
	}
	model, _ = page.Update(loaded)
	page = model

	detailView := page.View()
	if !strings.Contains(detailView, "History Entry") || !strings.Contains(detailView, "Reference Quotes") {
		t.Fatalf("detail view missing expected sections:\n%s", detailView)
	}
	if !strings.Contains(detailView, "Keep the response in the history entry.") {
		t.Fatalf("detail view missing response:\n%s", detailView)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	page = model
	if !page.detail {
		t.Fatal("detail view should ignore list-only keybindings")
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model
	if page.detail {
		t.Fatal("detail should close on esc")
	}
}

func TestHistoryPageOpensDeleteForSelectedEntries(t *testing.T) {
	t.Parallel()

	engine, _ := newHistoryTestEngine(t)
	page := NewHistoryPage(engine, 120, 40)
	entry := core.RecallHistorySummary{
		ID:        42,
		Question:  "Delete me",
		Response:  "Stored response",
		CreatedAt: time.Now(),
	}

	model, _ := page.Update(HistoryLoadedMsg{Entries: []core.RecallHistorySummary{entry}})
	page = model

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	page = model
	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	page = model
	if cmd == nil {
		t.Fatal("delete command = nil")
	}
	msg := cmd()
	open, ok := msg.(OpenDeleteRecallHistoryMsg)
	if !ok {
		t.Fatalf("delete msg type = %T, want OpenDeleteRecallHistoryMsg", msg)
	}
	if len(open.Entries) != 1 || open.Entries[0].ID != entry.ID {
		t.Fatalf("delete entries = %+v, want selected history entry", open.Entries)
	}
}

func TestHistoryPageSaveAsQuoteStatusInDetail(t *testing.T) {
	t.Parallel()

	page := NewHistoryPage(nil, 120, 40)
	page.detail = true
	page.entry = &core.RecallHistoryEntry{
		ID:        1,
		Question:  "How do I save this?",
		Response:  "Use the save action.",
		CreatedAt: time.Now(),
	}

	model, _ := page.Update(RecallHistoryQuoteSavedMsg{Quote: &core.Quote{ID: 1}, Err: nil})
	page = model

	view := page.View()
	if !strings.Contains(view, "Saved history entry as quote.") {
		t.Fatalf("history detail missing save status:\n%s", view)
	}
}

func newHistoryTestEngine(t *testing.T) (*core.Engine, core.Quote) {
	t.Helper()

	store, err := db.Open(filepath.Join(t.TempDir(), "history.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	engine := core.New(store, core.DefaultSettings())
	profile := &core.UserProfile{UserID: "user-1", DisplayName: "Alice"}
	engine.UpdateUserProfile(profile)

	quoteID, err := store.InsertQuote("History reference quote", db.QuoteIdentity{
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
	})
	if err != nil {
		t.Fatalf("insert history quote: %v", err)
	}

	return engine, core.Quote{
		ID:           quoteID,
		GlobalID:     "quote-1",
		AuthorUserID: profile.UserID,
		AuthorName:   profile.DisplayName,
		SourceUserID: profile.UserID,
		SourceName:   profile.DisplayName,
		Content:      "History reference quote",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
