package logger

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "logs/app.log", cfg.OutputPath)
	assert.Equal(t, DefaultMaxSize, cfg.MaxSize)
	assert.Equal(t, DefaultMaxBackups, cfg.MaxBackups)
	assert.Equal(t, DefaultMaxAge, cfg.MaxAge)
	assert.True(t, cfg.Compress)
	assert.False(t, cfg.Debug)
}

func TestProvideLogger(t *testing.T) {
	// Test with default config
	cfg := DefaultConfig()
	logger, err := ProvideLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Test with invalid config
	invalidCfg := Config{
		Level:      "invalid",
		OutputPath: filepath.Join("non-existent-dir", strings.Repeat("a", 1000), "test.log"),
	}
	logger, err = ProvideLogger(invalidCfg)
	assert.Error(t, err)
	assert.Nil(t, logger)
}
