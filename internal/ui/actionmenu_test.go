package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewActionMenu verifies the constructor creates a valid action menu
func TestNewActionMenu(t *testing.T) {
	menu := NewActionMenu()
	if menu == nil {
		t.Fatal("NewActionMenu() returned nil")
	}
	if menu.Visible() {
		t.Error("new ActionMenu should not be visible")
	}
	if len(menu.Actions()) == 0 {
		t.Error("new ActionMenu should have default actions")
	}
}

// TestActionMenuShow verifies Show makes the menu visible
func TestActionMenuShow(t *testing.T) {
	menu := NewActionMenu()
	item := &ListItem{ID: "test", Title: "Test Item", Description: "Test Description"}

	menu.Show(item)

	if !menu.Visible() {
		t.Error("ActionMenu should be visible after Show()")
	}
	if menu.Item() != item {
		t.Error("ActionMenu.Item() should return the item passed to Show()")
	}
	if menu.Selected() != 0 {
		t.Error("ActionMenu.Selected() should be 0 after Show()")
	}
}

// TestActionMenuHide verifies Hide makes the menu invisible
func TestActionMenuHide(t *testing.T) {
	menu := NewActionMenu()
	item := &ListItem{ID: "test", Title: "Test Item"}
	menu.Show(item)

	menu.Hide()

	if menu.Visible() {
		t.Error("ActionMenu should not be visible after Hide()")
	}
	if menu.Item() != nil {
		t.Error("ActionMenu.Item() should return nil after Hide()")
	}
	if menu.Selected() != 0 {
		t.Error("ActionMenu.Selected() should be 0 after Hide()")
	}
}

// TestActionMenuMoveDown verifies MoveDown navigates correctly
func TestActionMenuMoveDown(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	initial := menu.Selected()
	menu.MoveDown()

	if menu.Selected() != initial+1 {
		t.Errorf("after MoveDown(), Selected() = %d, want %d", menu.Selected(), initial+1)
	}
}

// TestActionMenuMoveUp verifies MoveUp navigates correctly
func TestActionMenuMoveUp(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})
	menu.MoveDown() // Move to second item

	menu.MoveUp()

	if menu.Selected() != 0 {
		t.Errorf("after MoveUp(), Selected() = %d, want 0", menu.Selected())
	}
}

// TestActionMenuMoveDownAtBoundary verifies MoveDown stops at last item
func TestActionMenuMoveDownAtBoundary(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	// Move to last item
	for i := 0; i < len(menu.Actions())+5; i++ {
		menu.MoveDown()
	}

	expected := len(menu.Actions()) - 1
	if menu.Selected() != expected {
		t.Errorf("Selected() at boundary = %d, want %d", menu.Selected(), expected)
	}
}

// TestActionMenuMoveUpAtBoundary verifies MoveUp stops at first item
func TestActionMenuMoveUpAtBoundary(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	// Try to move up from first item
	menu.MoveUp()

	if menu.Selected() != 0 {
		t.Errorf("Selected() at boundary = %d, want 0", menu.Selected())
	}
}

// TestActionMenuMoveDownEmpty verifies MoveDown handles empty actions
func TestActionMenuMoveDownEmpty(t *testing.T) {
	menu := NewActionMenu()
	menu.SetActions(nil)
	menu.Show(&ListItem{ID: "test"})

	// Should not panic
	menu.MoveDown()
	if menu.Selected() != 0 {
		t.Errorf("Selected() on empty menu = %d, want 0", menu.Selected())
	}
}

// TestActionMenuMoveUpEmpty verifies MoveUp handles empty actions
func TestActionMenuMoveUpEmpty(t *testing.T) {
	menu := NewActionMenu()
	menu.SetActions(nil)
	menu.Show(&ListItem{ID: "test"})

	// Should not panic
	menu.MoveUp()
	if menu.Selected() != 0 {
		t.Errorf("Selected() on empty menu = %d, want 0", menu.Selected())
	}
}

// TestActionMenuSelectedAction verifies SelectedAction returns correct action
func TestActionMenuSelectedAction(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	action := menu.SelectedAction()
	if action == nil {
		t.Fatal("SelectedAction() returned nil")
	}
	if action.ID != menu.Actions()[0].ID {
		t.Errorf("SelectedAction().ID = %q, want %q", action.ID, menu.Actions()[0].ID)
	}
}

// TestActionMenuSelectedActionEmpty verifies SelectedAction returns nil for empty menu
func TestActionMenuSelectedActionEmpty(t *testing.T) {
	menu := NewActionMenu()
	menu.SetActions(nil)
	menu.Show(&ListItem{ID: "test"})

	action := menu.SelectedAction()
	if action != nil {
		t.Error("SelectedAction() should return nil for empty menu")
	}
}

// TestActionMenuSetActions verifies SetActions updates actions
func TestActionMenuSetActions(t *testing.T) {
	menu := NewActionMenu()
	newActions := []Action{
		{ID: "custom1", Label: "Custom 1"},
		{ID: "custom2", Label: "Custom 2"},
	}

	menu.SetActions(newActions)

	if len(menu.Actions()) != 2 {
		t.Errorf("Actions() length = %d, want 2", len(menu.Actions()))
	}
	if menu.Actions()[0].ID != "custom1" {
		t.Errorf("Actions()[0].ID = %q, want %q", menu.Actions()[0].ID, "custom1")
	}
}

// TestActionMenuSetActionsClampsSelection verifies SetActions clamps selection
func TestActionMenuSetActionsClampsSelection(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})
	menu.MoveDown()
	menu.MoveDown()

	// Set fewer actions
	menu.SetActions([]Action{{ID: "only"}})

	if menu.Selected() != 0 {
		t.Errorf("Selected() after SetActions with fewer items = %d, want 0", menu.Selected())
	}
}

