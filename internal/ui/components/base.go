package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// BaseWidget provides common widget functionality
type BaseWidget struct {
	Width   int
	Height  int
	Focused bool
}

// Init implements Widget interface
func (w *BaseWidget) Init() tea.Cmd {
	return nil
}

// Update implements Widget interface
func (w *BaseWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	return w, nil
}

// Focus implements Focusable
func (w *BaseWidget) Focus() {
	w.Focused = true
}

// Blur implements Focusable
func (w *BaseWidget) Blur() {
	w.Focused = false
}

// IsFocused returns whether the widget is focused
func (w *BaseWidget) IsFocused() bool {
	return w.Focused
}

// SetSize implements Sizable
func (w *BaseWidget) SetSize(width, height int) {
	w.Width = width
	w.Height = height
}

// GetStyle returns the appropriate style based on focus state
func (w *BaseWidget) GetStyle() lipgloss.Style {
	if w.Focused {
		return styles.Focused
	}
	return styles.Base
}

// GetDimensions implements Sizable
func (w *BaseWidget) GetDimensions() (width, height int) {
	return w.Width, w.Height
}

// BaseInit provides a base initialization implementation for widgets
func (w *BaseWidget) BaseInit() tea.Cmd {
	return nil
}

// BaseUpdate provides a base update implementation for widgets
func (w *BaseWidget) BaseUpdate(msg tea.Msg) (Widget, tea.Cmd) {
	return w, nil
}

// View implements Widget interface with a default empty view
func (w *BaseWidget) View() string {
	return ""
}
