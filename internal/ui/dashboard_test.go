package ui

import (
	"testing"

	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	logger, _ := testlogger.NewTestLogger(t, "dashboard")
	dash := NewDashboard(logger)

	// Test initial state
	assert.False(t, dash.showHelp)
	assert.False(t, dash.debug)
	assert.NotNil(t, dash.container)
	assert.NotNil(t, dash.logger)
}
