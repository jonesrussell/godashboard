package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// BaseWidget provides common widget functionality
type BaseWidget struct {
	width   int
	height  int
	focused bool
}

// Focus implements Focusable
func (w *BaseWidget) Focus() {
	w.focused = true
}

// Blur implements Focusable
func (w *BaseWidget) Blur() {
	w.focused = false
}

// IsFocused returns whether the widget is focused
func (w *BaseWidget) IsFocused() bool {
	return w.focused
}

// SetSize implements Widget
func (w *BaseWidget) SetSize(width, height int) {
	w.width = width
	w.height = height
}

// GetStyle returns the appropriate style based on focus state
func (w *BaseWidget) GetStyle() lipgloss.Style {
	if w.focused {
		return styles.Focused
	}
	return styles.Base
}

// GetDimensions returns the widget dimensions
func (w *BaseWidget) GetDimensions() (width, height int) {
	return w.width, w.height
}

// DefaultInit provides a default Init implementation
func (w *BaseWidget) DefaultInit() tea.Cmd {
	return nil
}

// DefaultUpdate provides a default Update implementation
func (w *BaseWidget) DefaultUpdate(msg tea.Msg) (Widget, tea.Cmd) {
	return w, nil
}
