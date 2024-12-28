// Package container provides widget container and layout management
package container

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// GridConfig defines the configuration for a grid cell
type GridConfig struct {
	Row      int
	Col      int
	RowSpan  int
	ColSpan  int
	MinWidth int
}

// WidgetEntry represents a widget with its grid configuration
type WidgetEntry struct {
	Widget  components.Widget
	Config  GridConfig
	Focused bool
}

// Container manages multiple widgets in a grid layout
type Container struct {
	// Layout properties
	width  int
	height int
	rows   int
	cols   int

	// Widget management
	widgets    []WidgetEntry
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
		widgets:      make([]WidgetEntry, 0),
		focusIndex:   -1,
		styleCache:   styles.NewStyleCache(),
		contentCache: make(map[string]string),
		needsRefresh: true,
	}
}

// AddWidget adds a widget to the container with default configuration
func (c *Container) AddWidget(w components.Widget) {
	c.AddWidgetWithConfig(w, GridConfig{
		Row:      len(c.widgets) / c.cols,
		Col:      len(c.widgets) % c.cols,
		RowSpan:  1,
		ColSpan:  1,
		MinWidth: 20,
	})
}

// AddWidgetWithConfig adds a widget with specific grid configuration
func (c *Container) AddWidgetWithConfig(w components.Widget, config GridConfig) {
	entry := WidgetEntry{
		Widget:  w,
		Config:  config,
		Focused: len(c.widgets) == 0,
	}
	c.widgets = append(c.widgets, entry)
	if c.focusIndex == -1 {
		c.focusIndex = 0
	}
	c.needsRefresh = true
}

// Init implements tea.Model
func (c *Container) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, entry := range c.widgets {
		if cmd := entry.Widget.Init(); cmd != nil {
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
			for i := range c.widgets {
				c.widgets[i].Focused = i == c.focusIndex
			}
			c.needsRefresh = true
			return c, nil
		}

		// Forward message to focused widget
		if c.focusIndex >= 0 && c.focusIndex < len(c.widgets) {
			entry := &c.widgets[c.focusIndex]
			newWidget, cmd := entry.Widget.Update(msg)
			if w, ok := newWidget.(components.Widget); ok {
				entry.Widget = w
				c.needsRefresh = true
				return c, cmd
			}
		}

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.needsRefresh = true

		// Calculate grid dimensions
		c.updateGridSizes()

		// Forward size message to widgets
		for i := range c.widgets {
			entry := &c.widgets[i]
			if newWidget, cmd := entry.Widget.Update(msg); cmd != nil {
				if w, ok := newWidget.(components.Widget); ok {
					entry.Widget = w
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

	// Build grid view
	var grid [][]string
	for row := 0; row < c.rows; row++ {
		var rowContent []string
		for col := 0; col < c.cols; col++ {
			// Find widget for this cell
			var content string
			for _, entry := range c.widgets {
				if c.isWidgetInCell(entry, row, col) {
					style := c.styleCache.GetContentStyle(
						c.getWidgetWidth(entry),
						c.getWidgetHeight(entry),
					)
					if entry.Focused {
						style = style.BorderForeground(styles.Primary)
					}
					content = style.Render(entry.Widget.View())
					break
				}
			}
			if content == "" {
				style := c.styleCache.GetContentStyle(
					c.width/c.cols-4,
					c.height/c.rows-2,
				)
				content = style.Render("")
			}
			rowContent = append(rowContent, content)
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

// updateGridSizes calculates and updates widget sizes based on grid configuration
func (c *Container) updateGridSizes() {
	baseWidth := (c.width - 4) / c.cols
	baseHeight := (c.height - 4) / c.rows

	for i := range c.widgets {
		entry := &c.widgets[i]
		width := baseWidth*entry.Config.ColSpan - 4
		height := baseHeight*entry.Config.RowSpan - 2

		// Ensure minimum width
		if width < entry.Config.MinWidth {
			width = entry.Config.MinWidth
		}

		entry.Widget.SetSize(width, height)
	}
}

// getWidgetWidth returns the width for a widget based on its configuration
func (c *Container) getWidgetWidth(entry WidgetEntry) int {
	baseWidth := (c.width - 4) / c.cols
	width := baseWidth*entry.Config.ColSpan - 4
	if width < entry.Config.MinWidth {
		width = entry.Config.MinWidth
	}
	return width
}

// getWidgetHeight returns the height for a widget based on its configuration
func (c *Container) getWidgetHeight(entry WidgetEntry) int {
	baseHeight := (c.height - 4) / c.rows
	return baseHeight*entry.Config.RowSpan - 2
}

// isWidgetInCell checks if a widget occupies the given grid cell
func (c *Container) isWidgetInCell(entry WidgetEntry, row, col int) bool {
	return row >= entry.Config.Row &&
		row < entry.Config.Row+entry.Config.RowSpan &&
		col >= entry.Config.Col &&
		col < entry.Config.Col+entry.Config.ColSpan
}

// FocusedWidget returns the currently focused widget
func (c *Container) FocusedWidget() components.Widget {
	if c.focusIndex >= 0 && c.focusIndex < len(c.widgets) {
		return c.widgets[c.focusIndex].Widget
	}
	return nil
}

// HandleFocusKey handles focus-related key events
func (c *Container) HandleFocusKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyTab:
		c.focusIndex = (c.focusIndex + 1) % len(c.widgets)
		for i := range c.widgets {
			c.widgets[i].Focused = i == c.focusIndex
		}
		c.needsRefresh = true
		return true
	case tea.KeyShiftTab:
		c.focusIndex--
		if c.focusIndex < 0 {
			c.focusIndex = len(c.widgets) - 1
		}
		for i := range c.widgets {
			c.widgets[i].Focused = i == c.focusIndex
		}
		c.needsRefresh = true
		return true
	}
	return false
}
