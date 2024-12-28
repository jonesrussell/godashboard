package integration

import (
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/testutil"
	"github.com/jonesrussell/dashboard/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppInitialization verifies that the application initializes correctly
func TestAppInitialization(t *testing.T) {
	// Initialize logger
	cfg := logger.DefaultConfig()
	cfg.OutputPath = t.TempDir() + "/test.log" // Use test-specific log file
	log, err := logger.ProvideLogger(cfg)
	require.NoError(t, err)

	// Log test start
	log.Info("Starting app initialization test",
		logger.NewField("test", "initialization"),
		logger.NewField("logger", "configured"))

	// Initialize dashboard
	dashboard := ui.NewDashboard()
	require.NotNil(t, dashboard)
	log.Info("Dashboard initialized")

	// Create UI test helper
	ui := testutil.NewUITest(t, dashboard).
		WithSize(100, 40).
		Init()

	// Test initial state
	ui.AssertViewContains("Dashboard").
		AssertViewContains("Press ? for help")
	log.Info("Initial state verified")

	// Test help toggle
	ui.SendKey("?")
	ui.AssertViewContains("toggle help")
	log.Info("Help toggle tested")

	// Test quit command
	ui.SendKey("q")
	ui.AssertHasCommands() // Should have quit command
	log.Info("Quit command tested")
}

// TestAppLogging verifies that logging works throughout the application
func TestAppLogging(t *testing.T) {
	// Create a test-specific log file
	logPath := t.TempDir() + "/app.log"
	cfg := logger.DefaultConfig()
	cfg.OutputPath = logPath
	log, err := logger.ProvideLogger(cfg)
	require.NoError(t, err)

	// Log some test messages
	log.Info("Application starting", logger.NewField("test", true))
	log.Debug("Debug message")
	log.Warn("Warning message", logger.NewField("count", 42))

	// Verify log file exists and contains our messages
	logContents, err := testutil.ReadLogFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, logContents, "Application starting")
	assert.Contains(t, logContents, "test\":true")
	assert.Contains(t, logContents, "Warning message")
	assert.Contains(t, logContents, "count\":42")
}

// TestAppResize verifies that the application handles terminal resizing
func TestAppResize(t *testing.T) {
	// Initialize dashboard
	dashboard := ui.NewDashboard()

	// Create UI test helper
	ui := testutil.NewUITest(t, dashboard).
		Init()

	// Test different window sizes
	sizes := []struct {
		width  int
		height int
	}{
		{80, 24},  // Standard terminal
		{120, 40}, // Large terminal
		{40, 10},  // Small terminal
	}

	for _, size := range sizes {
		t.Run("size", func(t *testing.T) {
			ui.WithSize(size.width, size.height).
				SendWindowSize()

			// View should contain content regardless of size
			view := dashboard.View()
			assert.NotEmpty(t, view)
			assert.Contains(t, view, "Dashboard")
		})
	}
}
