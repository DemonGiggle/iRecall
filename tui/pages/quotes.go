package pages

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// QuotesLoadedMsg carries one page of quotes fetched from the DB.
type QuotesLoadedMsg struct {
	Quotes     []core.Quote
	TotalCount int64
	Offset     int
	Err        error
}

// QuotesPage lists all stored quotes with their tags.
type QuotesPage struct {
	engine     *core.Engine
	quoteList  quoteListWidget
	loading    bool
	errMsg     string
	totalCount int64
	offset     int
	pageSize   int
	width      int
	height     int
}

func NewQuotesPage(engine *core.Engine, width, height int) QuotesPage {
	page := QuotesPage{
		engine:    engine,
		loading:   true,
		width:     width,
		height:    height,
		pageSize:  quoteListPageSize,
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
			p.errMsg = ""
			p.totalCount = msg.TotalCount
			p.offset = msg.Offset
			p.quoteList.SetQuotes(msg.Quotes)
			p.quoteList.SetTitle(p.title())
		}
		return p, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			return p, func() tea.Msg { return OpenQuoteEditorMsg{Mode: QuoteEditorModeAdd} }
		case "i":
			return p, func() tea.Msg { return OpenQuoteImportMsg{} }
		case "r":
			if p.loading {
				return p, nil
			}
			p.loading = true
			p.errMsg = ""
			return p, p.loadQuotes()
		case "left":
			if !p.loading && !p.quoteList.isDetail() && p.canMovePrevPage() {
				p.loading = true
				p.errMsg = ""
				p.offset -= p.pageSize
				return p, p.loadQuotes()
			}
		case "right":
			if !p.loading && !p.quoteList.isDetail() && p.canMoveNextPage() {
				p.loading = true
				p.errMsg = ""
				p.offset += p.pageSize
				return p, p.loadQuotes()
			}
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
	pageHelp := "ctrl+n: Add Quote   i: Import   r: Refresh   ←/→: Result Page   pgup/pgdn: Scroll   tab/shift+tab: Switch Page"

	panel := p.quoteList.View(true, "", p.pageNavigationHint())
	switch {
	case p.loading:
		panel = styles.PanelActive.Width(p.width - 4).Height(p.quoteList.currentBodyHeight() + 5).Render(
			styles.Bold.Foreground(styles.ColorAccent).Render(p.title()) + "\n" +
				styles.Muted.Render("  Loading quotes...") + "\n\n" +
				styles.HelpBar.Render(joinQuoteListHelp(p.pageNavigationHint(), quoteListEntryActions)),
		)
	case p.errMsg != "":
		panel = styles.PanelActive.Width(p.width - 4).Height(p.quoteList.currentBodyHeight() + 5).Render(
			styles.Bold.Foreground(styles.ColorAccent).Render(p.title()) + "\n" +
				styles.StatusErr.Render("  "+p.errMsg) + "\n\n" +
				styles.HelpBar.Render(joinQuoteListHelp(p.pageNavigationHint(), quoteListEntryActions)),
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
	offset := p.offset
	pageSize := p.pageSize
	return func() tea.Msg {
		totalCount, err := engine.CountQuotes(context.Background())
		if err != nil {
			return QuotesLoadedMsg{Err: err}
		}
		offset = clampPageOffset(totalCount, pageSize, offset)
		quotes, err := engine.ListQuotesPage(context.Background(), pageSize, offset)
		return QuotesLoadedMsg{Quotes: quotes, TotalCount: totalCount, Offset: offset, Err: err}
	}
}

func (p *QuotesPage) ApplyQuoteUpdate(updated core.Quote) {
	p.quoteList.ApplyQuoteUpdate(updated)
}

func (p *QuotesPage) RemoveQuotes(ids []int64) {
	p.quoteList.RemoveQuotes(ids)
	if len(ids) > 0 {
		p.totalCount -= int64(len(ids))
		if p.totalCount < 0 {
			p.totalCount = 0
		}
	}
	p.quoteList.SetTitle(p.title())
}

func (p QuotesPage) title() string {
	return formatPagedTitle("Stored Quotes", p.totalCount, p.pageSize, p.offset)
}

func (p QuotesPage) pageNavigationHint() string {
	if p.quoteList.isDetail() {
		return ""
	}
	return fmt.Sprintf("←/→: Result Page %d/%d", max(1, currentPage(p.offset, p.pageSize)), max(1, totalPages(p.totalCount, p.pageSize)))
}

func (p QuotesPage) canMovePrevPage() bool {
	return p.offset > 0
}

func (p QuotesPage) canMoveNextPage() bool {
	return p.offset+p.pageSize < int(p.totalCount)
}
