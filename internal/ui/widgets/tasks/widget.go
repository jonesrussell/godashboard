package tasks

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// Todo represents a single todo item
type Todo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo item
func NewTodo(content string) *Todo {
	now := time.Now()
	return &Todo{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleDone toggles the done status of the todo
func (t *Todo) ToggleDone() {
	t.Done = !t.Done
	t.UpdatedAt = time.Now()
}

// Widget represents the tasks widget
type Widget struct {
	todos    []*Todo
	selected int
	width    int
	height   int
	style    lipgloss.Style
	focused  bool
	debug    bool
}

// New creates a new tasks widget
func New() *Widget {
	return &Widget{
		todos: []*Todo{
			NewTodo("Implement dashboard"),
			NewTodo("Add more widgets"),
			NewTodo("Write tests"),
		},
		style: styles.ContentStyle,
	}
}

// Init implements components.Widget
func (w *Widget) Init() tea.Cmd {
	return nil
}

// Update implements components.Widget
func (w *Widget) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if w.selected > 0 {
				w.selected--
			}
		case "down", "j":
			if w.selected < len(w.todos)-1 {
				w.selected++
			}
		case " ":
			if w.selected < len(w.todos) {
				w.todos[w.selected].ToggleDone()
			}
		}
	}
	return w, nil
}

// View implements components.Widget
func (w *Widget) View() string {
	if w.width == 0 || w.height == 0 {
		return ""
	}

	var b strings.Builder

	// Title
	title := "Tasks"
	if w.focused {
		title = fmt.Sprintf("> %s", title)
	}
	b.WriteString(lipgloss.NewStyle().Bold(true).Render(title))
	b.WriteRune('\n')

	// Tasks
	for i, todo := range w.todos {
		// Task style
		style := lipgloss.NewStyle()
		if i == w.selected && w.focused {
			style = style.Background(styles.Primary).Foreground(lipgloss.Color("#ffffff"))
		}

		// Task status
		status := "[ ]"
		if todo.Done {
			status = "[✓]"
		}

		// Format task line
		taskLine := fmt.Sprintf("%s %s", status, todo.Content)
		b.WriteString(style.Render(taskLine))
		b.WriteRune('\n')
	}

	return w.style.Render(b.String())
}

// Focus implements components.Focusable
func (w *Widget) Focus() {
	w.focused = true
}

// Blur implements components.Focusable
func (w *Widget) Blur() {
	w.focused = false
}

// EnableDebug enables debug output
func (w *Widget) EnableDebug() {
	w.debug = true
}

// SetSize implements components.Widget
func (w *Widget) SetSize(width, height int) {
	w.width = width
	w.height = height
	w.style = w.style.Copy().Width(width).Height(height)
}
