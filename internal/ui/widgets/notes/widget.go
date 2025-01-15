package notes

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/dashboard/internal/ui/components"
	"github.com/jonesrussell/dashboard/internal/ui/styles"
)

// Widget represents the notes widget
type Widget struct {
	components.BaseWidget
	client    *Client
	notes     []Note
	selected  int
	loading   bool
	lastError error
}

// New creates a new notes widget
func New(opts ...ClientOption) *Widget {
	return &Widget{
		client:   NewClient(opts...),
		notes:    make([]Note, 0),
		selected: 0,
	}
}

// Init implements components.Widget
func (w *Widget) Init() tea.Cmd {
	return w.fetchNotes
}

// Update implements components.Widget
func (w *Widget) Update(msg tea.Msg) (components.Widget, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !w.IsFocused() {
			return w, nil
		}
		switch msg.String() {
		case "up", "k":
			if w.selected > 0 {
				w.selected--
			}
		case "down", "j":
			if w.selected < len(w.notes)-1 {
				w.selected++
			}
		case " ":
			if w.selected >= 0 && w.selected < len(w.notes) {
				return w, w.toggleNote(w.notes[w.selected].ID)
			}
		case "d":
			if w.selected >= 0 && w.selected < len(w.notes) {
				return w, w.deleteNote(w.notes[w.selected].ID)
			}
		case "n":
			return w, w.createNote
		}
	case notesMsg:
		w.notes = msg
		w.loading = false
		w.lastError = nil
	case errorMsg:
		w.lastError = msg
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
	b.WriteString(styles.Title.Render("Notes"))
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

	// Notes
	if len(w.notes) == 0 {
		subtleStyle := lipgloss.NewStyle().Foreground(styles.Subtle)
		b.WriteString(subtleStyle.Render("No notes"))
		b.WriteString("\n\n")
		b.WriteString(subtleStyle.Render("Press 'n' to create a new note"))
	} else {
		for i, note := range w.notes {
			// Note style
			noteStyle := lipgloss.NewStyle()
			if i == w.selected && w.IsFocused() {
				noteStyle = styles.Selected
			}

			// Note status
			status := "[ ]"
			if note.CompletedAt != nil && !note.CompletedAt.IsZero() {
				status = "[✓]"
			}

			// Format note line
			noteLine := fmt.Sprintf("%s %s", status, note.Title)
			if note.Description != "" {
				noteLine += fmt.Sprintf(" - %s", note.Description)
			}
			b.WriteString(noteStyle.Render(noteLine))
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

// Commands
func (w *Widget) fetchNotes() tea.Msg {
	w.loading = true
	notes, err := w.client.ListNotes()
	if err != nil {
		return errorMsg(err)
	}
	return notesMsg(notes)
}

func (w *Widget) toggleNote(id string) tea.Cmd {
	return func() tea.Msg {
		note := w.notes[w.selected]
		input := NoteInput{
			Title:       note.Title,
			Description: note.Description,
		}

		now := time.Now()
		if note.CompletedAt == nil {
			note.CompletedAt = &now
		} else {
			note.CompletedAt = nil
		}

		_, err := w.client.UpdateNote(id, input)
		if err != nil {
			return errorMsg(err)
		}
		return w.fetchNotes()
	}
}

func (w *Widget) deleteNote(id string) tea.Cmd {
	return func() tea.Msg {
		if err := w.client.DeleteNote(id); err != nil {
			return errorMsg(err)
		}
		return w.fetchNotes()
	}
}

func (w *Widget) createNote() tea.Msg {
	input := NoteInput{
		Title: "New Note",
	}
	_, err := w.client.CreateNote(input)
	if err != nil {
		return errorMsg(err)
	}
	return w.fetchNotes()
}
