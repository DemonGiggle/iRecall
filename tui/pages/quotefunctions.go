package pages

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

const quoteListEntryActions = "↑/↓: Move   x: Select   a: Select all   u: Deselect all   e: Edit   d: Delete   s: Share"
const quoteDetailEntryActions = "enter/esc: Back   ↑/↓: Scroll   pgup/pgdn: Page   x: Select   e: Edit   d: Delete   s: Share"

type quoteListActionKind int

const (
	quoteListActionNone quoteListActionKind = iota
	quoteListActionEdit
	quoteListActionDelete
	quoteListActionShare
)

type quoteListAction struct {
	kind   quoteListActionKind
	quote  *core.Quote
	quotes []core.Quote
}

type quoteSelection struct {
	cursor   int
	selected map[int64]bool
}

func newQuoteSelection() quoteSelection {
	return quoteSelection{selected: map[int64]bool{}}
}

func (s *quoteSelection) clamp(quotes []core.Quote) {
	if len(quotes) == 0 {
		s.cursor = 0
		s.selected = map[int64]bool{}
		return
	}
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor >= len(quotes) {
		s.cursor = len(quotes) - 1
	}

	valid := make(map[int64]bool, len(quotes))
	for _, q := range quotes {
		valid[q.ID] = true
	}
	for id := range s.selected {
		if !valid[id] {
			delete(s.selected, id)
		}
	}
}

func (s *quoteSelection) move(delta int, quotes []core.Quote) {
	if len(quotes) == 0 {
		return
	}
	s.cursor += delta
	if s.cursor < 0 {
		s.cursor = 0
	}
	if s.cursor >= len(quotes) {
		s.cursor = len(quotes) - 1
	}
}

func (s *quoteSelection) current(quotes []core.Quote) *core.Quote {
	if len(quotes) == 0 || s.cursor < 0 || s.cursor >= len(quotes) {
		return nil
	}
	q := quotes[s.cursor]
	return &q
}

func (s *quoteSelection) toggleCurrent(quotes []core.Quote) {
	q := s.current(quotes)
	if q == nil {
		return
	}
	if s.selected[q.ID] {
		delete(s.selected, q.ID)
		return
	}
	s.selected[q.ID] = true
}

func (s *quoteSelection) clear() {
	s.selected = map[int64]bool{}
}

func (s *quoteSelection) selectAll(quotes []core.Quote) {
	s.selected = make(map[int64]bool, len(quotes))
	for _, q := range quotes {
		s.selected[q.ID] = true
	}
}

func (s *quoteSelection) selectedQuotes(quotes []core.Quote) []core.Quote {
	out := make([]core.Quote, 0, len(s.selected))
	for _, q := range quotes {
		if s.selected[q.ID] {
			out = append(out, q)
		}
	}
	if len(out) > 0 {
		return out
	}
	if q := s.current(quotes); q != nil {
		return []core.Quote{*q}
	}
	return nil
}

func (s *quoteSelection) selectedCount() int {
	return len(s.selected)
}

type quoteListWidget struct {
	title          string
	panelWidth     int
	bodyHeight     int
	listViewport   viewport.Model
	detailViewport viewport.Model
	quotes         []core.Quote
	selection      quoteSelection
	detail         bool
}

func newQuoteListWidget(title string, panelWidth, bodyHeight int) quoteListWidget {
	w := quoteListWidget{
		title:     title,
		selection: newQuoteSelection(),
	}
	w.SetSize(panelWidth, bodyHeight)
	return w
}

func (w *quoteListWidget) SetTitle(title string) {
	w.title = title
}

func (w *quoteListWidget) SetSize(panelWidth, bodyHeight int) {
	w.panelWidth = panelWidth
	if w.panelWidth < 22 {
		w.panelWidth = 22
	}
	w.bodyHeight = max(3, bodyHeight)
	innerWidth := max(20, w.panelWidth-2)
	if w.listViewport.Width == 0 && w.listViewport.Height == 0 {
		w.listViewport = viewport.New(innerWidth, w.bodyHeight)
		w.detailViewport = viewport.New(innerWidth, w.bodyHeight)
	} else {
		w.listViewport.Width = innerWidth
		w.listViewport.Height = w.bodyHeight
		w.detailViewport.Width = innerWidth
		w.detailViewport.Height = w.bodyHeight
	}
	w.refresh()
}

func (w *quoteListWidget) SetQuotes(quotes []core.Quote) {
	w.quotes = quotes
	w.selection.clamp(w.quotes)
	if len(w.quotes) == 0 {
		w.detail = false
	}
	w.refresh()
}

func (w *quoteListWidget) ClearQuotes() {
	w.quotes = nil
	w.selection.clear()
	w.detail = false
	w.refresh()
}

