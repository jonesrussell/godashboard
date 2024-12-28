// Package testlogger provides a centralized logger setup for tests
package testlogger

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger"
)

// TestLogger extends logger.Logger with test-specific functionality
type TestLogger struct {
	logger.Logger
	outputPath string
}

// GetOutputPath returns the path to the log file
func (t *TestLogger) GetOutputPath() string {
	return t.outputPath
}

// Setup creates a new test logger with the given name
func Setup(t *testing.T, name string) *TestLogger {
	outputPath := filepath.Join(t.TempDir(), name+".log")
	cfg := logger.Config{
		Level:      "debug",
		OutputPath: outputPath,
		MaxSize:    1,
		MaxBackups: 0,
		MaxAge:     1,
		Compress:   false,
		Debug:      true,
	}
	log, err := logger.NewZapLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}
	return &TestLogger{
		Logger:     log,
		outputPath: outputPath,
	}
}

// SetupBenchmark creates a new test logger optimized for benchmarks
func SetupBenchmark() *TestLogger {
	cfg := logger.Config{
		Level:      "error", // Minimize logging during benchmarks
		OutputPath: "benchmark.log",
		MaxSize:    1,
		MaxBackups: 0,
		MaxAge:     1,
		Compress:   false,
		Debug:      false,
	}
	log, err := logger.NewZapLogger(cfg)
	if err != nil {
		panic(err) // Acceptable in benchmarks
	}
	return &TestLogger{
		Logger:     log,
		outputPath: "benchmark.log",
	}
}
