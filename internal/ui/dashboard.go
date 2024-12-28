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
	"github.com/jonesrussell/dashboard/internal/ui/container"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
	"github.com/jonesrussell/dashboard/internal/ui/widgets/sysinfo"
	"github.com/jonesrussell/dashboard/internal/ui/widgets/tasks"
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

// DefaultKeyMap defines the default keyboard shortcuts for the dashboard
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

	// Cached styles
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
	styleCache  *styles.StyleCache

	// Content cache
	headerContent string
	footerContent string
	contentCache  map[string]string
	needsRefresh  bool

	// Widget container
	container *container.Container
}

// EnableDebug enables debug output
func (d *Dashboard) EnableDebug() {
	d.debug = true
	d.container.EnableDebug()
	if widget := d.container.FocusedWidget(); widget != nil {
		if debuggable, ok := widget.(interface{ EnableDebug() }); ok {
			debuggable.EnableDebug()
		}
	}
}

// NewDashboard creates a new dashboard instance
func NewDashboard(log logger.Logger) *Dashboard {
	d := &Dashboard{
		keys:     DefaultKeyMap,
		help:     help.New(),
		showHelp: false,
		debug:    false,
		logger:   log,

		// Initialize cached styles
		headerStyle: styles.HeaderStyle,
		footerStyle: styles.FooterStyle,
		styleCache:  styles.NewStyleCache(),

		// Initialize content cache
		contentCache: make(map[string]string),
		needsRefresh: true,

		// Initialize widget container with 2x2 grid
		container: container.New(2, 2, log),
	}

	// Pre-render static content
	d.headerContent = d.headerStyle.Render("Dashboard")
	d.footerContent = d.footerStyle.Render("Press ? for help")

	// Add system info widget
	sysInfo := sysinfo.New()
	d.container.AddWidgetWithConfig(sysInfo, container.GridConfig{
		Row:      0,
		Col:      0,
		RowSpan:  1,
		ColSpan:  1,
		MinWidth: 40,
	})

	// Add tasks widget
	tasks := tasks.New()
	d.container.AddWidgetWithConfig(tasks, container.GridConfig{
		Row:      0,
		Col:      1,
		RowSpan:  2,
		ColSpan:  1,
		MinWidth: 40,
	})

	return d
}

// Init implements tea.Model
func (d *Dashboard) Init() tea.Cmd {
	return d.container.Init()
}

// Update implements tea.Model
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.Quit):
			return d, tea.Quit
		case key.Matches(msg, d.keys.Help):
			d.showHelp = !d.showHelp
			d.needsRefresh = true
			return d, nil
		case msg.String() == "d":
			d.debug = !d.debug
			d.needsRefresh = true
			return d, nil
		case key.Matches(msg, d.keys.Tab):
			// Let container handle tab navigation
			return d.container.Update(msg)
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.needsRefresh = true
		return d, nil
	}

	// Forward other messages to container
	return d.container.Update(msg)
}

// View implements tea.Model
func (d *Dashboard) View() string {
	var b strings.Builder

	// Header
	header := "Dashboard"
	if d.debug {
		header += " Debug: ON"
	}
	b.WriteString(d.headerStyle.Render(header))
	b.WriteRune('\n')

	// Main content area with padding
	contentWidth := d.width - 4
	contentHeight := d.height - 6

	if d.showHelp {
		helpContent := "Help\n\n" + d.help.View(d.keys)
		helpStyle := d.styleCache.GetContentStyle(contentWidth, contentHeight)
		b.WriteString(helpStyle.Render(helpContent))
	} else {
		containerStyle := d.styleCache.GetContentStyle(contentWidth, contentHeight)
		b.WriteString(containerStyle.Render(d.container.View()))
	}
	b.WriteRune('\n')

	// Footer
	b.WriteString(d.footerStyle.Render("Press ? for help"))

	return b.String()
}

// AddWidget adds a widget to the container
func (d *Dashboard) AddWidget(w components.Widget) {
	d.container.AddWidget(w)
	d.needsRefresh = true
}
