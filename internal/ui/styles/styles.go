// Package styles provides UI styling constants and utilities for the dashboard
package styles

import "github.com/charmbracelet/lipgloss"

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

// WithSize returns a copy of the style with the given dimensions
func WithSize(style lipgloss.Style, width, height int) lipgloss.Style {
	return style.Copy().
		Width(width).
		Height(height)
}
