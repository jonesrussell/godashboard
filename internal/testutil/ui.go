package testutil

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// UITest provides utilities for testing Bubbletea components
type UITest struct {
	t      *testing.T
	model  tea.Model
	Cmds   []tea.Cmd
	width  int
	height int
}

// NewUITest creates a new UI test helper
func NewUITest(t *testing.T, model tea.Model) *UITest {
	return &UITest{
		t:      t,
		model:  model,
		width:  DefaultTestWindowWidth,
		height: DefaultTestWindowHeight,
	}
}

// WithSize sets the window size for the test
func (u *UITest) WithSize(width, height int) *UITest {
	u.width = width
	u.height = height
	return u
}

// Init initializes the model and captures commands
func (u *UITest) Init() *UITest {
	cmd := u.model.Init()
	if cmd != nil {
		u.Cmds = []tea.Cmd{cmd}
	}
	return u
}

// Send sends a message to the model and captures the result
func (u *UITest) Send(msg tea.Msg) *UITest {
	var cmd tea.Cmd
	u.model, cmd = u.model.Update(msg)
	if cmd != nil {
		u.Cmds = append(u.Cmds, cmd)
	}
	return u
}

// SendWindowSize sends a window size message
func (u *UITest) SendWindowSize() *UITest {
	return u.Send(tea.WindowSizeMsg{
		Width:  u.width,
		Height: u.height,
	})
}

// SendKey sends a key press message
func (u *UITest) SendKey(key string) *UITest {
	return u.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
}

// SendKeyType sends a special key type message
func (u *UITest) SendKeyType(keyType tea.KeyType) *UITest {
	return u.Send(tea.KeyMsg{Type: keyType})
}

// AssertView asserts that the current view matches the expected output
func (u *UITest) AssertView(expected string) *UITest {
	assert.Equal(u.t, expected, u.model.View())
	return u
}

// AssertViewContains asserts that the current view contains the expected string
func (u *UITest) AssertViewContains(expected string) *UITest {
	assert.Contains(u.t, u.model.View(), expected)
	return u
}

// AssertNoCommands asserts that no commands are pending
func (u *UITest) AssertNoCommands() *UITest {
	assert.Empty(u.t, u.Cmds)
	return u
}

// AssertHasCommands asserts that there are pending commands
func (u *UITest) AssertHasCommands() *UITest {
	assert.NotEmpty(u.t, u.Cmds)
	return u
}

// Model returns the current model state
func (u *UITest) Model() tea.Model {
	return u.model
}
