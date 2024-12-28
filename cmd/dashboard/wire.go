//go:build wireinject
// +build wireinject

package main

import (
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
	cfg := logger.Config{
		Level:      "debug",
		OutputPath: "logs/dashboard.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
		Debug:      true,
	}
	return logger.NewZapLogger(cfg)
}
