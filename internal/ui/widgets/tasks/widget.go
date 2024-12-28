// Package tasks provides a widget for managing tasks in the dashboard.
// It implements the Widget interface and handles task creation, deletion,
// and status toggling through a terminal user interface.
package tasks

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// Widget represents the tasks widget
type Widget struct {
	components.BaseWidget
	client    *Client
	tasks     []Task
	selected  int
	loading   bool
	lastError error
}

// New creates a new tasks widget
func New(opts ...ClientOption) *Widget {
	return &Widget{
		client:   NewClient(opts...),
		selected: -1,
	}
}

// Init implements components.Widget
func (w *Widget) Init() tea.Cmd {
	return w.fetchTasks
}

// Update implements components.Widget
func (w *Widget) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	if !w.IsFocused() {
		return w, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if w.selected > 0 {
				w.selected--
			}
		case "down", "j":
			if w.selected < len(w.tasks)-1 {
				w.selected++
			}
		case " ":
			if len(w.tasks) > 0 {
				id := w.tasks[w.selected].ID
				return w, w.toggleTask(id)
			}
		case "d":
			if len(w.tasks) > 0 {
				id := w.tasks[w.selected].ID
				return w, w.deleteTask(id)
			}
		}
	case tasksMsg:
		w.tasks = []Task(msg)
		w.loading = false
		if w.selected >= len(w.tasks) {
			w.selected = len(w.tasks) - 1
		}
	case errorMsg:
		w.lastError = error(msg)
		w.loading = false
	case loadingMsg:
		w.loading = bool(msg)
	}

	return w, nil
}

// View implements components.Widget
func (w *Widget) View() string {
	width, height := w.GetDimensions()
	var b strings.Builder
	b.Grow(width * height)

	// Title
	b.WriteString(styles.Title.Render("Tasks"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	// Loading state
	if w.loading {
		loadingStyle := lipgloss.NewStyle().Foreground(styles.Subtle)
		b.WriteString(loadingStyle.Render("Loading..."))
		return w.GetStyle().Width(width).Height(height).Render(b.String())
	}

	// Error state
	if w.lastError != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
		b.WriteString(errorStyle.Render(w.lastError.Error()))
		return w.GetStyle().Width(width).Height(height).Render(b.String())
	}

	// Tasks
	if len(w.tasks) == 0 {
		subtleStyle := lipgloss.NewStyle().Foreground(styles.Subtle)
		b.WriteString(subtleStyle.Render("No tasks"))
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'n' to create a new task"))
	} else {
		for i, task := range w.tasks {
			// Task style
			taskStyle := lipgloss.NewStyle()
			if i == w.selected && w.IsFocused() {
				taskStyle = styles.Selected
			}

			// Task status
			status := "[ ]"
			if task.CompletedAt != nil && !task.CompletedAt.IsZero() {
				status = "[✓]"
			}

			// Format task line
			taskLine := fmt.Sprintf("%s %s", status, task.Title)
			if task.Description != "" {
				taskLine += fmt.Sprintf(" - %s", task.Description)
			}
			b.WriteString(taskStyle.Render(taskLine))
			b.WriteRune('\n')
		}
	}

	// Help text
	if w.IsFocused() {
		b.WriteString("\n")
		helpStyle := lipgloss.NewStyle().Foreground(styles.Subtle)
		b.WriteString(helpStyle.Render("↑/↓: select • space: toggle • n: new • d: delete"))
	}

	return w.GetStyle().Width(width).Height(height).Render(b.String())
}

// Message types
type tasksMsg []Task
type errorMsg error
type loadingMsg bool

// Commands
func (w *Widget) fetchTasks() tea.Msg {
	w.loading = true
	tasks, err := w.client.ListTasks()
	if err != nil {
		return errorMsg(err)
	}
	return tasksMsg(tasks)
}

func (w *Widget) toggleTask(id string) tea.Cmd {
	return func() tea.Msg {
		task := w.tasks[w.selected]
		input := TaskInput{
			Title:       task.Title,
			Description: task.Description,
		}

		now := time.Now()
		if task.CompletedAt == nil {
			task.CompletedAt = &now
		} else {
			task.CompletedAt = nil
		}

		if _, err := w.client.UpdateTask(id, input); err != nil {
			return errorMsg(err)
		}
		return w.fetchTasks()
	}
}

func (w *Widget) deleteTask(id string) tea.Cmd {
	return func() tea.Msg {
		err := w.client.DeleteTask(id)
		if err != nil {
			return errorMsg(err)
		}
		return w.fetchTasks()
	}
}
