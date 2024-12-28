// Package styles provides UI styling constants and utilities for the dashboard
package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Colors defines the color palette for the dashboard
var (
	// Primary is the main accent color
	Primary = lipgloss.Color("#2196F3")
	// Secondary is used for highlights and accents
	Secondary = lipgloss.Color("#FFB74D")
	// Subtle is used for less prominent elements
	Subtle = lipgloss.Color("#4A4A4A")
)

// Pre-defined styles for common use cases
var (
	// Base style for all content
	Base = lipgloss.NewStyle().
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Primary)

	// Focused style for active elements
	Focused = Base.Copy().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(Primary)

	// Header style for section headers
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		MarginLeft(2)

	// Footer style for bottom text
	Footer = lipgloss.NewStyle().
		Foreground(Subtle).
		MarginLeft(2)

	// Title style for widget titles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary)

	// Selected style for highlighted items
	Selected = lipgloss.NewStyle().
			Background(Primary).
			Foreground(lipgloss.Color("#ffffff"))
)

// StyleCache provides caching for computed styles
type StyleCache struct {
	styles map[string]lipgloss.Style
}

// NewStyleCache creates a new style cache
func NewStyleCache() *StyleCache {
	return &StyleCache{
		styles: make(map[string]lipgloss.Style),
	}
}

// Get retrieves a cached style by key
func (c *StyleCache) Get(key string) (lipgloss.Style, bool) {
	style, ok := c.styles[key]
	return style, ok
}

// Set stores a style in the cache
func (c *StyleCache) Set(key string, style lipgloss.Style) {
	c.styles[key] = style
}

// GetFocusedStyle returns the focused style for a widget
func (c *StyleCache) GetFocusedStyle(width, height int) lipgloss.Style {
	key := fmt.Sprintf("focused_%d_%d", width, height)
	if style, ok := c.Get(key); ok {
		return style
	}
	style := WithSize(Focused, width, height)
	c.Set(key, style)
	return style
}

// GetContentStyle returns the base style for a widget
func (c *StyleCache) GetContentStyle(width, height int) lipgloss.Style {
	key := fmt.Sprintf("content_%d_%d", width, height)
	if style, ok := c.Get(key); ok {
		return style
	}
	style := WithSize(Base, width, height)
	c.Set(key, style)
	return style
}

// WithSize returns a style with the specified dimensions
func WithSize(style lipgloss.Style, width, height int) lipgloss.Style {
	return style.Width(width).Height(height)
}
