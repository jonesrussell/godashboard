// Package testutil provides testing utilities and helpers
package testutil

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/stretchr/testify/assert"
)

// LogEntry represents a parsed JSON log entry
type LogEntry struct {
	Level     string                 `json:"level"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"msg"`
	Fields    map[string]interface{} `json:"-"`
}

// NewTestLogger creates a new logger for testing
// For benchmarks, it uses a temporary file in the system temp directory
func NewTestLogger(tb testing.TB, name string) (logger.Logger, string) {
	tb.Helper()
	var logPath string
	if b, ok := tb.(*testing.B); ok {
		// For benchmarks, use system temp dir
		logPath = filepath.Join(os.TempDir(), "bench-"+b.Name()+"-"+name+".log")
	} else {
		// For tests, use test-specific temp dir
		logPath = filepath.Join(tb.TempDir(), name+".log")
	}

	cfg := logger.Config{
		Level:      "debug",
		OutputPath: logPath,
		MaxSize:    1,
		MaxBackups: 0,
		MaxAge:     1,
		Compress:   false,
		Debug:      true,
	}
	log, err := logger.NewZapLogger(cfg)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		_ = os.Remove(logPath)
	})
	return log, logPath
}

// ParseLogEntry parses a JSON log entry from a buffer
func ParseLogEntry(t *testing.T, line []byte) LogEntry {
	t.Helper()
	var entry LogEntry
	err := json.Unmarshal(line, &entry)
	assert.NoError(t, err, "failed to parse log entry")

	// Parse remaining fields into Fields map
	var raw map[string]interface{}
	err = json.Unmarshal(line, &raw)
	assert.NoError(t, err, "failed to parse raw log entry")

	// Remove known fields
	delete(raw, "level")
	delete(raw, "timestamp")
	delete(raw, "msg")
	entry.Fields = raw

	return entry
}

// AssertLogEntry asserts that a log entry matches expected values
func AssertLogEntry(t *testing.T, buf *bytes.Buffer, expectedLevel, expectedMsg string, expectedFields map[string]interface{}) {
	t.Helper()
	entry := ParseLogEntry(t, buf.Bytes())
	assert.Equal(t, expectedLevel, entry.Level)
	assert.Equal(t, expectedMsg, entry.Message)
	for k, v := range expectedFields {
		assert.Equal(t, v, entry.Fields[k])
	}
}

// ClearBuffer clears a bytes buffer and returns it
func ClearBuffer(buf *bytes.Buffer) *bytes.Buffer {
	buf.Reset()
	return buf
}

// ReadLogFile reads the contents of a log file
func ReadLogFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
