package main

import "github.com/charmbracelet/lipgloss"

// Warm palette inspired by the companion website.
var (
	colorPrimary    = lipgloss.Color("#C4A882")
	colorForeground = lipgloss.Color("#E8E0D4")
	colorMuted      = lipgloss.Color("#8A8078")
	colorDimmed     = lipgloss.Color("#5C564F")
	colorBorder     = lipgloss.Color("#3A3632")

	// Tag colors, one per stack category.
	colorFrontend = lipgloss.Color("#E8B87A")
	colorBackend  = lipgloss.Color("#82B8A8")
	colorTools    = lipgloss.Color("#A88EC4")
)

func newStyles(width int, r *lipgloss.Renderer) styles {
	contentWidth := min(width-8, 72)

	return styles{
		// Tab bar
		activeTab: r.NewStyle().
			Bold(true).
			Foreground(colorForeground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 3),

		inactiveTab: r.NewStyle().
			Foreground(colorDimmed).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 3),

		// Content area
		content: r.NewStyle().
			Width(contentWidth),

		paragraph: r.NewStyle().
			Foreground(colorMuted).
			Width(contentWidth),

		sectionTitle: r.NewStyle().
			Foreground(colorForeground).
			Bold(true).
			MarginBottom(1),

		heading: r.NewStyle().
			Foreground(colorPrimary).
			Bold(true),

		card: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2).
			Width(contentWidth),

		cardTitle: r.NewStyle().
			Foreground(colorForeground).
			Bold(true),

		cardSubtitle: r.NewStyle().
			Foreground(colorMuted),

		cardPeriod: r.NewStyle().
			Foreground(colorDimmed),

		cardDescription: r.NewStyle().
			Foreground(colorMuted).
			Width(contentWidth - 8).
			MarginTop(1),

		selectedCard: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2).
			Width(contentWidth),

		tagFrontend: r.NewStyle().
			Foreground(colorFrontend),

		tagBackend: r.NewStyle().
			Foreground(colorBackend),

		tagTools: r.NewStyle().
			Foreground(colorTools),

		tagProject: r.NewStyle().
			Foreground(colorPrimary),

		// Status bar
		statusBar: r.NewStyle().
			Foreground(colorDimmed).
			Width(contentWidth).
			Align(lipgloss.Center),

		statusBarKey: r.NewStyle().
			Foreground(colorMuted).
			Bold(true),

		statusBarValue: r.NewStyle().
			Foreground(colorDimmed),

		statusBarLine: r.NewStyle().
			Foreground(colorBorder).
			Width(contentWidth),

		// Contact
		contactLabel: r.NewStyle().
			Foreground(colorDimmed).
			Width(10),

		contactValue: r.NewStyle().
			Foreground(colorForeground),

		contactCursor: r.NewStyle().
			Foreground(colorPrimary),

		projectStatus: r.NewStyle().
			Foreground(colorPrimary),

		// Stats
		statValue: r.NewStyle().
			Foreground(colorPrimary).
			Bold(true),

		statLabel: r.NewStyle().
			Foreground(colorDimmed),
	}
}

type styles struct {
	activeTab       lipgloss.Style
	inactiveTab     lipgloss.Style
	content         lipgloss.Style
	paragraph       lipgloss.Style
	sectionTitle    lipgloss.Style
	heading         lipgloss.Style
	card            lipgloss.Style
	cardTitle       lipgloss.Style
	cardSubtitle    lipgloss.Style
	cardPeriod      lipgloss.Style
	cardDescription lipgloss.Style
	selectedCard    lipgloss.Style
	tagFrontend     lipgloss.Style
	tagBackend      lipgloss.Style
	tagTools        lipgloss.Style
	tagProject      lipgloss.Style
	statusBar       lipgloss.Style
	statusBarKey    lipgloss.Style
	statusBarValue  lipgloss.Style
	statusBarLine   lipgloss.Style
	contactLabel    lipgloss.Style
	contactValue    lipgloss.Style
	contactCursor   lipgloss.Style
	projectStatus   lipgloss.Style
	statValue       lipgloss.Style
	statLabel       lipgloss.Style
}
