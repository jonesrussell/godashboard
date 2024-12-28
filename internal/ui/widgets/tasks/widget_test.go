package tasks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskCompletionStatus(t *testing.T) {
	widget := New()

	// Setup test cases
	zeroTime := time.Time{} // zero value
	now := time.Now()

	testCases := []struct {
		name           string
		completedAt    *time.Time
		expectedStatus string
	}{
		{
			name:           "nil completion time",
			completedAt:    nil,
			expectedStatus: "[ ]",
		},
		{
			name:           "zero completion time",
			completedAt:    &zeroTime,
			expectedStatus: "[ ]",
		},
		{
			name:           "valid completion time",
			completedAt:    &now,
			expectedStatus: "[✓]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test task
			widget.tasks = []Task{{
				ID:          "test-id",
				Title:       "Test Task",
				CompletedAt: tc.completedAt,
			}}

			// Render the widget
			view := widget.View()

			// Check if the correct status is rendered
			assert.Contains(t, view, tc.expectedStatus,
				"Task should show correct completion status")
		})
	}
}

func TestTaskToggle(t *testing.T) {
	// Mock client implementation needed
	// TODO: Add tests for task toggling behavior
}
