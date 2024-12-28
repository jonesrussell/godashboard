// Package assertions provides test assertion utilities
package assertions

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// LogEntry represents a parsed JSON log entry
type LogEntry struct {
	Level     string                 `json:"level"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"msg"`
	Fields    map[string]interface{} `json:"-"`
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
