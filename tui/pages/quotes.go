package pages

import (
	"context"
	"fmt"

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
	engine    *core.Engine
	quoteList quoteListWidget
	loading   bool
	errMsg    string
	width     int
	height    int
}

func NewQuotesPage(engine *core.Engine, width, height int) QuotesPage {
	page := QuotesPage{
		engine:    engine,
		loading:   true,
		width:     width,
		height:    height,
		quoteList: newQuoteListWidget("Stored Quotes (0)", width-4, max(3, height-7)),
	}
	page.recalcLayout()
	return page
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
			p.quoteList.SetQuotes(msg.Quotes)
			p.quoteList.SetTitle(fmt.Sprintf("Stored Quotes (%d)", len(msg.Quotes)))
		}
		return p, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			return p, func() tea.Msg { return OpenQuoteEditorMsg{Mode: QuoteEditorModeAdd} }
		case "i":
			return p, func() tea.Msg { return OpenQuoteImportMsg{} }
		case "r":
			p.loading = true
			p.errMsg = ""
			return p, p.loadQuotes()
		}
	}

	action, cmd := p.quoteList.Update(msg)
	switch action.kind {
	case quoteListActionEdit:
		return p, func() tea.Msg {
			return OpenQuoteEditorMsg{Mode: QuoteEditorModeEdit, Quote: action.quote}
		}
	case quoteListActionDelete:
		return p, func() tea.Msg { return OpenDeleteQuotesMsg{Quotes: action.quotes} }
	case quoteListActionShare:
		return p, func() tea.Msg { return OpenQuoteShareMsg{Quotes: action.quotes} }
	}
	return p, cmd
}

func (p QuotesPage) View() string {
	pageHelp := "ctrl+n: Add Quote   i: Import   r: Refresh   pgup/pgdn: Page   tab/shift+tab: Switch Page"

	panel := p.quoteList.View(true, "", "")
	switch {
	case p.loading:
		panel = styles.PanelActive.Width(p.width - 4).Height(p.quoteList.currentBodyHeight() + 5).Render(
			styles.Bold.Foreground(styles.ColorAccent).Render("Stored Quotes (0)") + "\n" +
				styles.Muted.Render("  Loading quotes...") + "\n\n" +
				styles.HelpBar.Render(quoteListEntryActions),
		)
	case p.errMsg != "":
		panel = styles.PanelActive.Width(p.width - 4).Height(p.quoteList.currentBodyHeight() + 5).Render(
			styles.Bold.Foreground(styles.ColorAccent).Render("Stored Quotes (0)") + "\n" +
				styles.StatusErr.Render("  "+p.errMsg) + "\n\n" +
				styles.HelpBar.Render(quoteListEntryActions),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.HelpBar.Render(pageHelp),
		panel,
	)
}

func (p *QuotesPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.recalcLayout()
}

func (p *QuotesPage) recalcLayout() {
	p.quoteList.SetSize(p.width-4, max(3, p.height-7))
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

func (p *QuotesPage) ApplyQuoteUpdate(updated core.Quote) {
	p.quoteList.ApplyQuoteUpdate(updated)
}

func (p *QuotesPage) RemoveQuotes(ids []int64) {
	p.quoteList.RemoveQuotes(ids)
	p.quoteList.SetTitle(fmt.Sprintf("Stored Quotes (%d)", len(p.quoteList.quotes)))
}
