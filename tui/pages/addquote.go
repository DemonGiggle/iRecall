package pages

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gigol/irecall/core"
	"github.com/gigol/irecall/tui/styles"
)

type QuoteEditorMode int

const (
	QuoteEditorModeAdd QuoteEditorMode = iota
	QuoteEditorModeEdit
)

// OpenQuoteEditorMsg tells the app router to show the quote editor overlay.
type OpenQuoteEditorMsg struct {
	Mode  QuoteEditorMode
	Quote *core.Quote
}

// CloseQuoteEditorMsg tells the app router to dismiss the quote editor overlay.
type CloseQuoteEditorMsg struct {
	SavedQuote *core.Quote
}

// QuoteEditorDoneMsg signals the result of an add/edit operation.
type QuoteEditorDoneMsg struct {
	Quote    *core.Quote
	Warnings []string
	Err      error
}

// QuoteRefineDoneMsg signals the result of a draft refinement request.
type QuoteRefineDoneMsg struct {
	Refined string
	Err     error
}

// QuoteEditorPage is a modal overlay for adding or editing a quote.
type QuoteEditorPage struct {
	engine              *core.Engine
	textarea            textarea.Model
	spinner             spinner.Model
	mode                QuoteEditorMode
	editingID           int64
	busy                bool
	preview             bool
	statusMsg           string
	isErr               bool
	clearAt             time.Time
	original            string
	refined             string
	busyLabel           string
	attachmentMode      bool
	attachmentInput     textinput.Model
	retainedAttachments []core.QuoteAttachment
	newImages           []core.ImageInput
	attachmentCursor    int

	width  int
	height int
}

func NewQuoteEditorPage(engine *core.Engine, width, height int) QuoteEditorPage {
	ta := textarea.New()
	ta.Placeholder = "Type or paste your note here. Multi-line input supported."
	ta.Focus()
	ta.SetWidth(width - 12)
	ta.SetHeight(6)
	ta.CharLimit = 10000

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorAccent)
	pathInput := textinput.New()
	pathInput.Placeholder = "/path/to/image.png"
	pathInput.CharLimit = 4096

	return QuoteEditorPage{
		engine:          engine,
		textarea:        ta,
		spinner:         sp,
		attachmentInput: pathInput,
		width:           width,
		height:          height,
		mode:            QuoteEditorModeAdd,
	}
}

func (p QuoteEditorPage) Init() tea.Cmd {
	return textarea.Blink
}

