// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewConfirmDialog verifies the constructor.
func TestNewConfirmDialog(t *testing.T) {
	d := NewConfirmDialog()
	if d == nil {
		t.Fatal("Expected non-nil ConfirmDialog")
	}
	if d.Visible() {
		t.Error("Expected new dialog to be hidden")
	}
	if d.confirmLabel != "Confirm" {
		t.Errorf("Expected default confirm label 'Confirm', got '%s'", d.confirmLabel)
	}
	if d.cancelLabel != "Cancel" {
		t.Errorf("Expected default cancel label 'Cancel', got '%s'", d.cancelLabel)
	}
}

// TestConfirmDialogShow verifies Show makes the dialog visible.
func TestConfirmDialogShow(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Test Title", "Test message")

	if !d.Visible() {
		t.Error("Expected dialog to be visible after Show")
	}
	if d.Title() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", d.Title())
	}
	if d.Message() != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", d.Message())
	}
	// Should default to cancel for safety
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 (cancel), got %d", d.Selected())
	}
}

// TestConfirmDialogShowWithData verifies ShowWithData stores data.
func TestConfirmDialogShowWithData(t *testing.T) {
	d := NewConfirmDialog()
	data := &ListItem{ID: "test", Title: "Test Item"}
	d.ShowWithData("Title", "Message", data)

	if !d.Visible() {
		t.Error("Expected dialog to be visible")
	}
	if d.Data() != data {
		t.Error("Expected data to match")
	}
}

// TestConfirmDialogShowDanger verifies ShowDanger sets danger mode.
func TestConfirmDialogShowDanger(t *testing.T) {
	d := NewConfirmDialog()
	d.ShowDanger("Delete?", "This cannot be undone", nil)

	if !d.Visible() {
		t.Error("Expected dialog to be visible")
	}
	if !d.dangerMode {
		t.Error("Expected danger mode to be set")
	}
}

// TestConfirmDialogHide verifies Hide closes the dialog.
func TestConfirmDialogHide(t *testing.T) {
	d := NewConfirmDialog()
	d.ShowDanger("Title", "Message", "data")
	d.SetForceOption(true)
	d.forceSelected = true
	d.Hide()

	if d.Visible() {
		t.Error("Expected dialog to be hidden after Hide")
	}
	if d.dangerMode {
		t.Error("Expected danger mode to be reset")
	}
	if d.forceOption {
		t.Error("Expected force option to be reset")
	}
	if d.forceSelected {
		t.Error("Expected force selected to be reset")
	}
	if d.Data() != nil {
		t.Error("Expected data to be nil after Hide")
	}
}

// TestConfirmDialogMoveLeft verifies MoveLeft navigation.
func TestConfirmDialogMoveLeft(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	// Default is 1 (cancel)
	d.MoveLeft()
	if d.Selected() != 0 {
		t.Errorf("Expected selected 0 after MoveLeft, got %d", d.Selected())
	}
	// Should not go below 0
	d.MoveLeft()
	if d.Selected() != 0 {
		t.Errorf("Expected selected 0 at boundary, got %d", d.Selected())
	}
}

// TestConfirmDialogMoveRight verifies MoveRight navigation.
func TestConfirmDialogMoveRight(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.selected = 0
	d.MoveRight()
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 after MoveRight, got %d", d.Selected())
	}
	// Should not go above 1
	d.MoveRight()
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 at boundary, got %d", d.Selected())
	}
}

// TestConfirmDialogToggleForce verifies ToggleForce toggle.
func TestConfirmDialogToggleForce(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")

	// Toggle without force option enabled should do nothing
	d.ToggleForce()
	if d.ForceSelected() {
		t.Error("Expected force to remain false when option not enabled")
	}

	// Enable force option
	d.SetForceOption(true)
	d.ToggleForce()
	if !d.ForceSelected() {
		t.Error("Expected force to be true after toggle")
	}
	d.ToggleForce()
	if d.ForceSelected() {
		t.Error("Expected force to be false after second toggle")
	}
}

