package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary  = lipgloss.Color("#7C3AED") // violet
	ColorAccent   = lipgloss.Color("#A78BFA") // light violet
	ColorMuted    = lipgloss.Color("#6B7280") // gray
	ColorSuccess  = lipgloss.Color("#10B981") // green
	ColorError    = lipgloss.Color("#EF4444") // red
	ColorWarning  = lipgloss.Color("#F59E0B") // amber
	ColorFg       = lipgloss.Color("#F9FAFB") // near white
	ColorBorder   = lipgloss.Color("#374151") // dark gray
	ColorSelected = lipgloss.Color("#1F2937") // darker bg for active

	// Base styles
	Bold   = lipgloss.NewStyle().Bold(true)
	Muted  = lipgloss.NewStyle().Foreground(ColorMuted)
	Accent = lipgloss.NewStyle().Foreground(ColorAccent)

	// App title bar
	TitleBar = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	// Active / inactive tab indicators
	TabActive = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Underline(true)

	TabInactive = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Panel borders
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	PanelActive = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	// Section headers inside panels
	SectionHeader = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			MarginBottom(1)

	// Status / help bar at bottom
	HelpBar = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1)

	// Status messages
	StatusOK  = lipgloss.NewStyle().Foreground(ColorSuccess)
	StatusErr = lipgloss.NewStyle().Foreground(ColorError)

	// Quote items in reference list
	QuoteItem = lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingLeft(1)

	QuoteNumber = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// Modal overlay
	Modal = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2)

	// Form labels
	FormLabel = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Width(20)

	// Button (focused)
	ButtonFocused = lipgloss.NewStyle().
			Foreground(ColorFg).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPrimary).
			Background(ColorPrimary).
			Padding(0, 2).
			Bold(true)

	// Button (unfocused)
	ButtonNormal = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 2)
)
