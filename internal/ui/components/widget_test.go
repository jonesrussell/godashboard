package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWidgetInterface(t *testing.T) {
	t.Run("basic widget functionality", func(t *testing.T) {
		// Create a mock widget
		widget := NewMockWidget().
			WithContent("Test Content").
			WithInitCmd(func() tea.Msg { return nil })

		// Test initialization
		cmd := widget.Init()
		assert.NotNil(t, cmd)

		// Test view
		assert.Equal(t, "Test Content", widget.View())

		// Test size setting
		widget.SetSize(100, 50)
		assert.Equal(t, 100, widget.width)
		assert.Equal(t, 50, widget.height)
	})

	t.Run("widget update handling", func(t *testing.T) {
		// Create a mock widget with custom update function
		updateCalled := false
		w := NewMockWidget()
		w.WithUpdateFn(func(msg tea.Msg) (tea.Model, tea.Cmd) {
			updateCalled = true
			return w, nil
		})

		// Send a message
		_, _ = w.Update(nil)
		assert.True(t, updateCalled)
	})

	t.Run("widget with UI test helper", func(t *testing.T) {
		// Create a mock widget
		w := NewMockWidget().WithContent("Initial Content")
		w.WithUpdateFn(func(msg tea.Msg) (tea.Model, tea.Cmd) {
			return w, nil
		})

		// Create UI test helper
		ui := testutil.NewUITest(t, w).
			WithSize(80, 24).
			Init()

		// Test initial state
		ui.AssertView("Initial Content").
			AssertNoCommands()

		// Test window size handling
		ui.SendWindowSize()

		// Test key handling
		ui.SendKey("q")
	})

	t.Run("widget command chain", func(t *testing.T) {
		// Create a widget that returns commands
		cmdExecuted := false
		w := NewMockWidget()
		w.WithInitCmd(func() tea.Msg {
			cmdExecuted = true
			return nil
		}).WithUpdateFn(func(msg tea.Msg) (tea.Model, tea.Cmd) {
			return w, nil
		})

		// Create UI test helper
		ui := testutil.NewUITest(t, w).
			Init()

		// Verify command was captured
		ui.AssertHasCommands()

		// Execute pending commands
		for _, cmd := range ui.Cmds {
			if cmd != nil {
				cmd()
			}
		}

		assert.True(t, cmdExecuted)
	})

	t.Run("widget content updates", func(t *testing.T) {
		// Create a widget that updates its content
		w := NewMockWidget().WithContent("Initial")
		w.WithUpdateFn(func(msg tea.Msg) (tea.Model, tea.Cmd) {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.String() == "u" {
					return NewMockWidget().WithContent("Updated"), nil
				}
			}
			return w, nil
		})

		// Create UI test helper
		ui := testutil.NewUITest(t, w).
			Init()

		// Test initial content
		ui.AssertView("Initial")

		// Send update key and verify content change
		ui.SendKey("u").
			AssertView("Updated")
	})
}

func TestWidgetSizing(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "normal size",
			width:      80,
			height:     24,
			wantWidth:  80,
			wantHeight: 24,
		},
		{
			name:       "minimum size",
			width:      1,
			height:     1,
			wantWidth:  1,
			wantHeight: 1,
		},
		{
			name:       "large size",
			width:      200,
			height:     100,
			wantWidth:  200,
			wantHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := NewMockWidget()
			widget.SetSize(tt.width, tt.height)

			assert.Equal(t, tt.wantWidth, widget.width)
			assert.Equal(t, tt.wantHeight, widget.height)
		})
	}
}
