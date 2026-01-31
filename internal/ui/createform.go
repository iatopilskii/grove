// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CreateFormField identifies which field is currently focused.
type CreateFormField int

const (
	// FieldBranch is the branch name input field.
	FieldBranch CreateFormField = iota
	// FieldPath is the worktree path input field.
	FieldPath
	// FieldCreateNewBranch is the checkbox for creating a new branch.
	FieldCreateNewBranch
)

// CreateFormResult contains the data from a completed form.
type CreateFormResult struct {
	Branch       string
	Path         string
	CreateBranch bool
}

// CreateFormSubmittedMsg is sent when the form is submitted.
type CreateFormSubmittedMsg struct {
	Result CreateFormResult
}

// CreateFormCancelledMsg is sent when the form is cancelled.
type CreateFormCancelledMsg struct{}

// CreateForm is a modal form for creating a new worktree.
type CreateForm struct {
	visible      bool
	focused      CreateFormField
	branch       string
	path         string
	createBranch bool
	width        int
	height       int
	cursorPos    int // cursor position within the current input field
	errorMessage string
}

// NewCreateForm creates a new worktree creation form.
func NewCreateForm() *CreateForm {
	return &CreateForm{
		createBranch: true, // Default to creating a new branch
	}
}

// Visible returns whether the form is currently visible.
func (f *CreateForm) Visible() bool {
	return f.visible
}

// Show makes the form visible and resets its fields.
func (f *CreateForm) Show() {
	f.visible = true
	f.focused = FieldBranch
	f.branch = ""
	f.path = ""
	f.createBranch = true
	f.cursorPos = 0
	f.errorMessage = ""
}

// Hide hides the form.
func (f *CreateForm) Hide() {
	f.visible = false
	f.errorMessage = ""
}

// SetSize sets the form dimensions.
func (f *CreateForm) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// Branch returns the current branch name input value.
func (f *CreateForm) Branch() string {
	return f.branch
}

// Path returns the current path input value.
func (f *CreateForm) Path() string {
	return f.path
}

// CreateBranchEnabled returns whether the "create new branch" option is enabled.
func (f *CreateForm) CreateBranchEnabled() bool {
	return f.createBranch
}

// Focused returns the currently focused field.
func (f *CreateForm) Focused() CreateFormField {
	return f.focused
}

// SetError sets an error message to display on the form.
func (f *CreateForm) SetError(msg string) {
	f.errorMessage = msg
}

// Error returns the current error message.
func (f *CreateForm) Error() string {
	return f.errorMessage
}

// focusNext moves focus to the next field.
func (f *CreateForm) focusNext() {
	switch f.focused {
	case FieldBranch:
		f.focused = FieldPath
		f.cursorPos = len(f.path)
	case FieldPath:
		f.focused = FieldCreateNewBranch
		f.cursorPos = 0
	case FieldCreateNewBranch:
		f.focused = FieldBranch
		f.cursorPos = len(f.branch)
	}
}

// focusPrev moves focus to the previous field.
func (f *CreateForm) focusPrev() {
	switch f.focused {
	case FieldBranch:
		f.focused = FieldCreateNewBranch
		f.cursorPos = 0
	case FieldPath:
		f.focused = FieldBranch
		f.cursorPos = len(f.branch)
	case FieldCreateNewBranch:
		f.focused = FieldPath
		f.cursorPos = len(f.path)
	}
}

// validate checks if the form input is valid.
func (f *CreateForm) validate() bool {
	if f.branch == "" && f.createBranch {
		f.errorMessage = "Branch name is required"
		return false
	}
	if f.branch == "" && !f.createBranch {
		f.errorMessage = "Existing branch name is required"
		return false
	}
	if f.path == "" {
		f.errorMessage = "Path is required"
		return false
	}
	f.errorMessage = ""
	return true
}

// submit validates and submits the form.
func (f *CreateForm) submit() tea.Cmd {
	if !f.validate() {
		return nil
	}

	result := CreateFormResult{
		Branch:       f.branch,
		Path:         f.path,
		CreateBranch: f.createBranch,
	}

	f.Hide()

	return func() tea.Msg {
		return CreateFormSubmittedMsg{Result: result}
	}
}

// insertChar inserts a character at the current cursor position.
func (f *CreateForm) insertChar(char rune) {
	switch f.focused {
	case FieldBranch:
		if f.cursorPos > len(f.branch) {
			f.cursorPos = len(f.branch)
		}
		f.branch = f.branch[:f.cursorPos] + string(char) + f.branch[f.cursorPos:]
		f.cursorPos++
	case FieldPath:
		if f.cursorPos > len(f.path) {
			f.cursorPos = len(f.path)
		}
		f.path = f.path[:f.cursorPos] + string(char) + f.path[f.cursorPos:]
		f.cursorPos++
	}
}

