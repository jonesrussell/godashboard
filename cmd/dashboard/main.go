// Package main is the entry point for the dashboard application
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/logger/types"
	"github.com/jonesrussell/dashboard/internal/ui"
)

func main() {
	// Initialize logger
	log, err := logger.New(types.Config{
		Level:      "info",
		OutputPath: "dashboard.log",
		MaxSize:    10,    // 10MB
		MaxBackups: 3,     // Keep 3 backups
		MaxAge:     7,     // 7 days
		Compress:   true,  // Compress old files
		Debug:      false, // Production mode
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Create and run the dashboard
	dash := ui.NewDashboard(log)
	p := tea.NewProgram(dash,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running dashboard: %v\n", err)
		os.Exit(1)
	}
}