func (w *quoteListWidget) ApplyQuoteUpdate(updated core.Quote) {
	for i := range w.quotes {
		if w.quotes[i].ID == updated.ID {
			w.quotes[i] = updated
			w.refresh()
			return
		}
	}
}

func (w *quoteListWidget) RemoveQuotes(ids []int64) {
	if len(ids) == 0 || len(w.quotes) == 0 {
		return
	}
	remove := idsSet(ids)
	filtered := w.quotes[:0]
	for _, q := range w.quotes {
		if !remove[q.ID] {
			filtered = append(filtered, q)
		}
	}
	w.quotes = filtered
	w.selection.clamp(w.quotes)
	if len(w.quotes) == 0 {
		w.detail = false
	}
	w.refresh()
}

func (w *quoteListWidget) Update(msg tea.Msg) (quoteListAction, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if w.detail {
				w.detail = false
				w.refresh()
				return quoteListAction{}, nil
			}
			if w.selection.current(w.quotes) != nil {
				w.detail = true
				w.refresh()
			}
			return quoteListAction{}, nil
		case "esc":
			if w.detail {
				w.detail = false
				w.refresh()
				return quoteListAction{}, nil
			}
		case "up":
			if w.detail {
				break
			}
			w.selection.move(-1, w.quotes)
			w.refresh()
			return quoteListAction{}, nil
		case "down":
			if w.detail {
				break
			}
			w.selection.move(1, w.quotes)
			w.refresh()
			return quoteListAction{}, nil
		case "x":
			w.selection.toggleCurrent(w.quotes)
			w.refresh()
			return quoteListAction{}, nil
		case "a":
			if w.detail {
				return quoteListAction{}, nil
			}
			w.selection.selectAll(w.quotes)
			w.refresh()
			return quoteListAction{}, nil
		case "u":
			if w.detail {
				return quoteListAction{}, nil
			}
			w.selection.clear()
			w.refresh()
			return quoteListAction{}, nil
		case "e":
			if q := w.selection.current(w.quotes); q != nil {
				quote := *q
				return quoteListAction{kind: quoteListActionEdit, quote: &quote}, nil
			}
		case "d":
			selected := w.selection.selectedQuotes(w.quotes)
			if len(selected) > 0 {
				return quoteListAction{kind: quoteListActionDelete, quotes: selected}, nil
			}
		case "s":
			selected := w.selection.selectedQuotes(w.quotes)
			if len(selected) > 0 {
				return quoteListAction{kind: quoteListActionShare, quotes: selected}, nil
			}
		}
	}

	var cmd tea.Cmd
	if w.detail {
		w.detailViewport, cmd = w.detailViewport.Update(msg)
		return quoteListAction{}, cmd
	}
	w.listViewport, cmd = w.listViewport.Update(msg)
	return quoteListAction{}, cmd
}

func (w quoteListWidget) View(focused bool, inactiveHelp, navigationHint string) string {
	body := w.listViewport.View()
	actionHelp := quoteListEntryActions
	if w.detail {
		body = w.detailViewport.View()
		actionHelp = quoteDetailEntryActions
	}
	return renderQuoteListPanel(
		w.title,
		body,
		w.panelWidth,
		w.currentBodyHeight(),
		focused,
		inactiveHelp,
		navigationHint,
		actionHelp,
	)
}

func (w quoteListWidget) currentQuote() *core.Quote {
	return w.selection.current(w.quotes)
}

func (w quoteListWidget) selectedCount() int {
	return w.selection.selectedCount()
}

func (w quoteListWidget) currentCursor() int {
	return w.selection.cursor
}

func (w quoteListWidget) isDetail() bool {
	return w.detail
}

func (w quoteListWidget) yOffset() int {
	return w.listViewport.YOffset
}

func (w *quoteListWidget) refresh() {
	w.selection.clamp(w.quotes)
	w.listViewport.SetContent(renderQuoteFunctionList(w.quotes, w.selection, w.listViewport.Width))
	w.ensureCursorVisible()
	w.detailViewport.SetContent(w.renderDetail())
	if w.detail {
		w.detailViewport.GotoTop()
	}
}

func (w *quoteListWidget) ensureCursorVisible() {
	start, end := quoteEntryLineRange(w.quotes, w.selection.cursor)
	if start < 0 || end < 0 || w.listViewport.Height <= 0 {
		w.listViewport.SetYOffset(0)
		return
	}
	if start < w.listViewport.YOffset {
		w.listViewport.SetYOffset(start)
		return
	}
	bottom := w.listViewport.YOffset + w.listViewport.Height - 1
	if end > bottom {
		w.listViewport.SetYOffset(end - w.listViewport.Height + 1)
	}
}

