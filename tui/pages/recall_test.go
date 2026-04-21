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
	page.refQuotes.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta", "gamma", "delta"}},
		{ID: 2, Content: "second quote"},
	})

	if page.focus != focusInput {
		t.Fatalf("initial focus = %v, want %v", page.focus, focusInput)
	}
	if !page.input.Focused() {
		t.Fatal("input should start focused")
	}

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	if page.refQuotes.currentCursor() != 0 {
		t.Fatalf("cursor moved while input focused = %d, want 0", page.refQuotes.currentCursor())
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
	if page.refQuotes.currentCursor() != 1 {
		t.Fatalf("cursor after down on reference quotes = %d, want 1", page.refQuotes.currentCursor())
	}

	view := page.View()
	if !containsAllText(view, "ctrl+j: Focus input", "↑/↓: Move", "a: Select all", "u: Deselect all", "s: Share", "Reference Quotes", "alpha", "beta", "gamma", "+1 more") {
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
	page.refQuotes.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	})

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

func TestRecallPageReferenceQuotesCanBulkSelect(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)
	page.refQuotes.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
	})

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	page = model
	if got := page.refQuotes.selectedCount(); got != 2 {
		t.Fatalf("selected after a = %d, want 2", got)
	}

	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")})
	page = model
	if got := page.refQuotes.selectedCount(); got != 0 {
		t.Fatalf("selected after u = %d, want 0", got)
	}
}

func TestRecallPageReferenceQuotesEnterShowsSharedDetail(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)
	page.refQuotes.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote", Tags: []string{"alpha", "beta"}},
	})

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model

	if !page.refQuotes.isDetail() {
		t.Fatal("reference quotes detail = false, want true")
	}
	view := page.View()
	if !containsAllText(view, "Quote Information", "Quote [1]", "first quote", "alpha", "beta", "enter/esc: Back") {
		t.Fatalf("recall reference detail missing shared widget content:\n%s", view)
	}
}

func TestRecallPageReferenceQuotesScrollWithCursor(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 80, 20)
	page.refQuotes.SetQuotes([]core.Quote{
		{ID: 1, Content: "first quote"},
		{ID: 2, Content: "second quote"},
		{ID: 3, Content: "third quote"},
		{ID: 4, Content: "fourth quote"},
	})

	model, _ := page.Update(tea.KeyMsg{Type: tea.KeyCtrlJ})
	page = model
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model
	model, _ = page.Update(tea.KeyMsg{Type: tea.KeyDown})
	page = model

	if page.refQuotes.currentCursor() != 2 {
		t.Fatalf("cursor = %d, want 2", page.refQuotes.currentCursor())
	}
	if page.refQuotes.yOffset() == 0 {
		t.Fatalf("reference quotes y offset = %d, want > 0 after scrolling", page.refQuotes.yOffset())
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

func TestRecallPageSaveAsQuoteStatus(t *testing.T) {
	t.Parallel()

	page := NewRecallPage(nil, 120, 40)

	model, cmd := page.Update(RecallQuoteSavedMsg{Quote: &core.Quote{ID: 1}, Err: nil})
	page = model
	if cmd == nil {
		t.Fatal("notice command = nil")
	}
	msg := cmd()
	open, ok := msg.(OpenNoticeMsg)
	if !ok {
		t.Fatalf("notice msg type = %T, want OpenNoticeMsg", msg)
	}
	if open.Title != "Recall Saved as Quote" {
		t.Fatalf("notice title = %q, want Recall Saved as Quote", open.Title)
	}

	view := page.View()
	if !containsAllText(view, "Saved recall as quote.") {
		t.Fatalf("recall view missing save status:\n%s", view)
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
