// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewCreateForm verifies the constructor.
func TestNewCreateForm(t *testing.T) {
	form := NewCreateForm()
	if form == nil {
		t.Fatal("NewCreateForm returned nil")
	}
	if form.Visible() {
		t.Error("New form should not be visible")
	}
	if !form.CreateBranchEnabled() {
		t.Error("New form should default to creating new branch")
	}
}

// TestCreateFormShow verifies Show makes form visible and resets fields.
func TestCreateFormShow(t *testing.T) {
	form := NewCreateForm()
	form.branch = "old-value"
	form.path = "old-path"
	form.createBranch = false
	form.errorMessage = "old error"

	form.Show()

	if !form.Visible() {
		t.Error("Form should be visible after Show")
	}
	if form.Branch() != "" {
		t.Errorf("Branch should be reset, got '%s'", form.Branch())
	}
	if form.Path() != "" {
		t.Errorf("Path should be reset, got '%s'", form.Path())
	}
	if !form.CreateBranchEnabled() {
		t.Error("CreateBranch should be reset to true")
	}
	if form.Focused() != FieldBranch {
		t.Error("Focus should be on branch field")
	}
	if form.Error() != "" {
		t.Errorf("Error should be reset, got '%s'", form.Error())
	}
}

// TestCreateFormHide verifies Hide makes form invisible.
func TestCreateFormHide(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.errorMessage = "some error"

	form.Hide()

	if form.Visible() {
		t.Error("Form should not be visible after Hide")
	}
	if form.Error() != "" {
		t.Error("Error should be cleared after Hide")
	}
}

// TestCreateFormSetSize verifies size setting.
func TestCreateFormSetSize(t *testing.T) {
	form := NewCreateForm()
	form.SetSize(100, 50)

	if form.width != 100 {
		t.Errorf("Expected width 100, got %d", form.width)
	}
	if form.height != 50 {
		t.Errorf("Expected height 50, got %d", form.height)
	}
}

// TestCreateFormSetError verifies error setting.
func TestCreateFormSetError(t *testing.T) {
	form := NewCreateForm()
	form.SetError("test error")

	if form.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", form.Error())
	}
}

// TestCreateFormFocusNext verifies tab navigation.
func TestCreateFormFocusNext(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	// Start at branch
	if form.Focused() != FieldBranch {
		t.Fatal("Should start at FieldBranch")
	}

	form.focusNext()
	if form.Focused() != FieldPath {
		t.Error("Should move to FieldPath")
	}

	form.focusNext()
	if form.Focused() != FieldCreateNewBranch {
		t.Error("Should move to FieldCreateNewBranch")
	}

	form.focusNext()
	if form.Focused() != FieldBranch {
		t.Error("Should wrap to FieldBranch")
	}
}

// TestCreateFormFocusPrev verifies shift-tab navigation.
func TestCreateFormFocusPrev(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	form.focusPrev()
	if form.Focused() != FieldCreateNewBranch {
		t.Error("Should move to FieldCreateNewBranch")
	}

	form.focusPrev()
	if form.Focused() != FieldPath {
		t.Error("Should move to FieldPath")
	}

	form.focusPrev()
	if form.Focused() != FieldBranch {
		t.Error("Should move to FieldBranch")
	}
}

// TestCreateFormInsertChar verifies character insertion.
func TestCreateFormInsertChar(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	// Insert into branch field
	form.insertChar('a')
	form.insertChar('b')
	form.insertChar('c')

	if form.Branch() != "abc" {
		t.Errorf("Expected branch 'abc', got '%s'", form.Branch())
	}

	// Move to path field and insert
	form.focusNext()
	form.insertChar('x')
	form.insertChar('y')

	if form.Path() != "xy" {
		t.Errorf("Expected path 'xy', got '%s'", form.Path())
	}
}

// TestCreateFormDeleteChar verifies backspace.
func TestCreateFormDeleteChar(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	form.branch = "abc"
	form.cursorPos = 3

	form.deleteChar()
	if form.Branch() != "ab" {
		t.Errorf("Expected branch 'ab', got '%s'", form.Branch())
	}

	form.deleteChar()
	if form.Branch() != "a" {
		t.Errorf("Expected branch 'a', got '%s'", form.Branch())
	}

	// Test delete at beginning (should do nothing)
	form.deleteChar()
	form.deleteChar() // This one is at position 0
	if form.Branch() != "" {
		t.Errorf("Expected empty branch, got '%s'", form.Branch())
	}
}

