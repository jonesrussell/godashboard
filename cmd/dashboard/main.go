package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/dashboard/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewDashboard(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
