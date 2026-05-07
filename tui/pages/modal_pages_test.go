package pages

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gigol/irecall/core"
)

func TestNoticePageLifecycleAndView(t *testing.T) {
	t.Parallel()

	page := NewNoticePage(80, 24)
	page.Reset("  ", "  finished successfully  ")
	page.SetSize(96, 30)

	view := page.View()
	if !containsAll(view, "Done", "finished successfully", "enter/esc: Close") {
		t.Fatalf("notice view missing expected content:\n%s", view)
	}

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if cmd == nil {
		t.Fatal("notice close command = nil")
	}
	if _, ok := cmd().(CloseNoticeMsg); !ok {
		t.Fatalf("notice close msg type = %T, want CloseNoticeMsg", cmd())
	}
}

func TestDeleteQuotesPageLifecycleAndSummary(t *testing.T) {
	t.Parallel()

	engine := newShareTestEngine(t)
	first, err := engine.AddQuote(context.Background(), "first quote slated for deletion")
	if err != nil {
		t.Fatalf("AddQuote(first) error = %v", err)
	}
	second, err := engine.AddQuote(context.Background(), "second quote slated for deletion")
	if err != nil {
		t.Fatalf("AddQuote(second) error = %v", err)
	}

	page := NewDeleteQuotesPage(engine, 80, 24)
	page.Reset([]core.Quote{*first, *second})
	page.SetSize(96, 30)

	if got := page.summary(); !strings.Contains(got, "Delete 2 quotes?") || !strings.Contains(got, "first quote slated for deletion") {
		t.Fatalf("summary() = %q, want delete summary with quote preview", got)
	}
	if !containsAll(page.View(), "Delete Quotes", "This action cannot be undone.", "Confirm delete") {
		t.Fatalf("delete quotes view missing expected content:\n%s", page.View())
	}

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if cmd == nil {
		t.Fatal("delete quotes command = nil")
	}
	done, ok := cmd().(DeleteQuotesDoneMsg)
	if !ok {
		t.Fatalf("delete quotes done msg type = %T, want DeleteQuotesDoneMsg", cmd())
	}
	if len(done.DeletedIDs) != 2 {
		t.Fatalf("delete quotes done ids = %+v, want both quote ids", done.DeletedIDs)
	}

	model, cmd = page.Update(done)
	page = model
	if cmd == nil {
		t.Fatal("delete quotes close command = nil")
	}
	closeMsg, ok := cmd().(CloseDeleteQuotesMsg)
	if !ok {
		t.Fatalf("delete quotes close msg type = %T, want CloseDeleteQuotesMsg", cmd())
	}
	if len(closeMsg.DeletedIDs) != 2 {
		t.Fatalf("delete quotes close ids = %+v, want both quote ids", closeMsg.DeletedIDs)
	}

	quotes, err := engine.ListQuotes(context.Background())
	if err != nil {
		t.Fatalf("ListQuotes() error = %v", err)
	}
	if len(quotes) != 0 {
		t.Fatalf("quotes after delete = %+v, want empty", quotes)
	}

	page.Reset(nil)
	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model
	if cmd == nil {
		t.Fatal("delete quotes esc command = nil")
	}
	if _, ok := cmd().(CloseDeleteQuotesMsg); !ok {
		t.Fatalf("delete quotes esc msg type = %T, want CloseDeleteQuotesMsg", cmd())
	}
}

