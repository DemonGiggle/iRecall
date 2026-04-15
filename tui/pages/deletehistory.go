package pages

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

type CloseDeleteRecallHistoryMsg struct {
	DeletedIDs []int64
}

type DeleteRecallHistoryDoneMsg struct {
	DeletedIDs []int64
	Err        error
}

type DeleteRecallHistoryPage struct {
	engine    *core.Engine
	entries   []core.RecallHistorySummary
	busy      bool
	statusMsg string
	width     int
	height    int
}

func NewDeleteRecallHistoryPage(engine *core.Engine, width, height int) DeleteRecallHistoryPage {
	return DeleteRecallHistoryPage{engine: engine, width: width, height: height}
}

func (p DeleteRecallHistoryPage) Init() tea.Cmd { return nil }

func (p DeleteRecallHistoryPage) Update(msg tea.Msg) (DeleteRecallHistoryPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseDeleteRecallHistoryMsg{} }
			}
		case "enter":
			if !p.busy && len(p.entries) > 0 {
				p.busy = true
				p.statusMsg = ""
				return p, p.deleteHistory()
			}
		}
	case DeleteRecallHistoryDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			return p, nil
		}
		return p, func() tea.Msg { return CloseDeleteRecallHistoryMsg{DeletedIDs: msg.DeletedIDs} }
	}
	return p, nil
}

func (p DeleteRecallHistoryPage) View() string {
	title := styles.Bold.Foreground(styles.ColorError).Render(" Delete History ")
	help := styles.HelpBar.Render("enter: Confirm delete   esc: Cancel")
	if p.busy {
		help = styles.HelpBar.Render("Deleting...")
	}

	var status string
	if p.statusMsg != "" {
		status = styles.StatusErr.Render(p.statusMsg)
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
		styles.QuoteItem.Render(p.summary()),
		"",
		styles.Muted.Render("This action cannot be undone."),
		"",
		status,
		help,
	)

	modalW := p.width - 20
	if modalW < 50 {
		modalW = 50
	}

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			styles.Modal.Width(modalW).Render(inner),
		),
	)
}

func (p *DeleteRecallHistoryPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *DeleteRecallHistoryPage) Reset(entries []core.RecallHistorySummary) {
	p.entries = append([]core.RecallHistorySummary(nil), entries...)
	p.busy = false
	p.statusMsg = ""
}

func (p *DeleteRecallHistoryPage) deleteHistory() tea.Cmd {
	engine := p.engine
	ids := recallHistoryIDs(p.entries)
	return func() tea.Msg {
		err := engine.DeleteRecallHistory(context.Background(), ids)
		return DeleteRecallHistoryDoneMsg{DeletedIDs: ids, Err: err}
	}
}

func (p DeleteRecallHistoryPage) summary() string {
	switch len(p.entries) {
	case 0:
		return "No history selected."
	case 1:
		return fmt.Sprintf("Delete this history entry?\n\n%s", truncateQuoteContent(p.entries[0].Question, 180))
	default:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Delete %d history entries?\n\n", len(p.entries)))
		limit := len(p.entries)
		if limit > 3 {
			limit = 3
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, truncateQuoteContent(p.entries[i].Question, 120)))
		}
		if len(p.entries) > limit {
			sb.WriteString(fmt.Sprintf("...and %d more", len(p.entries)-limit))
		}
		return strings.TrimRight(sb.String(), "\n")
	}
}

func recallHistoryIDs(entries []core.RecallHistorySummary) []int64 {
	ids := make([]int64, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.ID)
	}
	return ids
}
