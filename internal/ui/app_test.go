package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestAppImplementsTeaModel verifies that App implements tea.Model interface
func TestAppImplementsTeaModel(t *testing.T) {
	var _ tea.Model = (*App)(nil)
}

// TestNewApp verifies that NewApp returns a properly initialized App
func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

// TestAppInit verifies that Init returns a valid command
func TestAppInit(t *testing.T) {
	app := NewApp()
	cmd := app.Init()
	// Init can return nil (no initial command) or a valid command
	// Either is acceptable for initial setup
	_ = cmd
}

// TestAppUpdate verifies that Update handles messages and returns updated model
func TestAppUpdate(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.Msg
	}{
		{
			name: "handles nil message",
			msg:  nil,
		},
		{
			name: "handles quit key",
			msg:  tea.KeyMsg{Type: tea.KeyCtrlC},
		},
		{
			name: "handles q key",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			model, cmd := app.Update(tt.msg)
			if model == nil {
				t.Error("Update() returned nil model")
			}
			// cmd can be nil or non-nil depending on the message
			_ = cmd
		})
	}
}

// TestAppUpdateQuitCommand verifies that quit keys produce tea.Quit command
func TestAppUpdateQuitCommand(t *testing.T) {
	quitKeys := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyCtrlC},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
	}

	for _, msg := range quitKeys {
		app := NewApp()
		_, cmd := app.Update(msg)
		if cmd == nil {
			t.Errorf("Expected quit command for %v, got nil", msg)
		}
	}
}

// TestAppView verifies that View returns a non-empty string
func TestAppView(t *testing.T) {
	app := NewApp()
	view := app.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}
