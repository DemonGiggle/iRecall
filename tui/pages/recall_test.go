package pages

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
)

func TestRecallPageFocusJumpAndReferenceHints(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)
	page.quotes = []core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	}
	page.refreshReferencePanel()

	if page.focus != focusInput {
		t.Fatalf("initial focus = %v, want %v", page.focus, focusInput)
	}
	if !page.input.Focused() {
		t.Fatal("input should start focused")
	}

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	if page.quoteFns.cursor != 0 {
		t.Fatalf("cursor moved while input focused = %d, want 0", page.quoteFns.cursor)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model
	if page.focus != focusReferenceQuotes {
		t.Fatalf("focus after ctrl+j = %v, want %v", page.focus, focusReferenceQuotes)
	}
	if page.input.Focused() {
		t.Fatal("input should be blurred when reference quotes are focused")
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	if page.quoteFns.cursor != 1 {
		t.Fatalf("cursor after down on reference quotes = %d, want 1", page.quoteFns.cursor)
	}

	view := page.View()
	if !containsAllText(view, "ctrl+j: Focus input", "↑/↓: Move", "x: Select", "s: Share", "Reference Quotes") {
		t.Fatalf("reference panel hints missing expected content:\n%s", view)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model
	if page.focus != focusInput {
		t.Fatalf("focus after second ctrl+j = %v, want %v", page.focus, focusInput)
	}
	if !page.input.Focused() {
		t.Fatal("input should be focused after returning from reference quotes")
	}
}

func TestRecallPageShareRequiresReferenceFocus(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)
	page.quotes = []core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	}
	page.refreshReferencePanel()

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	page = model
	if cmd != nil {
		if _, ok := cmd().(OpenQuoteShareMsg); ok {
			t.Fatal("share command while input focused = OpenQuoteShareMsg, want no share")
		}
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model

	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	page = model
	if cmd == nil {
		t.Fatal("share command on reference focus = nil, want command")
	}
	msg := cmd()
	open, ok := msg.(OpenQuoteShareMsg)
	if !ok {
		t.Fatalf("msg type = %T, want OpenQuoteShareMsg", msg)
	}
	if len(open.Quotes) != 1 || open.Quotes[0].ID != 1 {
		t.Fatalf("shared quotes = %+v, want current quote 1", open.Quotes)
	}
}

func TestRecallPageResponseShowsQuestionContext(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)
	page.question = "Ask about memory"
	page.respBuf = "Here is the answer."
	page.updateResponsePanel()

	view := page.View()
	if !containsAllText(view, "Question:", "Ask about memory", "Here is the answer.") {
		t.Fatalf("recall response missing question context:\n%s", view)
	}
}

func containsAllText(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
