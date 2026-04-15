package pages

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

type HistoryLoadedMsg struct {
	Entries []core.RecallHistorySummary
	Err     error
}

type HistoryDetailLoadedMsg struct {
	Entry *core.RecallHistoryEntry
	Err   error
}

type OpenDeleteRecallHistoryMsg struct {
	Entries []core.RecallHistorySummary
}

type RecallHistoryQuoteSavedMsg struct {
	Quote *core.Quote
	Err   error
}

type historyListSelection struct {
	cursor   int
	selected map[int64]bool
}

func newHistoryListSelection() historyListSelection {
	return historyListSelection{selected: map[int64]bool{}}
}

func (s *historyListSelection) clamp(entries []core.RecallHistorySummary) {
	if len(entries) == 0 {
		s.cursor = 0
		s.selected = map[int64]bool{}
		return
	}
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor >= len(entries) {
		s.cursor = len(entries) - 1
	}
	valid := make(map[int64]bool, len(entries))
	for _, entry := range entries {
		valid[entry.ID] = true
	}
	for id := range s.selected {
		if !valid[id] {
			delete(s.selected, id)
		}
	}
}

func (s *historyListSelection) move(delta int, entries []core.RecallHistorySummary) {
	if len(entries) == 0 {
		return
	}
	s.cursor += delta
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor >= len(entries) {
		s.cursor = len(entries) - 1
	}
}

func (s *historyListSelection) current(entries []core.RecallHistorySummary) *core.RecallHistorySummary {
	if len(entries) == 0 || s.cursor < 0 || s.cursor >= len(entries) {
		return nil
	}
	entry := entries[s.cursor]
	return &entry
}

func (s *historyListSelection) toggleCurrent(entries []core.RecallHistorySummary) {
	entry := s.current(entries)
	if entry == nil {
		return
	}
	if s.selected[entry.ID] {
		delete(s.selected, entry.ID)
		return
	}
	s.selected[entry.ID] = true
}

func (s *historyListSelection) selectAll(entries []core.RecallHistorySummary) {
	s.selected = make(map[int64]bool, len(entries))
	for _, entry := range entries {
		s.selected[entry.ID] = true
	}
}

func (s *historyListSelection) clear() {
	s.selected = map[int64]bool{}
}

func (s *historyListSelection) selectedEntries(entries []core.RecallHistorySummary) []core.RecallHistorySummary {
	out := make([]core.RecallHistorySummary, 0, len(s.selected))
	for _, entry := range entries {
		if s.selected[entry.ID] {
			out = append(out, entry)
		}
	}
	if len(out) > 0 {
		return out
	}
	if entry := s.current(entries); entry != nil {
		return []core.RecallHistorySummary{*entry}
	}
	return nil
}

type historyFocus int

const (
	historyFocusDetail historyFocus = iota
	historyFocusReferenceQuotes
)

type HistoryPage struct {
	engine *core.Engine

	listViewport   viewport.Model
	detailViewport viewport.Model
	refViewport    viewport.Model

	entries       []core.RecallHistorySummary
	entry         *core.RecallHistoryEntry
	selection     historyListSelection
	quoteFns      quoteSelection
	detail        bool
	detailLoading bool
	loading       bool
	focus         historyFocus
	errMsg        string
	statusMsg     string
	statusErr     bool

	width  int
	height int
}

func NewHistoryPage(engine *core.Engine, width, height int) HistoryPage {
	page := HistoryPage{
		engine:    engine,
		selection: newHistoryListSelection(),
		quoteFns:  newQuoteSelection(),
		loading:   true,
		width:     width,
		height:    height,
		focus:     historyFocusDetail,
	}
	page.recalcLayout()
	return page
}

func (p HistoryPage) Init() tea.Cmd {
	return p.loadHistory()
}

