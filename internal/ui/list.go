// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListItem represents a single item in the list.
type ListItem struct {
	ID          string
	Title       string
	Description string
}

// List is a scrollable list component.
type List struct {
	items    []ListItem
	selected int
	width    int
	height   int
	offsetX  int // X position on screen for mouse handling
	offsetY  int // Y position on screen for mouse handling
}

// NewList creates a new list with the given items.
func NewList(items []ListItem) *List {
	return &List{
		items:    items,
		selected: 0,
	}
}

// Items returns all items in the list.
func (l *List) Items() []ListItem {
	return l.items
}

// SetItems replaces the items in the list.
func (l *List) SetItems(items []ListItem) {
	l.items = items
	// Clamp selection to valid range
	if len(items) == 0 {
		l.selected = 0
	} else if l.selected >= len(items) {
		l.selected = len(items) - 1
	}
}

// Selected returns the index of the currently selected item.
func (l *List) Selected() int {
	return l.selected
}

// SetSelected sets the selected index with bounds checking.
func (l *List) SetSelected(index int) {
	if len(l.items) == 0 {
		l.selected = 0
		return
	}
	if index < 0 {
		l.selected = 0
		return
	}
	if index >= len(l.items) {
		l.selected = len(l.items) - 1
		return
	}
	l.selected = index
}

// SelectedItem returns the currently selected item, or nil if the list is empty.
func (l *List) SelectedItem() *ListItem {
	if len(l.items) == 0 || l.selected < 0 || l.selected >= len(l.items) {
		return nil
	}
	return &l.items[l.selected]
}

// MoveDown moves the selection down by one.
func (l *List) MoveDown() {
	if len(l.items) == 0 {
		return
	}
	if l.selected < len(l.items)-1 {
		l.selected++
	}
}

// MoveUp moves the selection up by one.
func (l *List) MoveUp() {
	if len(l.items) == 0 {
		return
	}
	if l.selected > 0 {
		l.selected--
	}
}

// PageDown moves the selection down by one page (based on visible height).
func (l *List) PageDown() {
	if len(l.items) == 0 {
		return
	}
	pageSize := l.height
	if pageSize < 1 {
		pageSize = 1
	}
	l.selected += pageSize
	if l.selected >= len(l.items) {
		l.selected = len(l.items) - 1
	}
}

// PageUp moves the selection up by one page (based on visible height).
func (l *List) PageUp() {
	if len(l.items) == 0 {
		return
	}
	pageSize := l.height
	if pageSize < 1 {
		pageSize = 1
	}
	l.selected -= pageSize
	if l.selected < 0 {
		l.selected = 0
	}
}

// SetSize sets the list dimensions for rendering.
func (l *List) SetSize(width, height int) {
	l.width = width
	l.height = height
}

// SetOffset sets the screen position of the list for mouse handling.
func (l *List) SetOffset(x, y int) {
	l.offsetX = x
	l.offsetY = y
}

// IsInBounds checks if the given screen coordinates are within the list bounds.
func (l *List) IsInBounds(x, y int) bool {
	return x >= l.offsetX && x < l.offsetX+l.width &&
		y >= l.offsetY && y < l.offsetY+l.height
}

// Update handles input messages for the list.
func (l *List) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			l.MoveDown()
		case tea.KeyUp:
			l.MoveUp()
		case tea.KeyPgDown:
			l.PageDown()
		case tea.KeyPgUp:
			l.PageUp()
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'j':
					l.MoveDown()
				case 'k':
					l.MoveUp()
				}
			}
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonLeft:
			// Handle click to select item
			if len(l.items) > 0 && l.IsInBounds(msg.X, msg.Y) {
				// Calculate which item was clicked
				clickedIndex := msg.Y - l.offsetY
				if clickedIndex >= 0 && clickedIndex < len(l.items) {
					l.SetSelected(clickedIndex)
				}
			}
		case tea.MouseButtonWheelDown:
			l.MoveDown()
		case tea.MouseButtonWheelUp:
			l.MoveUp()
		}
	}
	return nil
}

// View renders the list.
func (l *List) View() string {
	if len(l.items) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#888888", Dark: "#666666"}).
			Italic(true)
		return emptyStyle.Render("No items")
	}

	// Calculate effective width
	effectiveWidth := l.width - 2
	if effectiveWidth < 0 {
		effectiveWidth = 0
	}

	// Define styles
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
		Bold(true).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"}).
		Padding(0, 1)

	// Apply width if set
	if effectiveWidth > 0 {
		selectedStyle = selectedStyle.Width(effectiveWidth)
		normalStyle = normalStyle.Width(effectiveWidth)
	}

	// Selection indicator
	selectedPrefix := "▸ "
	normalPrefix := "  "

	var lines []string
	for i, item := range l.items {
		if i == l.selected {
			lines = append(lines, selectedStyle.Render(selectedPrefix+item.Title))
		} else {
			lines = append(lines, normalStyle.Render(normalPrefix+item.Title))
		}
	}

	return strings.Join(lines, "\n")
}