// TestCreateFormValidate verifies validation.
func TestCreateFormValidate(t *testing.T) {
	tests := []struct {
		name         string
		branch       string
		path         string
		createBranch bool
		expectValid  bool
		expectError  string
	}{
		{
			name:         "valid new branch",
			branch:       "feature",
			path:         "/path/to/worktree",
			createBranch: true,
			expectValid:  true,
		},
		{
			name:         "valid existing branch",
			branch:       "main",
			path:         "/path/to/worktree",
			createBranch: false,
			expectValid:  true,
		},
		{
			name:         "empty branch new",
			branch:       "",
			path:         "/path/to/worktree",
			createBranch: true,
			expectValid:  false,
			expectError:  "Branch name is required",
		},
		{
			name:         "empty branch existing",
			branch:       "",
			path:         "/path/to/worktree",
			createBranch: false,
			expectValid:  false,
			expectError:  "Existing branch name is required",
		},
		{
			name:         "empty path",
			branch:       "feature",
			path:         "",
			createBranch: true,
			expectValid:  false,
			expectError:  "Path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := NewCreateForm()
			form.branch = tt.branch
			form.path = tt.path
			form.createBranch = tt.createBranch

			valid := form.validate()
			if valid != tt.expectValid {
				t.Errorf("validate() = %v, want %v", valid, tt.expectValid)
			}
			if !valid && form.Error() != tt.expectError {
				t.Errorf("Error = '%s', want '%s'", form.Error(), tt.expectError)
			}
		})
	}
}

// TestCreateFormUpdateEscape verifies Escape key closes form.
func TestCreateFormUpdateEscape(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	cmd := form.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if form.Visible() {
		t.Error("Form should be hidden after Escape")
	}

	// Check command returns cancelled message
	if cmd == nil {
		t.Fatal("Expected command from Escape")
	}
	msg := cmd()
	if _, ok := msg.(CreateFormCancelledMsg); !ok {
		t.Errorf("Expected CreateFormCancelledMsg, got %T", msg)
	}
}

// TestCreateFormUpdateTab verifies Tab navigation.
func TestCreateFormUpdateTab(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	form.Update(tea.KeyMsg{Type: tea.KeyTab})
	if form.Focused() != FieldPath {
		t.Error("Tab should move to path field")
	}

	form.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if form.Focused() != FieldBranch {
		t.Error("Shift+Tab should move back to branch field")
	}
}

// TestCreateFormUpdateBackspace verifies backspace handling.
func TestCreateFormUpdateBackspace(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.branch = "abc"
	form.cursorPos = 3

	form.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	if form.Branch() != "ab" {
		t.Errorf("Expected branch 'ab', got '%s'", form.Branch())
	}
}

// TestCreateFormUpdateArrows verifies cursor movement.
func TestCreateFormUpdateArrows(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.branch = "abc"
	form.cursorPos = 3

	form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if form.cursorPos != 2 {
		t.Errorf("Expected cursor at 2, got %d", form.cursorPos)
	}

	form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if form.cursorPos != 0 {
		t.Errorf("Expected cursor at 0, got %d", form.cursorPos)
	}

	// Left at beginning should stay at 0
	form.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if form.cursorPos != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", form.cursorPos)
	}

	form.Update(tea.KeyMsg{Type: tea.KeyRight})
	if form.cursorPos != 1 {
		t.Errorf("Expected cursor at 1, got %d", form.cursorPos)
	}
}

// TestCreateFormUpdateSpace verifies space toggles checkbox.
func TestCreateFormUpdateSpace(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.focused = FieldCreateNewBranch

	// Toggle checkbox
	form.Update(tea.KeyMsg{Type: tea.KeySpace})
	if form.CreateBranchEnabled() {
		t.Error("Space should toggle createBranch to false")
	}

	form.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !form.CreateBranchEnabled() {
		t.Error("Space should toggle createBranch back to true")
	}
}

// TestCreateFormUpdateRunes verifies character input.
func TestCreateFormUpdateRunes(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h', 'i'}})

	if form.Branch() != "hi" {
		t.Errorf("Expected branch 'hi', got '%s'", form.Branch())
	}
}

// TestCreateFormUpdateEnterValid verifies Enter submits valid form.
func TestCreateFormUpdateEnterValid(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.branch = "feature"
	form.path = "/path/to/worktree"
	form.createBranch = true

	cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if form.Visible() {
		t.Error("Form should be hidden after valid submit")
	}

	if cmd == nil {
		t.Fatal("Expected command from Enter")
	}
	msg := cmd()
	submitted, ok := msg.(CreateFormSubmittedMsg)
	if !ok {
		t.Fatalf("Expected CreateFormSubmittedMsg, got %T", msg)
	}
	if submitted.Result.Branch != "feature" {
		t.Errorf("Expected branch 'feature', got '%s'", submitted.Result.Branch)
	}
	if submitted.Result.Path != "/path/to/worktree" {
		t.Errorf("Expected path '/path/to/worktree', got '%s'", submitted.Result.Path)
	}
	if !submitted.Result.CreateBranch {
		t.Error("Expected CreateBranch true")
	}
}

