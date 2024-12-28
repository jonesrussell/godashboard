package logger

import "github.com/google/wire"

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		OutputPath: "logs/app.log",
		MaxSize:    10, // 10MB
		MaxBackups: 5,  // 5 backups
		MaxAge:     30, // 30 days
		Compress:   true,
		Debug:      false,
	}
}

// ProvideLogger creates a new logger instance
func ProvideLogger(cfg Config) (Logger, error) {
	return NewZapLogger(cfg)
}

// ProviderSet is the wire provider set for logger
var ProviderSet = wire.NewSet(
	DefaultConfig,
	ProvideLogger,
)
