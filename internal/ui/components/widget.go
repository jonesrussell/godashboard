// Package components provides reusable UI components for the dashboard
package components

import tea "github.com/charmbracelet/bubbletea"

// Widget represents a dashboard widget
type Widget interface {
	// Core Bubbletea interface
	Init() tea.Cmd
	Update(msg tea.Msg) (Widget, tea.Cmd)
	View() string

	// Size management
	SetSize(width, height int)
	GetDimensions() (width, height int)

	// Focus management
	Focus()
	Blur()
	IsFocused() bool
}

// Focusable represents a component that can receive focus
type Focusable interface {
	Focus()
	Blur()
	IsFocused() bool
}

// Sizable represents a component that can be resized
type Sizable interface {
	SetSize(width, height int)
	GetDimensions() (width, height int)
}
