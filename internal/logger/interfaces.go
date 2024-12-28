// Package logger provides logging functionality
package logger

import (
	"github.com/jonesrussell/dashboard/internal/logger/types"
)

// Re-export types for backward compatibility
type (
	Field  = types.Field
	Logger = types.Logger
)

// NewField creates a new log field
func NewField(key string, value interface{}) Field {
	return types.NewField(key, value)
}
