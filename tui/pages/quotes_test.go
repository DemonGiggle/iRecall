package pages

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
)

func TestQuotesPageShareUsesSelectedOrCurrentQuote(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quotes = []core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	}
	page.quoteFns.clamp(page.quotes)

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	page = model
	if cmd == nil {
		t.Fatal("share command = nil, want command")
	}
	msg := cmd()
	open, ok := msg.(OpenQuoteShareMsg)
	if !ok {
		t.Fatalf("msg type = %T, want OpenQuoteShareMsg", msg)
	}
	if len(open.Quotes) != 1 || open.Quotes[0].ID != 1 {
		t.Fatalf("shared quotes = %+v, want current quote 1", open.Quotes)
	}

	page.quoteFns.selected[2] = true
	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	page = model
	if cmd == nil {
		t.Fatal("share command with selection = nil, want command")
	}
	msg = cmd()
	open, ok = msg.(OpenQuoteShareMsg)
	if !ok {
		t.Fatalf("msg type = %T, want OpenQuoteShareMsg", msg)
	}
	if len(open.Quotes) != 1 || open.Quotes[0].ID != 2 {
		t.Fatalf("shared quotes = %+v, want selected quote 2", open.Quotes)
	}

	if !containsAllText(page.View(), "s: Share") {
		t.Fatalf("quotes page help missing share hint:\n%s", page.View())
	}
}

func TestQuotesPageCanOpenAddQuoteModal(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyCtrlN})
	page = model
	if cmd == nil {
		t.Fatal("add quote command = nil, want command")
	}

	msg := cmd()
	open, ok := msg.(OpenQuoteEditorMsg)
	if !ok {
		t.Fatalf("msg type = %T, want OpenQuoteEditorMsg", msg)
	}
	if open.Mode != QuoteEditorModeAdd {
		t.Fatalf("editor mode = %v, want %v", open.Mode, QuoteEditorModeAdd)
	}

	if !containsAllText(page.View(), "ctrl+n: Add Quote") {
		t.Fatalf("quotes page help missing add quote hint:\n%s", page.View())
	}
}

func TestQuotesPageCanOpenImportModal(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("i")})
	page = model
	if cmd == nil {
		t.Fatal("import command = nil, want command")
	}

	msg := cmd()
	if _, ok := msg.(OpenQuoteImportMsg); !ok {
		t.Fatalf("msg type = %T, want OpenQuoteImportMsg", msg)
	}

	if !containsAllText(page.View(), "i: Import") {
		t.Fatalf("quotes page help missing import hint:\n%s", page.View())
	}
}
