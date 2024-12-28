//go:build wireinject
// +build wireinject

package main

import (
	"os"

	"github.com/google/wire"
	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/ui"
)

// InitializeDashboard creates a new dashboard instance with all dependencies wired up
func InitializeDashboard() (*ui.Dashboard, error) {
	wire.Build(
		provideLogger,
		ui.NewDashboard,
	)
	return nil, nil
}

func provideLogger() (logger.Logger, error) {
	cfg := logger.DefaultConfig()

	// Override with environment variables if present
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Level = level
	}
	if path := os.Getenv("LOG_PATH"); path != "" {
		cfg.OutputPath = path
	}
	if debug := os.Getenv("DEBUG"); debug == "true" {
		cfg.Debug = true
	}

	return logger.NewZapLogger(cfg)
}
