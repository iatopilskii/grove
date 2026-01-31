// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmDialog is a modal dialog that asks for user confirmation.
type ConfirmDialog struct {
	visible       bool
	title         string
	message       string
	confirmLabel  string
	cancelLabel   string
	dangerMode    bool
	forceOption   bool
	forceSelected bool
	selected      int // 0 = confirm, 1 = cancel
	data          interface{}
	width         int
	height        int
}

// NewConfirmDialog creates a new confirmation dialog.
func NewConfirmDialog() *ConfirmDialog {
	return &ConfirmDialog{
		confirmLabel: "Confirm",
		cancelLabel:  "Cancel",
	}
}

// Visible returns whether the dialog is currently visible.
func (d *ConfirmDialog) Visible() bool {
	return d.visible
}

// Title returns the dialog title.
func (d *ConfirmDialog) Title() string {
	return d.title
}

// Message returns the dialog message.
func (d *ConfirmDialog) Message() string {
	return d.message
}

// Data returns the associated data (e.g., the item being confirmed for deletion).
func (d *ConfirmDialog) Data() interface{} {
	return d.data
}

// ForceSelected returns whether the force option is selected.
func (d *ConfirmDialog) ForceSelected() bool {
	return d.forceSelected
}

// HasForceOption returns whether the dialog has a force option.
func (d *ConfirmDialog) HasForceOption() bool {
	return d.forceOption
}

// Selected returns the currently selected button (0 = confirm, 1 = cancel).
func (d *ConfirmDialog) Selected() int {
	return d.selected
}

// Show displays the confirmation dialog with the given title and message.
func (d *ConfirmDialog) Show(title, message string) {
	d.visible = true
	d.title = title
	d.message = message
	d.selected = 1 // Default to cancel for safety
	d.forceSelected = false
	d.data = nil
}

// ShowWithData displays the confirmation dialog and stores associated data.
func (d *ConfirmDialog) ShowWithData(title, message string, data interface{}) {
	d.Show(title, message)
	d.data = data
}

// ShowDanger displays the confirmation dialog styled for dangerous actions.
func (d *ConfirmDialog) ShowDanger(title, message string, data interface{}) {
	d.ShowWithData(title, message, data)
	d.dangerMode = true
}

// SetForceOption enables or disables the force checkbox option.
func (d *ConfirmDialog) SetForceOption(enabled bool) {
	d.forceOption = enabled
}

// SetConfirmLabel sets the text for the confirm button.
func (d *ConfirmDialog) SetConfirmLabel(label string) {
	d.confirmLabel = label
}

// SetCancelLabel sets the text for the cancel button.
func (d *ConfirmDialog) SetCancelLabel(label string) {
	d.cancelLabel = label
}

// Hide closes the confirmation dialog.
func (d *ConfirmDialog) Hide() {
	d.visible = false
	d.dangerMode = false
	d.forceOption = false
	d.forceSelected = false
	d.data = nil
	d.selected = 1
}

