// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App is the main application model implementing tea.Model.
// It uses the Elm architecture with Init, Update, and View methods.
type App struct {
	// quitting indicates the application should exit
	quitting bool
	// tabs is the tab bar component
	tabs *Tabs
	// list is the list pane component
	list *List
	// details is the details pane component
	details *Details
	// width is the terminal width
	width int
	// height is the terminal height
	height int
}

// NewApp creates and returns a new App instance.
func NewApp() *App {
	// Sample worktree items for demonstration
	items := []ListItem{
		{ID: "main", Title: "main", Description: "Main worktree at /path/to/repo"},
		{ID: "feature-1", Title: "feature-1", Description: "Feature branch worktree at /path/to/repo-feature-1"},
		{ID: "bugfix-2", Title: "bugfix-2", Description: "Bugfix branch worktree at /path/to/repo-bugfix-2"},
	}

	list := NewList(items)
	details := NewDetails()

	// Initialize details with first item
	if len(items) > 0 {
		details.SetItem(list.SelectedItem())
	}

	return &App{
		tabs:    NewTabs(),
		list:    list,
		details: details,
	}
}

// Init initializes the application and returns an initial command.
// This is called once when the program starts.
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model accordingly.
// It returns the updated model and any command to execute.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.tabs.SetWidth(msg.Width)
		a.updatePaneSizes()
		return a, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			a.quitting = true
			return a, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			a.tabs.Update(msg)
			return a, nil
		case tea.KeyUp, tea.KeyDown:
			// Handle list navigation on Worktrees and Branches tabs
			if a.tabs.Active() == TabWorktrees || a.tabs.Active() == TabBranches {
				a.list.Update(msg)
				a.details.SetItem(a.list.SelectedItem())
			}
			return a, nil
		case tea.KeyRunes:
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'q':
					a.quitting = true
					return a, tea.Quit
				case 'j', 'k':
					// Handle vim-style navigation
					if a.tabs.Active() == TabWorktrees || a.tabs.Active() == TabBranches {
						a.list.Update(msg)
						a.details.SetItem(a.list.SelectedItem())
					}
					return a, nil
				}
			}
		}
	}
	return a, nil
}

// updatePaneSizes updates the sizes of list and details panes based on terminal size.
func (a *App) updatePaneSizes() {
	// Calculate available space after tabs and help text
	// Tabs take ~2 lines, help takes ~1 line, leave some margin
	availableHeight := a.height - 4
	if availableHeight < 0 {
		availableHeight = 0
	}

	// Split width between list and details (40% list, 60% details)
	listWidth := a.width * 40 / 100
	detailsWidth := a.width - listWidth - 1 // -1 for separator

	if listWidth < 0 {
		listWidth = 0
	}
	if detailsWidth < 0 {
		detailsWidth = 0
	}

	a.list.SetSize(listWidth, availableHeight)
	a.details.SetSize(detailsWidth, availableHeight)
}

// View renders the current state of the application as a string.
func (a *App) View() string {
	if a.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder

	// Render tab bar at top
	b.WriteString(a.tabs.View())
	b.WriteString("\n")

	// Render content area based on active tab
	switch a.tabs.Active() {
	case TabWorktrees, TabBranches:
		b.WriteString(a.renderTwoPaneLayout())
	case TabSettings:
		contentStyle := lipgloss.NewStyle().
			Padding(1, 2)
		content := "Settings content\n\nThis will show application settings."
		b.WriteString(contentStyle.Render(content))
	}

	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"})
	b.WriteString(helpStyle.Render("↑/↓: navigate • Tab/Shift+Tab: switch tabs • q: quit"))

	return b.String()
}

// renderTwoPaneLayout renders the list and details side by side.
func (a *App) renderTwoPaneLayout() string {
	listView := a.list.View()
	detailsView := a.details.View()

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, listView, " ", detailsView)
}