func TestDeleteRecallHistoryPageLifecycleAndSummary(t *testing.T) {
	t.Parallel()

	engine, quote := newHistoryTestEngine(t)
	first, err := engine.SaveRecallHistory(context.Background(), "first history entry", "first response", []core.Quote{quote})
	if err != nil {
		t.Fatalf("SaveRecallHistory(first) error = %v", err)
	}
	second, err := engine.SaveRecallHistory(context.Background(), "second history entry", "second response", []core.Quote{quote})
	if err != nil {
		t.Fatalf("SaveRecallHistory(second) error = %v", err)
	}

	summaries, err := engine.ListRecallHistory(context.Background())
	if err != nil {
		t.Fatalf("ListRecallHistory() error = %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("history summary count = %d, want 2", len(summaries))
	}

	page := NewDeleteRecallHistoryPage(engine, 80, 24)
	page.Reset(summaries)
	page.SetSize(96, 30)

	if got := page.summary(); !strings.Contains(got, "Delete 2 history entries?") || !strings.Contains(got, "first history entry") {
		t.Fatalf("summary() = %q, want delete history summary", got)
	}
	if !containsAll(page.View(), "Delete History", "This action cannot be undone.", "Confirm delete") {
		t.Fatalf("delete history view missing expected content:\n%s", page.View())
	}

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyEnter})
	page = model
	if cmd == nil {
		t.Fatal("delete history command = nil")
	}
	done, ok := cmd().(DeleteRecallHistoryDoneMsg)
	if !ok {
		t.Fatalf("delete history done msg type = %T, want DeleteRecallHistoryDoneMsg", cmd())
	}
	if len(done.DeletedIDs) != 2 {
		t.Fatalf("delete history done ids = %+v, want both history ids", done.DeletedIDs)
	}

	model, cmd = page.Update(done)
	page = model
	if cmd == nil {
		t.Fatal("delete history close command = nil")
	}
	closeMsg, ok := cmd().(CloseDeleteRecallHistoryMsg)
	if !ok {
		t.Fatalf("delete history close msg type = %T, want CloseDeleteRecallHistoryMsg", cmd())
	}
	if len(closeMsg.DeletedIDs) != 2 {
		t.Fatalf("delete history close ids = %+v, want both history ids", closeMsg.DeletedIDs)
	}

	remaining, err := engine.ListRecallHistory(context.Background())
	if err != nil {
		t.Fatalf("ListRecallHistory() after delete error = %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("history after delete = %+v, want empty", remaining)
	}

	page.Reset([]core.RecallHistorySummary{{ID: first.ID, Question: first.Question}})
	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model
	if cmd == nil {
		t.Fatal("delete history esc command = nil")
	}
	if _, ok := cmd().(CloseDeleteRecallHistoryMsg); !ok {
		t.Fatalf("delete history esc msg type = %T, want CloseDeleteRecallHistoryMsg", cmd())
	}

	if ids := recallHistoryIDs([]core.RecallHistorySummary{{ID: first.ID}, {ID: second.ID}}); len(ids) != 2 || ids[0] != first.ID || ids[1] != second.ID {
		t.Fatalf("recallHistoryIDs() = %v, want [%d %d]", ids, first.ID, second.ID)
	}
}

func TestQuoteSharePageViewResizeValidationAndClose(t *testing.T) {
	t.Parallel()

	engine := newShareTestEngine(t)
	quotes := make([]core.Quote, 0, 4)
	for i := 0; i < 4; i++ {
		quote, err := engine.AddQuote(context.Background(), strings.Repeat("shared quote ", i+1))
		if err != nil {
			t.Fatalf("AddQuote(%d) error = %v", i, err)
		}
		quotes = append(quotes, *quote)
	}

	page := NewQuoteSharePage(engine, 80, 24)
	page.Reset(quotes)
	page.SetSize(96, 30)

	msg := page.Init()()
	model, _ := page.Update(msg)
	page = model
	if !containsAll(page.View(), "Share Quotes", "Export Payload", "Save To") {
		t.Fatalf("share view missing expected content:\n%s", page.View())
	}
	if got := page.summary(); !strings.Contains(got, "Exporting 4 quote(s):") || !strings.Contains(got, "...and 1 more") {
		t.Fatalf("summary() = %q, want truncated export summary", got)
	}

	model, cmd := page.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	page = model
	if cmd != nil {
		t.Fatal("save command with empty path = non-nil, want validation failure")
	}
	if !strings.Contains(page.statusMsg, "Enter a file path") || !page.isErr {
		t.Fatalf("share validation status = %q err=%v, want missing-path error", page.statusMsg, page.isErr)
	}

	path := filepath.Join(t.TempDir(), "exports", "quotes.json")
	page.pathInput.SetValue(path)
	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	page = model
	if cmd == nil {
		t.Fatal("save command = nil")
	}
	model, _ = page.Update(cmd())
	page = model
	if !strings.Contains(page.statusMsg, path) || page.isErr {
		t.Fatalf("share save status = %q err=%v, want success for %q", page.statusMsg, page.isErr, path)
	}

	model, cmd = page.Update(tea.KeyMsg{Type: tea.KeyEsc})
	page = model
	if cmd == nil {
		t.Fatal("share close command = nil")
	}
	if _, ok := cmd().(CloseQuoteShareMsg); !ok {
		t.Fatalf("share close msg type = %T, want CloseQuoteShareMsg", cmd())
	}
}
