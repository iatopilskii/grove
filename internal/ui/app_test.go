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
		{TabWorktrees, "main"}, // List shows worktree names
		{TabBranches, "main"},  // Branches tab also shows list
		{TabSettings, "Settings content"},
	}

	for _, tt := range tests {
		t.Run(tt.tab.String(), func(t *testing.T) {
			app := NewApp()
			app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
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

// TestAppHasList verifies App has a list component
func TestAppHasList(t *testing.T) {
	app := NewApp()
	if app.list == nil {
		t.Error("App should have a list component")
	}
}

// TestAppHasDetails verifies App has a details component
func TestAppHasDetails(t *testing.T) {
	app := NewApp()
	if app.details == nil {
		t.Error("App should have a details component")
	}
}

// TestAppViewShowsList verifies View includes list pane
func TestAppViewShowsList(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	view := app.View()

	// List should show items when on Worktrees tab
	if !strings.Contains(view, "▸") {
		t.Error("View() should show list selection indicator")
	}
}

// TestAppViewShowsDetails verifies View includes details pane
func TestAppViewShowsDetails(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	view := app.View()

	// Details should show border (rounded corners)
	hasBorder := strings.Contains(view, "─") || strings.Contains(view, "│") ||
		strings.Contains(view, "╭") || strings.Contains(view, "╰")
	if !hasBorder {
		t.Error("View() should show details pane border")
	}
}

// TestAppListNavigationDown verifies arrow key navigation in list
func TestAppListNavigationDown(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	initial := app.list.Selected()
	app.Update(tea.KeyMsg{Type: tea.KeyDown})
	after := app.list.Selected()

	if after != initial+1 {
		t.Errorf("after KeyDown, list selection = %d, want %d", after, initial+1)
	}
}

// TestAppListNavigationUp verifies arrow key navigation in list
func TestAppListNavigationUp(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Move down first, then up
	app.Update(tea.KeyMsg{Type: tea.KeyDown})
	app.Update(tea.KeyMsg{Type: tea.KeyUp})

	if app.list.Selected() != 0 {
		t.Errorf("after KeyDown then KeyUp, list selection = %d, want 0", app.list.Selected())
	}
}

// TestAppDetailsUpdatesWithSelection verifies details pane updates when list selection changes
func TestAppDetailsUpdatesWithSelection(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Get initial details item
	initial := app.details.Item()
	if initial == nil {
		t.Fatal("details should have an item when list has selection")
	}

	// Move selection
	app.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Details should update
	after := app.details.Item()
	if after == nil {
		t.Fatal("details should still have an item after selection change")
	}
	if initial.ID == after.ID {
		t.Error("details item should change when list selection changes")
	}
}

// TestAppViewRendersBothPanes verifies both panes are rendered in view
func TestAppViewRendersBothPanes(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	view := app.View()

	// Should have multiple lines showing list items
	lines := strings.Split(view, "\n")
	if len(lines) < 5 {
		t.Errorf("View() should have multiple lines for two-pane layout, got %d lines", len(lines))
	}
}

// TestAppListNavigationPageDown verifies PageDown key navigation in list
func TestAppListNavigationPageDown(t *testing.T) {
	app := NewApp()
	// Set up with many items for page navigation testing
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item", Description: "Desc"}
	}
	app.list.SetItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 10}) // height determines page size

	initial := app.list.Selected()
	app.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	after := app.list.Selected()

	if after <= initial {
		t.Errorf("after KeyPgDown, list selection = %d, should be greater than %d", after, initial)
	}
}

// TestAppListNavigationPageUp verifies PageUp key navigation in list
func TestAppListNavigationPageUp(t *testing.T) {
	app := NewApp()
	// Set up with many items for page navigation testing
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item", Description: "Desc"}
	}
	app.list.SetItems(items)
	app.list.SetSelected(15)                              // Start from further down
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 10}) // height determines page size

	initial := app.list.Selected()
	app.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	after := app.list.Selected()

	if after >= initial {
		t.Errorf("after KeyPgUp, list selection = %d, should be less than %d", after, initial)
	}
}

// TestAppPageNavigationUpdatesDetails verifies details pane updates after page navigation
func TestAppPageNavigationUpdatesDetails(t *testing.T) {
	app := NewApp()
	// Set up with many items
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{
			ID:          string(rune('a' + i)),
			Title:       "Item " + string(rune('A'+i)),
			Description: "Description " + string(rune('A'+i)),
		}
	}
	app.list.SetItems(items)
	app.details.SetItem(app.list.SelectedItem())
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 5})

	// Get initial details item
	initialItem := app.details.Item()
	if initialItem == nil {
		t.Fatal("details should have an item")
	}

	// PageDown should update details
	app.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	afterItem := app.details.Item()
	if afterItem == nil {
		t.Fatal("details should still have an item after page navigation")
	}
	if initialItem.ID == afterItem.ID {
		t.Error("details item should change after PageDown navigation")
	}
}

// TestAppPageNavigationOnlyOnListTabs verifies PageUp/PageDown only work on Worktrees/Branches tabs
func TestAppPageNavigationOnlyOnListTabs(t *testing.T) {
	app := NewApp()
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	app.list.SetItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 5})

	// Switch to Settings tab
	app.tabs.SetActive(TabSettings)

	initial := app.list.Selected()
	app.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	after := app.list.Selected()

	// Selection should not change on Settings tab
	if after != initial {
		t.Errorf("PageDown on Settings tab should not change selection, got %d want %d", after, initial)
	}
}