// SetSize sets the dialog dimensions.
func (d *ConfirmDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// MoveLeft moves the selection to the left (toward confirm).
func (d *ConfirmDialog) MoveLeft() {
	if d.selected > 0 {
		d.selected--
	}
}

// MoveRight moves the selection to the right (toward cancel).
func (d *ConfirmDialog) MoveRight() {
	if d.selected < 1 {
		d.selected++
	}
}

// ToggleForce toggles the force option checkbox.
func (d *ConfirmDialog) ToggleForce() {
	if d.forceOption {
		d.forceSelected = !d.forceSelected
	}
}

// ConfirmDialogResultMsg is sent when the dialog is confirmed.
type ConfirmDialogResultMsg struct {
	Confirmed bool
	Force     bool
	Data      interface{}
}

// Update handles input messages for the confirmation dialog.
func (d *ConfirmDialog) Update(msg tea.Msg) tea.Cmd {
	if !d.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			d.Hide()
			return func() tea.Msg {
				return ConfirmDialogResultMsg{Confirmed: false}
			}
		case tea.KeyEnter:
			confirmed := d.selected == 0
			force := d.forceSelected
			data := d.data
			d.Hide()
			return func() tea.Msg {
				return ConfirmDialogResultMsg{
					Confirmed: confirmed,
					Force:     force,
					Data:      data,
				}
			}
		case tea.KeyLeft, tea.KeyShiftTab:
			d.MoveLeft()
		case tea.KeyRight, tea.KeyTab:
			d.MoveRight()
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'h':
					d.MoveLeft()
				case 'l':
					d.MoveRight()
				case 'y':
					// Quick confirm
					force := d.forceSelected
					data := d.data
					d.Hide()
					return func() tea.Msg {
						return ConfirmDialogResultMsg{
							Confirmed: true,
							Force:     force,
							Data:      data,
						}
					}
				case 'n':
					// Quick cancel
					d.Hide()
					return func() tea.Msg {
						return ConfirmDialogResultMsg{Confirmed: false}
					}
				case ' ':
					// Space toggles force when enabled
					d.ToggleForce()
				case 'f':
					// 'f' also toggles force when enabled
					d.ToggleForce()
				}
			}
		}
	}
	return nil
}

// View renders the confirmation dialog.
func (d *ConfirmDialog) View() string {
	if !d.visible {
		return ""
	}

	// Title styling
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)

	if d.dangerMode {
		titleStyle = titleStyle.Foreground(Colors.Error)
	} else {
		titleStyle = titleStyle.Foreground(Colors.Text)
	}

	// Message styling
	messageStyle := lipgloss.NewStyle().
		Foreground(Colors.Text).
		MarginBottom(1)

	// Button styles
	confirmBtnStyle := lipgloss.NewStyle().
		Padding(0, 2)
	cancelBtnStyle := lipgloss.NewStyle().
		Padding(0, 2)

	// Apply selection styling
	if d.selected == 0 {
		if d.dangerMode {
			confirmBtnStyle = confirmBtnStyle.
				Background(Colors.Error).
				Foreground(Colors.OnError).
				Bold(true)
		} else {
			confirmBtnStyle = confirmBtnStyle.
				Background(Colors.Primary).
				Foreground(Colors.OnPrimary).
				Bold(true)
		}
		cancelBtnStyle = cancelBtnStyle.
			Foreground(Colors.TextMuted)
	} else {
		confirmBtnStyle = confirmBtnStyle.
			Foreground(Colors.TextMuted)
		cancelBtnStyle = cancelBtnStyle.
			Background(Colors.Primary).
			Foreground(Colors.OnPrimary).
			Bold(true)
	}

	var lines []string
	lines = append(lines, titleStyle.Render(d.title))
	lines = append(lines, messageStyle.Render(d.message))

	// Force option checkbox
	if d.forceOption {
		checkboxStyle := lipgloss.NewStyle().
			Foreground(Colors.Text).
			MarginBottom(1)

		checkbox := "[ ]"
		if d.forceSelected {
			checkbox = "[x]"
		}
		forceText := checkbox + " Force removal (ignore uncommitted changes)"
		lines = append(lines, checkboxStyle.Render(forceText))
	}

	// Buttons
	confirmBtn := confirmBtnStyle.Render(d.confirmLabel)
	cancelBtn := cancelBtnStyle.Render(d.cancelLabel)
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, confirmBtn, "  ", cancelBtn)
	lines = append(lines, buttons)

	// Help text using centralized style
	helpStyle := Styles.Help.MarginTop(1)
	lines = append(lines, helpStyle.Render("y/n: quick answer • ←/→: select • Enter: confirm • Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Box style with thin border and consistent padding
	boxStyle := Styles.Box.Padding(Padding.Small, Padding.Medium)
	if d.dangerMode {
		boxStyle = boxStyle.BorderForeground(Colors.Error)
	}

	return boxStyle.Render(content)
}
