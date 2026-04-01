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

// OpenDeleteQuotesMsg tells the app router to show the delete confirmation overlay.
type OpenDeleteQuotesMsg struct {
	Quotes []core.Quote
}

// CloseDeleteQuotesMsg tells the app router to dismiss the delete confirmation overlay.
type CloseDeleteQuotesMsg struct {
	DeletedIDs []int64
}

// DeleteQuotesDoneMsg is emitted after a delete operation completes.
type DeleteQuotesDoneMsg struct {
	DeletedIDs []int64
	Err        error
}

// DeleteQuotesPage confirms and executes quote deletion.
type DeleteQuotesPage struct {
	engine    *core.Engine
	quotes    []core.Quote
	busy      bool
	statusMsg string
	width     int
	height    int
}

func NewDeleteQuotesPage(engine *core.Engine, width, height int) DeleteQuotesPage {
	return DeleteQuotesPage{
		engine: engine,
		width:  width,
		height: height,
	}
}

func (p DeleteQuotesPage) Init() tea.Cmd {
	return nil
}

func (p DeleteQuotesPage) Update(msg tea.Msg) (DeleteQuotesPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseDeleteQuotesMsg{} }
			}
		case "enter":
			if !p.busy && len(p.quotes) > 0 {
				p.busy = true
				p.statusMsg = ""
				return p, p.deleteQuotes()
			}
		}

	case DeleteQuotesDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			return p, nil
		}
		return p, func() tea.Msg { return CloseDeleteQuotesMsg{DeletedIDs: msg.DeletedIDs} }
	}

	return p, nil
}

func (p DeleteQuotesPage) View() string {
	title := styles.Bold.Foreground(styles.ColorError).Render(" Delete Quotes ")
	copy := styles.QuoteItem.Render(p.summary())
	help := styles.HelpBar.Render("Enter: Confirm delete   Esc: Cancel")
	if p.busy {
		help = styles.HelpBar.Render("Deleting...")
	}

	var status string
	if p.statusMsg != "" {
		status = styles.StatusErr.Render(p.statusMsg)
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
		copy,
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

func (p *DeleteQuotesPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *DeleteQuotesPage) Reset(quotes []core.Quote) {
	p.quotes = append([]core.Quote(nil), quotes...)
	p.busy = false
	p.statusMsg = ""
}

func (p *DeleteQuotesPage) deleteQuotes() tea.Cmd {
	engine := p.engine
	ids := quoteIDs(p.quotes)
	return func() tea.Msg {
		err := engine.DeleteQuotes(context.Background(), ids)
		return DeleteQuotesDoneMsg{DeletedIDs: ids, Err: err}
	}
}

func (p DeleteQuotesPage) summary() string {
	switch len(p.quotes) {
	case 0:
		return "No quotes selected."
	case 1:
		return fmt.Sprintf("Delete this quote?\n\n%s", truncateQuoteContent(p.quotes[0].Content, 180))
	default:
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Delete %d quotes?\n\n", len(p.quotes)))
		limit := len(p.quotes)
		if limit > 3 {
			limit = 3
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, truncateQuoteContent(p.quotes[i].Content, 120)))
		}
		if len(p.quotes) > limit {
			sb.WriteString(fmt.Sprintf("...and %d more", len(p.quotes)-limit))
		}
		return strings.TrimRight(sb.String(), "\n")
	}
}

func truncateQuoteContent(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func quoteIDs(quotes []core.Quote) []int64 {
	ids := make([]int64, 0, len(quotes))
	for _, q := range quotes {
		ids = append(ids, q.ID)
	}
	return ids
}
