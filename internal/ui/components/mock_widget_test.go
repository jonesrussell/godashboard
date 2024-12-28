package components

import tea "github.com/charmbracelet/bubbletea"

// mockWidget is a test implementation of the Widget interface
type mockWidget struct {
	width    int
	height   int
	content  string
	initCmd  tea.Cmd
	updateFn func(msg tea.Msg) (tea.Model, tea.Cmd)
}

// NewMockWidget creates a new mock widget for testing
func NewMockWidget() *mockWidget {
	return &mockWidget{
		content: "Mock Widget",
		updateFn: func(msg tea.Msg) (tea.Model, tea.Cmd) {
			return nil, nil
		},
		initCmd: nil,
	}
}

// WithContent sets the mock widget's content
func (w *mockWidget) WithContent(content string) *mockWidget {
	w.content = content
	return w
}

// WithInitCmd sets the mock widget's Init command
func (w *mockWidget) WithInitCmd(cmd tea.Cmd) *mockWidget {
	w.initCmd = cmd
	return w
}

// WithUpdateFn sets the mock widget's Update function
func (w *mockWidget) WithUpdateFn(fn func(msg tea.Msg) (tea.Model, tea.Cmd)) *mockWidget {
	w.updateFn = fn
	return w
}

// Init implements tea.Model
func (w *mockWidget) Init() tea.Cmd {
	return w.initCmd
}

// Update implements tea.Model
func (w *mockWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if w.updateFn != nil {
		return w.updateFn(msg)
	}
	return w, nil
}

// View implements tea.Model
func (w *mockWidget) View() string {
	return w.content
}

// SetSize implements Widget
func (w *mockWidget) SetSize(width, height int) {
	w.width = width
	w.height = height
}
