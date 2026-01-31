// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ilatopilskij/gwt/internal/git"
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
	// actionMenu is the action menu modal
	actionMenu *ActionMenu
	// feedback is the feedback message component
	feedback *Feedback
	// createForm is the worktree creation form modal
	createForm *CreateForm
	// width is the terminal width
	width int
	// height is the terminal height
	height int
	// worktrees stores the git worktrees
	worktrees []git.Worktree
	// gitError stores any error from git operations
	gitError error
	// repoPath is the path to the git repository
	repoPath string
}

// NewApp creates and returns a new App instance.
// It attempts to load worktrees from the current directory.
func NewApp() *App {
	return NewAppWithPath("")
}

// NewAppWithPath creates a new App instance for a specific path.
// If path is empty, uses the current working directory.
func NewAppWithPath(path string) *App {
	app := &App{
		tabs:       NewTabs(),
		list:       NewList(nil),
		details:    NewDetails(),
		actionMenu: NewActionMenu(),
		feedback:   NewFeedback(),
		createForm: NewCreateForm(),
		repoPath:   path,
	}

	// Determine the repository path
	if path == "" {
		var err error
		path, err = git.GetCurrentDirectory()
		if err != nil {
			app.gitError = err
			return app
		}
		app.repoPath = path
	}

	// Load worktrees
	app.loadWorktrees()

	return app
}

// NewAppWithItems creates a new App instance with predefined items.
// This is primarily used for testing.
func NewAppWithItems(items []ListItem) *App {
	list := NewList(items)
	details := NewDetails()

	// Initialize details with first item
	if len(items) > 0 {
		details.SetItem(list.SelectedItem())
	}

	return &App{
		tabs:       NewTabs(),
		list:       list,
		details:    details,
		actionMenu: NewActionMenu(),
		feedback:   NewFeedback(),
		createForm: NewCreateForm(),
	}
}

// loadWorktrees loads git worktrees from the repository and updates the list.
func (a *App) loadWorktrees() {
	worktrees, err := git.ListWorktrees(a.repoPath)
	if err != nil {
		a.gitError = err
		a.worktrees = nil
		a.list.SetItems(nil)
		return
	}

	a.worktrees = worktrees
	a.gitError = nil

	// Convert worktrees to list items
	items := make([]ListItem, len(worktrees))
	for i, wt := range worktrees {
		items[i] = worktreeToListItem(wt)
	}

	a.list.SetItems(items)

	// Initialize details with first item
	if len(items) > 0 {
		a.details.SetItem(a.list.SelectedItem())
	}
}

// worktreeToListItem converts a git.Worktree to a ListItem.
func worktreeToListItem(wt git.Worktree) ListItem {
	var description string
	if wt.IsBare {
		description = "Bare repository at " + wt.Path
	} else if wt.IsDetached {
		description = "Detached HEAD at " + wt.Path
	} else if wt.Branch != "" {
		description = "Branch: " + wt.Branch + "\nPath: " + wt.Path
	} else {
		description = "Path: " + wt.Path
	}

	return ListItem{
		ID:          wt.Path,
		Title:       wt.Name(),
		Description: description,
	}
}

// Worktrees returns the list of git worktrees.
func (a *App) Worktrees() []git.Worktree {
	return a.worktrees
}

// GitError returns any error from git operations.
func (a *App) GitError() error {
	return a.gitError
}

// IsInGitRepo returns true if the app is running in a git repository.
func (a *App) IsInGitRepo() bool {
	return a.gitError == nil && !git.IsNotGitRepoError(a.gitError)
}

// RefreshWorktrees reloads the worktree list from git.
func (a *App) RefreshWorktrees() {
	a.loadWorktrees()
}

// Init initializes the application and returns an initial command.
// This is called once when the program starts.
func (a *App) Init() tea.Cmd {
	return tea.EnableMouseCellMotion
}

// Update handles incoming messages and updates the model accordingly.
// It returns the updated model and any command to execute.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle action execution results and form submissions
	switch msg := msg.(type) {
	case ActionExecutedMsg:
		return a.handleActionExecuted(msg)
	case ClearFeedbackMsg:
		a.feedback.Update(msg)
		return a, nil
	case CreateFormSubmittedMsg:
		return a.handleCreateFormSubmitted(msg)
	case CreateFormCancelledMsg:
		// Form was cancelled, nothing to do
		return a, nil
	}

	// If create form is visible, route all key events to it
	if a.createForm.Visible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// Allow Ctrl+C to quit even with form open
			if keyMsg.Type == tea.KeyCtrlC {
				a.quitting = true
				return a, tea.Quit
			}
			cmd := a.createForm.Update(keyMsg)
			return a, cmd
		}
	}

	// If action menu is visible, route all key events to it
	if a.actionMenu.Visible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// Allow Ctrl+C to quit even with menu open
			if keyMsg.Type == tea.KeyCtrlC {
				a.quitting = true
				return a, tea.Quit
			}
			cmd := a.actionMenu.Update(keyMsg)
			return a, cmd
		}
	}

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
		case tea.KeyEnter:
			// Open action menu on Worktrees or Branches tabs
			if a.tabs.Active() == TabWorktrees || a.tabs.Active() == TabBranches {
				if item := a.list.SelectedItem(); item != nil {
					a.actionMenu.Show(item)
				}
			}
			return a, nil
		case tea.KeyEsc:
			// Escape cancels action menu (if visible)
			if a.actionMenu.Visible() {
				a.actionMenu.Hide()
			}
			return a, nil
		case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
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
				case 'n':
					// Open create form on Worktrees tab
					if a.tabs.Active() == TabWorktrees && !git.IsNotGitRepoError(a.gitError) {
						a.createForm.Show()
					}
					return a, nil
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

	case tea.MouseMsg:
		// Handle mouse events
		if msg.Y == 0 {
			// Click on tab bar row
			a.tabs.Update(msg)
		} else if a.tabs.Active() == TabWorktrees || a.tabs.Active() == TabBranches {
			// Handle mouse in list pane
			if a.list.IsInBounds(msg.X, msg.Y) || msg.Button == tea.MouseButtonWheelDown || msg.Button == tea.MouseButtonWheelUp {
				a.list.Update(msg)
				a.details.SetItem(a.list.SelectedItem())
			}
		}
		return a, nil
	}
	return a, nil
}

