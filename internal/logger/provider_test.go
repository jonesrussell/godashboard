package logger

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger/types"
	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	defaultLogger, err := ProvideLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, defaultLogger)
	t.Cleanup(func() {
		if defaultLogger != nil {
			if err := defaultLogger.Close(); err != nil {
				t.Logf("Failed to close default logger: %v", err)
			}
		}
	})

	// Test logging with provided logger
	defaultLogger.Info("test message", types.NewField("test", true))

	// Test with test logger
	testLogger, _ := testlogger.NewTestLogger(t, "provider-test")
	assert.NotNil(t, testLogger)

	// Test logging with test logger
	testLogger.Info("test message", types.NewField("test", true))

	// Test with invalid config
	invalidCfg := types.Config{
		Level:      "invalid",
		OutputPath: filepath.Join("non-existent-dir", strings.Repeat("a", 1000), "test.log"),
	}
	invalidLogger, err := ProvideLogger(invalidCfg)
	assert.Error(t, err)
	assert.Nil(t, invalidLogger)
}

func TestNew(t *testing.T) {
	logger, err := New(types.Config{
		Level:      "debug",
		OutputPath: filepath.Join(t.TempDir(), "test.log"),
	})
	require.NoError(t, err)
	require.NotNil(t, logger)
	defer func() {
		require.NoError(t, logger.Close())
	}()
}
