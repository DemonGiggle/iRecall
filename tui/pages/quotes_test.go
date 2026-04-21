package pages

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
)

func TestQuotesPageShareUsesSelectedOrCurrentQuote(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	})
	page.quoteList.SetTitle("Stored Quotes (2)")

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

	page.quoteList.selection.selected[2] = true
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

	if !containsAllText(page.View(), "Stored Quotes", "a: Select all", "u: Deselect all", "s: Share") {
		t.Fatalf("quotes page widget missing shared quote actions:\n%s", page.View())
	}
}

func TestQuotesPageListShowsTagPreviewOnly(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta", "gamma", "delta", "epsilon"}},
	})
	page.quoteList.SetTitle("Stored Quotes (1)")

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
	page.quoteList.SetQuotes([]core.Quote{
		{ID: 1, Content: "This is a very long quote entry that should be truncated in the list view so every row stays compact and easy to scan.", Tags: []string{"alpha"}},
	})
	page.quoteList.SetTitle("Stored Quotes (1)")

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

func TestQuotesPageEnterShowsQuoteInformation(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta", "gamma", "delta"}},
	})
	page.quoteList.SetTitle("Stored Quotes (1)")

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if cmd != nil {
		t.Fatalf("enter command = %v, want nil", cmd)
	}
	view := page.View()
	if !page.quoteList.isDetail() {
		t.Fatal("detail = false, want true")
	}
	if !containsAllText(view, "Quote Information", "Quote [1]", "first quote", "alpha", "beta", "gamma", "delta", "enter/esc: Back", "↑/↓: Scroll", "pgup/pgdn: Page") {
		t.Fatalf("quotes page detail view missing expected information:\n%s", view)
	}
}

func TestQuotesPageDetailSupportsArrowScrolling(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 80, 12)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{
		{
			ID:      1,
			Content: strings.Join([]string{"line 1", "line 2", "line 3", "line 4", "line 5", "line 6", "line 7", "line 8", "line 9", "line 10", "line 11", "line 12"}, "\n"),
			Tags:    []string{"alpha"},
		},
	})
	page.quoteList.SetTitle("Stored Quotes (1)")

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if !page.quoteList.isDetail() {
		t.Fatal("detail = false, want true")
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	if page.quoteList.detailViewport.YOffset == 0 {
		t.Fatalf("detail viewport y offset = %d, want > 0 after down", page.quoteList.detailViewport.YOffset)
	}

	beforePageDown := page.quoteList.detailViewport.YOffset
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	page = model
	if page.quoteList.detailViewport.YOffset <= beforePageDown {
		t.Fatalf("detail viewport y offset after pgdown = %d, want > %d", page.quoteList.detailViewport.YOffset, beforePageDown)
	}

	beforePageUp := page.quoteList.detailViewport.YOffset
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	page = model
	if page.quoteList.detailViewport.YOffset >= beforePageUp {
		t.Fatalf("detail viewport y offset after pgup = %d, want < %d", page.quoteList.detailViewport.YOffset, beforePageUp)
	}
}

func TestQuotesPageHelpAppearsBeforeStoredQuotesPanel(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 120, 40)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{{ID: 1, Content: "first quote"}})
	page.quoteList.SetTitle("Stored Quotes (1)")

	view := page.View()
	if strings.Index(view, "ctrl+n: Add Quote") > strings.Index(view, "Stored Quotes") {
		t.Fatalf("quotes page actions should appear before stored quotes section:\n%s", view)
	}
}

func TestQuotesPageCursorMovementScrollsViewport(t *testing.T) {
	t.Parallel()

	page := NewQuotesPage(nil, 80, 12)
	page.loading = false
	page.quoteList.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
		{ID: 3, Content: "third quote"},
		{ID: 4, Content: "fourth quote"},
	})
	page.quoteList.SetTitle("Stored Quotes (4)")

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model

	if page.quoteList.currentCursor() != 2 {
		t.Fatalf("cursor = %d, want 2", page.quoteList.currentCursor())
	}
	if page.quoteList.yOffset() == 0 {
		t.Fatalf("viewport y offset = %d, want > 0 after scrolling", page.quoteList.yOffset())
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