func (p HistoryPage) Update(msg tea.Msg) (HistoryPage, tea.Cmd) {
	switch msg := msg.(type) {
	case HistoryLoadedMsg:
		p.loading = false
		if msg.Err != nil {
			p.errMsg = "Error loading history: " + msg.Err.Error()
			return p, nil
		}
		p.entries = msg.Entries
		p.selection.clamp(p.entries)
		p.listViewport.SetContent(p.renderList())
		return p, nil

	case HistoryDetailLoadedMsg:
		p.detailLoading = false
		if msg.Err != nil {
			p.errMsg = "Error loading history entry: " + msg.Err.Error()
			return p, nil
		}
		p.entry = msg.Entry
		p.focus = historyFocusDetail
		p.quoteFns.clear()
		p.quoteFns.clamp(p.entry.Quotes)
		p.refreshDetail()
		return p, nil

	case tea.KeyMsg:
		if !p.detail {
			switch msg.String() {
			case "r":
				p.loading = true
				p.errMsg = ""
				return p, p.loadHistory()
			case "up":
				p.selection.move(-1, p.entries)
				p.listViewport.SetContent(p.renderList())
				return p, nil
			case "down":
				p.selection.move(1, p.entries)
				p.listViewport.SetContent(p.renderList())
				return p, nil
			case "x":
				p.selection.toggleCurrent(p.entries)
				p.listViewport.SetContent(p.renderList())
				return p, nil
			case "a":
				p.selection.selectAll(p.entries)
				p.listViewport.SetContent(p.renderList())
				return p, nil
			case "u":
				p.selection.clear()
				p.listViewport.SetContent(p.renderList())
				return p, nil
			case "d":
				selected := p.selection.selectedEntries(p.entries)
				if len(selected) > 0 {
					return p, func() tea.Msg { return OpenDeleteRecallHistoryMsg{Entries: selected} }
				}
			case "enter":
				if current := p.selection.current(p.entries); current != nil {
					p.detail = true
					p.detailLoading = true
					p.errMsg = ""
					p.entry = nil
					return p, p.loadHistoryDetail(current.ID)
				}
			}
		} else {
			switch msg.String() {
			case "esc", "enter":
				p.detail = false
				p.detailLoading = false
				p.focus = historyFocusDetail
				p.quoteFns.clear()
				return p, nil
			case "ctrl+j":
				if p.focus == historyFocusDetail {
					p.focus = historyFocusReferenceQuotes
				} else {
					p.focus = historyFocusDetail
				}
				return p, nil
			case "ctrl+s":
				if p.entry != nil {
					p.statusMsg = ""
					p.statusErr = false
					return p, p.saveHistoryAsQuote()
				}
			case "up":
				if p.focus == historyFocusReferenceQuotes {
					p.quoteFns.move(-1, p.currentQuotes())
					p.refreshReferenceQuotes()
					return p, nil
				}
			case "down":
				if p.focus == historyFocusReferenceQuotes {
					p.quoteFns.move(1, p.currentQuotes())
					p.refreshReferenceQuotes()
					return p, nil
				}
			case "x":
				if p.focus == historyFocusReferenceQuotes {
					p.quoteFns.toggleCurrent(p.currentQuotes())
					p.refreshReferenceQuotes()
					return p, nil
				}
			case "e":
				if p.focus == historyFocusReferenceQuotes {
					if q := p.quoteFns.current(p.currentQuotes()); q != nil {
						quote := *q
						return p, func() tea.Msg {
							return OpenQuoteEditorMsg{Mode: QuoteEditorModeEdit, Quote: &quote}
						}
					}
				}
			case "d":
				if p.focus == historyFocusReferenceQuotes {
					selected := p.quoteFns.selectedQuotes(p.currentQuotes())
					if len(selected) > 0 {
						return p, func() tea.Msg { return OpenDeleteQuotesMsg{Quotes: selected} }
					}
				}
			case "s":
				if p.focus == historyFocusReferenceQuotes {
					selected := p.quoteFns.selectedQuotes(p.currentQuotes())
					if len(selected) > 0 {
						return p, func() tea.Msg { return OpenQuoteShareMsg{Quotes: selected} }
					}
				}
			}
		}
	}

	switch msg := msg.(type) {
	case RecallHistoryQuoteSavedMsg:
		if msg.Err != nil {
			p.statusMsg = "Error saving history as quote: " + msg.Err.Error()
			p.statusErr = true
		} else {
			p.statusMsg = "Saved history entry as quote."
			p.statusErr = false
		}
		return p, nil
	}

	var cmd tea.Cmd
	if !p.detail {
		p.listViewport, cmd = p.listViewport.Update(msg)
		return p, cmd
	}
	if p.focus == historyFocusDetail {
		p.detailViewport, cmd = p.detailViewport.Update(msg)
		return p, cmd
	}
	p.refViewport, cmd = p.refViewport.Update(msg)
	return p, cmd
}

