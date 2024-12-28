// Package sysinfo provides a widget for displaying system information
package sysinfo

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Widget represents the system information widget
type Widget struct {
	components.BaseWidget
	cpuUsage    float64
	memoryUsage float64
	diskUsage   float64
}

// New creates a new system information widget
func New() *Widget {
	return &Widget{}
}

// Init implements components.Widget
func (w *Widget) Init() tea.Cmd {
	return w.tick()
}

// Update implements components.Widget
func (w *Widget) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case systemInfoMsg:
		w.cpuUsage = msg.cpu
		w.memoryUsage = msg.memory
		w.diskUsage = msg.disk
		return w, w.tick()
	}
	return w, nil
}

// View implements components.Widget
func (w *Widget) View() string {
	width, height := w.GetDimensions()
	var b strings.Builder
	b.Grow(width * height)

	// Title
	b.WriteString(styles.Title.Render("System Information"))
	b.WriteString("\n\n")

	// Calculate bar width (minimum 10 characters)
	barWidth := width - 20
	if barWidth < 10 {
		barWidth = 10
	}

	// Format system info with bars
	cpuBar := createUsageBar(w.cpuUsage, barWidth)
	memBar := createUsageBar(w.memoryUsage, barWidth)
	diskBar := createUsageBar(w.diskUsage, barWidth)

	// CPU
	b.WriteString(styles.Title.Render("CPU"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%.1f%% ", w.cpuUsage))
	b.WriteString(cpuBar)
	b.WriteString("\n\n")

	// Memory
	b.WriteString(styles.Title.Render("Memory"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%.1f%% ", w.memoryUsage))
	b.WriteString(memBar)
	b.WriteString("\n\n")

	// Disk
	b.WriteString(styles.Title.Render("Disk"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%.1f%% ", w.diskUsage))
	b.WriteString(diskBar)

	return w.GetStyle().Width(width).Height(height).Render(b.String())
}

// SetSize implements components.Widget
func (w *Widget) SetSize(width, height int) {
	w.width = width
	w.height = height
}

// Focus implements components.Focusable
func (w *Widget) Focus() {
	w.focused = true
}

// Blur implements components.Focusable
func (w *Widget) Blur() {
	w.focused = false
}

// createUsageBar creates a progress bar for the given percentage
func createUsageBar(percent float64, width int) string {
	// Ensure valid percentage
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	// Calculate filled and empty portions
	filled := int(float64(width) * percent / 100)
	if filled > width {
		filled = width
	}
	empty := width - filled

	// Create the bar with colors
	bar := lipgloss.NewStyle().Foreground(styles.Primary).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(styles.Subtle).Render(strings.Repeat("░", empty))

	return bar
}

// systemInfoMsg carries system information updates
type systemInfoMsg struct {
	cpu    float64
	memory float64
	disk   float64
}

// tick returns a command that waits for the update interval
func (w *Widget) tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return updateSystemInfoMsg{}
	})
}

// updateSystemInfoMsg triggers a system info update
type updateSystemInfoMsg struct{}

// updateSystemInfo returns a command that updates system information
func (w *Widget) updateSystemInfo() tea.Msg {
	// Get CPU usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		cpuPercent = []float64{0}
	}

	// Get memory usage
	memInfo, err := mem.VirtualMemory()
	memPercent := 0.0
	if err == nil {
		memPercent = memInfo.UsedPercent
	}

	// Get disk usage
	diskInfo, err := disk.Usage("/")
	diskPercent := 0.0
	if err == nil {
		diskPercent = diskInfo.UsedPercent
	}

	return systemInfoMsg{
		cpu:    cpuPercent[0],
		memory: memPercent,
		disk:   diskPercent,
	}
}
