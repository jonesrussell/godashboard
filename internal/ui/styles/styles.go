// Package styles provides UI styling constants and utilities for the dashboard
package styles

import "github.com/charmbracelet/lipgloss"

var (
	// HeaderStyle is the style for the dashboard header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("69")).
			MarginLeft(2)

	// FooterStyle is the style for the dashboard footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)
)
