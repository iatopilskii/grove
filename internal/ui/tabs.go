// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab in the tab bar.
type Tab int

const (
	// TabWorktrees shows the worktrees list.
	TabWorktrees Tab = iota
	// TabBranches shows the branches list.
	TabBranches
	// TabSettings shows the settings view.
	TabSettings
)

// String returns the display name of the tab.
func (t Tab) String() string {
	switch t {
	case TabWorktrees:
		return "Worktrees"
	case TabBranches:
		return "Branches"
	case TabSettings:
		return "Settings"
	default:
		return "Unknown"
	}
}

// TabCount is the total number of tabs.
const TabCount = 3

// TabPosition represents the screen position of a tab for mouse handling.
type TabPosition struct {
	Tab    Tab
	StartX int
	EndX   int
}

// Tabs is the tab bar model.
type Tabs struct {
	active Tab
	width  int
}

// NewTabs creates a new tab bar with Worktrees as the default active tab.
func NewTabs() *Tabs {
	return &Tabs{
		active: TabWorktrees,
	}
}

// Active returns the currently active tab.
func (t *Tabs) Active() Tab {
	return t.active
}

// SetActive sets the active tab.
func (t *Tabs) SetActive(tab Tab) {
	if tab >= 0 && tab < TabCount {
		t.active = tab
	}
}

// Next moves to the next tab, wrapping around to the first tab.
func (t *Tabs) Next() {
	t.active = (t.active + 1) % TabCount
}

// Prev moves to the previous tab, wrapping around to the last tab.
func (t *Tabs) Prev() {
	t.active = (t.active - 1 + TabCount) % TabCount
}

// SetWidth sets the available width for rendering.
func (t *Tabs) SetWidth(w int) {
	t.width = w
}

// GetTabPositions calculates the screen positions of each tab.
func (t *Tabs) GetTabPositions() []TabPosition {
	positions := make([]TabPosition, TabCount)
	currentX := 0

	// Each tab has padding of 2 on each side ("  TabName  ")
	for i := Tab(0); i < TabCount; i++ {
		tabWidth := len(i.String()) + 4 // 2 padding on each side
		positions[i] = TabPosition{
			Tab:    i,
			StartX: currentX,
			EndX:   currentX + tabWidth,
		}
		currentX += tabWidth
	}

	return positions
}

// Update handles key messages for tab navigation.
func (t *Tabs) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			t.Next()
		case tea.KeyShiftTab:
			t.Prev()
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft && msg.Y == 0 {
			// Click is on the tab row, determine which tab was clicked
			positions := t.GetTabPositions()
			for _, pos := range positions {
				if msg.X >= pos.StartX && msg.X < pos.EndX {
					t.SetActive(pos.Tab)
					break
				}
			}
		}
	}
	return nil
}

// View renders the tab bar.
func (t *Tabs) View() string {
	// Base tab style with consistent padding
	tabStyle := lipgloss.NewStyle().
		Padding(0, Padding.Medium)

	// Active tab: subtle underline effect via primary color text
	activeStyle := tabStyle.
		Foreground(Colors.Primary).
		Bold(true)

	// Inactive tab: muted color
	inactiveStyle := tabStyle.
		Foreground(Colors.TextMuted)

	// Build tab bar
	var tabs []string
	for i := Tab(0); i < TabCount; i++ {
		if i == t.active {
			tabs = append(tabs, activeStyle.Render(i.String()))
		} else {
			tabs = append(tabs, inactiveStyle.Render(i.String()))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	// Add thin border below tabs using single-line character
	border := lipgloss.NewStyle().
		Foreground(Colors.Border).
		Render(strings.Repeat(Borders.Thin.Bottom, max(t.width, lipgloss.Width(row))))

	return row + "\n" + border
}
