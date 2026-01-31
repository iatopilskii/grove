package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ilatopilskij/gwt/internal/git"
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
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree at /path/to/repo"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch worktree at /path/to/repo-feature-1"},
		{ID: "bugfix-2", Title: "bugfix-2", Description: "Bugfix branch worktree at /path/to/repo-bugfix-2"},
	}
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
			app := NewAppWithItems(sampleItems)
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
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree at /path/to/repo"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
		{ID: "bugfix-2", Title: "bugfix-2", Description: "Bugfix branch"},
	}
	app := NewAppWithItems(sampleItems)
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
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
	}
	app := NewAppWithItems(sampleItems)
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
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
	}
	app := NewAppWithItems(sampleItems)
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

// TestAppMouseClickListItem verifies clicking on list item selects it
func TestAppMouseClickListItem(t *testing.T) {
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
		{ID: "bugfix-2", Title: "bugfix-2", Description: "Bugfix branch"},
	}
	app := NewAppWithItems(sampleItems)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Click on second item in the list
	// List starts after tabs (2 lines for tabs + border)
	// Each list item is 1 line
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      10, // Inside list pane width
		Y:      4,  // Tabs (2) + first item = 3, second item = 4
	})

	// Selection should have changed
	if app.list.Selected() != 1 {
		t.Errorf("after mouse click on second item, list.Selected() = %d, want 1", app.list.Selected())
	}
}

// TestAppMouseClickTab verifies clicking on tab switches to it
func TestAppMouseClickTab(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Click on Settings tab (rightmost)
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      32, // In Settings tab area
		Y:      0,  // Tab row
	})

	if app.tabs.Active() != TabSettings {
		t.Errorf("after click on Settings tab, Active() = %v, want TabSettings", app.tabs.Active())
	}
}

// TestAppMouseWheelInList verifies mouse wheel scrolls list
func TestAppMouseWheelInList(t *testing.T) {
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
	}
	app := NewAppWithItems(sampleItems)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	initial := app.list.Selected()

	app.Update(tea.MouseMsg{
		Type:   tea.MouseWheelDown,
		Button: tea.MouseButtonWheelDown,
		X:      10,
		Y:      5, // Inside list area
	})

	if app.list.Selected() != initial+1 {
		t.Errorf("after mouse wheel down, list.Selected() = %d, want %d", app.list.Selected(), initial+1)
	}
}

// TestAppMouseOutsideBounds verifies clicking outside bounds doesn't crash
func TestAppMouseOutsideBounds(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Click outside terminal bounds - should not panic
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      500,
		Y:      500,
	})

	// Click with negative coordinates - should not panic
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      -10,
		Y:      -10,
	})
}

// TestAppMouseDetailsPane verifies clicking on details pane doesn't crash
func TestAppMouseDetailsPane(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	initial := app.list.Selected()

	// Click in details pane area (right side)
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      80, // In details pane area
		Y:      5,
	})

	// Selection should not change
	if app.list.Selected() != initial {
		t.Errorf("click in details pane should not change list selection, got %d want %d", app.list.Selected(), initial)
	}
}

// TestAppMouseDetailsUpdateAfterClick verifies details updates after mouse selection
func TestAppMouseDetailsUpdateAfterClick(t *testing.T) {
	sampleItems := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch"},
	}
	app := NewAppWithItems(sampleItems)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	initialItem := app.details.Item()
	if initialItem == nil {
		t.Fatal("details should have an item")
	}

	// Click on second item
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      10,
		Y:      4, // Second item
	})

	afterItem := app.details.Item()
	if afterItem == nil {
		t.Fatal("details should still have an item")
	}
	if initialItem.ID == afterItem.ID {
		t.Error("details should update after mouse click on different item")
	}
}

