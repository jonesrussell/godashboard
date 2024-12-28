package integration

import (
	"os"
	"testing"

	"github.com/jonesrussell/dashboard/internal/logger"
	"github.com/jonesrussell/dashboard/internal/testutil"
	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
	"github.com/jonesrussell/dashboard/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppInitialization verifies that the application initializes correctly
func TestAppInitialization(t *testing.T) {
	log, _ := testlogger.NewTestLogger(t, "init")
	t.Cleanup(func() {
		if err := log.Close(); err != nil {
			t.Logf("Failed to close logger: %v", err)
		}
	})

	// Initialize dashboard
	dashboard := ui.NewDashboard(log)

	// Create UI test helper
	ui := testutil.NewUITest(t, dashboard).
		WithSize(80, 24)

	// Send window size before init to ensure proper layout
	ui.SendWindowSize()

	// Initialize and verify no pending commands
	ui.Init()

	// Test initial state
	ui.AssertViewContains("Dashboard")
	log.Info("Initial state tested")

	// Test help toggle
	ui.SendKey("?")
	ui.AssertViewContains("Help")
	log.Info("Help toggle tested")

	// Test quit command
	ui.SendKey("q")
	ui.AssertHasCommands() // Should have quit command
	log.Info("Quit command tested")
}

// TestAppLogging verifies that logging works throughout the application
func TestAppLogging(t *testing.T) {
	log, logPath := testlogger.NewTestLogger(t, "logging")
	defer func() {
		require.NoError(t, log.Close())
	}()

	// Log some test messages
	log.Info("Application starting", logger.NewField("test", true))
	log.Debug("Debug message")
	log.Warn("Warning message", logger.NewField("count", 42))

	// Verify log file exists and contains our messages
	logContents, err := os.ReadFile(logPath)
	require.NoError(t, err)
	logStr := string(logContents)
	assert.Contains(t, logStr, "Application starting")
	assert.Contains(t, logStr, "test\":true")
	assert.Contains(t, logStr, "Warning message")
	assert.Contains(t, logStr, "count\":42")
}

// TestAppResize verifies that the application handles terminal resizing
func TestAppResize(t *testing.T) {
	log, _ := testlogger.NewTestLogger(t, "resize")
	defer func() {
		require.NoError(t, log.Close())
	}()

	// Initialize dashboard
	dashboard := ui.NewDashboard(log)

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

func TestApp_Start(t *testing.T) {
	log, _ := testlogger.NewTestLogger(t, "app-test")
	defer func() {
		require.NoError(t, log.Close())
	}()
}

func TestApp_Stop(t *testing.T) {
	log, _ := testlogger.NewTestLogger(t, "app-test")
	defer func() {
		require.NoError(t, log.Close())
	}()
}