// TestConfirmDialogUpdateEscape verifies Escape key closes dialog.
func TestConfirmDialogUpdateEscape(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if d.Visible() {
		t.Error("Expected dialog to be hidden after Escape")
	}
	if cmd == nil {
		t.Error("Expected command to be returned")
	}

	// Execute the command and check the message
	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if result.Confirmed {
		t.Error("Expected Confirmed false for Escape")
	}
}

// TestConfirmDialogUpdateEnterConfirm verifies Enter key on confirm.
func TestConfirmDialogUpdateEnterConfirm(t *testing.T) {
	d := NewConfirmDialog()
	d.ShowWithData("Title", "Message", "testdata")
	d.selected = 0 // Select confirm

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if d.Visible() {
		t.Error("Expected dialog to be hidden after Enter")
	}
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if !result.Confirmed {
		t.Error("Expected Confirmed true")
	}
	if result.Data != "testdata" {
		t.Error("Expected data to be passed")
	}
}

// TestConfirmDialogUpdateEnterCancel verifies Enter key on cancel.
func TestConfirmDialogUpdateEnterCancel(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.selected = 1 // Select cancel

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if d.Visible() {
		t.Error("Expected dialog to be hidden after Enter")
	}

	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if result.Confirmed {
		t.Error("Expected Confirmed false for cancel")
	}
}

// TestConfirmDialogUpdateArrowKeys verifies arrow key navigation.
func TestConfirmDialogUpdateArrowKeys(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.selected = 1

	d.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if d.Selected() != 0 {
		t.Errorf("Expected selected 0 after Left, got %d", d.Selected())
	}

	d.Update(tea.KeyMsg{Type: tea.KeyRight})
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 after Right, got %d", d.Selected())
	}
}

// TestConfirmDialogUpdateTabKeys verifies Tab key navigation.
func TestConfirmDialogUpdateTabKeys(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.selected = 0

	d.Update(tea.KeyMsg{Type: tea.KeyTab})
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 after Tab, got %d", d.Selected())
	}

	d.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if d.Selected() != 0 {
		t.Errorf("Expected selected 0 after Shift+Tab, got %d", d.Selected())
	}
}

// TestConfirmDialogUpdateHLKeys verifies h/l vim navigation.
func TestConfirmDialogUpdateHLKeys(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.selected = 1

	d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if d.Selected() != 0 {
		t.Errorf("Expected selected 0 after 'h', got %d", d.Selected())
	}

	d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if d.Selected() != 1 {
		t.Errorf("Expected selected 1 after 'l', got %d", d.Selected())
	}
}

// TestConfirmDialogUpdateYKey verifies 'y' quick confirm.
func TestConfirmDialogUpdateYKey(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if d.Visible() {
		t.Error("Expected dialog to be hidden after 'y'")
	}

	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if !result.Confirmed {
		t.Error("Expected Confirmed true for 'y'")
	}
}

// TestConfirmDialogUpdateNKey verifies 'n' quick cancel.
func TestConfirmDialogUpdateNKey(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if d.Visible() {
		t.Error("Expected dialog to be hidden after 'n'")
	}

	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if result.Confirmed {
		t.Error("Expected Confirmed false for 'n'")
	}
}

// TestConfirmDialogUpdateSpaceToggleForce verifies space toggles force.
func TestConfirmDialogUpdateSpaceToggleForce(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.SetForceOption(true)

	d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	if !d.ForceSelected() {
		t.Error("Expected force to be true after space")
	}
}

// TestConfirmDialogUpdateFToggleForce verifies 'f' toggles force.
func TestConfirmDialogUpdateFToggleForce(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Title", "Message")
	d.SetForceOption(true)

	d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	if !d.ForceSelected() {
		t.Error("Expected force to be true after 'f'")
	}
}

