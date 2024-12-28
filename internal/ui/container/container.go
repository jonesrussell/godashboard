// Package container provides a flexible container for organizing UI widgets
package container

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/components"
)

const (
	// Minimum dimensions
	minCellWidth  = 20
	minCellHeight = 10
	minWidth      = 40
	minHeight     = 20
)

// Container manages the layout and interaction of multiple widgets
type Container struct {
	width    int
	height   int
	focused  bool
	entries  []WidgetEntry
	selected int
}

// WidgetEntry represents a widget and its layout properties
type WidgetEntry struct {
	Widget    components.Widget
	Row       int
	Col       int
	RowSpan   int
	ColSpan   int
	MinWidth  int
	MinHeight int
}

// New creates a new container instance
func New() *Container {
	return &Container{
		entries:  make([]WidgetEntry, 0),
		selected: -1,
	}
}

// AddWidget adds a widget to the container
func (c *Container) AddWidget(widget components.Widget, row, col, rowSpan, colSpan int) {
	entry := WidgetEntry{
		Widget:    widget,
		Row:       row,
		Col:       col,
		RowSpan:   rowSpan,
		ColSpan:   colSpan,
		MinWidth:  minCellWidth,
		MinHeight: minCellHeight,
	}
	c.entries = append(c.entries, entry)
}

// Init implements tea.Model
func (c *Container) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, entry := range c.entries {
		if cmd := entry.Widget.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update implements components.Widget
func (c *Container) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("tab"))) {
			c.focusNext()
			return c, nil
		}
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.updateWidgetSizes()
	}

	// Update focused widget
	if c.selected >= 0 && c.selected < len(c.entries) {
		if widget, cmd := c.entries[c.selected].Widget.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
			c.entries[c.selected].Widget = widget
		}
	}

	return c, tea.Batch(cmds...)
}

// View implements tea.Model
func (c *Container) View() string {
	if !c.hasValidDimensions() {
		return "Window too small"
	}

	rows := make([]string, 0)
	for row := 0; row < c.getRowCount(); row++ {
		rowContent := make([]string, 0)
		for col := 0; col < c.getColCount(); col++ {
			widget := c.getWidgetAt(row, col)
			if widget != nil {
				rowContent = append(rowContent, widget.View())
			} else {
				rowContent = append(rowContent, "")
			}
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowContent...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// SetSize implements components.Widget
func (c *Container) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.updateWidgetSizes()
}

// GetDimensions implements components.Widget
func (c *Container) GetDimensions() (width, height int) {
	return c.width, c.height
}

// Focus implements components.Widget
func (c *Container) Focus() {
	c.focused = true
	if c.selected == -1 && len(c.entries) > 0 {
		c.selected = 0
		c.entries[0].Widget.Focus()
	}
}

// Blur implements components.Widget
func (c *Container) Blur() {
	c.focused = false
	if c.selected >= 0 && c.selected < len(c.entries) {
		c.entries[c.selected].Widget.Blur()
	}
}

// IsFocused implements components.Widget
func (c *Container) IsFocused() bool {
	return c.focused
}

func (c *Container) hasValidDimensions() bool {
	return c.width >= minWidth && c.height >= minHeight
}

func (c *Container) focusNext() {
	if len(c.entries) == 0 {
		return
	}

	if c.selected >= 0 {
		c.entries[c.selected].Widget.Blur()
	}

	c.selected = (c.selected + 1) % len(c.entries)
	c.entries[c.selected].Widget.Focus()
}

func (c *Container) updateWidgetSizes() {
	cellWidth := c.width / c.getColCount()
	cellHeight := c.height / c.getRowCount()

	if cellWidth < minCellWidth || cellHeight < minCellHeight {
		return
	}

	for i := range c.entries {
		width := cellWidth * c.entries[i].ColSpan
		height := cellHeight * c.entries[i].RowSpan
		c.entries[i].Widget.SetSize(width, height)
	}
}

func (c *Container) getWidgetAt(row, col int) components.Widget {
	for _, entry := range c.entries {
		if entry.Row <= row && row < entry.Row+entry.RowSpan &&
			entry.Col <= col && col < entry.Col+entry.ColSpan {
			return entry.Widget
		}
	}
	return nil
}

func (c *Container) getRowCount() int {
	maxRow := 0
	for _, entry := range c.entries {
		if entry.Row+entry.RowSpan > maxRow {
			maxRow = entry.Row + entry.RowSpan
		}
	}
	return maxRow
}

func (c *Container) getColCount() int {
	maxCol := 0
	for _, entry := range c.entries {
		if entry.Col+entry.ColSpan > maxCol {
			maxCol = entry.Col + entry.ColSpan
		}
	}
	return maxCol
}