func (p QuoteEditorPage) Update(msg tea.Msg) (QuoteEditorPage, tea.Cmd) {
	var cmds []tea.Cmd

	if !p.clearAt.IsZero() && time.Now().After(p.clearAt) {
		p.statusMsg = ""
		p.clearAt = time.Time{}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if p.attachmentMode {
			switch msg.String() {
			case "esc", "ctrl+o":
				p.attachmentMode = false
				p.attachmentInput.Blur()
				p.textarea.Focus()
				return p, nil
			case "enter":
				path := strings.TrimSpace(p.attachmentInput.Value())
				if path == "" {
					return p, nil
				}
				if len(p.retainedAttachments)+len(p.newImages) >= core.MaxQuoteImages {
					p.statusMsg = fmt.Sprintf("A quote can have at most %d images.", core.MaxQuoteImages)
					p.isErr = true
					return p, nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					p.statusMsg = "Error: " + err.Error()
					p.isErr = true
					return p, nil
				}
				p.newImages = append(p.newImages, core.ImageInput{Filename: filepath.Base(path), Data: data})
				p.attachmentCursor = len(p.retainedAttachments) + len(p.newImages) - 1
				p.attachmentInput.SetValue("")
				p.statusMsg = "Image staged. It will be copied when you save."
				p.isErr = false
				return p, nil
			case "up":
				if p.attachmentCursor > 0 {
					p.attachmentCursor--
				}
				return p, nil
			case "down":
				if p.attachmentCursor+1 < len(p.retainedAttachments)+len(p.newImages) {
					p.attachmentCursor++
				}
				return p, nil
			case "delete", "ctrl+d":
				if p.attachmentCursor < len(p.retainedAttachments) {
					p.retainedAttachments = append(p.retainedAttachments[:p.attachmentCursor], p.retainedAttachments[p.attachmentCursor+1:]...)
				} else {
					i := p.attachmentCursor - len(p.retainedAttachments)
					if i >= 0 && i < len(p.newImages) {
						p.newImages = append(p.newImages[:i], p.newImages[i+1:]...)
					}
				}
				if total := len(p.retainedAttachments) + len(p.newImages); total == 0 {
					p.attachmentCursor = 0
				} else if p.attachmentCursor >= total {
					p.attachmentCursor = total - 1
				}
				return p, nil
			}
			var cmd tea.Cmd
			p.attachmentInput, cmd = p.attachmentInput.Update(msg)
			return p, cmd
		}
		if p.preview {
			switch msg.String() {
			case "enter":
				p.textarea.SetValue(p.refined)
				p.preview = false
				p.original = ""
				p.refined = ""
				p.statusMsg = "Refined draft applied. Review and keep editing."
				p.isErr = false
				p.clearAt = time.Now().Add(3 * time.Second)
				p.textarea.Focus()
				return p, nil
			case "esc":
				p.preview = false
				p.original = ""
				p.refined = ""
				p.statusMsg = "Refined draft discarded."
				p.isErr = false
				p.clearAt = time.Now().Add(2 * time.Second)
				p.textarea.Focus()
				return p, nil
			}
		}

		switch msg.String() {
		case "esc":
			if !p.busy {
				return p, func() tea.Msg { return CloseQuoteEditorMsg{} }
			}
		case "ctrl+s":
			if p.busy {
				break
			}
			content := strings.TrimSpace(p.textarea.Value())
			if content == "" {
				p.statusMsg = "Nothing to save."
				p.isErr = true
				p.clearAt = time.Now().Add(2 * time.Second)
				break
			}
			p.busy = true
			p.busyLabel = "Saving quote and generating tags..."
			p.statusMsg = p.busyLabel
			p.isErr = false
			cmds = append(cmds, p.spinner.Tick, p.persistQuote(content))
		case "ctrl+r":
			if p.busy {
				break
			}
			content := strings.TrimSpace(p.textarea.Value())
			if content == "" {
				p.statusMsg = "Nothing to refine."
				p.isErr = true
				p.clearAt = time.Now().Add(2 * time.Second)
				break
			}
			p.busy = true
			p.busyLabel = "Refining draft..."
			p.statusMsg = p.busyLabel
			p.isErr = false
			cmds = append(cmds, p.spinner.Tick, p.refineQuote(content))
		case "ctrl+o":
			if !p.busy {
				p.attachmentMode = true
				p.textarea.Blur()
				p.attachmentInput.Focus()
				return p, textinput.Blink
			}
		}

	case QuoteEditorDoneMsg:
		p.busy = false
		p.busyLabel = ""
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.clearAt = time.Now().Add(4 * time.Second)
		} else {
			p.statusMsg = "Saved."
			if len(msg.Warnings) > 0 {
				p.statusMsg += " " + strings.Join(msg.Warnings, " ")
			}
			p.isErr = false
			p.clearAt = time.Now().Add(2 * time.Second)
			cmds = append(cmds, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return CloseQuoteEditorMsg{SavedQuote: msg.Quote}
			}))
		}

	case QuoteRefineDoneMsg:
		p.busy = false
		p.busyLabel = ""
		if msg.Err != nil {
			p.statusMsg = "Error: " + msg.Err.Error()
			p.isErr = true
			p.clearAt = time.Now().Add(4 * time.Second)
		} else {
			p.preview = true
			p.original = p.textarea.Value()
			p.refined = msg.Refined
			p.statusMsg = ""
			p.isErr = false
			p.clearAt = time.Time{}
			p.textarea.Blur()
		}

	case spinner.TickMsg:
		if p.busy {
			var cmd tea.Cmd
			p.spinner, cmd = p.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	p.textarea, cmd = p.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p QuoteEditorPage) View() string {
	helpLine := "  ctrl+s: Save   ctrl+r: Refine   ctrl+o: Images   esc: Cancel"
	if p.busy {
		helpLine = "  " + p.spinner.View() + " " + p.busyLabel
	} else if p.preview {
		helpLine = "  enter: Accept Refined Draft   esc: Reject and Continue Editing"
	}

	var statusLine string
	if p.statusMsg != "" {
		if p.isErr {
			statusLine = "\n  " + styles.StatusErr.Render(p.statusMsg)
		} else {
			statusLine = "\n  " + styles.StatusOK.Render(p.statusMsg)
		}
	}

	hint := styles.Muted.Render("  Tags will be regenerated automatically by the LLM.")
	title := " Add Quote "
	if p.mode == QuoteEditorModeEdit {
		title = " Edit Quote "
	}

	body := p.textarea.View()
	if p.attachmentMode {
		body = p.attachmentManagerView()
		hint = styles.Muted.Render("  Paste or drag an image path, then press enter. Images are copied only on Save.")
		helpLine = "  enter: Stage Path   ↑/↓: Select   delete/ctrl+d: Remove   esc: Back"
	}
	if p.preview {
		body = p.previewComparisonView()
		hint = styles.Muted.Render("  Compare your current draft with the suggested rewrite before applying it.")
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
		"\n",
		body,
		"\n",
		hint,
		statusLine,
		"\n",
		styles.HelpBar.Render(helpLine),
	)

	modalW := p.width - 8
	if modalW < 40 {
		modalW = 40
	}

	modal := styles.Modal.Width(modalW).Render(inner)

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left,
			styles.Bold.Foreground(styles.ColorPrimary).Render(title),
			modal,
		),
	)
}