// TestActionMenuUpdateEscape verifies Escape hides the menu
func TestActionMenuUpdateEscape(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	menu.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if menu.Visible() {
		t.Error("ActionMenu should be hidden after Escape key")
	}
}

// TestActionMenuUpdateArrowKeys verifies arrow keys navigate
func TestActionMenuUpdateArrowKeys(t *testing.T) {
	tests := []struct {
		name     string
		key      tea.KeyType
		expected int
	}{
		{"KeyDown", tea.KeyDown, 1},
		{"KeyUp first item", tea.KeyUp, 0}, // Should stay at 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := NewActionMenu()
			menu.Show(&ListItem{ID: "test"})
			if tt.key == tea.KeyUp {
				// First move down to test up
				menu.MoveDown()
				menu.Update(tea.KeyMsg{Type: tea.KeyUp})
			} else {
				menu.Update(tea.KeyMsg{Type: tt.key})
			}
			if menu.Selected() != tt.expected {
				t.Errorf("Selected() = %d, want %d", menu.Selected(), tt.expected)
			}
		})
	}
}

// TestActionMenuUpdateJKKeys verifies j/k keys navigate
func TestActionMenuUpdateJKKeys(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	// j should move down
	menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if menu.Selected() != 1 {
		t.Errorf("after 'j', Selected() = %d, want 1", menu.Selected())
	}

	// k should move up
	menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if menu.Selected() != 0 {
		t.Errorf("after 'k', Selected() = %d, want 0", menu.Selected())
	}
}

// TestActionMenuUpdateEnter verifies Enter executes selected action
func TestActionMenuUpdateEnter(t *testing.T) {
	menu := NewActionMenu()
	item := &ListItem{ID: "test", Title: "Test"}
	menu.Show(item)

	cmd := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Menu should be hidden after Enter
	if menu.Visible() {
		t.Error("ActionMenu should be hidden after Enter key")
	}

	// Should return a command that produces ActionExecutedMsg
	if cmd == nil {
		t.Fatal("Update(Enter) should return a command")
	}

	// Execute the command
	msg := cmd()
	execMsg, ok := msg.(ActionExecutedMsg)
	if !ok {
		t.Fatalf("command returned %T, want ActionExecutedMsg", msg)
	}
	if execMsg.Action == nil {
		t.Error("ActionExecutedMsg.Action should not be nil")
	}
	if execMsg.Item != item {
		t.Error("ActionExecutedMsg.Item should be the menu item")
	}
}

// TestActionMenuUpdateWhenHidden verifies Update does nothing when hidden
func TestActionMenuUpdateWhenHidden(t *testing.T) {
	menu := NewActionMenu()
	// Menu is hidden by default

	cmd := menu.Update(tea.KeyMsg{Type: tea.KeyDown})

	if cmd != nil {
		t.Error("Update when hidden should return nil")
	}
}

// TestActionMenuView verifies View renders correctly
func TestActionMenuView(t *testing.T) {
	menu := NewActionMenu()
	item := &ListItem{ID: "test", Title: "Test Worktree"}
	menu.Show(item)

	view := menu.View()

	if view == "" {
		t.Error("View() returned empty string when visible")
	}

	// Should contain item title
	if !strings.Contains(view, "Test Worktree") {
		t.Error("View() should contain item title")
	}

	// Should contain action labels
	for _, action := range menu.Actions() {
		if !strings.Contains(view, action.Label) {
			t.Errorf("View() should contain action label %q", action.Label)
		}
	}

	// Should contain selection indicator
	if !strings.Contains(view, "▸") {
		t.Error("View() should contain selection indicator")
	}

	// Should contain help text
	if !strings.Contains(view, "Esc") {
		t.Error("View() should contain Esc in help text")
	}
}

// TestActionMenuViewWhenHidden verifies View returns empty when hidden
func TestActionMenuViewWhenHidden(t *testing.T) {
	menu := NewActionMenu()
	// Menu is hidden by default

	view := menu.View()

	if view != "" {
		t.Errorf("View() when hidden returned %q, want empty string", view)
	}
}

// TestActionMenuViewShowsSelectedDescription verifies description shown for selected
func TestActionMenuViewShowsSelectedDescription(t *testing.T) {
	menu := NewActionMenu()
	menu.Show(&ListItem{ID: "test"})

	view := menu.View()

	// First action's description should be visible
	firstAction := menu.Actions()[0]
	if firstAction.Description != "" && !strings.Contains(view, firstAction.Description) {
		t.Errorf("View() should contain selected action's description %q", firstAction.Description)
	}
}

// TestDefaultWorktreeActions verifies default actions are set
func TestDefaultWorktreeActions(t *testing.T) {
	actions := defaultWorktreeActions()

	if len(actions) == 0 {
		t.Fatal("defaultWorktreeActions() returned empty list")
	}

	// Check that essential actions exist
	hasOpen := false
	hasDelete := false
	for _, a := range actions {
		if a.ID == "open" {
			hasOpen = true
		}
		if a.ID == "delete" {
			hasDelete = true
		}
	}

	if !hasOpen {
		t.Error("defaultWorktreeActions() should include 'open' action")
	}
	if !hasDelete {
		t.Error("defaultWorktreeActions() should include 'delete' action")
	}
}

// TestActionMenuSetSize verifies SetSize sets dimensions
func TestActionMenuSetSize(t *testing.T) {
	menu := NewActionMenu()
	menu.SetSize(100, 50)

	if menu.width != 100 || menu.height != 50 {
		t.Errorf("SetSize(100, 50) resulted in width=%d, height=%d", menu.width, menu.height)
	}
}