// TestAppMouseOnSettingsTab verifies mouse in list doesn't affect selection on Settings tab
func TestAppMouseOnSettingsTab(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app.tabs.SetActive(TabSettings)

	initial := app.list.Selected()

	// Click in what would be list area
	app.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      10,
		Y:      5,
	})

	// Selection should not change on Settings tab
	if app.list.Selected() != initial {
		t.Errorf("mouse click should not change list selection on Settings tab, got %d want %d", app.list.Selected(), initial)
	}
}

// TestAppHasActionMenu verifies App has an action menu component
func TestAppHasActionMenu(t *testing.T) {
	app := NewApp()
	if app.actionMenu == nil {
		t.Error("App should have an action menu component")
	}
}

// TestAppHasFeedback verifies App has a feedback component
func TestAppHasFeedback(t *testing.T) {
	app := NewApp()
	if app.feedback == nil {
		t.Error("App should have a feedback component")
	}
}

// TestAppEnterOpensActionMenu verifies Enter key opens action menu on Worktrees tab
func TestAppEnterOpensActionMenu(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Press Enter
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !app.actionMenu.Visible() {
		t.Error("Enter key should open action menu")
	}
	if app.actionMenu.Item() == nil {
		t.Error("Action menu should have the selected item")
	}
}

// TestAppEnterOpensActionMenuOnBranchesTab verifies Enter works on Branches tab
func TestAppEnterOpensActionMenuOnBranchesTab(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app.tabs.SetActive(TabBranches)

	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !app.actionMenu.Visible() {
		t.Error("Enter key should open action menu on Branches tab")
	}
}

// TestAppEnterDoesNotOpenOnSettingsTab verifies Enter doesn't open menu on Settings
func TestAppEnterDoesNotOpenOnSettingsTab(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app.tabs.SetActive(TabSettings)

	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if app.actionMenu.Visible() {
		t.Error("Enter key should not open action menu on Settings tab")
	}
}

// TestAppEscapeClosesActionMenu verifies Escape closes action menu
func TestAppEscapeClosesActionMenu(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !app.actionMenu.Visible() {
		t.Fatal("Action menu should be visible after Enter")
	}

	// Press Escape
	app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if app.actionMenu.Visible() {
		t.Error("Escape key should close action menu")
	}
}

// TestAppActionMenuRoutesKeys verifies keys go to action menu when visible
func TestAppActionMenuRoutesKeys(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	initialSelection := app.actionMenu.Selected()

	// Press Down in action menu
	app.Update(tea.KeyMsg{Type: tea.KeyDown})

	if app.actionMenu.Selected() != initialSelection+1 {
		t.Errorf("Down key should navigate action menu, selection = %d, want %d", app.actionMenu.Selected(), initialSelection+1)
	}
}

// TestAppActionMenuEnterExecutesAction verifies Enter in menu executes action
func TestAppActionMenuEnterExecutesAction(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Press Enter to execute action
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Menu should close
	if app.actionMenu.Visible() {
		t.Error("Action menu should close after executing action")
	}

	// Should return a command
	if cmd == nil {
		t.Error("Executing action should return a command")
	}
}

