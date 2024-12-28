// Package ui implements the terminal user interface for the dashboard
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
	"github.com/jonesrussell/dashboard/internal/ui/widgets/sysinfo"
	"github.com/jonesrussell/dashboard/internal/ui/widgets/tasks"
)

// Dashboard messages
type dashboardMsg int

const (
	msgNone dashboardMsg = iota
	msgToggleHelp
	msgToggleDebug
)

// KeyMap defines the keybindings for the dashboard
type KeyMap struct {
	Quit  key.Binding
	Help  key.Binding
	Tab   key.Binding
	Enter key.Binding
}

// ShortHelp implements help.KeyMap
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp implements help.KeyMap
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
		{k.Tab, k.Enter},
	}
}

// DefaultKeyMap defines the default keyboard shortcuts
var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next widget"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

// Dashboard represents the main application model
type Dashboard struct {
	keys     KeyMap
	help     help.Model
	width    int
	height   int
	showHelp bool
	debug    bool
	logger   logger.Logger

	// Widgets
	sysInfo components.Widget
	tasks   components.Widget
}

// NewDashboard creates a new dashboard instance
func NewDashboard(log logger.Logger) *Dashboard {
	return &Dashboard{
		keys:     DefaultKeyMap,
		help:     help.New(),
		showHelp: false,
		debug:    false,
		logger:   log,
		sysInfo:  sysinfo.New(),
		tasks:    tasks.New(),
	}
}

// Init implements tea.Model
func (d *Dashboard) Init() tea.Cmd {
	return tea.Batch(
		d.sysInfo.Init(),
		d.tasks.Init(),
	)
}

// Update implements tea.Model
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.Quit):
			return d, tea.Quit
		case key.Matches(msg, d.keys.Help):
			d.showHelp = !d.showHelp
			return d, nil
		case msg.String() == "d":
			d.debug = !d.debug
			return d, nil
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

		// Calculate widget sizes
		contentWidth := d.width - 4
		contentHeight := d.height - 6
		halfWidth := contentWidth / 2

		// Update widget sizes
		d.sysInfo.SetSize(halfWidth-1, contentHeight)
		d.tasks.SetSize(halfWidth-1, contentHeight)
	}

	// Update widgets
	if sysInfo, cmd := d.sysInfo.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
		d.sysInfo = sysInfo.(components.Widget)
	}
	if tasks, cmd := d.tasks.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
		d.tasks = tasks.(components.Widget)
	}

	return d, tea.Batch(cmds...)
}

// View implements tea.Model
func (d *Dashboard) View() string {
	var b strings.Builder

	// Header
	header := "Dashboard"
	if d.debug {
		header += " Debug: ON"
	}
	b.WriteString(styles.Header.Render(header))
	b.WriteRune('\n')

	// Main content area
	contentWidth := d.width - 4
	contentHeight := d.height - 6

	if d.showHelp {
		helpContent := "Help\n\n" + d.help.View(d.keys)
		b.WriteString(styles.WithSize(styles.Base, contentWidth, contentHeight).Render(helpContent))
	} else {
		// Layout widgets side by side
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			d.sysInfo.View(),
			d.tasks.View(),
		)
		b.WriteString(content)
	}
	b.WriteRune('\n')

	// Footer
	b.WriteString(styles.Footer.Render("Press ? for help"))

	return b.String()
}
