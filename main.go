package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
)

var tabNames = []string{"about", "experience", "stack", "projects", "contact"}

type tab int

const (
	tabAbout tab = iota
	tabExperience
	tabStack
	tabProjects
	tabContact
)

type clearClipboardMsg struct{}

type model struct {
	width          int
	height         int
	activeTab      tab
	styles         styles
	renderer       *lipgloss.Renderer
	cursor         int
	expanded       int
	konamiBuffer   string
	showEasterEgg  bool
	booting        bool
	boot           bootModel
	clipboardLabel string
	clipboardURL   string
}

func initialModel(renderer *lipgloss.Renderer) model {
	return model{
		width:     80,
		height:    24,
		activeTab: tabAbout,
		renderer:  renderer,
		styles:    newStyles(80, 24, renderer),
		expanded:  -1,
		booting:   true,
		boot:      initialBootModel(renderer),
	}
}

func (m model) Init() tea.Cmd {
	return m.boot.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.booting {
		if wsm, ok := msg.(tea.WindowSizeMsg); ok {
			m.width = wsm.Width
			m.height = wsm.Height
			m.styles = newStyles(wsm.Width, wsm.Height, m.renderer)
			m.boot.width = wsm.Width
			m.boot.height = wsm.Height
		}
		var cmd tea.Cmd
		m.boot, cmd = m.boot.Update(msg)
		if m.boot.done {
			m.booting = false
			return m, nil
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles = newStyles(msg.Width, msg.Height, m.renderer)
		return m, nil

	case clearClipboardMsg:
		m.clipboardLabel = ""
		m.clipboardURL = ""
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Easter egg tracking
		m.konamiBuffer += key
		if len(m.konamiBuffer) > 40 {
			m.konamiBuffer = m.konamiBuffer[len(m.konamiBuffer)-40:]
		}
		if strings.Contains(m.konamiBuffer, "upupdowndownleftrightleftrightba") {
			m.showEasterEgg = true
			m.konamiBuffer = ""
			return m, nil
		}
		if strings.HasSuffix(m.konamiBuffer, "coffee") {
			m.showEasterEgg = true
			m.konamiBuffer = ""
			return m, nil
		}

		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.showEasterEgg {
				m.showEasterEgg = false
				return m, nil
			}
			if m.expanded >= 0 {
				m.expanded = -1
				return m, nil
			}
			return m, tea.Quit

		case "tab", "l", "right":
			if m.showEasterEgg {
				return m, nil
			}
			m.activeTab = tab((int(m.activeTab) + 1) % len(tabNames))
			m.cursor = 0
			m.expanded = -1
			return m, nil

		case "shift+tab", "h", "left":
			if m.showEasterEgg {
				return m, nil
			}
			m.activeTab = tab((int(m.activeTab) - 1 + len(tabNames)) % len(tabNames))
			m.cursor = 0
			m.expanded = -1
			return m, nil

		case "j", "down":
			if !m.showEasterEgg {
				m = m.moveCursorDown()
			}
			return m, nil

		case "k", "up":
			if !m.showEasterEgg {
				m = m.moveCursorUp()
			}
			return m, nil

		case "enter":
			if m.showEasterEgg {
				m.showEasterEgg = false
				return m, nil
			}
			if m.activeTab == tabContact && m.cursor < len(contacts) {
				c := contacts[m.cursor]
				m.clipboardLabel = c.Label
				m.clipboardURL = c.URL
				return m, clearClipboardAfter(3 * time.Second)
			}
			m = m.toggleExpand()
			return m, nil

		case "?":
			m.showEasterEgg = !m.showEasterEgg
			return m, nil

		case "1", "2", "3", "4", "5":
			if !m.showEasterEgg {
				idx := int(key[0] - '1')
				if idx < len(tabNames) {
					m.activeTab = tab(idx)
					m.cursor = 0
					m.expanded = -1
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m model) moveCursorDown() model {
	max := m.maxCursorItems()
	if m.cursor < max-1 {
		m.cursor++
	}
	return m
}

func (m model) moveCursorUp() model {
	if m.cursor > 0 {
		m.cursor--
	}
	return m
}

func (m model) maxCursorItems() int {
	switch m.activeTab {
	case tabExperience:
		return len(experience)
	case tabProjects:
		return len(projects)
	case tabContact:
		return len(contacts)
	default:
		return 0
	}
}

func (m model) toggleExpand() model {
	switch m.activeTab {
	case tabExperience, tabProjects:
		if m.expanded == m.cursor {
			m.expanded = -1
		} else {
			m.expanded = m.cursor
		}
	}
	return m
}

func clearClipboardAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearClipboardMsg{}
	})
}

// osc52Copy returns the OSC 52 escape sequence to copy text to the client's clipboard.
func osc52Copy(text string) string {
	return fmt.Sprintf("\x1b]52;c;%s\x1b\\", base64.StdEncoding.EncodeToString([]byte(text)))
}

// View renders the full screen.
// Layout: centered tabs at top third, content in middle, status bar pinned to bottom.
func (m model) View() string {
	if m.booting {
		return m.boot.View()
	}
	if m.showEasterEgg {
		return m.viewEasterEgg()
	}

	tabs := m.viewTabs()
	content := m.viewContent()
	status := m.viewStatusBar()

	tabsCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, tabs)
	contentCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, content)
	statusCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, status)

	// Fixed layout:
	// - 2 blank lines from top
	// - tabs
	// - 2 blank lines
	// - content (fills remaining space)
	// - status bar pinned to bottom with 1 line padding
	topPad := 2
	gapAfterTabs := 2
	bottomPad := 1

	tabsHeight := lipgloss.Height(tabsCentered)
	statusHeight := lipgloss.Height(statusCentered)

	// Calculate content area height (fixed regardless of content)
	contentAreaHeight := m.height - topPad - tabsHeight - gapAfterTabs - statusHeight - bottomPad
	if contentAreaHeight < 1 {
		contentAreaHeight = 1
	}

	// Place content within its fixed-height area, aligned to top
	contentArea := lipgloss.Place(
		m.width, contentAreaHeight,
		lipgloss.Center, lipgloss.Top,
		contentCentered,
	)

	var sections []string
	sections = append(sections, strings.Repeat("\n", topPad))
	sections = append(sections, tabsCentered)
	sections = append(sections, strings.Repeat("\n", gapAfterTabs))
	sections = append(sections, contentArea)
	sections = append(sections, statusCentered)

	full := strings.Join(sections, "\n")

	out := lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, full)

	// Overlay "copied" notification
	if m.clipboardURL != "" {
		notif := m.renderClipboardNotification()
		notifCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, notif)
		// Replace line near the bottom (above status bar) with the notification
		lines := strings.Split(out, "\n")
		notifLine := m.height - 4
		if notifLine > 0 && notifLine < len(lines) {
			lines[notifLine] = notifCentered
		}
		out = strings.Join(lines, "\n")
		// Prepend OSC 52 so the terminal copies the URL to clipboard
		out = osc52Copy(m.clipboardURL) + out
	}

	return out
}

