// Package styles provides UI styling constants and utilities for the dashboard
package styles

import "github.com/charmbracelet/lipgloss"

const (
	// DefaultHeaderWidth is the default width for header elements
	DefaultHeaderWidth = 30
)

// Colors
var (
	Primary   = lipgloss.Color("#2196F3")
	Secondary = lipgloss.Color("#FFB74D")
	Border    = lipgloss.Color("#3C3C3C")
	Text      = lipgloss.Color("#FFFFFF")
	Subtle    = lipgloss.Color("#4A4A4A")
	Highlight = lipgloss.Color("#82AAFF")
)

// BaseStyle defines the base styling for all UI elements
var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(Border).
	Foreground(Text).
	Padding(1)

// HeaderStyle defines the styling for header elements
var HeaderStyle = BaseStyle.Copy().
	Bold(true).
	Foreground(Primary).
	BorderStyle(lipgloss.DoubleBorder()).
	BorderForeground(Primary).
	Align(lipgloss.Center).
	Width(DefaultHeaderWidth)

// ContentStyle defines the styling for content areas
var ContentStyle = BaseStyle.Copy().
	BorderForeground(Border).
	BorderStyle(lipgloss.RoundedBorder())

// FooterStyle defines the styling for footer elements
var FooterStyle = BaseStyle.Copy().
	BorderStyle(lipgloss.HiddenBorder()).
	Foreground(Subtle).
	Align(lipgloss.Center)

// FocusedStyle defines the styling for focused elements
var FocusedStyle = ContentStyle.Copy().
	BorderForeground(Primary).
	BorderStyle(lipgloss.DoubleBorder())

// StyleCache provides cached styles for different dimensions
type StyleCache struct {
	contentStyles map[string]lipgloss.Style
	focusedStyles map[string]lipgloss.Style
}

// NewStyleCache creates a new style cache
func NewStyleCache() *StyleCache {
	return &StyleCache{
		contentStyles: make(map[string]lipgloss.Style),
		focusedStyles: make(map[string]lipgloss.Style),
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

// GetFocusedStyle returns a cached focused style for the given dimensions
func (c *StyleCache) GetFocusedStyle(width, height int) lipgloss.Style {
	key := styleKey(width, height)
	if style, ok := c.focusedStyles[key]; ok {
		return style
	}

	style := FocusedStyle.Copy().
		Width(width).
		Height(height)
	c.focusedStyles[key] = style
	return style
}

// styleKey generates a cache key for dimensions
func styleKey(width, height int) string {
	return string(rune(width)) + "x" + string(rune(height))
}