// handleActionExecuted processes an action that was executed from the menu.
func (a *App) handleActionExecuted(msg ActionExecutedMsg) (tea.Model, tea.Cmd) {
	if msg.Action == nil {
		return a, nil
	}

	// Execute the action and show feedback
	switch msg.Action.ID {
	case "open":
		// For now, just show success feedback
		// In a real implementation, this would open a terminal at the worktree path
		cmd := a.feedback.ShowSuccess("Opened worktree: " + msg.Item.Title)
		return a, cmd
	case "cd":
		// For now, just show success feedback
		// In a real implementation, this would copy the path to clipboard
		cmd := a.feedback.ShowSuccess("Path copied to clipboard")
		return a, cmd
	case "delete":
		// For now, just show info feedback
		// In a real implementation, this would prompt for confirmation
		cmd := a.feedback.ShowInfo("Delete action: " + msg.Item.Title + " (not implemented)")
		return a, cmd
	default:
		cmd := a.feedback.ShowError("Unknown action: " + msg.Action.ID)
		return a, cmd
	}
}

// handleCreateFormSubmitted processes the submitted create worktree form.
func (a *App) handleCreateFormSubmitted(msg CreateFormSubmittedMsg) (tea.Model, tea.Cmd) {
	opts := git.AddWorktreeOptions{
		Path:         msg.Result.Path,
		Branch:       msg.Result.Branch,
		CreateBranch: msg.Result.CreateBranch,
	}

	err := git.AddWorktree(a.repoPath, opts)
	if err != nil {
		cmd := a.feedback.ShowError("Failed to create worktree: " + err.Error())
		return a, cmd
	}

	// Refresh the worktree list
	a.loadWorktrees()

	cmd := a.feedback.ShowSuccess("Created worktree: " + msg.Result.Branch + " at " + msg.Result.Path)
	return a, cmd
}

// CreateForm returns the create form component for testing.
func (a *App) CreateForm() *CreateForm {
	return a.createForm
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
	a.list.SetOffset(0, 3) // List starts at Y=3 (after tabs and border, which take 2 lines + 1 newline)
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
		// Show error if not in a git repository
		if git.IsNotGitRepoError(a.gitError) {
			b.WriteString(a.renderGitError())
		} else {
			b.WriteString(a.renderTwoPaneLayout())
		}
	case TabSettings:
		contentStyle := lipgloss.NewStyle().
			Padding(1, 2)
		content := "Settings content\n\nThis will show application settings."
		b.WriteString(contentStyle.Render(content))
	}

	b.WriteString("\n\n")

	// Show feedback message if visible
	if a.feedback.Visible() {
		b.WriteString(a.feedback.View())
		b.WriteString("\n\n")
	}

	// Help text using centralized style
	helpText := "↑/↓: navigate • Enter: action • n: new worktree • PgUp/PgDn: scroll page • Tab/Shift+Tab: switch tabs • q: quit"
	b.WriteString(Styles.Help.Render(helpText))

	// If action menu is visible, render it as an overlay
	if a.actionMenu.Visible() {
		b.WriteString("\n\n")
		b.WriteString(a.actionMenu.View())
	}

	// If create form is visible, render it as an overlay
	if a.createForm.Visible() {
		b.WriteString("\n\n")
		b.WriteString(a.createForm.View())
	}

	return b.String()
}

// renderTwoPaneLayout renders the list and details side by side.
func (a *App) renderTwoPaneLayout() string {
	listView := a.list.View()
	detailsView := a.details.View()

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, listView, " ", detailsView)
}

// renderGitError renders an error message for git-related errors.
func (a *App) renderGitError() string {
	errorStyle := lipgloss.NewStyle().
		Padding(2, 4).
		Foreground(Colors.Error)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(Colors.Error)

	var b strings.Builder
	b.WriteString(titleStyle.Render("Not a Git Repository"))
	b.WriteString("\n\n")
	b.WriteString("This directory is not part of a git repository.")
	b.WriteString("\n\n")
	b.WriteString("Please run this application from within a git repository to manage worktrees.")
	b.WriteString("\n\n")
	b.WriteString("To initialize a git repository, run:")
	b.WriteString("\n")
	b.WriteString("  git init")

	return errorStyle.Render(b.String())
}
