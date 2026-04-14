package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var bootLines = []string{
	"establishing connection...",
	"loading profile...",
	"arkadiusz.juszczyk — tech lead & product developer",
	"initializing interface...",
	"ready.",
}

type bootTickMsg time.Time

func bootTick() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg {
		return bootTickMsg(t)
	})
}

type bootModel struct {
	width     int
	height    int
	lineIndex int
	ready     bool
	done      bool
	renderer  *lipgloss.Renderer
}

func initialBootModel(r *lipgloss.Renderer) bootModel {
	return bootModel{renderer: r}
}

func (m bootModel) Init() tea.Cmd {
	return bootTick()
}

func (m bootModel) Update(msg tea.Msg) (bootModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case bootTickMsg:
		_ = msg
		m.lineIndex++
		if m.lineIndex >= len(bootLines) {
			m.ready = true
			return m, nil
		}
		return m, bootTick()

	case tea.KeyMsg:
		if m.ready {
			m.done = true
		} else {
			// Fast-forward if user presses key before animation finishes
			m.lineIndex = len(bootLines)
			m.ready = true
		}
		return m, nil
	}

	return m, nil
}

func (m bootModel) View() string {
	r := m.renderer
	if r == nil {
		r = lipgloss.DefaultRenderer()
	}
	prompt := r.NewStyle().Foreground(colorDimmed)
	text := r.NewStyle().Foreground(colorMuted)
	highlight := r.NewStyle().Foreground(colorPrimary)

	var lines []string
	for i := 0; i < m.lineIndex && i < len(bootLines); i++ {
		line := bootLines[i]
		p := prompt.Render("> ")
		if i == 2 {
			// The profile line gets highlighted
			lines = append(lines, p+highlight.Render(line))
		} else {
			lines = append(lines, p+text.Render(line))
		}
	}

	// Blinking cursor while booting
	if m.lineIndex < len(bootLines) {
		lines = append(lines, prompt.Render("> ")+text.Render("█"))
	}

	content := strings.Join(lines, "\n")

	// Prompt once boot is complete
	if m.ready {
		hint := r.NewStyle().
			Foreground(colorDimmed).
			MarginTop(2).
			Render("press any key to continue")
		content = content + "\n" + hint
	}

	content = fmt.Sprintf("\n%s\n", content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Center,
		r.NewStyle().MarginLeft(4).Render(content),
	)
}
