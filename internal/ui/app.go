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
	// width is the terminal width
	width int
	// height is the terminal height
	height int
}

// NewApp creates and returns a new App instance.
func NewApp() *App {
	return &App{
		tabs: NewTabs(),
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
		return a, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			a.quitting = true
			return a, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			a.tabs.Update(msg)
			return a, nil
		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
				a.quitting = true
				return a, tea.Quit
			}
		}
	}
	return a, nil
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
	contentStyle := lipgloss.NewStyle().
		Padding(1, 2)

	var content string
	switch a.tabs.Active() {
	case TabWorktrees:
		content = "Worktrees content\n\nThis will show the list of git worktrees."
	case TabBranches:
		content = "Branches content\n\nThis will show the list of branches."
	case TabSettings:
		content = "Settings content\n\nThis will show application settings."
	}

	b.WriteString(contentStyle.Render(content))
	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"})
	b.WriteString(helpStyle.Render("Tab/Shift+Tab: switch tabs • q: quit"))

	return b.String()
}
