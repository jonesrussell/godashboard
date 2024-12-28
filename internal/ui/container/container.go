// Package container provides widget container and layout management
package container

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// Calculate grid dimensions with spacing
	cellWidth := (c.width - 4 - (c.cols-1)*2) / c.cols   // Account for spacing between cells
	cellHeight := (c.height - 4 - (c.rows-1)*1) / c.rows // Account for spacing between rows

	fmt.Printf("Grid dimensions: cols=%d, rows=%d\n", c.cols, c.rows)
	fmt.Printf("Cell dimensions: width=%d, height=%d\n", cellWidth, cellHeight)

	// Build grid view
	var grid [][]string
	for row := 0; row < c.rows; row++ {
		var rowContent []string
		for col := 0; col < c.cols; col++ {
			// Find widget for this cell
			var content string
			for _, entry := range c.widgets {
				if c.isWidgetInCell(entry, row, col) {
					width := cellWidth*entry.Config.ColSpan + (entry.Config.ColSpan-1)*2
					height := cellHeight*entry.Config.RowSpan + (entry.Config.RowSpan-1)*1

					// Ensure minimum width
					if width < entry.Config.MinWidth {
						width = entry.Config.MinWidth
					}

					fmt.Printf("Widget at [%d,%d]: width=%d, height=%d, focused=%v\n",
						row, col, width, height, entry.Focused)

					var style lipgloss.Style
					if entry.Focused {
						style = c.styleCache.GetFocusedStyle(width, height)
					} else {
						style = c.styleCache.GetContentStyle(width, height)
					}

					content = style.Render(entry.Widget.View())
					break
				}
			}
			if content == "" {
				fmt.Printf("Empty cell at [%d,%d]\n", row, col)
				style := c.styleCache.GetContentStyle(cellWidth, cellHeight)
				content = style.Render("")
			}
			rowContent = append(rowContent, content)

			// Add spacing between columns
			if col < c.cols-1 {
				rowContent = append(rowContent, strings.Repeat(" ", 2))
			}
		}
		grid = append(grid, rowContent)

		// Add spacing between rows
		if row < c.rows-1 {
			grid = append(grid, []string{strings.Repeat("\n", 1)})
		}
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
	// Ensure minimum container dimensions
	if c.width < 40 {
		c.width = 40
	}
	if c.height < 20 {
		c.height = 20
	}

	// Calculate cell dimensions with spacing
	cellWidth := (c.width - 4 - (c.cols-1)*2) / c.cols
	if cellWidth < 20 {
		cellWidth = 20
	}

	cellHeight := (c.height - 4 - (c.rows-1)*1) / c.rows
	if cellHeight < 10 {
		cellHeight = 10
	}

	for i := range c.widgets {
		entry := &c.widgets[i]
		width := cellWidth*entry.Config.ColSpan + (entry.Config.ColSpan-1)*2
		height := cellHeight*entry.Config.RowSpan + (entry.Config.RowSpan-1)*1

		// Ensure minimum width
		if width < entry.Config.MinWidth {
			width = entry.Config.MinWidth
		}

		entry.Widget.SetSize(width, height)
	}
}

// getWidgetWidth returns the width for a widget based on its configuration
func (c *Container) getWidgetWidth(entry WidgetEntry) int {
	cellWidth := (c.width - 4 - (c.cols-1)*2) / c.cols
	if cellWidth < 20 {
		cellWidth = 20
	}

	width := cellWidth*entry.Config.ColSpan + (entry.Config.ColSpan-1)*2
	if width < entry.Config.MinWidth {
		width = entry.Config.MinWidth
	}
	return width
}

// getWidgetHeight returns the height for a widget based on its configuration
func (c *Container) getWidgetHeight(entry WidgetEntry) int {
	cellHeight := (c.height - 4 - (c.rows-1)*1) / c.rows
	if cellHeight < 10 {
		cellHeight = 10
	}

	return cellHeight*entry.Config.RowSpan + (entry.Config.RowSpan-1)*1
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
