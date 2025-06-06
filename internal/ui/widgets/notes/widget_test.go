package notes

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/dashboard/internal/testutil/testlogger"
	"github.com/stretchr/testify/assert"
)

func TestNotesWidget(t *testing.T) {
	log, _ := testlogger.NewTestLogger(t, "notes-test")
	w := New(log)

	t.Run("initialization", func(t *testing.T) {
		assert.NotNil(t, w)
		assert.Equal(t, 0, w.selected)
		assert.False(t, w.loading)
		assert.Nil(t, w.lastError)
	})

	t.Run("view states", func(t *testing.T) {
		t.Run("normal view", func(t *testing.T) {
			w := New(log)
			w.notes = []Note{{
				ID:        "1",
				Content:   "Test Note",
				Done:      false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}}

			view := w.View()
			assert.Contains(t, view, "Test Note")
		})

		t.Run("error view", func(t *testing.T) {
			w := New(log)
			w.lastError = assert.AnError
			view := w.View()
			assert.Contains(t, view, assert.AnError.Error())
		})
	})

	t.Run("focus handling", func(t *testing.T) {
		assert.False(t, w.IsFocused())
		w.Focus()
		assert.True(t, w.IsFocused())
		w.Blur()
		assert.False(t, w.IsFocused())
	})

	t.Run("navigation", func(t *testing.T) {
		w := New(log)
		w.notes = []Note{
			{ID: "1", Content: "Note 1"},
			{ID: "2", Content: "Note 2"},
		}
		w.Focus()

		tests := []struct {
			name     string
			key      string
			expected int
		}{
			{"move down with j", "j", 1},
			{"move down with down", "down", 1},
			{"move up with k", "k", 0},
			{"move up with up", "up", 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				model, cmd := w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
				widget := model.(*Widget)
				assert.Equal(t, tt.expected, widget.selected)
				assert.Nil(t, cmd)
			})
		}
	})

	t.Run("commands", func(t *testing.T) {
		t.Run("fetch notes", func(t *testing.T) {
			msg := w.fetchNotes()
			assert.True(t, w.loading)
			assert.NotNil(t, msg)
			_, ok := msg.(errorMsg)
			assert.True(t, ok, "expected error message type")
		})

		t.Run("note operations", func(t *testing.T) {
			w := New(log)
			w.notes = []Note{{
				ID:      "1",
				Content: "Test Note",
				Done:    false,
			}}
			w.selected = 0

			tests := []struct {
				name   string
				getCmd func() tea.Msg
			}{
				{"toggle note", func() tea.Msg { return w.toggleNote("1")() }},
				{"delete note", func() tea.Msg { return w.deleteNote("1")() }},
				{"create note", func() tea.Msg { return w.createNote() }},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					msg := tt.getCmd()
					assert.NotNil(t, msg)
					_, ok := msg.(errorMsg)
					assert.True(t, ok, "expected error message type")
				})
			}
		})
	})

	t.Run("key commands", func(t *testing.T) {
		w := New(log)
		w.notes = []Note{{
			ID:      "1",
			Content: "Test Note",
			Done:    false,
		}}
		w.Focus()

		tests := []struct {
			name    string
			key     string
			wantCmd bool
		}{
			{"space to toggle", " ", true},
			{"d to delete", "d", true},
			{"n to create", "n", true},
			{"invalid key", "x", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, cmd := w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
				if tt.wantCmd {
					assert.NotNil(t, cmd)
				} else {
					assert.Nil(t, cmd)
				}
			})
		}
	})
}
