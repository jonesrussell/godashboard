package logger

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonesrussell/dashboard/internal/testutil/assertions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Level:      "info",
				OutputPath: "test.log",
				MaxSize:    1,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				Debug:      false,
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			cfg: Config{
				Level:      "info",
				OutputPath: filepath.Join("non-existent-dir", strings.Repeat("a", 1000), "test.log"),
				MaxSize:    1,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				Debug:      false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up test files
			defer os.Remove(tt.cfg.OutputPath)

			logger, err := NewZapLogger(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create an encoder
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// Create a core that writes to our buffer
	core := zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.DebugLevel)

	// Create a logger
	logger := &zapLogger{
		logger: zap.New(core),
	}

	tests := []struct {
		name       string
		logFunc    func(string, ...Field)
		level      string
		msg        string
		fields     []Field
		wantFields map[string]interface{}
	}{
		{
			name:    "debug level",
			logFunc: logger.Debug,
			level:   "debug",
			msg:     "debug message",
			fields:  []Field{NewField("key", "value")},
			wantFields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "info level",
			logFunc: logger.Info,
			level:   "info",
			msg:     "info message",
			fields:  []Field{NewField("number", 42)},
			wantFields: map[string]interface{}{
				"number": float64(42),
			},
		},
		{
			name:    "warn level",
			logFunc: logger.Warn,
			level:   "warn",
			msg:     "warn message",
			fields:  []Field{NewField("bool", true)},
			wantFields: map[string]interface{}{
				"bool": true,
			},
		},
		{
			name:    "error level",
			logFunc: logger.Error,
			level:   "error",
			msg:     "error message",
			fields:  []Field{NewField("error", "failed")},
			wantFields: map[string]interface{}{
				"error": "failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.msg, tt.fields...)
			assertions.AssertLogEntry(t, &buf, tt.level, tt.msg, tt.wantFields)
		})
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.InfoLevel)
	logger := &zapLogger{logger: zap.New(core)}

	// Create logger with fields
	fields := []Field{
		NewField("service", "test"),
		NewField("version", "1.0"),
	}
	loggerWithFields := logger.WithFields(fields...)

	// Log a message
	loggerWithFields.Info("test message")

	// Assert log entry
	wantFields := map[string]interface{}{
		"service": "test",
		"version": "1.0",
	}
	assertions.AssertLogEntry(t, &buf, "info", "test message", wantFields)
}

func TestWithContext(t *testing.T) {
	var buf bytes.Buffer
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.InfoLevel)
	logger := &zapLogger{logger: zap.New(core)}

	// Create context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "123")
	loggerWithCtx := logger.WithContext(ctx)

	// Log a message
	loggerWithCtx.Info("test message")

	// Assert log entry
	wantFields := map[string]interface{}{
		"request_id": "123",
	}
	assertions.AssertLogEntry(t, &buf, "info", "test message", wantFields)
}

func TestLogRotation(t *testing.T) {
	// Create test directory
	testDir := t.TempDir()
	logPath := filepath.Join(testDir, "test.log")

	// Create a scope for the logger to ensure it's cleaned up
	func() {
		cfg := Config{
			Level:      "info",
			OutputPath: logPath,
			MaxSize:    1, // 1MB
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   true,
		}

		// Create logger
		logger, err := NewZapLogger(cfg)
		require.NoError(t, err)

		// Write enough logs to trigger rotation
		for i := 0; i < 100000; i++ {
			logger.Info("test message", NewField("count", i))
		}

		// Sync and close the logger
		if zl, ok := logger.(*zapLogger); ok {
			require.NoError(t, zl.logger.Sync())
		}
	}()

	// Give the OS a moment to release file handles
	time.Sleep(100 * time.Millisecond)

	// Check if log file exists
	_, err := os.Stat(logPath)
	assert.NoError(t, err)

	// Check if backup file exists (may be compressed)
	backupExists := false
	files, err := os.ReadDir(testDir)
	require.NoError(t, err)
	for _, file := range files {
		if file.Name() != "test.log" && filepath.Ext(file.Name()) != ".gz" {
			backupExists = true
			break
		}
	}
	assert.True(t, backupExists, "backup file should exist")
}
