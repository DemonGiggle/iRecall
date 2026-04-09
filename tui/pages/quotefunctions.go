package pages

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

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

func renderQuoteFunctionList(quotes []core.Quote, selection quoteSelection, innerW int, showTags bool) string {
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

		if showTags {
			if len(q.Tags) > 0 {
				tagStr := previewTags(q.Tags, 3)
				sb.WriteString(styles.Muted.Render("    Tags: ") + styles.Accent.Render(tagStr) + "\n")
			} else {
				sb.WriteString(styles.Muted.Render("    Tags: (none)") + "\n")
			}
		}

		if i < len(quotes)-1 {
			sb.WriteString(sep + "\n")
		}
	}
	return sb.String()
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
