// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Action represents an action that can be performed on an item.
type Action struct {
	ID          string
	Label       string
	Description string
}

// ActionMenu is a modal dialog that displays available actions for an item.
type ActionMenu struct {
	visible  bool
	item     *ListItem
	actions  []Action
	selected int
	width    int
	height   int
}

// NewActionMenu creates a new action menu.
func NewActionMenu() *ActionMenu {
	return &ActionMenu{
		actions: defaultWorktreeActions(),
	}
}

// defaultWorktreeActions returns the default actions available for worktrees.
func defaultWorktreeActions() []Action {
	return []Action{
		{ID: "open", Label: "Open", Description: "Open worktree in new terminal"},
		{ID: "cd", Label: "Copy Path", Description: "Copy worktree path to clipboard"},
		{ID: "delete", Label: "Delete", Description: "Remove this worktree"},
	}
}

// Visible returns whether the action menu is currently visible.
func (m *ActionMenu) Visible() bool {
	return m.visible
}

// Show makes the action menu visible for the given item.
func (m *ActionMenu) Show(item *ListItem) {
	m.visible = true
	m.item = item
	m.selected = 0
}

// Hide hides the action menu.
func (m *ActionMenu) Hide() {
	m.visible = false
	m.item = nil
	m.selected = 0
}

// Item returns the item the action menu is showing actions for.
func (m *ActionMenu) Item() *ListItem {
	return m.item
}

// Actions returns the list of available actions.
func (m *ActionMenu) Actions() []Action {
	return m.actions
}

// SetActions sets the list of available actions.
func (m *ActionMenu) SetActions(actions []Action) {
	m.actions = actions
	if m.selected >= len(actions) {
		m.selected = 0
	}
}

// Selected returns the index of the currently selected action.
func (m *ActionMenu) Selected() int {
	return m.selected
}

// SelectedAction returns the currently selected action, or nil if none.
func (m *ActionMenu) SelectedAction() *Action {
	if len(m.actions) == 0 || m.selected < 0 || m.selected >= len(m.actions) {
		return nil
	}
	return &m.actions[m.selected]
}

// MoveUp moves the selection up by one.
func (m *ActionMenu) MoveUp() {
	if len(m.actions) == 0 {
		return
	}
	if m.selected > 0 {
		m.selected--
	}
}

// MoveDown moves the selection down by one.
func (m *ActionMenu) MoveDown() {
	if len(m.actions) == 0 {
		return
	}
	if m.selected < len(m.actions)-1 {
		m.selected++
	}
}

// SetSize sets the action menu dimensions.
func (m *ActionMenu) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// ActionExecutedMsg is sent when an action is executed.
type ActionExecutedMsg struct {
	Action *Action
	Item   *ListItem
}

// Update handles input messages for the action menu.
func (m *ActionMenu) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.Hide()
		case tea.KeyEnter:
			// Execute the selected action
			if action := m.SelectedAction(); action != nil {
				item := m.item
				m.Hide()
				return func() tea.Msg {
					return ActionExecutedMsg{Action: action, Item: item}
				}
			}
		case tea.KeyUp:
			m.MoveUp()
		case tea.KeyDown:
			m.MoveDown()
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'k':
					m.MoveUp()
				case 'j':
					m.MoveDown()
				}
			}
		}
	}
	return nil
}

// View renders the action menu.
func (m *ActionMenu) View() string {
	if !m.visible {
		return ""
	}

	// Colors
	borderColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	titleColor := lipgloss.AdaptiveColor{Light: "#333333", Dark: "#EEEEEE"}
	selectedBg := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	selectedFg := lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}
	normalFg := lipgloss.AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"}
	descColor := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		MarginBottom(1)

	title := "Actions"
	if m.item != nil {
		title = "Actions: " + m.item.Title
	}

	// Action items
	selectedStyle := lipgloss.NewStyle().
		Background(selectedBg).
		Foreground(selectedFg).
		Bold(true).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(normalFg).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(descColor).
		Italic(true).
		PaddingLeft(3)

	var lines []string
	lines = append(lines, titleStyle.Render(title))

	for i, action := range m.actions {
		var line string
		if i == m.selected {
			line = selectedStyle.Render("▸ " + action.Label)
			if action.Description != "" {
				line += "\n" + descStyle.Render(action.Description)
			}
		} else {
			line = normalStyle.Render("  " + action.Label)
		}
		lines = append(lines, line)
	}

	// Add help text
	helpStyle := lipgloss.NewStyle().
		Foreground(descColor).
		MarginTop(1)
	lines = append(lines, helpStyle.Render("↑/↓: navigate • Enter: select • Esc: cancel"))

	content := strings.Join(lines, "\n")

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)

	return boxStyle.Render(content)
}
