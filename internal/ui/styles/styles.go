// Package styles provides UI styling constants and utilities for the dashboard
package styles

import "github.com/charmbracelet/lipgloss"

const (
	// DefaultHeaderWidth is the default width for header elements
	DefaultHeaderWidth = 30
)

// Primary is the main color used in the UI theme
var Primary = lipgloss.Color("#2196F3")

// Secondary is the complementary color used in the UI theme
var Secondary = lipgloss.Color("#FFB74D")

// BaseStyle defines the base styling for all UI elements
var BaseStyle = lipgloss.NewStyle().
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder())

// HeaderStyle defines the styling for header elements
var HeaderStyle = BaseStyle.Copy().
	Bold(true).
	Foreground(Primary).
	BorderForeground(Primary).
	Width(DefaultHeaderWidth)

// ContentStyle defines the styling for content areas
var ContentStyle = BaseStyle.Copy().
	BorderForeground(Secondary)

// FooterStyle defines the styling for footer elements
var FooterStyle = BaseStyle.Copy().
	BorderStyle(lipgloss.HiddenBorder()).
	Align(lipgloss.Center)

// StyleCache provides cached styles for different dimensions
type StyleCache struct {
	contentStyles map[string]lipgloss.Style
}

// NewStyleCache creates a new style cache
func NewStyleCache() *StyleCache {
	return &StyleCache{
		contentStyles: make(map[string]lipgloss.Style),
	}
}

// GetContentStyle returns a cached content style for the given dimensions
func (c *StyleCache) GetContentStyle(width, height int) lipgloss.Style {
	key := styleKey(width, height)
	if style, ok := c.contentStyles[key]; ok {
		return style
	}

	style := ContentStyle.Copy().
		Width(width).
		Height(height)
	c.contentStyles[key] = style
	return style
}

// styleKey generates a cache key for dimensions
func styleKey(width, height int) string {
	return string(rune(width)) + "x" + string(rune(height))
}
