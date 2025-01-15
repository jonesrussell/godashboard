// Package main is the entry point for the dashboard application
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/ui"
)

func main() {
	// Initialize logger first
	log, err := logger.New(logger.DefaultConfig())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer log.Close()

	// Test logger
	log.Debug("Logger initialized",
		logger.NewField("config", logger.DefaultConfig()),
	)
	log.Info("Starting dashboard application")

	// Initialize dashboard
	dash := ui.NewDashboard(log)
	log.Debug("Dashboard initialized")

	// Create and start program
	p := tea.NewProgram(dash)
	log.Debug("Starting Bubbletea program")

	if _, err := p.Run(); err != nil {
		log.Error("Error running program", logger.NewField("error", err))
		os.Exit(1)
	}
	log.Info("Program exiting normally")
}
