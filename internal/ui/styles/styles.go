package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("#7E57C2")
	Secondary = lipgloss.Color("#FFB74D")
	Error     = lipgloss.Color("#EF5350")
	Success   = lipgloss.Color("#66BB6A")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Primary)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Padding(0, 1)

	// Content styles
	ContentStyle = BaseStyle.Copy().
			Padding(1, 2)

	// Footer styles
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