// TestConfirmDialogUpdateWhenHidden verifies no action when hidden.
func TestConfirmDialogUpdateWhenHidden(t *testing.T) {
	d := NewConfirmDialog()
	// Dialog is hidden by default

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("Expected nil command when dialog is hidden")
	}
}

// TestConfirmDialogSetSize verifies dimension setting.
func TestConfirmDialogSetSize(t *testing.T) {
	d := NewConfirmDialog()
	d.SetSize(80, 20)

	if d.width != 80 {
		t.Errorf("Expected width 80, got %d", d.width)
	}
	if d.height != 20 {
		t.Errorf("Expected height 20, got %d", d.height)
	}
}

// TestConfirmDialogSetLabels verifies label setting.
func TestConfirmDialogSetLabels(t *testing.T) {
	d := NewConfirmDialog()
	d.SetConfirmLabel("Delete")
	d.SetCancelLabel("Keep")

	if d.confirmLabel != "Delete" {
		t.Errorf("Expected confirm label 'Delete', got '%s'", d.confirmLabel)
	}
	if d.cancelLabel != "Keep" {
		t.Errorf("Expected cancel label 'Keep', got '%s'", d.cancelLabel)
	}
}

// TestConfirmDialogView verifies the view renders correctly.
func TestConfirmDialogView(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Delete Worktree?", "This will remove the worktree directory.")

	view := d.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
	if !strings.Contains(view, "Delete Worktree?") {
		t.Error("Expected view to contain title")
	}
	if !strings.Contains(view, "remove the worktree") {
		t.Error("Expected view to contain message")
	}
	if !strings.Contains(view, "Confirm") {
		t.Error("Expected view to contain confirm button")
	}
	if !strings.Contains(view, "Cancel") {
		t.Error("Expected view to contain cancel button")
	}
}

// TestConfirmDialogViewWhenHidden verifies empty view when hidden.
func TestConfirmDialogViewWhenHidden(t *testing.T) {
	d := NewConfirmDialog()
	view := d.View()
	if view != "" {
		t.Error("Expected empty view when hidden")
	}
}

// TestConfirmDialogViewWithForce verifies force checkbox appears.
func TestConfirmDialogViewWithForce(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Delete?", "Message")
	d.SetForceOption(true)

	view := d.View()
	if !strings.Contains(view, "[ ]") && !strings.Contains(view, "[x]") {
		t.Error("Expected view to contain force checkbox")
	}
	if !strings.Contains(view, "Force removal") {
		t.Error("Expected view to contain force label")
	}
}

// TestConfirmDialogViewForceChecked verifies force checkbox checked state.
func TestConfirmDialogViewForceChecked(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Delete?", "Message")
	d.SetForceOption(true)
	d.forceSelected = true

	view := d.View()
	if !strings.Contains(view, "[x]") {
		t.Error("Expected view to contain checked checkbox")
	}
}

// TestConfirmDialogHasForceOption verifies HasForceOption getter.
func TestConfirmDialogHasForceOption(t *testing.T) {
	d := NewConfirmDialog()
	if d.HasForceOption() {
		t.Error("Expected HasForceOption false by default")
	}

	d.SetForceOption(true)
	if !d.HasForceOption() {
		t.Error("Expected HasForceOption true after SetForceOption")
	}
}

// TestConfirmDialogForceWithEnter verifies force state is passed on confirm.
func TestConfirmDialogForceWithEnter(t *testing.T) {
	d := NewConfirmDialog()
	d.Show("Delete?", "Message")
	d.SetForceOption(true)
	d.forceSelected = true
	d.selected = 0 // Confirm

	cmd := d.Update(tea.KeyMsg{Type: tea.KeyEnter})
	msg := cmd()
	result, ok := msg.(ConfirmDialogResultMsg)
	if !ok {
		t.Fatal("Expected ConfirmDialogResultMsg")
	}
	if !result.Force {
		t.Error("Expected Force true in result")
	}
}
