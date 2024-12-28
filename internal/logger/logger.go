package logger

import "github.com/jonesrussell/dashboard/internal/logger/types"

// New creates a new logger instance with the given configuration
func New(cfg types.Config) (types.Logger, error) {
	return NewZapLogger(cfg)
}
