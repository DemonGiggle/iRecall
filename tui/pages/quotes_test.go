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

func TestQuotesPageListShowsTagPreviewOnly(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quotes = []core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta", "gamma", "delta", "epsilon"}},
	}
	page.quoteFns.clamp(page.quotes)
	page.viewport.SetContent(page.renderQuotes())

	view := page.View()
	if !containsAllText(view, "alpha", "beta", "gamma", "+2 more") {
		t.Fatalf("quotes page preview missing expected compact tags:\n%s", view)
	}
	if containsAllText(view, "delta", "epsilon") {
		t.Fatalf("quotes page preview should hide extra tags:\n%s", view)
	}
}

func TestQuotesPageListTruncatesLongContent(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 60, 40)
	page.loading = false
	page.quotes = []core.Quote{
		{ID: 1, Content: "This is a very long quote entry that should be truncated in the list view so every row stays compact and easy to scan.", Tags: []string{"alpha"}},
	}
	page.quoteFns.clamp(page.quotes)
	page.viewport.SetContent(page.renderQuotes())

	view := page.View()
	if !containsAllText(view, "This is a very long") {
		t.Fatalf("quotes page preview missing expected content prefix:\n%s", view)
	}
	if containsAllText(view, "every row stays compact and easy to scan.") {
		t.Fatalf("quotes page preview should truncate long content:\n%s", view)
	}
	if !containsAllText(view, "…") {
		t.Fatalf("quotes page preview should show ellipsis for truncated content:\n%s", view)
	}
}

func TestQuotesPageEnterShowsDetailView(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quotes = []core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta", "gamma", "delta"}},
	}
	page.quoteFns.clamp(page.quotes)
	page.viewport.SetContent(page.renderQuotes())

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model

	if !page.detail {
		t.Fatal("detail = false, want true")
	}
	view := page.View()
	if !containsAllText(view, "Quote [1]", "alpha", "beta", "gamma", "delta", "enter/esc: Back to list") {
		t.Fatalf("quotes page detail view missing expected full data:\n%s", view)
	}
	if !containsAllText(view, "first quote") {
		t.Fatalf("quotes page detail view missing full quote content:\n%s", view)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model
	if page.detail {
		t.Fatal("detail = true after esc, want false")
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
