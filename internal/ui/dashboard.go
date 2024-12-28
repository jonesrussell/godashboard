// Package ui implements the terminal user interface for the dashboard
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
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
	}

	// Pre-render static content
	d.headerContent = d.headerStyle.Render("Dashboard")
	d.footerContent = d.footerStyle.Render("Press ? for help")

	return d
}

// Init implements tea.Model
func (d *Dashboard) Init() tea.Cmd {
	return nil
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
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.needsRefresh = true
	}

	return d, nil
}

// View implements tea.Model
func (d *Dashboard) View() string {
	// Return cached content if available and no refresh needed
	if !d.needsRefresh {
		if d.showHelp {
			return d.getHelpView()
		}
		return d.getNormalView()
	}

	// Update cache and return new content
	d.needsRefresh = false
	if d.showHelp {
		return d.getHelpView()
	}
	return d.getNormalView()
}

// getHelpView returns the help view (not cached since help model handles its own rendering)
func (d *Dashboard) getHelpView() string {
	var b strings.Builder
	b.Grow(d.width * d.height)

	b.WriteString(d.headerContent)
	b.WriteByte('\n')

	contentStyle := d.styleCache.GetContentStyle(d.width-4, d.height-6)
	b.WriteString(contentStyle.Render("Welcome to the dashboard!"))
	b.WriteByte('\n')

	b.WriteString(d.help.View(d.keys))
	return b.String()
}

// getNormalView returns the normal view (with caching)
func (d *Dashboard) getNormalView() string {
	key := fmt.Sprintf("%dx%d", d.width, d.height)
	if content, ok := d.contentCache[key]; ok {
		return content
	}

	var b strings.Builder
	b.Grow(d.width * d.height)

	b.WriteString(d.headerContent)
	b.WriteByte('\n')

	contentStyle := d.styleCache.GetContentStyle(d.width-4, d.height-6)
	b.WriteString(contentStyle.Render("Welcome to the dashboard!"))
	b.WriteByte('\n')

	b.WriteString(d.footerContent)

	content := b.String()
	d.contentCache[key] = content
	return content
}