// TestCreateFormUpdateEnterInvalid verifies Enter with invalid form shows error.
func TestCreateFormUpdateEnterInvalid(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.branch = ""
	form.path = "/path"

	cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !form.Visible() {
		t.Error("Form should remain visible with invalid input")
	}
	if form.Error() == "" {
		t.Error("Error should be set for invalid input")
	}
	if cmd != nil {
		t.Error("Should not return command for invalid submit")
	}
}

// TestCreateFormUpdateWhenHidden verifies update does nothing when hidden.
func TestCreateFormUpdateWhenHidden(t *testing.T) {
	form := NewCreateForm()
	// Don't call Show()

	cmd := form.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Error("Should not return command when hidden")
	}
}

// TestCreateFormView verifies the view renders correctly.
func TestCreateFormView(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	view := form.View()

	if view == "" {
		t.Error("View should not be empty when visible")
	}
	if !strings.Contains(view, "Create New Worktree") {
		t.Error("View should contain title")
	}
	if !strings.Contains(view, "Branch name:") {
		t.Error("View should contain branch label")
	}
	if !strings.Contains(view, "Worktree path:") {
		t.Error("View should contain path label")
	}
	if !strings.Contains(view, "Create new branch") {
		t.Error("View should contain checkbox label")
	}
	if !strings.Contains(view, "Tab:") {
		t.Error("View should contain help text")
	}
}

// TestCreateFormViewWhenHidden verifies view is empty when hidden.
func TestCreateFormViewWhenHidden(t *testing.T) {
	form := NewCreateForm()
	// Don't call Show()

	view := form.View()

	if view != "" {
		t.Error("View should be empty when hidden")
	}
}

// TestCreateFormViewShowsError verifies error is shown in view.
func TestCreateFormViewShowsError(t *testing.T) {
	form := NewCreateForm()
	form.Show()
	form.SetError("Test error message")

	view := form.View()

	if !strings.Contains(view, "Test error message") {
		t.Error("View should contain error message")
	}
}

// TestCreateFormViewCheckboxState verifies checkbox rendering.
func TestCreateFormViewCheckboxState(t *testing.T) {
	form := NewCreateForm()
	form.Show()

	// Checked state
	view := form.View()
	if !strings.Contains(view, "[✓]") {
		t.Error("View should show checked checkbox")
	}

	// Unchecked state
	form.createBranch = false
	view = form.View()
	if !strings.Contains(view, "[ ]") {
		t.Error("View should show unchecked checkbox")
	}
}

// TestCreateFormRenderInputWithCursor verifies cursor rendering.
func TestCreateFormRenderInputWithCursor(t *testing.T) {
	form := NewCreateForm()

	tests := []struct {
		name     string
		text     string
		pos      int
		expected string
	}{
		{
			name:     "empty text",
			text:     "",
			pos:      0,
			expected: "│",
		},
		{
			name:     "cursor at start",
			text:     "abc",
			pos:      0,
			expected: "│abc",
		},
		{
			name:     "cursor in middle",
			text:     "abc",
			pos:      1,
			expected: "a│bc",
		},
		{
			name:     "cursor at end",
			text:     "abc",
			pos:      3,
			expected: "abc│",
		},
		{
			name:     "cursor beyond end",
			text:     "abc",
			pos:      5,
			expected: "abc│",
		},
		{
			name:     "negative position",
			text:     "abc",
			pos:      -1,
			expected: "│abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := form.renderInputWithCursor(tt.text, tt.pos)
			if result != tt.expected {
				t.Errorf("renderInputWithCursor(%q, %d) = %q, want %q",
					tt.text, tt.pos, result, tt.expected)
			}
		})
	}
}

// TestCreateFormFieldConstants verifies field constants are distinct.
func TestCreateFormFieldConstants(t *testing.T) {
	fields := []CreateFormField{FieldBranch, FieldPath, FieldCreateNewBranch}
	seen := make(map[CreateFormField]bool)

	for _, f := range fields {
		if seen[f] {
			t.Errorf("Field %v is duplicated", f)
		}
		seen[f] = true
	}
}
