package pages

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// QuotesLoadedMsg carries all quotes fetched from the DB.
type QuotesLoadedMsg struct {
	Quotes []core.Quote
	Err    error
}

// QuotesPage lists all stored quotes with their tags.
type QuotesPage struct {
	engine   *core.Engine
	viewport viewport.Model
	quotes   []core.Quote
	quoteFns quoteSelection
	loading  bool
	errMsg   string
	width    int
	height   int
}

func NewQuotesPage(engine *core.Engine, width, height int) QuotesPage {
	vp := viewport.New(width-4, height-4)
	return QuotesPage{
		engine:   engine,
		viewport: vp,
		quoteFns: newQuoteSelection(),
		loading:  true,
		width:    width,
		height:   height,
	}
}

func (p QuotesPage) Init() tea.Cmd {
	return p.loadQuotes()
}

func (p QuotesPage) Update(msg tea.Msg) (QuotesPage, tea.Cmd) {
	switch msg := msg.(type) {
	case QuotesLoadedMsg:
		p.loading = false
		if msg.Err != nil {
			p.errMsg = "Error loading quotes: " + msg.Err.Error()
		} else {
			p.quotes = msg.Quotes
			p.quoteFns.clamp(p.quotes)
			p.viewport.SetContent(p.renderQuotes())
		}
		return p, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			p.loading = true
			p.errMsg = ""
			return p, p.loadQuotes()
		case "up":
			p.quoteFns.move(-1, p.quotes)
			p.viewport.SetContent(p.renderQuotes())
			return p, nil
		case "down":
			p.quoteFns.move(1, p.quotes)
			p.viewport.SetContent(p.renderQuotes())
			return p, nil
		case "x":
			p.quoteFns.toggleCurrent(p.quotes)
			p.viewport.SetContent(p.renderQuotes())
			return p, nil
		case "e":
			if q := p.quoteFns.current(p.quotes); q != nil {
				quote := *q
				return p, func() tea.Msg {
					return OpenQuoteEditorMsg{Mode: QuoteEditorModeEdit, Quote: &quote}
				}
			}
		case "d":
			selected := p.quoteFns.selectedQuotes(p.quotes)
			if len(selected) > 0 {
				return p, func() tea.Msg { return OpenDeleteQuotesMsg{Quotes: selected} }
			}
		}
	}

	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p QuotesPage) View() string {
	helpLine := styles.HelpBar.Render("↑/↓: Move   X: Select   E: Edit   D: Delete   R: Refresh   PgUp/PgDn: Page")

	var body string
	switch {
	case p.loading:
		body = styles.Muted.Render("  Loading quotes...")
	case p.errMsg != "":
		body = styles.StatusErr.Render("  " + p.errMsg)
	case len(p.quotes) == 0:
		body = styles.Muted.Render("  No quotes yet. Press Ctrl+N on the Recall page to add one.")
	default:
		body = p.viewport.View()
	}

	return styles.Panel.Width(p.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			styles.SectionHeader.Render(fmt.Sprintf("Stored Quotes (%d)", len(p.quotes))),
			body,
			"",
			helpLine,
		),
	)
}

func (p *QuotesPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.viewport.Width = width - 6
	p.viewport.Height = height - 6
	if len(p.quotes) > 0 {
		p.viewport.SetContent(p.renderQuotes())
	}
}

// Reload refreshes the quote list from the DB.
func (p *QuotesPage) Reload() tea.Cmd {
	p.loading = true
	p.errMsg = ""
	return p.loadQuotes()
}

func (p *QuotesPage) loadQuotes() tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		quotes, err := engine.ListQuotes(context.Background())
		return QuotesLoadedMsg{Quotes: quotes, Err: err}
	}
}

func (p *QuotesPage) renderQuotes() string {
	return renderQuoteFunctionList(p.quotes, p.quoteFns, p.viewport.Width-2, true)
}

func (p *QuotesPage) ApplyQuoteUpdate(updated core.Quote) {
	for i := range p.quotes {
		if p.quotes[i].ID == updated.ID {
			p.quotes[i] = updated
			p.viewport.SetContent(p.renderQuotes())
			return
		}
	}
}

func (p *QuotesPage) RemoveQuotes(ids []int64) {
	if len(ids) == 0 || len(p.quotes) == 0 {
		return
	}
	remove := idsSet(ids)
	filtered := p.quotes[:0]
	for _, q := range p.quotes {
		if !remove[q.ID] {
			filtered = append(filtered, q)
		}
	}
	p.quotes = filtered
	p.quoteFns.clamp(p.quotes)
	p.viewport.SetContent(p.renderQuotes())
}
