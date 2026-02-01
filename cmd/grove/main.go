// Package main is the entry point for the Git Worktree TUI application.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/iatopilskii/grove/internal/ui"
)

func main() {
	// Load and apply configuration from ~/.config/grove/config.yaml
	// Invalid config falls back to defaults; missing file is not an error
	if err := ui.LoadAndApplyTheme(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: theme config error: %v (using defaults)\n", err)
	}

	app := ui.NewApp()
	p := tea.NewProgram(app)

	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	// Check if we should cd to a target path (after worktree creation)
	if finalApp, ok := m.(*ui.App); ok {
		if targetPath := finalApp.TargetPath(); targetPath != "" {
			fmt.Println(targetPath)
			os.Exit(2)
		}
	}
}
