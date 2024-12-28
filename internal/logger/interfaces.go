// Package logger provides logging functionality
package logger

import (
	"github.com/jonesrussell/dashboard/internal/logger/types"
)

// Field represents a key-value pair for structured logging
type Field = types.Field

// Logger represents the interface for logging operations
type Logger = types.Logger

// NewField creates a new log field
func NewField(key string, value interface{}) Field {
	return types.NewField(key, value)
}