func (w quoteListWidget) renderDetail() string {
	q := w.selection.current(w.quotes)
	if q == nil {
		return styles.Muted.Render("  No quote selected.")
	}

	var parts []string
	parts = append(parts, styles.SectionHeader.Render("Quote Information"))
	parts = append(parts, styles.Accent.Render(fmt.Sprintf("Quote [%d]", w.selection.cursor+1)))
	parts = append(parts, lipgloss.NewStyle().Width(max(20, w.panelWidth-6)).Render(q.Content))

	if !q.IsOwnedByMe && q.SourceName != "" {
		parts = append(parts, styles.Muted.Render("From: ")+styles.Accent.Render(q.SourceName))
	}

	if len(q.Tags) > 0 {
		parts = append(parts, styles.Muted.Render("Tags: ")+styles.Accent.Render(strings.Join(q.Tags, "  ·  ")))
	} else {
		parts = append(parts, styles.Muted.Render("Tags: (none)"))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (w quoteListWidget) currentBodyHeight() int {
	if w.detail {
		return w.detailViewport.Height
	}
	return w.listViewport.Height
}

func renderQuoteListPanel(title, body string, panelWidth, bodyHeight int, focused bool, inactiveHelp, navigationHint, actionHelp string) string {
	label := styles.Accent.Render(title)
	help := inactiveHelp
	panelStyle := styles.Panel
	if focused {
		label = styles.Bold.Foreground(styles.ColorAccent).Render(title)
		help = joinQuoteListHelp(navigationHint, actionHelp)
		panelStyle = styles.PanelActive
	}
	if bodyHeight < 3 {
		bodyHeight = 3
	}
	return panelStyle.Width(panelWidth).Height(bodyHeight + 5).Render(
		label + "\n" + body + "\n\n" + styles.HelpBar.Render(help),
	)
}

func joinQuoteListHelp(navigationHint, actionHelp string) string {
	parts := make([]string, 0, 2)
	if navigationHint != "" {
		parts = append(parts, navigationHint)
	}
	if actionHelp != "" {
		parts = append(parts, actionHelp)
	}
	return strings.Join(parts, "   ")
}

func renderQuoteFunctionList(quotes []core.Quote, selection quoteSelection, innerW int) string {
	if len(quotes) == 0 {
		return styles.Muted.Render("No quotes available.")
	}
	if innerW < 20 {
		innerW = 20
	}

	var sb strings.Builder
	sep := styles.Muted.Render(strings.Repeat("─", innerW))

	for i, q := range quotes {
		prefix := "  "
		if i == selection.cursor {
			prefix = styles.Accent.Render("> ")
		}

		check := "[ ]"
		if selection.selected[q.ID] {
			check = "[x]"
		}
		check = styles.QuoteNumber.Render(check)
		number := styles.QuoteNumber.Render(fmt.Sprintf("[%d]", i+1))
		content := truncateQuotePreview(q.Content, innerW-10)
		sb.WriteString(prefix + check + " " + number + " " + content + "\n")

		if !q.IsOwnedByMe && q.SourceName != "" {
			sb.WriteString(styles.Muted.Render("    From: ") + styles.Accent.Render(q.SourceName) + "\n")
		}

		if len(q.Tags) > 0 {
			tagStr := previewTags(q.Tags, 3)
			sb.WriteString(styles.Muted.Render("    Tags: ") + styles.Accent.Render(tagStr) + "\n")
		} else {
			sb.WriteString(styles.Muted.Render("    Tags: (none)") + "\n")
		}

		if i < len(quotes)-1 {
			sb.WriteString(sep + "\n")
		}
	}
	return sb.String()
}

func quoteEntryLineRange(quotes []core.Quote, index int) (start, end int) {
	if index < 0 || index >= len(quotes) {
		return -1, -1
	}
	line := 0
	for i, q := range quotes {
		itemStart := line
		line++ // quote preview
		if !q.IsOwnedByMe && q.SourceName != "" {
			line++
		}
		line++ // tags line
		itemEnd := line - 1
		if i == index {
			return itemStart, itemEnd
		}
		if i < len(quotes)-1 {
			line++ // separator
		}
	}
	return -1, -1
}

func previewTags(tags []string, limit int) string {
	if len(tags) == 0 {
		return ""
	}
	if limit <= 0 || len(tags) <= limit {
		return strings.Join(tags, "  ·  ")
	}
	return strings.Join(tags[:limit], "  ·  ") + fmt.Sprintf("  ·  +%d more", len(tags)-limit)
}

func truncateQuotePreview(content string, width int) string {
	if width < 8 {
		width = 8
	}
	flat := strings.Join(strings.Fields(content), " ")
	if lipgloss.Width(flat) <= width {
		return flat
	}
	if width <= 1 {
		return "…"
	}
	truncated := []rune(flat)
	if len(truncated) > width-1 {
		truncated = truncated[:width-1]
	}
	return strings.TrimSpace(string(truncated)) + "…"
}

func idsSet(ids []int64) map[int64]bool {
	out := make(map[int64]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out
}

func selectedIDsFromMap(m map[int64]bool) []int64 {
	ids := make([]int64, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}
