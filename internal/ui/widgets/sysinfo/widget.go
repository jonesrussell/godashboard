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

// Widget displays system information like CPU, memory, and disk usage
type Widget struct {
	width  int
	height int
	style  lipgloss.Style

	// System info
	cpuUsage    float64
	memoryUsage float64
	diskUsage   float64

	// Update ticker
	updateInterval time.Duration
}

// New creates a new system information widget
func New() *Widget {
	return &Widget{
		style:          styles.ContentStyle,
		updateInterval: 2 * time.Second,
	}
}

// Init implements tea.Model
func (w *Widget) Init() tea.Cmd {
	return tea.Batch(
		w.updateSystemInfo,
		w.tick(),
	)
}

// Update implements components.Widget
func (w *Widget) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case systemInfoMsg:
		w.cpuUsage = msg.cpu
		w.memoryUsage = msg.memory
		w.diskUsage = msg.disk
		return w, w.tick()

	case tea.WindowSizeMsg:
		w.SetSize(msg.Width, msg.Height)
	}

	return w, nil
}

// View implements components.Widget
func (w *Widget) View() string {
	var b strings.Builder
	b.Grow(w.width * w.height)

	// Format system info
	b.WriteString("System Information\n\n")
	b.WriteString(fmt.Sprintf("CPU Usage:    %.1f%%\n", w.cpuUsage))
	b.WriteString(fmt.Sprintf("Memory Usage: %.1f%%\n", w.memoryUsage))
	b.WriteString(fmt.Sprintf("Disk Usage:   %.1f%%\n", w.diskUsage))

	return w.style.Width(w.width).Height(w.height).Render(b.String())
}

// SetSize implements components.Widget
func (w *Widget) SetSize(width, height int) {
	w.width = width
	w.height = height
}

// systemInfoMsg carries system information updates
type systemInfoMsg struct {
	cpu    float64
	memory float64
	disk   float64
}

// tick returns a command that waits for the update interval
func (w *Widget) tick() tea.Cmd {
	return tea.Tick(w.updateInterval, func(time.Time) tea.Msg {
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
