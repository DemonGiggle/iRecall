package pages

import (
	"context"
	"fmt"
	"strings"

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
	loading  bool
	errMsg   string
	width    int
	height   int
}

func NewQuotesPage(engine *core.Engine, width, height int) QuotesPage {
	vp := viewport.New(width-4, height-4)
	return QuotesPage{
		engine:  engine,
		viewport: vp,
		loading: true,
		width:   width,
		height:  height,
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
			p.viewport.SetContent(p.renderQuotes())
		}
		return p, nil

	case tea.KeyMsg:
		// Reload on 'r'.
		if msg.String() == "r" {
			p.loading = true
			p.errMsg = ""
			return p, p.loadQuotes()
		}
	}

	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

func (p QuotesPage) View() string {
	helpLine := styles.HelpBar.Render("↑/↓: Scroll   PgUp/PgDn: Page   R: Refresh   Tab: Switch page")

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
	innerW := p.viewport.Width - 2
	if innerW < 20 {
		innerW = 20
	}

	var sb strings.Builder
	sep := styles.Muted.Render(strings.Repeat("─", innerW))

	for i, q := range p.quotes {
		// Quote number + content.
		num := styles.QuoteNumber.Render(fmt.Sprintf("[%d]", i+1))
		content := lipgloss.NewStyle().Width(innerW).Render(q.Content)
		sb.WriteString(num + " " + content + "\n")

		// Tags line.
		if len(q.Tags) > 0 {
			tagStr := strings.Join(q.Tags, "  ·  ")
			sb.WriteString(styles.Muted.Render("    Tags: ") + styles.Accent.Render(tagStr) + "\n")
		} else {
			sb.WriteString(styles.Muted.Render("    Tags: (none)") + "\n")
		}

		if i < len(p.quotes)-1 {
			sb.WriteString(sep + "\n")
		}
	}
	return sb.String()
}