// TestAppActionExecutedShowsFeedback verifies action execution shows feedback
func TestAppActionExecutedShowsFeedback(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Send an ActionExecutedMsg directly
	action := &Action{ID: "open", Label: "Open"}
	item := &ListItem{ID: "test", Title: "Test Worktree"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	if !app.feedback.Visible() {
		t.Error("Feedback should be visible after action executed")
	}
}

// TestAppCtrlCQuitsEvenWithMenuOpen verifies Ctrl+C quits even with menu open
func TestAppCtrlCQuitsEvenWithMenuOpen(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Press Ctrl+C
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if !app.quitting {
		t.Error("Ctrl+C should set quitting to true even with menu open")
	}
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}

// TestAppViewShowsActionMenu verifies View includes action menu when visible
func TestAppViewShowsActionMenu(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	view := app.View()

	// Should contain action menu elements
	if !strings.Contains(view, "Actions") {
		t.Error("View should show action menu title")
	}
	if !strings.Contains(view, "Open") {
		t.Error("View should show Open action")
	}
	if !strings.Contains(view, "Esc") {
		t.Error("View should show Esc hint in action menu")
	}
}

// TestAppViewShowsFeedback verifies View includes feedback when visible
func TestAppViewShowsFeedback(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show feedback
	app.feedback.ShowSuccess("Test message")

	view := app.View()

	if !strings.Contains(view, "Test message") {
		t.Error("View should show feedback message")
	}
}

// TestAppViewHelpIncludesEnter verifies help text includes Enter key
func TestAppViewHelpIncludesEnter(t *testing.T) {
	app := NewApp()
	view := app.View()

	if !strings.Contains(view, "Enter") {
		t.Error("Help text should include Enter key hint")
	}
}

// TestAppClearFeedbackMsg verifies ClearFeedbackMsg clears feedback
func TestAppClearFeedbackMsg(t *testing.T) {
	app := NewApp()
	app.feedback.ShowSuccess("Test")

	app.Update(ClearFeedbackMsg{})

	if app.feedback.Visible() {
		t.Error("Feedback should be cleared after ClearFeedbackMsg")
	}
}

// TestAppActionMenuJKNavigation verifies j/k work in action menu
func TestAppActionMenuJKNavigation(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open action menu
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Press j to move down
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if app.actionMenu.Selected() != 1 {
		t.Errorf("'j' should move action menu selection down, got %d", app.actionMenu.Selected())
	}

	// Press k to move up
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if app.actionMenu.Selected() != 0 {
		t.Errorf("'k' should move action menu selection up, got %d", app.actionMenu.Selected())
	}
}

// TestAppEnterWithEmptyList verifies Enter does nothing with empty list
func TestAppEnterWithEmptyList(t *testing.T) {
	app := NewApp()
	app.list.SetItems(nil) // Empty list
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	app.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if app.actionMenu.Visible() {
		t.Error("Enter should not open action menu when list is empty")
	}
}

// TestAppHandleActionExecutedOpen verifies open action shows success
func TestAppHandleActionExecutedOpen(t *testing.T) {
	app := NewApp()

	action := &Action{ID: "open", Label: "Open"}
	item := &ListItem{ID: "test", Title: "Test"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	if app.feedback.Type() != FeedbackSuccess {
		t.Errorf("Open action should show success feedback, got %v", app.feedback.Type())
	}
}

// TestAppHandleActionExecutedDelete verifies delete action shows confirmation dialog
func TestAppHandleActionExecutedDelete(t *testing.T) {
	app := NewApp()

	action := &Action{ID: "delete", Label: "Delete"}
	item := &ListItem{ID: "/path/to/worktree", Title: "Test"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	if !app.confirmDialog.Visible() {
		t.Error("Delete action should show confirmation dialog")
	}
	if app.confirmDialog.Title() != "Delete Worktree?" {
		t.Errorf("Expected title 'Delete Worktree?', got '%s'", app.confirmDialog.Title())
	}
	if !app.confirmDialog.HasForceOption() {
		t.Error("Delete confirmation should have force option")
	}
}

// TestAppHandleActionExecutedUnknown verifies unknown action shows error
func TestAppHandleActionExecutedUnknown(t *testing.T) {
	app := NewApp()

	action := &Action{ID: "unknown", Label: "Unknown"}
	item := &ListItem{ID: "test", Title: "Test"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	if app.feedback.Type() != FeedbackError {
		t.Errorf("Unknown action should show error feedback, got %v", app.feedback.Type())
	}
}

// TestAppHandleActionExecutedNilAction verifies nil action is handled
func TestAppHandleActionExecutedNilAction(t *testing.T) {
	app := NewApp()

	// Should not panic
	app.Update(ActionExecutedMsg{Action: nil, Item: nil})

	if app.feedback.Visible() {
		t.Error("Nil action should not show feedback")
	}
}

// TestNewAppWithItems verifies NewAppWithItems creates app with custom items
func TestNewAppWithItems(t *testing.T) {
	items := []ListItem{
		{ID: "test1", Title: "Test 1", Description: "Desc 1"},
		{ID: "test2", Title: "Test 2", Description: "Desc 2"},
	}
	app := NewAppWithItems(items)

	if app == nil {
		t.Fatal("NewAppWithItems returned nil")
	}
	if len(app.list.Items()) != 2 {
		t.Errorf("Expected 2 items, got %d", len(app.list.Items()))
	}
	if app.list.Items()[0].ID != "test1" {
		t.Errorf("Expected first item ID 'test1', got '%s'", app.list.Items()[0].ID)
	}
}

// TestNewAppWithItemsInitializesDetails verifies details is set for first item
func TestNewAppWithItemsInitializesDetails(t *testing.T) {
	items := []ListItem{
		{ID: "first", Title: "First", Description: "First Description"},
		{ID: "second", Title: "Second", Description: "Second Description"},
	}
	app := NewAppWithItems(items)

	item := app.details.Item()
	if item == nil {
		t.Fatal("Details should have an item")
	}
	if item.ID != "first" {
		t.Errorf("Expected first item in details, got %s", item.ID)
	}
}

// TestNewAppWithEmptyItems verifies app works with empty items
func TestNewAppWithEmptyItems(t *testing.T) {
	app := NewAppWithItems(nil)

	if app == nil {
		t.Fatal("NewAppWithItems returned nil")
	}
	if len(app.list.Items()) != 0 {
		t.Errorf("Expected 0 items, got %d", len(app.list.Items()))
	}
}

// TestAppWorktreesGetter verifies Worktrees() returns worktrees
func TestAppWorktreesGetter(t *testing.T) {
	app := NewApp()
	// In a git repo, Worktrees() should return some worktrees
	// The exact count depends on the test environment
	_ = app.Worktrees() // Just ensure it doesn't panic
}

// TestAppGitErrorGetter verifies GitError() returns error state
func TestAppGitErrorGetter(t *testing.T) {
	app := NewApp()
	// In a git repo, GitError() should be nil
	// We just test that the method works
	_ = app.GitError() // Just ensure it doesn't panic
}

// TestAppIsInGitRepo verifies IsInGitRepo() works
func TestAppIsInGitRepo(t *testing.T) {
	app := NewApp()
	// Since tests run in a git repo, this should return true
	if !app.IsInGitRepo() {
		t.Skip("Test must be run in a git repository")
	}
}

// TestAppRefreshWorktrees verifies RefreshWorktrees works
func TestAppRefreshWorktrees(t *testing.T) {
	app := NewApp()
	initialCount := len(app.Worktrees())

	app.RefreshWorktrees()

	// Count should be same after refresh (no worktrees were added/removed)
	if len(app.Worktrees()) != initialCount {
		t.Errorf("Worktree count changed after refresh: %d -> %d", initialCount, len(app.Worktrees()))
	}
}

// TestAppViewShowsGitError verifies View shows error for non-git directory
func TestAppViewShowsGitError(t *testing.T) {
	app := NewAppWithItems(nil)
	// Simulate a git error using actual NotGitRepoError
	app.gitError = &git.NotGitRepoError{Path: "/tmp/test"}
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	view := app.View()

	if !strings.Contains(view, "Not a Git Repository") {
		t.Error("View should show 'Not a Git Repository' error message")
	}
}

// TestAppViewShowsWorktreeList verifies View shows worktree list in git repo
func TestAppViewShowsWorktreeList(t *testing.T) {
	app := NewApp()
	if !app.IsInGitRepo() {
		t.Skip("Test must be run in a git repository")
	}

	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	view := app.View()

	// Should show the selection indicator and not show git error
	if !strings.Contains(view, "▸") {
		t.Error("View should show list selection indicator")
	}
	if strings.Contains(view, "Not a Git Repository") {
		t.Error("View should not show git error in a git repository")
	}
}

// TestAppHasCreateForm verifies App has createForm component
func TestAppHasCreateForm(t *testing.T) {
	app := NewApp()
	if app.CreateForm() == nil {
		t.Error("App should have createForm component")
	}
}

// TestAppNKeyOpensCreateForm verifies 'n' key opens create form on Worktrees tab
func TestAppNKeyOpensCreateForm(t *testing.T) {
	app := NewApp()
	if !app.IsInGitRepo() {
		t.Skip("Test must be run in a git repository")
	}
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app.tabs.SetActive(TabWorktrees)

	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if !app.createForm.Visible() {
		t.Error("'n' key should open create form on Worktrees tab")
	}
}

// TestAppNKeyDoesNotOpenOnNonWorktreesTabs verifies 'n' doesn't open form on other tabs
func TestAppNKeyDoesNotOpenOnNonWorktreesTabs(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app.tabs.SetActive(TabBranches)

	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if app.createForm.Visible() {
		t.Error("'n' key should not open create form on Branches tab")
	}

	app.tabs.SetActive(TabSettings)
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if app.createForm.Visible() {
		t.Error("'n' key should not open create form on Settings tab")
	}
}

// TestAppNKeyDoesNotOpenOnGitError verifies 'n' doesn't open form when git error
func TestAppNKeyDoesNotOpenOnGitError(t *testing.T) {
	app := NewAppWithItems(nil)
	app.gitError = &git.NotGitRepoError{Path: "/tmp/test"}
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if app.createForm.Visible() {
		t.Error("'n' key should not open create form when not in git repo")
	}
}

// TestAppCreateFormRoutesKeys verifies keys go to create form when visible
func TestAppCreateFormRoutesKeys(t *testing.T) {
	app := NewApp()
	if !app.IsInGitRepo() {
		t.Skip("Test must be run in a git repository")
	}
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open create form
	app.createForm.Show()

	// Type in the form
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a', 'b', 'c'}})

	if app.createForm.Branch() != "abc" {
		t.Errorf("Keys should be routed to create form, branch = '%s'", app.createForm.Branch())
	}
}

// TestAppCreateFormEscapeCloses verifies Escape closes create form
func TestAppCreateFormEscapeCloses(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open create form
	app.createForm.Show()

	// Press Escape
	app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if app.createForm.Visible() {
		t.Error("Escape should close create form")
	}
}

// TestAppCtrlCQuitsEvenWithFormOpen verifies Ctrl+C quits even with form open
func TestAppCtrlCQuitsEvenWithFormOpen(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open create form
	app.createForm.Show()

	// Press Ctrl+C
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if !app.quitting {
		t.Error("Ctrl+C should set quitting to true even with form open")
	}
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}

// TestAppViewShowsCreateForm verifies View includes create form when visible
func TestAppViewShowsCreateForm(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Open create form
	app.createForm.Show()

	view := app.View()

	if !strings.Contains(view, "Create New Worktree") {
		t.Error("View should show create form title")
	}
	if !strings.Contains(view, "Branch name:") {
		t.Error("View should show branch field")
	}
}

// TestAppViewHelpIncludesNewKey verifies help text includes 'n' key
func TestAppViewHelpIncludesNewKey(t *testing.T) {
	app := NewApp()
	view := app.View()

	if !strings.Contains(view, "n: new worktree") {
		t.Error("Help text should include 'n: new worktree' hint")
	}
}

// TestAppCreateFormCancelledMsg verifies cancel message is handled
func TestAppCreateFormCancelledMsg(t *testing.T) {
	app := NewApp()
	app.createForm.Show()

	// Should not panic
	app.Update(CreateFormCancelledMsg{})

	// Form should be hidden (handled in the form itself)
	// Just verify the message doesn't cause issues
}

// TestAppCreateFormSubmittedSuccess verifies successful form submission
func TestAppCreateFormSubmittedSuccess(t *testing.T) {
	app := NewApp()
	if !app.IsInGitRepo() {
		t.Skip("Test must be run in a git repository")
	}

	// Note: We can't easily test actual worktree creation without modifying the git repo
	// So we'll just verify the handler doesn't panic and shows appropriate feedback

	// Send a form submission (this will fail due to invalid path, but tests the handler)
	app.Update(CreateFormSubmittedMsg{
		Result: CreateFormResult{
			Branch:       "test-branch",
			Path:         "/nonexistent/path",
			CreateBranch: true,
		},
	})

	if !app.feedback.Visible() {
		t.Error("Form submission should show feedback")
	}
}

// TestAppCreateFormTabNavigation verifies Tab key in form
func TestAppCreateFormTabNavigation(t *testing.T) {
	app := NewApp()
	app.createForm.Show()

	// Press Tab
	app.Update(tea.KeyMsg{Type: tea.KeyTab})

	if app.createForm.Focused() != FieldPath {
		t.Error("Tab should move focus to path field")
	}
}

// TestAppHasConfirmDialog verifies App has confirmDialog component
func TestAppHasConfirmDialog(t *testing.T) {
	app := NewApp()
	if app.ConfirmDialog() == nil {
		t.Error("App should have confirmDialog component")
	}
}

// TestAppConfirmDialogRoutesKeys verifies keys go to confirm dialog when visible
func TestAppConfirmDialogRoutesKeys(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show confirm dialog
	app.confirmDialog.Show("Test", "Message")

	// Press Left to move to confirm
	app.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if app.confirmDialog.Selected() != 0 {
		t.Error("Keys should be routed to confirm dialog")
	}
}

// TestAppConfirmDialogEscapeCloses verifies Escape closes confirm dialog
func TestAppConfirmDialogEscapeCloses(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show confirm dialog
	app.confirmDialog.Show("Test", "Message")

	// Press Escape
	app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if app.confirmDialog.Visible() {
		t.Error("Escape should close confirm dialog")
	}
}

// TestAppCtrlCQuitsEvenWithConfirmDialogOpen verifies Ctrl+C quits even with dialog open
func TestAppCtrlCQuitsEvenWithConfirmDialogOpen(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show confirm dialog
	app.confirmDialog.Show("Test", "Message")

	// Press Ctrl+C
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if !app.quitting {
		t.Error("Ctrl+C should set quitting to true even with dialog open")
	}
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}

// TestAppViewShowsConfirmDialog verifies View includes confirm dialog when visible
func TestAppViewShowsConfirmDialog(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show confirm dialog
	app.confirmDialog.Show("Delete Worktree?", "This will remove the worktree.")

	view := app.View()

	if !strings.Contains(view, "Delete Worktree?") {
		t.Error("View should show confirm dialog title")
	}
	if !strings.Contains(view, "remove the worktree") {
		t.Error("View should show confirm dialog message")
	}
}

// TestAppConfirmDialogResultMsgCancelled verifies cancelled confirmation
func TestAppConfirmDialogResultMsgCancelled(t *testing.T) {
	app := NewApp()

	// Should not panic and should not show feedback
	app.Update(ConfirmDialogResultMsg{Confirmed: false})

	if app.feedback.Visible() {
		t.Error("Cancelled confirmation should not show feedback")
	}
}

// TestAppConfirmDialogResultMsgConfirmedNoData verifies confirmed without data
func TestAppConfirmDialogResultMsgConfirmedNoData(t *testing.T) {
	app := NewApp()

	// Should not panic
	app.Update(ConfirmDialogResultMsg{Confirmed: true, Data: nil})

	// Nothing happens without valid data
}

// TestAppDeleteConfirmationFlow verifies the complete delete confirmation flow
func TestAppDeleteConfirmationFlow(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Trigger delete action
	action := &Action{ID: "delete", Label: "Delete"}
	item := &ListItem{ID: "/path/to/worktree", Title: "test-worktree"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	// Confirm dialog should be visible
	if !app.confirmDialog.Visible() {
		t.Fatal("Confirm dialog should be visible after delete action")
	}

	// Select confirm button (move left from cancel which is default)
	app.Update(tea.KeyMsg{Type: tea.KeyLeft})

	// Verify the data is stored
	if app.confirmDialog.Data() == nil {
		t.Error("Confirm dialog should have stored the item data")
	}
}

// TestAppDeleteWithForceOption verifies force option in delete
func TestAppDeleteWithForceOption(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Trigger delete action
	action := &Action{ID: "delete", Label: "Delete"}
	item := &ListItem{ID: "/path/to/worktree", Title: "test-worktree"}
	app.Update(ActionExecutedMsg{Action: action, Item: item})

	// Toggle force option
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

	if !app.confirmDialog.ForceSelected() {
		t.Error("'f' should toggle force option")
	}
}

// TestAppConfirmDialogQuickAnswer verifies quick y/n answers
func TestAppConfirmDialogQuickAnswer(t *testing.T) {
	app := NewApp()
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Show confirm dialog
	app.confirmDialog.Show("Test", "Message")

	// Press 'n' to cancel
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if app.confirmDialog.Visible() {
		t.Error("'n' should close confirm dialog")
	}
}

// TestAppPKeyTriggersPrune verifies 'p' key opens prune confirmation on Worktrees tab
func TestAppPKeyTriggersPrune(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Press 'p' to trigger prune
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	// Should show confirmation dialog for prune
	if !app.confirmDialog.Visible() {
		t.Error("'p' should show prune confirmation dialog on Worktrees tab")
	}
}

// TestAppPKeyDoesNotTriggerOnSettingsTab verifies 'p' doesn't work on Settings tab
func TestAppPKeyDoesNotTriggerOnSettingsTab(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Switch to Settings tab
	app.tabs.SetActive(TabSettings)

	// Press 'p' - should not trigger prune
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	if app.confirmDialog.Visible() {
		t.Error("'p' should not work on Settings tab")
	}
}

// TestAppPruneConfirmationFlow verifies the prune confirmation flow
func TestAppPruneConfirmationFlow(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Press 'p' to trigger prune
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	if !app.confirmDialog.Visible() {
		t.Fatal("Expected prune confirmation dialog to be visible")
	}

	// Check the dialog title
	view := app.confirmDialog.View()
	if !strings.Contains(view, "Prune") {
		t.Error("Confirmation dialog should mention 'Prune'")
	}
}

// TestAppPruneCancellation verifies prune can be cancelled
func TestAppPruneCancellation(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Press 'p' to trigger prune
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	// Press Escape to cancel
	app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if app.confirmDialog.Visible() {
		t.Error("Escape should close prune confirmation dialog")
	}
}

// TestAppViewHelpIncludesPrune verifies help text includes prune shortcut
func TestAppViewHelpIncludesPrune(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	view := app.View()
	if !strings.Contains(view, "p:") || !strings.Contains(view, "prune") {
		t.Error("Help text should include 'p: prune' hint")
	}
}

// TestAppPruneResultMsg verifies handling of prune result message
func TestAppPruneResultMsg(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Worktree 1", Description: "Description 1"},
	}
	app := NewAppWithItems(items)
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Send a prune result message (confirmed)
	app.Update(ConfirmDialogResultMsg{
		Confirmed: true,
		Data:      "prune",
	})

	// Should show feedback (success or error depending on git state)
	// Since we're not in a real git repo, it will likely show an error
	// but the message handling should work
}

// TestAppPKeyDoesNotTriggerWhenGitError verifies 'p' doesn't work when not in git repo
func TestAppPKeyDoesNotTriggerWhenGitError(t *testing.T) {
	app := NewApp() // Will have git error in non-git directory
	app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Simulate not being in a git repo by setting git error
	app.gitError = &git.NotGitRepoError{Path: "/tmp"}

	// Press 'p' - should not trigger prune
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	if app.confirmDialog.Visible() {
		t.Error("'p' should not work when there is a git error")
	}
}
