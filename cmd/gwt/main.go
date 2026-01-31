// Package main is the entry point for the Git Worktree TUI application.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ilatopilskij/gwt/internal/ui"
)

func main() {
	app := ui.NewApp()
	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
