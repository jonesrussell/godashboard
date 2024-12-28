// Package container provides widget container and layout management
package container

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// Container manages multiple widgets in a grid layout
type Container struct {
	// Layout properties
	width  int
	height int
	rows   int
	cols   int

	// Widget management
	widgets    []components.Widget
	focusIndex int
	styleCache *styles.StyleCache

	// Content cache
	contentCache map[string]string
	needsRefresh bool
}

// New creates a new widget container
func New(rows, cols int) *Container {
	return &Container{
		rows:         rows,
		cols:         cols,
		widgets:      make([]components.Widget, 0),
		focusIndex:   -1,
		styleCache:   styles.NewStyleCache(),
		contentCache: make(map[string]string),
		needsRefresh: true,
	}
}

// AddWidget adds a widget to the container
func (c *Container) AddWidget(w components.Widget) {
	c.widgets = append(c.widgets, w)
	if c.focusIndex == -1 {
		c.focusIndex = 0
	}
	c.needsRefresh = true
}

// Init implements tea.Model
func (c *Container) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, w := range c.widgets {
		if cmd := w.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update implements tea.Model
func (c *Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle tab navigation between widgets
		if msg.Type == tea.KeyTab {
			c.focusIndex = (c.focusIndex + 1) % len(c.widgets)
			c.needsRefresh = true
			return c, nil
		}

		// Forward message to focused widget
		if c.focusIndex >= 0 && c.focusIndex < len(c.widgets) {
			widget := c.widgets[c.focusIndex]
			newWidget, cmd := widget.Update(msg)
			if w, ok := newWidget.(components.Widget); ok {
				c.widgets[c.focusIndex] = w
				c.needsRefresh = true
				return c, cmd
			}
		}

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.needsRefresh = true

		// Update widget sizes
		cellWidth := (c.width - 4) / c.cols
		cellHeight := (c.height - 4) / c.rows

		for i, w := range c.widgets {
			w.SetSize(cellWidth, cellHeight)

			// Forward size message to widget
			if newWidget, cmd := w.Update(msg); cmd != nil {
				if w, ok := newWidget.(components.Widget); ok {
					c.widgets[i] = w
				}
			}
		}
	}

	return c, nil
}

// View implements tea.Model
func (c *Container) View() string {
	if len(c.widgets) == 0 {
		return "No widgets"
	}

	// Calculate grid layout
	cellWidth := (c.width - 4) / c.cols
	cellHeight := (c.height - 4) / c.rows

	// Build grid view
	var grid [][]string
	for row := 0; row < c.rows; row++ {
		var rowContent []string
		for col := 0; col < c.cols; col++ {
			idx := row*c.cols + col
			if idx < len(c.widgets) {
				widget := c.widgets[idx]
				style := c.styleCache.GetContentStyle(cellWidth, cellHeight)
				if idx == c.focusIndex {
					style = style.BorderForeground(styles.Primary)
				}
				rowContent = append(rowContent, style.Render(widget.View()))
			} else {
				style := c.styleCache.GetContentStyle(cellWidth, cellHeight)
				rowContent = append(rowContent, style.Render(""))
			}
		}
		grid = append(grid, rowContent)
	}

	// Render grid
	var b strings.Builder
	b.Grow(c.width * c.height)

	for _, row := range grid {
		for _, cell := range row {
			b.WriteString(cell)
		}
		b.WriteByte('\n')
	}

	return b.String()
}

// FocusedWidget returns the currently focused widget
func (c *Container) FocusedWidget() components.Widget {
	if c.focusIndex >= 0 && c.focusIndex < len(c.widgets) {
		return c.widgets[c.focusIndex]
	}
	return nil
}

// HandleFocusKey handles focus-related key events
func (c *Container) HandleFocusKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyTab:
		c.focusIndex = (c.focusIndex + 1) % len(c.widgets)
		c.needsRefresh = true
		return true
	case tea.KeyShiftTab:
		c.focusIndex--
		if c.focusIndex < 0 {
			c.focusIndex = len(c.widgets) - 1
		}
		c.needsRefresh = true
		return true
	}
	return false
}
