package logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonesrussell/dashboard/internal/logger/types"
	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	logger, logPath := testlogger.NewTestLogger(t, "log-levels")
	defer logger.Close()

	tests := []struct {
		name       string
		logFunc    func(string, ...types.Field)
		level      string
		msg        string
		fields     []types.Field
		wantFields map[string]interface{}
	}{
		{
			name:    "debug level",
			logFunc: logger.Debug,
			level:   "DEBUG",
			msg:     "debug message",
			fields:  []types.Field{types.NewField("key", "value")},
			wantFields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "info level",
			logFunc: logger.Info,
			level:   "INFO",
			msg:     "info message",
			fields:  []types.Field{types.NewField("number", 42)},
			wantFields: map[string]interface{}{
				"number": float64(42),
			},
		},
		{
			name:    "warn level",
			logFunc: logger.Warn,
			level:   "WARN",
			msg:     "warn message",
			fields:  []types.Field{types.NewField("bool", true)},
			wantFields: map[string]interface{}{
				"bool": true,
			},
		},
		{
			name:    "error level",
			logFunc: logger.Error,
			level:   "ERROR",
			msg:     "error message",
			fields:  []types.Field{types.NewField("error", "failed")},
			wantFields: map[string]interface{}{
				"error": "failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log file before each test case
			if err := os.Truncate(logPath, 0); err != nil {
				t.Fatalf("Failed to clear log file: %v", err)
			}

			tt.logFunc(tt.msg, tt.fields...)
			content, err := os.ReadFile(logPath)
			require.NoError(t, err)
			contentStr := string(content)
			assert.Contains(t, contentStr, tt.level)
			assert.Contains(t, contentStr, tt.msg)
			for k, v := range tt.wantFields {
				assert.Contains(t, contentStr, k)
				assert.Contains(t, contentStr, v)
			}
		})
	}
}

func TestWithFields(t *testing.T) {
	logger, logPath := testlogger.NewTestLogger(t, "with-fields")
	defer logger.Close()

	// Create logger with fields
	fields := []types.Field{
		types.NewField("service", "test"),
		types.NewField("version", "1.0"),
	}
	loggerWithFields := logger.WithFields(fields...)

	// Log a message
	loggerWithFields.Info("test message")

	// Verify log output
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	contentStr := string(content)
	assert.Contains(t, contentStr, "test message")
	assert.Contains(t, contentStr, "service")
	assert.Contains(t, contentStr, "test")
	assert.Contains(t, contentStr, "version")
	assert.Contains(t, contentStr, "1.0")
}

func TestWithContext(t *testing.T) {
	logger, logPath := testlogger.NewTestLogger(t, "with-context")
	defer logger.Close()

	// Create context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "123")
	loggerWithCtx := logger.WithContext(ctx)

	// Log a message
	loggerWithCtx.Info("test message")

	// Verify log output
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	contentStr := string(content)
	assert.Contains(t, contentStr, "test message")
	assert.Contains(t, contentStr, "request_id")
	assert.Contains(t, contentStr, "123")
}

func TestLogRotation(t *testing.T) {
	// Create test directory
	testDir := t.TempDir()
	logPath := filepath.Join(testDir, "test.log")

	// Create logger
	cfg := Config{
		Level:      "info",
		OutputPath: logPath,
		MaxSize:    1, // 1MB
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   true,
	}

	logger, err := NewZapLogger(cfg)
	require.NoError(t, err)

	// Ensure logger is closed
	defer func() {
		require.NoError(t, logger.Close())
		// Wait a moment for Windows to release the file handle
		time.Sleep(10 * time.Millisecond)
	}()

	// Write enough logs to trigger rotation
	for i := 0; i < 100000; i++ {
		logger.Info("test message", types.NewField("count", i))
	}

	// Check if log file exists
	_, err = os.Stat(logPath)
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
