package components

import tea "github.com/charmbracelet/bubbletea"

// Widget represents a dashboard widget
type Widget interface {
	// Init initializes the widget
	Init() tea.Cmd

	// Update handles messages and updates the widget state
	Update(msg tea.Msg) (Widget, tea.Cmd)

	// View renders the widget
	View() string

	// SetSize sets the widget dimensions
	SetSize(width, height int)
}
