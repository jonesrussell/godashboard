package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func BenchmarkDashboardInit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dashboard := NewDashboard()
		_ = dashboard.Init()
	}
}

func BenchmarkDashboardView(b *testing.B) {
	dashboard := NewDashboard()
	// Set a reasonable size
	dashboard.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = dashboard.View()
	}
}

func BenchmarkDashboardUpdate(b *testing.B) {
	dashboard := NewDashboard()
	msgs := []tea.Msg{
		tea.WindowSizeMsg{Width: 80, Height: 24},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, msg := range msgs {
			_, _ = dashboard.Update(msg)
		}
	}
}

func BenchmarkDashboardResize(b *testing.B) {
	dashboard := NewDashboard()
	sizes := []tea.WindowSizeMsg{
		{Width: 80, Height: 24},  // Standard
		{Width: 120, Height: 40}, // Large
		{Width: 40, Height: 10},  // Small
		{Width: 100, Height: 30}, // Medium
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, size := range sizes {
			_, _ = dashboard.Update(size)
			_ = dashboard.View()
		}
	}
}

func BenchmarkDashboardHelpToggle(b *testing.B) {
	dashboard := NewDashboard()
	// Initial setup
	dashboard.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Toggle help on
		_, _ = dashboard.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		_ = dashboard.View()
		// Toggle help off
		_, _ = dashboard.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		_ = dashboard.View()
	}
}