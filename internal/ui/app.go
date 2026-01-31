// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// App is the main application model implementing tea.Model.
// It uses the Elm architecture with Init, Update, and View methods.
type App struct {
	// quitting indicates the application should exit
	quitting bool
}

// NewApp creates and returns a new App instance.
func NewApp() *App {
	return &App{}
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			a.quitting = true
			return a, tea.Quit
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
	return "Git Worktree Manager\n\nPress 'q' to quit.\n"
}
