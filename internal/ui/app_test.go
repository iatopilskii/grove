package ui

import (
	"strings"
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
	if app.tabs == nil {
		t.Error("NewApp() did not initialize tabs")
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

// TestAppViewContainsTabs verifies that View shows all tab names
func TestAppViewContainsTabs(t *testing.T) {
	app := NewApp()
	view := app.View()

	tabNames := []string{"Worktrees", "Branches", "Settings"}
	for _, name := range tabNames {
		if !strings.Contains(view, name) {
			t.Errorf("View() does not contain tab name %q", name)
		}
	}
}

// TestAppViewShowsActiveTabContent verifies content updates based on active tab
func TestAppViewShowsActiveTabContent(t *testing.T) {
	tests := []struct {
		tab            Tab
		expectedPhrase string
	}{
		{TabWorktrees, "Worktrees content"},
		{TabBranches, "Branches content"},
		{TabSettings, "Settings content"},
	}

	for _, tt := range tests {
		t.Run(tt.tab.String(), func(t *testing.T) {
			app := NewApp()
			app.tabs.SetActive(tt.tab)
			view := app.View()
			if !strings.Contains(view, tt.expectedPhrase) {
				t.Errorf("View() with %v tab does not contain %q", tt.tab, tt.expectedPhrase)
			}
		})
	}
}

// TestAppUpdateTabKey verifies Tab key switches tabs
func TestAppUpdateTabKey(t *testing.T) {
	app := NewApp()

	// Initial tab should be Worktrees
	if app.tabs.Active() != TabWorktrees {
		t.Fatalf("expected initial tab to be TabWorktrees")
	}

	// Tab key should switch to next tab
	app.Update(tea.KeyMsg{Type: tea.KeyTab})
	if app.tabs.Active() != TabBranches {
		t.Errorf("after Tab key, active tab = %v, want TabBranches", app.tabs.Active())
	}
}

// TestAppUpdateShiftTabKey verifies Shift+Tab switches tabs backwards
func TestAppUpdateShiftTabKey(t *testing.T) {
	app := NewApp()

	// Shift+Tab should wrap to Settings
	app.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if app.tabs.Active() != TabSettings {
		t.Errorf("after Shift+Tab key, active tab = %v, want TabSettings", app.tabs.Active())
	}
}

// TestAppUpdateWindowSize verifies window size updates dimensions
func TestAppUpdateWindowSize(t *testing.T) {
	app := NewApp()

	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if app.width != 120 || app.height != 40 {
		t.Errorf("after WindowSizeMsg, got width=%d height=%d, want 120x40", app.width, app.height)
	}
}

// TestAppTabCycling verifies full cycle through tabs
func TestAppTabCycling(t *testing.T) {
	app := NewApp()

	expectedOrder := []Tab{TabBranches, TabSettings, TabWorktrees}
	for _, expected := range expectedOrder {
		app.Update(tea.KeyMsg{Type: tea.KeyTab})
		if app.tabs.Active() != expected {
			t.Errorf("tab cycling: got %v, want %v", app.tabs.Active(), expected)
		}
	}
}
