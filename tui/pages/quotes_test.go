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
