package pages

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

// OpenQuoteShareMsg tells the app router to show the quote share overlay.
type OpenQuoteShareMsg struct {
	Quotes []core.Quote
}

// CloseQuoteShareMsg tells the app router to dismiss the quote share overlay.
type CloseQuoteShareMsg struct{}

// QuoteShareLoadedMsg carries the exported share payload.
type QuoteShareLoadedMsg struct {
	Payload string
	Err     error
}

// QuoteShareSavedMsg carries the result of saving the share payload to disk.
type QuoteShareSavedMsg struct {
	Path string
	Err  error
}

// QuoteSharePage previews exported quote payloads and saves them to a file.
type QuoteSharePage struct {
	engine    *core.Engine
	quotes    []core.Quote
	payload   string
	pathInput textinput.Model
	viewport  viewport.Model
	busy      bool
	statusMsg string
	isErr     bool
	width     int
	height    int
}

func NewQuoteSharePage(engine *core.Engine, width, height int) QuoteSharePage {
	input := textinput.New()
	input.Placeholder = "/tmp/irecall-share.json"
	input.CharLimit = 4096
	input.Focus()

	vp := viewport.New(max(40, width-20), max(8, height-18))

	return QuoteSharePage{
		engine:    engine,
		pathInput: input,
		viewport:  vp,
		width:     width,
		height:    height,
	}
}

func (p QuoteSharePage) Init() tea.Cmd {
	return p.exportQuotes()
}

func (p QuoteSharePage) Update(msg tea.Msg) (QuoteSharePage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseQuoteShareMsg{} }
			}
		case "enter", "ctrl+s":
			if p.busy {
				return p, nil
			}
			path := strings.TrimSpace(p.pathInput.Value())
			if path == "" {
				p.statusMsg = "Enter a file path to save the share payload."
				p.isErr = true
				return p, nil
			}
			if strings.TrimSpace(p.payload) == "" {
				p.statusMsg = "Share payload is not ready yet."
				p.isErr = true
				return p, nil
			}
			p.busy = true
			p.statusMsg = ""
			return p, p.savePayload(path)
		}

	case QuoteShareLoadedMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.payload = ""
			p.viewport.SetContent("")
			return p, nil
		}
		p.payload = msg.Payload
		p.viewport.SetContent(msg.Payload)
		p.statusMsg = "Share payload ready. Save it to a file to send it to another user."
		p.isErr = false
		return p, nil

	case QuoteShareSavedMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			return p, nil
		}
		p.statusMsg = "Saved share payload to " + msg.Path
		p.isErr = false
		return p, nil
	}

	var cmd tea.Cmd
	p.pathInput, cmd = p.pathInput.Update(msg)
	p.viewport, _ = p.viewport.Update(msg)
	return p, cmd
}

func (p QuoteSharePage) View() string {
	title := styles.Bold.Foreground(styles.ColorPrimary).Render(" Share Quotes ")
	help := "  ctrl+s: Save export file   esc: Close"
	if p.busy {
		help = "  Exporting..."
	}

	status := ""
	if p.statusMsg != "" {
		if p.isErr {
			status = styles.StatusErr.Render(p.statusMsg)
		} else {
			status = styles.StatusOK.Render(p.statusMsg)
		}
	}

	summary := styles.QuoteItem.Render(p.summary())
	pathField := styles.PanelActive.Width(max(44, p.width-24)).Render(p.pathInput.View())
	payloadBox := styles.Panel.Width(max(44, p.width-24)).
		Height(max(10, p.height-18)).
		Render(styles.SectionHeader.Render("Export Payload") + "\n" + p.viewport.View())

	body := lipgloss.JoinVertical(lipgloss.Left,
		summary,
		"",
		styles.SectionHeader.Render("Save To"),
		pathField,
		styles.Muted.Render("  Export to a JSON file and transfer it manually to the recipient."),
		"",
		payloadBox,
		"",
		status,
		styles.HelpBar.Render(help),
	)

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			styles.Modal.Width(max(54, p.width-12)).Render(body),
		),
	)
}

func (p *QuoteSharePage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.pathInput.Width = max(36, width-28)
	p.viewport.Width = max(40, width-28)
	p.viewport.Height = max(8, height-22)
	if p.payload != "" {
		p.viewport.SetContent(p.payload)
	}
}

func (p *QuoteSharePage) Reset(quotes []core.Quote) {
	p.quotes = append([]core.Quote(nil), quotes...)
	p.payload = ""
	p.busy = true
	p.statusMsg = ""
	p.isErr = false
	p.pathInput.Focus()
	p.pathInput.SetValue("")
	p.viewport.SetContent("")
}

func (p QuoteSharePage) summary() string {
	if len(p.quotes) == 0 {
		return "No quotes selected."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Exporting %d quote(s):\n\n", len(p.quotes)))
	limit := len(p.quotes)
	if limit > 3 {
		limit = 3
	}
	for i := 0; i < limit; i++ {
		sb.WriteString(fmt.Sprintf("[%d] v%d %s\n", i+1, p.quotes[i].Version, truncateQuoteContent(p.quotes[i].Content, 100)))
	}
	if len(p.quotes) > limit {
		sb.WriteString(fmt.Sprintf("...and %d more", len(p.quotes)-limit))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func (p QuoteSharePage) exportQuotes() tea.Cmd {
	engine := p.engine
	ids := quoteIDs(p.quotes)
	return func() tea.Msg {
		payload, err := engine.ExportQuotes(context.Background(), ids)
		return QuoteShareLoadedMsg{Payload: string(payload), Err: err}
	}
}

func (p QuoteSharePage) savePayload(path string) tea.Cmd {
	payload := p.payload
	return func() tea.Msg {
		dir := filepath.Dir(path)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o700); err != nil {
				return QuoteShareSavedMsg{Err: fmt.Errorf("create share directory: %w", err)}
			}
		}
		if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
			return QuoteShareSavedMsg{Err: fmt.Errorf("write share file: %w", err)}
		}
		return QuoteShareSavedMsg{Path: path}
	}
}