// deleteChar deletes the character before the cursor.
func (f *CreateForm) deleteChar() {
	switch f.focused {
	case FieldBranch:
		if f.cursorPos > 0 && len(f.branch) > 0 {
			f.branch = f.branch[:f.cursorPos-1] + f.branch[f.cursorPos:]
			f.cursorPos--
		}
	case FieldPath:
		if f.cursorPos > 0 && len(f.path) > 0 {
			f.path = f.path[:f.cursorPos-1] + f.path[f.cursorPos:]
			f.cursorPos--
		}
	}
}

// Update handles input messages for the form.
func (f *CreateForm) Update(msg tea.Msg) tea.Cmd {
	if !f.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			f.Hide()
			return func() tea.Msg {
				return CreateFormCancelledMsg{}
			}
		case tea.KeyEnter:
			return f.submit()
		case tea.KeyTab:
			f.focusNext()
		case tea.KeyShiftTab:
			f.focusPrev()
		case tea.KeyBackspace:
			f.deleteChar()
		case tea.KeyLeft:
			if f.focused == FieldBranch || f.focused == FieldPath {
				if f.cursorPos > 0 {
					f.cursorPos--
				}
			}
		case tea.KeyRight:
			if f.focused == FieldBranch {
				if f.cursorPos < len(f.branch) {
					f.cursorPos++
				}
			} else if f.focused == FieldPath {
				if f.cursorPos < len(f.path) {
					f.cursorPos++
				}
			}
		case tea.KeySpace:
			if f.focused == FieldCreateNewBranch {
				f.createBranch = !f.createBranch
			} else {
				f.insertChar(' ')
			}
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					f.insertChar(r)
				}
			}
		}
	}

	return nil
}

// View renders the form.
func (f *CreateForm) View() string {
	if !f.visible {
		return ""
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true).
		MarginBottom(1)

	// Label style
	labelStyle := lipgloss.NewStyle().
		Foreground(Colors.TextMuted)

	// Input field style (unfocused)
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.TextMuted).
		Padding(0, 1).
		Width(40)

	// Input field style (focused)
	inputFocusedStyle := inputStyle.
		BorderForeground(Colors.Primary)

	// Checkbox style
	checkboxStyle := lipgloss.NewStyle().
		Foreground(Colors.Text)

	// Error style
	errorStyle := lipgloss.NewStyle().
		Foreground(Colors.Error).
		Bold(true)

	var lines []string
	lines = append(lines, titleStyle.Render("Create New Worktree"))

	// Branch name field
	branchLabel := "Branch name:"
	if !f.createBranch {
		branchLabel = "Existing branch:"
	}
	lines = append(lines, labelStyle.Render(branchLabel))

	branchValue := f.branch
	if f.focused == FieldBranch {
		// Show cursor
		branchValue = f.renderInputWithCursor(f.branch, f.cursorPos)
		lines = append(lines, inputFocusedStyle.Render(branchValue))
	} else {
		if branchValue == "" {
			branchValue = " "
		}
		lines = append(lines, inputStyle.Render(branchValue))
	}
	lines = append(lines, "")

	// Path field
	lines = append(lines, labelStyle.Render("Worktree path:"))

	pathValue := f.path
	if f.focused == FieldPath {
		// Show cursor
		pathValue = f.renderInputWithCursor(f.path, f.cursorPos)
		lines = append(lines, inputFocusedStyle.Render(pathValue))
	} else {
		if pathValue == "" {
			pathValue = " "
		}
		lines = append(lines, inputStyle.Render(pathValue))
	}
	lines = append(lines, "")

	// Create new branch checkbox
	checkbox := "[ ]"
	if f.createBranch {
		checkbox = "[✓]"
	}
	checkboxLine := checkbox + " Create new branch"
	if f.focused == FieldCreateNewBranch {
		lines = append(lines, checkboxStyle.Bold(true).Foreground(Colors.Primary).Render(checkboxLine))
	} else {
		lines = append(lines, checkboxStyle.Render(checkboxLine))
	}

	// Error message
	if f.errorMessage != "" {
		lines = append(lines, "")
		lines = append(lines, errorStyle.Render("✗ "+f.errorMessage))
	}

	// Help text
	lines = append(lines, "")
	lines = append(lines, Styles.Help.Render("Tab: next field • Space: toggle • Enter: create • Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Box style
	boxStyle := Styles.Box.Padding(Padding.Small, Padding.Medium)

	return boxStyle.Render(content)
}

// renderInputWithCursor renders text with a cursor at the given position.
func (f *CreateForm) renderInputWithCursor(text string, pos int) string {
	if pos > len(text) {
		pos = len(text)
	}
	if pos < 0 {
		pos = 0
	}

	cursor := "│"
	if text == "" {
		return cursor
	}

	before := text[:pos]
	after := text[pos:]

	return before + cursor + after
}