func (p HistoryPage) View() string {
	help := "↑/↓: Move   enter: View   x: Select   a: Select all   u: Deselect all   d: Delete   r: Refresh"
	body := styles.Muted.Render("  Loading history...")
	if p.errMsg != "" {
		body = styles.StatusErr.Render("  " + p.errMsg)
	} else if !p.loading {
		if len(p.entries) == 0 {
			body = styles.Muted.Render("  No recall history yet.")
		} else if !p.detail {
			body = p.listViewport.View()
		} else if p.detailLoading {
			body = styles.Muted.Render("  Loading history entry...")
		} else {
			help = "enter/esc: Back   ctrl+j: Toggle focus   x: Select quote   e: Edit quote   d: Delete quote   s: Share quote"
			if p.entry != nil {
				help = "enter/esc: Back   ctrl+s: Save Q/A as Quote   ctrl+j: Toggle focus   x: Select quote   e: Edit quote   d: Delete quote   s: Share quote"
			}
			body = p.detailView()
		}
	}

	status := ""
	if p.statusMsg != "" {
		if p.statusErr {
			status = styles.StatusErr.Render("  " + p.statusMsg)
		} else {
			status = styles.StatusOK.Render("  " + p.statusMsg)
		}
	}

	return styles.Panel.Width(p.width - 4).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			styles.SectionHeader.Render(fmt.Sprintf("History (%d)", len(p.entries))),
			body,
			status,
			"",
			styles.HelpBar.Render(help),
		),
	)
}

func (p *HistoryPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.recalcLayout()
	if !p.detail {
		p.listViewport.SetContent(p.renderList())
		return
	}
	p.refreshDetail()
}

func (p *HistoryPage) Reload() tea.Cmd {
	p.loading = true
	p.errMsg = ""
	return p.loadHistory()
}

func (p *HistoryPage) ApplyQuoteUpdate(updated core.Quote) {
	if p.entry == nil {
		return
	}
	for i := range p.entry.Quotes {
		if p.entry.Quotes[i].ID == updated.ID {
			p.entry.Quotes[i] = updated
			p.refreshDetail()
			return
		}
	}
}

func (p *HistoryPage) RemoveQuotes(ids []int64) {
	if p.entry == nil || len(ids) == 0 {
		return
	}
	remove := idsSet(ids)
	filtered := p.entry.Quotes[:0]
	for _, q := range p.entry.Quotes {
		if !remove[q.ID] {
			filtered = append(filtered, q)
		}
	}
	p.entry.Quotes = filtered
	p.quoteFns.clamp(p.entry.Quotes)
	p.refreshDetail()
}

func (p *HistoryPage) RemoveHistories(ids []int64) {
	if len(ids) == 0 || len(p.entries) == 0 {
		return
	}
	remove := idsSet(ids)
	filtered := p.entries[:0]
	for _, entry := range p.entries {
		if !remove[entry.ID] {
			filtered = append(filtered, entry)
		}
	}
	p.entries = filtered
	p.selection.clamp(p.entries)
	if p.entry != nil && remove[p.entry.ID] {
		p.entry = nil
		p.detail = false
		p.detailLoading = false
	}
	p.listViewport.SetContent(p.renderList())
}

func (p *HistoryPage) recalcLayout() {
	innerW := p.width - 6
	if innerW < 20 {
		innerW = 20
	}
	p.listViewport = viewport.New(innerW, max(4, p.height-8))

	remaining := p.height - 13
	detailH := remaining * 2 / 3
	refH := remaining - detailH
	if detailH < 5 {
		detailH = 5
	}
	if refH < 4 {
		refH = 4
	}
	p.detailViewport = viewport.New(innerW, detailH)
	p.refViewport = viewport.New(innerW, refH)
}

