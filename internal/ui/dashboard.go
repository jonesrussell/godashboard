// Package ui implements the terminal user interface for the dashboard
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/container"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
	"github.com/jonesrussell/dashboard/internal/ui/widgets/sysinfo"
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
		key.WithHelp("q", "quit"),
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

// NewDashboard creates a new dashboard instance
func NewDashboard() *Dashboard {
	d := &Dashboard{
		keys:     DefaultKeyMap,
		help:     help.New(),
		showHelp: false,

		// Initialize cached styles
		headerStyle: styles.HeaderStyle,
		footerStyle: styles.FooterStyle,
		styleCache:  styles.NewStyleCache(),

		// Initialize content cache
		contentCache: make(map[string]string),
		needsRefresh: true,

		// Initialize widget container with 2x2 grid
		container: container.New(2, 2),
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
		ColSpan:  2,
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
		case key.Matches(msg, d.keys.Tab):
			// Let container handle tab navigation
			if d.container.HandleFocusKey(msg) {
				d.needsRefresh = true
				return d, nil
			}
		}

		// Forward other key messages to container
		if !d.showHelp {
			if _, cmd := d.container.Update(msg); cmd != nil {
				return d, cmd
			}
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.needsRefresh = true

		// Update container size
		contentHeight := d.height - 6 // Account for header and footer
		if _, cmd := d.container.Update(tea.WindowSizeMsg{
			Width:  d.width - 4,
			Height: contentHeight,
		}); cmd != nil {
			return d, cmd
		}
	}

	return d, nil
}

// View implements tea.Model
func (d *Dashboard) View() string {
	var b strings.Builder
	b.Grow(d.width * d.height)

	// Log dimensions
	fmt.Printf("Dashboard dimensions: width=%d, height=%d\n", d.width, d.height)

	// Add header
	b.WriteString(d.headerContent)
	b.WriteByte('\n')

	// Add main content area
	if d.showHelp {
		contentStyle := d.styleCache.GetContentStyle(d.width-4, d.height-6)
		b.WriteString(contentStyle.Render("Welcome to the dashboard!"))
	} else {
		fmt.Printf("Container dimensions: width=%d, height=%d\n", d.width-4, d.height-6)
		b.WriteString(d.container.View())
	}
	b.WriteByte('\n')

	// Add footer with help
	if d.showHelp {
		b.WriteString(d.help.View(d.keys))
	} else {
		b.WriteString(d.footerContent)
	}

	return b.String()
}

// AddWidget adds a widget to the container
func (d *Dashboard) AddWidget(w components.Widget) {
	d.container.AddWidget(w)
	d.needsRefresh = true
}