func (m model) viewTabs() string {
	var rendered []string
	for i, name := range tabNames {
		shortcut := m.styles.statusBarKey.Render(fmt.Sprintf("%d", i+1))
		if tab(i) == m.activeTab {
			rendered = append(rendered, m.styles.activeTab.Render(shortcut+" "+name))
		} else {
			rendered = append(rendered, m.styles.inactiveTab.Render(shortcut+" "+name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Bottom, rendered...)
}

func (m model) viewContent() string {
	s := m.styles

	var content string
	switch m.activeTab {
	case tabAbout:
		content = m.viewAbout()
	case tabExperience:
		content = m.viewExperience()
	case tabStack:
		content = m.viewStack()
	case tabProjects:
		content = m.viewProjects()
	case tabContact:
		content = m.viewContact()
	}

	return s.content.Render(content)
}

func (m model) viewAbout() string {
	s := m.styles
	var sections []string

	// Name and title
	name := s.heading.Render("Arkadiusz Juszczyk")
	title := m.renderer.NewStyle().Foreground(colorMuted).Render("  tech lead & product developer")
	sections = append(sections, name+title)
	sections = append(sections, "")

	// Tagline
	sections = append(sections, s.sectionTitle.Render(`the engineer who asks "why" before "how"`))
	sections = append(sections, "")

	// Short bio
	shortBio := `I'm Arkadiusz, a tech lead and product developer based in Poland, working remotely with teams across the US and Europe. 10 years building software — and the bridges between the people who make it and the people who need it.`
	sections = append(sections, s.paragraph.Render(shortBio))
	sections = append(sections, "")

	// Stats
	stat := func(val, label string) string {
		return s.statValue.Render(val) + " " + s.statLabel.Render(label)
	}
	stats := stat("10+", "years") + "    " + stat("6+", "products shipped") + "    " + stat("3", "teams led")
	sections = append(sections, stats)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) viewExperience() string {
	s := m.styles
	var rows []string

	contentWidth := min(m.width-8, 72)
	periodWidth := 16
	leftWidth := contentWidth - periodWidth

	leftStyle := m.renderer.NewStyle().Width(leftWidth)
	periodStyle := m.renderer.NewStyle().
		Foreground(colorDimmed).
		Width(periodWidth).
		Align(lipgloss.Right)

	for i, job := range experience {
		cursor := "  "
		if i == m.cursor {
			cursor = s.contactCursor.Render("→ ")
		}

		title := s.cardTitle.Render(job.Role)
		company := s.cardSubtitle.Render(" at " + job.Company)

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftStyle.Render(cursor+title+company),
			periodStyle.Render(job.Period),
		)
		rows = append(rows, row)

		if i == m.expanded {
			desc := s.cardDescription.MarginLeft(4).Render(job.Description)
			rows = append(rows, desc)
			rows = append(rows, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) viewStack() string {
	s := m.styles
	_ = s
	var rows []string

	contentWidth := min(m.width-8, 72)
	labelWidth := 12

	labelStyle := m.renderer.NewStyle().
		Foreground(colorDimmed).
		Width(labelWidth)

	tagsStyle := m.renderer.NewStyle().
		Width(contentWidth - labelWidth)

	separator := m.renderer.NewStyle().Foreground(colorBorder).Render(" · ")

	renderGroup := func(label string, items []string, tagStyle lipgloss.Style) string {
		var colored []string
		for _, item := range items {
			colored = append(colored, tagStyle.Render(item))
		}
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render(label),
			tagsStyle.Render(strings.Join(colored, separator)),
		)
	}

	rows = append(rows, renderGroup("frontend", stackFrontend, m.styles.tagFrontend))
	rows = append(rows, "")
	rows = append(rows, renderGroup("backend", stackBackend, m.styles.tagBackend))
	rows = append(rows, "")
	rows = append(rows, renderGroup("tools", stackTools, m.styles.tagTools))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) viewProjects() string {
	s := m.styles
	var rows []string

	contentWidth := min(m.width-8, 72)
	yearWidth := 8
	leftWidth := contentWidth - yearWidth

	leftStyle := m.renderer.NewStyle().Width(leftWidth)
	yearStyle := m.renderer.NewStyle().
		Foreground(colorDimmed).
		Width(yearWidth).
		Align(lipgloss.Right)

	for i, project := range projects {
		cursor := "  "
		if i == m.cursor {
			cursor = s.contactCursor.Render("→ ")
		}

		title := s.cardTitle.Render(project.Title)
		status := s.projectStatus.Render("  " + project.Status)

		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftStyle.Render(cursor+title+status),
			yearStyle.Render(project.Year),
		)
		rows = append(rows, row)

		if i == m.expanded {
			desc := s.cardDescription.MarginLeft(4).Render(project.Description)
			var tags []string
			for _, t := range project.Tags {
				tags = append(tags, s.tagProject.Render("[ "+t+" ]"))
			}
			tagLine := m.renderer.NewStyle().MarginLeft(4).Render(lipgloss.JoinHorizontal(lipgloss.Bottom, tags...))
			rows = append(rows, desc)
			rows = append(rows, tagLine)
			rows = append(rows, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) viewContact() string {
	s := m.styles
	var rows []string

	rows = append(rows, s.heading.Render("let's connect"))
	rows = append(rows, "")
	rows = append(rows, s.paragraph.Render("open to ambitious projects and interesting conversations."))
	rows = append(rows, "")

	for i, c := range contacts {
		cursor := "  "
		if i == m.cursor {
			cursor = s.contactCursor.Render("→ ")
		}
		// OSC 8 hyperlink — renders as clickable in modern terminals
		link := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", c.URL, c.Value)
		row := cursor + s.contactLabel.Render(c.Label) + s.contactValue.Render(link)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) viewStatusBar() string {
	s := m.styles

	// Divider line
	contentWidth := min(m.width-8, 72)
	divider := s.statusBarLine.Render(strings.Repeat("─", contentWidth))

	// Hints
	var hints []string
	key := s.statusBarKey
	val := s.statusBarValue

	hints = append(hints, key.Render("←/→")+val.Render(" tabs"))
	switch m.activeTab {
	case tabExperience, tabProjects:
		hints = append(hints, key.Render("↑/↓")+val.Render(" navigate"))
		hints = append(hints, key.Render("enter")+val.Render(" expand"))
	case tabContact:
		hints = append(hints, key.Render("↑/↓")+val.Render(" navigate"))
		hints = append(hints, key.Render("enter")+val.Render(" copy link"))
	}
	hints = append(hints, key.Render("q")+val.Render(" quit"))

	hintLine := s.statusBar.Render(strings.Join(hints, "    "))

	return lipgloss.JoinVertical(lipgloss.Center, divider, hintLine)
}

func (m model) renderClipboardNotification() string {
	label := m.renderer.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Render("✓ copied")
	url := m.renderer.NewStyle().
		Foreground(colorMuted).
		Render(m.clipboardURL)
	return label + "  " + url
}

func (m model) viewEasterEgg() string {
	coffee := `
      ( (
       ) )
     ........
     |      |]
     \      /
      '----'

  built with go, bubble tea,
  and an unhealthy amount of coffee.

  press any key to go back.
`
	style := m.renderer.NewStyle().
		Foreground(colorPrimary).
		Align(lipgloss.Center).
		Width(m.width)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(coffee),
			)
}

func main() {
	host := flag.String("host", "::", "host to bind to (:: for dual-stack, 0.0.0.0 for IPv4 only)")
	port := flag.Int("port", 2222, "port to listen on")
	hostKey := flag.String("host-key", ".ssh/id_ed25519", "path to SSH host key (generated on first run)")
	flag.Parse()

	// Format address correctly for IPv6 (needs brackets)
	addr := fmt.Sprintf("[%s]:%d", *host, *port)
	if !strings.Contains(*host, ":") {
		addr = fmt.Sprintf("%s:%d", *host, *port)
	}

	srv, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(*hostKey),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				// Per-session renderer that respects the client's terminal capabilities
				renderer := bubbletea.MakeRenderer(s)
				return initialModel(renderer), []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		log.Fatalf("Could not create server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}
