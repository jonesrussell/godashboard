// Package types provides common types for logging
package types

import "context"

// Field represents a log field
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new log field
func NewField(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
	Close() error
}
