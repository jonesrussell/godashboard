package logger

import (
	"github.com/google/wire"
	"github.com/jonesrussell/dashboard/internal/logger/types"
)

const (
	// DefaultMaxSize is the default maximum size in megabytes of the log file before it gets rotated
	DefaultMaxSize = 10
	// DefaultMaxBackups is the default maximum number of old log files to retain
	DefaultMaxBackups = 5
	// DefaultMaxAge is the default maximum number of days to retain old log files
	DefaultMaxAge = 30
)

// DefaultConfig returns the default logger configuration
func DefaultConfig() types.Config {
	return types.Config{
		Level:      "info",
		OutputPath: "logs/app.log",
		MaxSize:    DefaultMaxSize,    // 10MB
		MaxBackups: DefaultMaxBackups, // 5 backups
		MaxAge:     DefaultMaxAge,     // 30 days
		Compress:   true,
		Debug:      false,
	}
}

// ProvideLogger creates a new logger instance
func ProvideLogger(cfg types.Config) (types.Logger, error) {
	return NewZapLogger(cfg)
}

// ProviderSet is the wire provider set for logger
var ProviderSet = wire.NewSet(
	DefaultConfig,
	ProvideLogger,
)
