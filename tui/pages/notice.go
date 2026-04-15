package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/tui/styles"
)

// OpenNoticeMsg tells the app router to show a notice modal.
type OpenNoticeMsg struct {
	Title   string
	Message string
}

// CloseNoticeMsg tells the app router to dismiss the notice modal.
type CloseNoticeMsg struct{}

// NoticePage is a simple modal used for explicit success notifications.
type NoticePage struct {
	title   string
	message string
	width   int
	height  int
}

func NewNoticePage(width, height int) NoticePage {
	return NoticePage{
		title:  "Done",
		width:  width,
		height: height,
	}
}

func (p *NoticePage) Reset(title, message string) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "Done"
	}
	p.title = title
	p.message = strings.TrimSpace(message)
}

func (p *NoticePage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p NoticePage) Init() tea.Cmd {
	return nil
}

func (p NoticePage) Update(msg tea.Msg) (NoticePage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc", " ":
			return p, func() tea.Msg { return CloseNoticeMsg{} }
		}
	}
	return p, nil
}

func (p NoticePage) View() string {
	body := lipgloss.JoinVertical(lipgloss.Left,
		styles.SectionHeader.Render(p.title),
		"",
		lipgloss.NewStyle().Foreground(styles.ColorFg).Width(maxInt(36, minInt(p.width-16, 72))).Render(p.message),
		"",
		styles.HelpBar.Render("enter/esc: Close"),
	)

	modalW := p.width - 20
	if modalW < 44 {
		modalW = 44
	}

	return styles.Modal.Width(modalW).Render(body)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