func (p *HistoryPage) refreshDetail() {
	if p.entry == nil {
		p.detailViewport.SetContent("")
		p.refViewport.SetContent("")
		return
	}
	var detail strings.Builder
	detail.WriteString(styles.Muted.Render("Created: ") + styles.Accent.Render(formatHistoryTime(p.entry.CreatedAt)) + "\n\n")
	detail.WriteString(styles.SectionHeader.Render("Question") + "\n")
	detail.WriteString(p.entry.Question + "\n\n")
	detail.WriteString(styles.SectionHeader.Render("Response") + "\n")
	detail.WriteString(p.entry.Response)
	p.detailViewport.SetContent(detail.String())
	p.refreshReferenceQuotes()
}

func (p *HistoryPage) refreshReferenceQuotes() {
	p.quoteFns.clamp(p.currentQuotes())
	p.refViewport.SetContent(renderQuoteFunctionList(p.currentQuotes(), p.quoteFns, p.refViewport.Width, true))
}

func (p HistoryPage) detailView() string {
	detailStyle := styles.Panel
	refStyle := styles.Panel
	refHelp := "ctrl+j: Focus detail"
	if p.focus == historyFocusDetail {
		detailStyle = styles.PanelActive
	} else {
		refStyle = styles.PanelActive
		refHelp = "ctrl+j: Focus detail   ↑/↓: Move   x: Select   e: Edit   d: Delete   s: Share"
	}
	detailBox := detailStyle.Width(p.width - 8).Height(p.detailViewport.Height + 3).Render(
		styles.Accent.Render("History Entry") + "\n" + p.detailViewport.View(),
	)
	refBox := refStyle.Width(p.width - 8).Height(p.refViewport.Height + 5).Render(
		styles.Accent.Render("Reference Quotes") + "\n" + p.refViewport.View() + "\n\n" + styles.HelpBar.Render(refHelp),
	)
	return lipgloss.JoinVertical(lipgloss.Left, detailBox, refBox)
}

func (p HistoryPage) renderList() string {
	if len(p.entries) == 0 {
		return styles.Muted.Render("No recall history yet.")
	}
	var sb strings.Builder
	sep := styles.Muted.Render(strings.Repeat("─", max(20, p.listViewport.Width)))
	for i, entry := range p.entries {
		prefix := "  "
		if i == p.selection.cursor {
			prefix = styles.Accent.Render("> ")
		}
		check := "[ ]"
		if p.selection.selected[entry.ID] {
			check = "[x]"
		}
		sb.WriteString(prefix + styles.QuoteNumber.Render(check) + " " +
			styles.QuoteNumber.Render(formatHistoryTime(entry.CreatedAt)) + " " +
			truncateQuotePreview(entry.Question, max(20, p.listViewport.Width-26)) + "\n")
		resp := strings.TrimSpace(entry.Response)
		if resp == "" {
			resp = "(empty response)"
		}
		sb.WriteString(styles.Muted.Render("    Response: ") + truncateQuotePreview(resp, max(20, p.listViewport.Width-16)) + "\n")
		if i < len(p.entries)-1 {
			sb.WriteString(sep + "\n")
		}
	}
	return sb.String()
}

func (p *HistoryPage) loadHistory() tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		entries, err := engine.ListRecallHistory(context.Background())
		return HistoryLoadedMsg{Entries: entries, Err: err}
	}
}

func (p *HistoryPage) loadHistoryDetail(id int64) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		entry, err := engine.GetRecallHistory(context.Background(), id)
		return HistoryDetailLoadedMsg{Entry: entry, Err: err}
	}
}

func (p *HistoryPage) saveHistoryAsQuote() tea.Cmd {
	engine := p.engine
	question := p.entry.Question
	response := p.entry.Response
	return func() tea.Msg {
		quote, err := engine.SaveRecallAsQuote(context.Background(), question, response, nil)
		return RecallHistoryQuoteSavedMsg{Quote: quote, Err: err}
	}
}

func (p *HistoryPage) currentQuotes() []core.Quote {
	if p.entry == nil {
		return nil
	}
	return p.entry.Quotes
}

func formatHistoryTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04")
}
