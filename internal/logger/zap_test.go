package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name      string
		config    types.Config
		wantError bool
	}{
		{
			name: "valid config",
			config: types.Config{
				Level:      "info",
				OutputPath: filepath.Join(t.TempDir(), "test.log"),
			},
			wantError: false,
		},
		{
			name: "invalid path",
			config: types.Config{
				Level:      "info",
				OutputPath: filepath.Join("non-existent-dir", strings.Repeat("a", 1000), "test.log"),
			},
			wantError: true,
		},
		{
			name: "invalid level",
			config: types.Config{
				Level:      "invalid",
				OutputPath: filepath.Join(t.TempDir(), "test.log"),
			},
			wantError: false, // Level defaults to info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewZapLogger(tt.config)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				if logger != nil {
					require.NoError(t, logger.Close())
				}
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	logger, logPath := setupTestLogger(t, "log-levels")
	defer func() {
		require.NoError(t, logger.Close())
	}()

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
			fields:  []types.Field{types.NewField("error", "test error")},
			wantFields: map[string]interface{}{
				"error": "test error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write log message
			tt.logFunc(tt.msg, tt.fields...)

			// Read and verify log content
			content, err := os.ReadFile(logPath)
			require.NoError(t, err)
			contentStr := string(content)

			// Verify log content
			assert.Contains(t, contentStr, tt.level)
			assert.Contains(t, contentStr, tt.msg)
			for k, v := range tt.wantFields {
				assert.Contains(t, contentStr, k)
				assert.Contains(t, contentStr, fmt.Sprint(v))
			}

			// Clear log file for next test
			err = os.WriteFile(logPath, []byte{}, 0644)
			require.NoError(t, err)
		})
	}
}

func TestWithFields(t *testing.T) {
	logger, logPath := setupTestLogger(t, "with-fields")
	defer func() {
		require.NoError(t, logger.Close())
	}()

	// Create logger with fields
	fields := []types.Field{
		types.NewField("service", "test"),
		types.NewField("version", "1.0"),
	}
	loggerWithFields := logger.WithFields(fields...)

	// Log message with additional fields
	loggerWithFields.Info("test message",
		types.NewField("extra", "field"),
	)

	// Verify log output
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	contentStr := string(content)

	// Check all fields are present
	assert.Contains(t, contentStr, "service")
	assert.Contains(t, contentStr, "test")
	assert.Contains(t, contentStr, "version")
	assert.Contains(t, contentStr, "1.0")
	assert.Contains(t, contentStr, "extra")
	assert.Contains(t, contentStr, "field")
}

func TestWithContext(t *testing.T) {
	logger, logPath := setupTestLogger(t, "with-context")
	defer func() {
		require.NoError(t, logger.Close())
	}()

	// Create context with request ID
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Create logger with context
	loggerWithContext := logger.WithContext(ctx)

	// Log message
	loggerWithContext.Info("test message")

	// Verify log output
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	contentStr := string(content)

	// Check request ID is present
	assert.Contains(t, contentStr, "request_id")
	assert.Contains(t, contentStr, "test-request-id")
}

func TestZapLogger_Close(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "zap-test-*.log")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(tmpFile.Name()))
	}()

	logger, err := NewZapLogger(types.Config{
		Level:      "debug",
		OutputPath: tmpFile.Name(),
	})
	require.NoError(t, err)

	// Test closing the logger
	require.NoError(t, logger.Close())
}

func TestZapLogger_WithFields(t *testing.T) {
	logger, err := NewZapLogger(types.Config{
		Level:      "debug",
		OutputPath: "test.log",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, logger.Close())
	}()

	// Create logger with fields
	fields := []types.Field{
		types.NewField("service", "test"),
		types.NewField("version", "1.0"),
	}
	loggerWithFields := logger.WithFields(fields...)

	// Log message with additional fields
	loggerWithFields.Info("test message",
		types.NewField("extra", "field"),
	)

	// Verify log output
	content, err := os.ReadFile("test.log")
	require.NoError(t, err)
	contentStr := string(content)

	// Check all fields are present
	assert.Contains(t, contentStr, "service")
	assert.Contains(t, contentStr, "test")
	assert.Contains(t, contentStr, "version")
	assert.Contains(t, contentStr, "1.0")
	assert.Contains(t, contentStr, "extra")
	assert.Contains(t, contentStr, "field")
}

func TestZapLogger_WithContext(t *testing.T) {
	logger, err := NewZapLogger(types.Config{
		Level:      "debug",
		OutputPath: "test.log",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, logger.Close())
	}()

	// Create context with request ID
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Create logger with context
	loggerWithContext := logger.WithContext(ctx)

	// Log message
	loggerWithContext.Info("test message")

	// Verify log output
	content, err := os.ReadFile("test.log")
	require.NoError(t, err)
	contentStr := string(content)

	// Check request ID is present
	assert.Contains(t, contentStr, "request_id")
	assert.Contains(t, contentStr, "test-request-id")
}

func TestZapLogger_LogLevels(t *testing.T) {
	logger, err := NewZapLogger(types.Config{
		Level:      "debug",
		OutputPath: "test.log",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, logger.Close())
	}()

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
			fields:  []types.Field{types.NewField("error", "test error")},
			wantFields: map[string]interface{}{
				"error": "test error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write log message
			tt.logFunc(tt.msg, tt.fields...)

			// Read and verify log content
			content, err := os.ReadFile("test.log")
			require.NoError(t, err)
			contentStr := string(content)

			// Verify log content
			assert.Contains(t, contentStr, tt.level)
			assert.Contains(t, contentStr, tt.msg)
			for k, v := range tt.wantFields {
				assert.Contains(t, contentStr, k)
				assert.Contains(t, contentStr, fmt.Sprint(v))
			}

			// Clear log file for next test
			err = os.WriteFile("test.log", []byte{}, 0644)
			require.NoError(t, err)
		})
	}
}

func setupTestLogger(t *testing.T, name string) (types.Logger, string) {
	logPath := filepath.Join(t.TempDir(), name+".log")
	logger, err := NewZapLogger(types.Config{
		Level:      "debug",
		OutputPath: logPath,
	})
	require.NoError(t, err)
	return logger, logPath
}