func (p *QuoteEditorPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.textarea.SetWidth(width - 12)
}

func (p *QuoteEditorPage) Reset(mode QuoteEditorMode, quote *core.Quote) {
	p.mode = mode
	p.textarea.Reset()
	p.textarea.Focus()
	p.statusMsg = ""
	p.isErr = false
	p.busy = false
	p.preview = false
	p.attachmentMode = false
	p.attachmentInput.SetValue("")
	p.retainedAttachments = nil
	p.newImages = nil
	p.attachmentCursor = 0
	p.original = ""
	p.refined = ""
	p.clearAt = time.Time{}
	p.editingID = 0
	if quote != nil {
		p.editingID = quote.ID
		p.textarea.SetValue(quote.Content)
		p.retainedAttachments = append([]core.QuoteAttachment(nil), quote.Attachments...)
	}
}

func (p *QuoteEditorPage) persistQuote(content string) tea.Cmd {
	engine := p.engine
	mode := p.mode
	editingID := p.editingID
	return func() tea.Msg {
		retainedIDs := make([]string, len(p.retainedAttachments))
		for i, attachment := range p.retainedAttachments {
			retainedIDs[i] = attachment.ID
		}
		if mode == QuoteEditorModeEdit {
			result, err := engine.UpdateQuoteWithImages(context.Background(), editingID, content, retainedIDs, p.newImages)
			if err != nil {
				return QuoteEditorDoneMsg{Err: err}
			}
			return QuoteEditorDoneMsg{Quote: &result.Quote, Warnings: result.Warnings}
		}
		result, err := engine.AddQuoteWithImages(context.Background(), content, p.newImages)
		if err != nil {
			return QuoteEditorDoneMsg{Err: err}
		}
		return QuoteEditorDoneMsg{Quote: &result.Quote, Warnings: result.Warnings}
	}
}

func (p QuoteEditorPage) attachmentManagerView() string {
	total := len(p.retainedAttachments) + len(p.newImages)
	var lines []string
	for i, attachment := range p.retainedAttachments {
		prefix := "  "
		if i == p.attachmentCursor {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s  %dx%d  %s", prefix, attachment.Filename, attachment.Width, attachment.Height, formatAttachmentBytes(attachment.Size)))
	}
	for i, image := range p.newImages {
		index := len(p.retainedAttachments) + i
		prefix := "  "
		if index == p.attachmentCursor {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s  staged  %s", prefix, image.Filename, formatAttachmentBytes(int64(len(image.Data)))))
	}
	if len(lines) == 0 {
		lines = append(lines, "  No images attached.")
	}
	return styles.SectionHeader.Render("Image path") + "\n" + p.attachmentInput.View() + "\n\n" + styles.SectionHeader.Render(fmt.Sprintf("Attachments (%d/%d)", total, core.MaxQuoteImages)) + "\n" + strings.Join(lines, "\n")
}

func formatAttachmentBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.1f KiB", float64(size)/1024)
	}
	return fmt.Sprintf("%.1f MiB", float64(size)/(1024*1024))
}

func (p *QuoteEditorPage) refineQuote(content string) tea.Cmd {
	engine := p.engine
	return func() tea.Msg {
		refined, err := engine.RefineQuoteDraft(context.Background(), content)
		return QuoteRefineDoneMsg{Refined: refined, Err: err}
	}
}

func (p QuoteEditorPage) previewComparisonView() string {
	panelWidth := (p.width - 20) / 2
	if panelWidth < 24 {
		panelWidth = 24
	}

	current := styles.PanelActive.Width(panelWidth).Render(
		styles.SectionHeader.Render("Current Draft") + "\n" + strings.TrimSpace(p.original),
	)
	refined := styles.Panel.Width(panelWidth).Render(
		styles.SectionHeader.Render("Refined Draft") + "\n" + strings.TrimSpace(p.refined),
	)

	if p.width < 90 {
		return lipgloss.JoinVertical(lipgloss.Left, current, "", refined)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, current, "  ", refined)
}
