package logger

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger/types"
	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
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
	defer logger.Close()

	// Test logging with provided logger
	logger.Info("test message", types.NewField("test", true))

	// Test with test logger
	testLogger, _ := testlogger.NewTestLogger(t, "provider-test")
	defer testLogger.Close()

	// Test logging with test logger
	testLogger.Info("test message", types.NewField("test", true))

	// Test with invalid config
	invalidCfg := types.Config{
		Level:      "invalid",
		OutputPath: filepath.Join("non-existent-dir", strings.Repeat("a", 1000), "test.log"),
	}
	logger, err = ProvideLogger(invalidCfg)
	assert.Error(t, err)
	assert.Nil(t, logger)
}
