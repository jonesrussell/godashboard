package styles

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// StyleCache caches Lipgloss styles for better performance
type StyleCache struct {
	mu     sync.RWMutex
	styles map[string]lipgloss.Style
}

// NewStyleCache creates a new style cache
func NewStyleCache() *StyleCache {
	return &StyleCache{
		styles: make(map[string]lipgloss.Style),
	}
}

// GetContentStyle returns a cached content style for the given dimensions
func (c *StyleCache) GetContentStyle(width, height int) lipgloss.Style {
	key := fmt.Sprintf("content-%dx%d", width, height)

	c.mu.RLock()
	if style, ok := c.styles[key]; ok {
		c.mu.RUnlock()
		return style
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("69"))

	c.styles[key] = style
	return style
}

// Clear clears the style cache
func (c *StyleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.styles = make(map[string]lipgloss.Style)
}
