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
	// Border is used for borders and dividers
	Border = lipgloss.Color("#3C3C3C")
)

// Base styles for the dashboard
var (
	// ContentStyle is the base style for content areas
	ContentStyle = lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary)

	// HeaderStyle is the style for the dashboard header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginLeft(2)

	// FooterStyle is the style for the dashboard footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			MarginLeft(2)
)
