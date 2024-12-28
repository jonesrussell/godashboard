// Package ui implements the terminal user interface for the dashboard
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
}

// NewDashboard creates a new dashboard instance
func NewDashboard() *Dashboard {
	return &Dashboard{
		keys:     DefaultKeyMap,
		help:     help.New(),
		showHelp: false,
	}
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
			return d, nil
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
	}

	return d, nil
}

// View implements tea.Model
func (d *Dashboard) View() string {
	// Create the main container
	doc := strings.Builder{}

	// Add header
	doc.WriteString(d.renderHeader())
	doc.WriteString("\n")

	// Add main content area
	doc.WriteString(d.renderContent())
	doc.WriteString("\n")

	// Add footer with help
	if d.showHelp {
		doc.WriteString(d.help.View(d.keys))
	} else {
		doc.WriteString(d.renderFooter())
	}

	return doc.String()
}

func (d *Dashboard) renderHeader() string {
	return styles.HeaderStyle.Render("Dashboard")
}

func (d *Dashboard) renderContent() string {
	return styles.ContentStyle.
		Width(d.width - 4).
		Height(d.height - 6).
		Render("Welcome to the dashboard!")
}

func (d *Dashboard) renderFooter() string {
	return styles.FooterStyle.Render("Press ? for help")
}
