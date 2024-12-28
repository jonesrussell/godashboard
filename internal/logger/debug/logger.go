// Package debug provides debug logging functionality for the dashboard
package debug

import (
	"fmt"
	"io"
)

// Logger provides debug logging functionality
type Logger struct {
	out io.Writer
}

// New creates a new debug logger
func New(out io.Writer) *Logger {
	if out == nil {
		out = io.Discard
	}
	return &Logger{out: out}
}

// Printf formats and writes a debug message
func (l *Logger) Printf(format string, args ...interface{}) {
	if l.out != nil {
		// Use fmt.Fprintf internally, but expose a cleaner API
		// This is okay because it's isolated in the debug package
		fmt.Fprintf(l.out, format+"\n", args...)
	}
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(out io.Writer) {
	if out == nil {
		out = io.Discard
	}
	l.out = out
}
