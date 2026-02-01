// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ilatopilskij/grove/internal/git"
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
	// confirmDialog is the confirmation dialog modal
	confirmDialog *ConfirmDialog
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
	// targetPath is the path to cd to after quitting (for shell wrapper)
	targetPath string
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
		tabs:          NewTabs(),
		list:          NewList(nil),
		details:       NewDetails(),
		actionMenu:    NewActionMenu(),
		feedback:      NewFeedback(),
		createForm:    NewCreateForm(),
		confirmDialog: NewConfirmDialog(),
		repoPath:      path,
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
		tabs:          NewTabs(),
		list:          list,
		details:       details,
		actionMenu:    NewActionMenu(),
		feedback:      NewFeedback(),
		createForm:    NewCreateForm(),
		confirmDialog: NewConfirmDialog(),
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

// worktreeToListItem converts a git.Worktree to a ListItem with status information.
func worktreeToListItem(wt git.Worktree) ListItem {
	// Get worktree status (modified/staged file counts)
	var modifiedCount, stagedCount, untrackedCount int
	if !wt.IsBare {
		status, err := git.GetWorktreeStatus(wt.Path)
		if err == nil && status != nil {
			modifiedCount = status.ModifiedCount
			stagedCount = status.StagedCount
			untrackedCount = status.UntrackedCount
		}
	}

	// Build metadata
	metadata := &WorktreeItemData{
		Path:           wt.Path,
		Branch:         wt.Branch,
		CommitHash:     wt.CommitHash,
		IsBare:         wt.IsBare,
		IsDetached:     wt.IsDetached,
		ModifiedCount:  modifiedCount,
		StagedCount:    stagedCount,
		UntrackedCount: untrackedCount,
	}

	// Build simple description for backwards compatibility
	var description string
	if wt.IsBare {
		description = "Bare repository"
	} else if wt.IsDetached {
		description = "Detached HEAD"
	} else if wt.Branch != "" {
		description = wt.Branch
	}

	return ListItem{
		ID:          wt.Path,
		Title:       wt.Name(),
		Description: description,
		Metadata:    metadata,
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
	case ConfirmDialogResultMsg:
		return a.handleConfirmDialogResult(msg)
	}

	// If confirm dialog is visible, route all key events to it
	if a.confirmDialog.Visible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// Allow Ctrl+C to quit even with dialog open
			if keyMsg.Type == tea.KeyCtrlC {
				a.quitting = true
				return a, tea.Quit
			}
			cmd := a.confirmDialog.Update(keyMsg)
			return a, cmd
		}
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
				case 'p':
					// Prune stale worktrees on Worktrees tab
					if a.tabs.Active() == TabWorktrees && !git.IsNotGitRepoError(a.gitError) {
						a.confirmDialog.SetConfirmLabel("Prune")
						a.confirmDialog.SetForceOption(false)
						a.confirmDialog.ShowWithData(
							"Prune Stale Worktrees?",
							"This will remove worktree entries whose directories no longer exist.",
							"prune",
						)
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
		// Open the worktree in a new terminal or provide cd command
		worktreePath := msg.Item.ID // ID is the worktree path
		opener := git.NewTerminalOpener()
		result, err := opener.OpenWorktree(worktreePath)
		if err != nil {
			cmd := a.feedback.ShowError("Failed to open worktree: " + err.Error())
			return a, cmd
		}

		// Show appropriate feedback based on result
		if result.Success {
			cmd := a.feedback.ShowSuccess(result.Message)
			return a, cmd
		}
		// Fallback: show the cd command to the user
		cmd := a.feedback.ShowInfo(result.Message)
		return a, cmd
	case "cd":
		// Get the cd command for the worktree
		worktreePath := msg.Item.ID
		cdCommand := git.GetCDCommand(worktreePath)
		cmd := a.feedback.ShowInfo("Copy: " + cdCommand)
		return a, cmd
	case "delete":
		// Show confirmation dialog for delete action
		a.confirmDialog.SetConfirmLabel("Delete")
		a.confirmDialog.SetForceOption(true)
		a.confirmDialog.ShowDanger(
			"Delete Worktree?",
			"This will remove the worktree '"+msg.Item.Title+"'.\nPath: "+msg.Item.ID,
			msg.Item,
		)
		return a, nil
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

	// Set target path and quit so shell wrapper can cd to it
	a.targetPath = msg.Result.Path
	a.quitting = true
	return a, tea.Quit
}

// handleConfirmDialogResult processes the result of a confirmation dialog.
func (a *App) handleConfirmDialogResult(msg ConfirmDialogResultMsg) (tea.Model, tea.Cmd) {
	if !msg.Confirmed {
		// User cancelled, nothing to do
		return a, nil
	}

	// Handle the confirmed action based on the data type
	if item, ok := msg.Data.(*ListItem); ok {
		// This is a worktree delete confirmation
		opts := git.RemoveWorktreeOptions{
			Path:  item.ID, // ID is the worktree path
			Force: msg.Force,
		}

		err := git.RemoveWorktree(a.repoPath, opts)
		if err != nil {
			cmd := a.feedback.ShowError("Failed to remove worktree: " + err.Error())
			return a, cmd
		}

		// Refresh the worktree list
		a.loadWorktrees()

		cmd := a.feedback.ShowSuccess("Removed worktree: " + item.Title)
		return a, cmd
	}

	// Handle prune confirmation
	if action, ok := msg.Data.(string); ok && action == "prune" {
		output, err := git.PruneWorktrees(a.repoPath)
		if err != nil {
			cmd := a.feedback.ShowError("Failed to prune worktrees: " + err.Error())
			return a, cmd
		}

		// Refresh the worktree list
		a.loadWorktrees()

		// Show success message
		message := "Pruned stale worktrees"
		if output != "" {
			message = "Pruned: " + output
		}
		cmd := a.feedback.ShowSuccess(message)
		return a, cmd
	}

	return a, nil
}

// ConfirmDialog returns the confirmation dialog component for testing.
func (a *App) ConfirmDialog() *ConfirmDialog {
	return a.confirmDialog
}

// CreateForm returns the create form component for testing.
func (a *App) CreateForm() *CreateForm {
	return a.createForm
}

// TargetPath returns the path to cd to after quitting (for shell wrapper).
func (a *App) TargetPath() string {
	return a.targetPath
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
		// Suppress output when targetPath is set (shell wrapper will handle cd)
		if a.targetPath != "" {
			return ""
		}
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
	helpText := "↑/↓: navigate • Enter: action • n: new worktree • p: prune • Tab: switch tabs • q: quit"
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

	// If confirm dialog is visible, render it as an overlay
	if a.confirmDialog.Visible() {
		b.WriteString("\n\n")
		b.WriteString(a.confirmDialog.View())
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
