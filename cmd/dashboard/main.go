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

const (
	logMaxSizeMB    = 10 // 10MB
	logMaxBackups   = 3  // Keep 3 backups
	logMaxAgeDays   = 7  // 7 days
	logDefaultLevel = "info"
)

func main() {
	// Initialize logger
	log, err := logger.New(types.Config{
		Level:      logDefaultLevel,
		OutputPath: "dashboard.log",
		MaxSize:    logMaxSizeMB,
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAgeDays,
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
		log.Error("Error running dashboard", logger.NewField("error", err))
		os.Exit(1)
	}

	log.Info("Dashboard shutdown complete")
}
