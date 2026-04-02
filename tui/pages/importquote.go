package pages

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

type OpenQuoteImportMsg struct{}

type CloseQuoteImportMsg struct {
	Reload bool
}

type QuoteImportDoneMsg struct {
	Result core.ImportResult
	Err    error
}

type QuoteImportPage struct {
	engine    *core.Engine
	pathInput textinput.Model
	busy      bool
	statusMsg string
	isErr     bool
	result    *core.ImportResult
	width     int
	height    int
}

func NewQuoteImportPage(engine *core.Engine, width, height int) QuoteImportPage {
	input := textinput.New()
	input.Placeholder = "/tmp/irecall-share.json"
	input.CharLimit = 4096
	input.Focus()
	return QuoteImportPage{
		engine:    engine,
		pathInput: input,
		width:     width,
		height:    height,
	}
}

func (p QuoteImportPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p QuoteImportPage) Update(msg tea.Msg) (QuoteImportPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseQuoteImportMsg{Reload: p.result != nil && !p.isErr} }
			}
		case "enter", "ctrl+s":
			if p.busy {
				return p, nil
			}
			path := strings.TrimSpace(p.pathInput.Value())
			if path == "" {
				p.statusMsg = "Enter a file path to import."
				p.isErr = true
				return p, nil
			}
			p.busy = true
			p.statusMsg = ""
			p.result = nil
			return p, p.importPayload(path)
		}
	case QuoteImportDoneMsg:
		p.busy = false
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.result = nil
			return p, nil
		}
		p.result = &msg.Result
		p.statusMsg = fmt.Sprintf("Imported quotes. inserted=%d updated=%d duplicates=%d stale=%d",
			msg.Result.Inserted, msg.Result.Updated, msg.Result.Duplicates, msg.Result.Stale)
		p.isErr = false
		return p, nil
	}

	var cmd tea.Cmd
	p.pathInput, cmd = p.pathInput.Update(msg)
	return p, cmd
}

func (p QuoteImportPage) View() string {
	title := styles.Bold.Foreground(styles.ColorPrimary).Render(" Import Quotes ")
	help := "  enter/ctrl+s: Import file   esc: Close"
	if p.busy {
		help = "  Importing..."
	}
	status := ""
	if p.statusMsg != "" {
		if p.isErr {
			status = styles.StatusErr.Render(p.statusMsg)
		} else {
			status = styles.StatusOK.Render(p.statusMsg)
		}
	}

	summary := styles.QuoteItem.Render("Import a quote share JSON file exported from another iRecall instance.")
	pathField := styles.PanelActive.Width(max(44, p.width-24)).Render(p.pathInput.View())

	var resultBox string
	if p.result != nil {
		resultBox = styles.Panel.Width(max(44, p.width-24)).Render(
			styles.SectionHeader.Render("Import Summary") + "\n" +
				fmt.Sprintf("Inserted: %d\nUpdated: %d\nDuplicates: %d\nStale: %d",
					p.result.Inserted, p.result.Updated, p.result.Duplicates, p.result.Stale),
		)
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		summary,
		"",
		styles.SectionHeader.Render("Import From"),
		pathField,
		styles.Muted.Render("  The file must contain the iRecall share envelope exported by another user."),
		"",
		resultBox,
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

func (p *QuoteImportPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.pathInput.Width = max(36, width-28)
}

func (p *QuoteImportPage) Reset() {
	p.pathInput.Focus()
	p.pathInput.SetValue("")
	p.busy = false
	p.statusMsg = ""
	p.isErr = false
	p.result = nil
}

func (p QuoteImportPage) importPayload(path string) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		payload, err := os.ReadFile(path)
		if err != nil {
			return QuoteImportDoneMsg{Err: fmt.Errorf("read import file: %w", err)}
		}
		result, err := engine.ImportSharedQuotes(context.Background(), payload)
		return QuoteImportDoneMsg{Result: result, Err: err}
	}
}
