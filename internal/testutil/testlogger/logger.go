// Package testlogger provides a simple logger for testing
package testlogger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jonesrussell/dashboard/internal/logger/types"
)

// testLogger implements Logger for testing
type testLogger struct {
	file   *os.File
	mu     sync.Mutex
	fields []types.Field
}

// Config holds logger configuration
type Config struct {
	Level      string
	OutputPath string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Debug      bool
}

// DefaultConfig returns default test logger config
func DefaultConfig(name string, tb testing.TB) Config {
	var logPath string
	if b, ok := tb.(*testing.B); ok {
		logPath = filepath.Join(os.TempDir(), "bench-"+b.Name()+"-"+name+".log")
	} else {
		logPath = filepath.Join(tb.TempDir(), name+".log")
	}

	return Config{
		Level:      "debug",
		OutputPath: logPath,
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
		Debug:      true,
	}
}

// NewTestLogger creates a new logger for testing
func NewTestLogger(tb testing.TB, name string) (types.Logger, string) {
	cfg := DefaultConfig(name, tb)

	// Create log directory if needed
	if err := os.MkdirAll(filepath.Dir(cfg.OutputPath), 0o755); err != nil {
		tb.Fatal(err)
	}

	// Open log file
	file, err := os.OpenFile(cfg.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		tb.Fatal(err)
	}

	logger := &testLogger{
		file: file,
	}

	// Register cleanup
	tb.Cleanup(func() {
		if err := logger.Close(); err != nil {
			tb.Logf("Failed to close logger: %v", err)
		}
		// Wait a moment for Windows to release the file handle
		time.Sleep(10 * time.Millisecond)
		_ = os.Remove(cfg.OutputPath)
	})

	return logger, cfg.OutputPath
}

func (l *testLogger) log(level, msg string, fields ...types.Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Combine logger fields with message fields
	allFields := make([]types.Field, 0, len(l.fields)+len(fields))
	allFields = append(allFields, l.fields...)
	allFields = append(allFields, fields...)

	// Create log entry
	entry := map[string]interface{}{
		"level":     level,
		"timestamp": time.Now().Format(time.RFC3339),
		"msg":       msg,
	}

	// Add fields
	for _, f := range allFields {
		entry[f.Key] = f.Value
	}

	// Write log entry
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log entry: %v\n", err)
	}
}

func (l *testLogger) Debug(msg string, fields ...types.Field) {
	l.log("DEBUG", msg, fields...)
}

func (l *testLogger) Info(msg string, fields ...types.Field) {
	l.log("INFO", msg, fields...)
}

func (l *testLogger) Warn(msg string, fields ...types.Field) {
	l.log("WARN", msg, fields...)
}

func (l *testLogger) Error(msg string, fields ...types.Field) {
	l.log("ERROR", msg, fields...)
}

func (l *testLogger) Fatal(msg string, fields ...types.Field) {
	l.log("FATAL", msg, fields...)
}

func (l *testLogger) WithFields(fields ...types.Field) types.Logger {
	newLogger := &testLogger{
		file:   l.file,
		fields: make([]types.Field, 0, len(l.fields)+len(fields)),
	}
	newLogger.fields = append(newLogger.fields, l.fields...)
	newLogger.fields = append(newLogger.fields, fields...)
	return newLogger
}

func (l *testLogger) WithContext(ctx context.Context) types.Logger {
	// Extract request ID from context if available
	if reqID, ok := ctx.Value("request_id").(string); ok {
		return l.WithFields(types.NewField("request_id", reqID))
	}
	return l
}

func (l *testLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
