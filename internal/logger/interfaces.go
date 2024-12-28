// Package logger provides structured logging capabilities for the dashboard
package logger

import "context"

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	WithFields(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

// Field represents a logging field key-value pair
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new logging field
func NewField(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
