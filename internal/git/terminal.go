// Package git provides git operations for the worktree manager.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// TerminalOpener provides functionality to open worktrees in a new terminal window.
type TerminalOpener struct {
	// terminalCmd is the terminal emulator command to use.
	// If empty, will auto-detect based on environment.
	terminalCmd string
}

// NewTerminalOpener creates a new TerminalOpener with auto-detection.
func NewTerminalOpener() *TerminalOpener {
	return &TerminalOpener{}
}

// NewTerminalOpenerWithCmd creates a new TerminalOpener with a specific terminal command.
func NewTerminalOpenerWithCmd(cmd string) *TerminalOpener {
	return &TerminalOpener{terminalCmd: cmd}
}

// OpenWorktreeResult contains the result of opening a worktree.
type OpenWorktreeResult struct {
	// Success indicates if the terminal was opened successfully.
	Success bool
	// Method describes how the worktree was opened (e.g., "terminal", "cd command").
	Method string
	// Message is a user-friendly message about the result.
	Message string
	// CDCommand is the cd command that can be used to switch to the worktree.
	CDCommand string
}

// OpenWorktree opens a new terminal window at the specified worktree path.
// Falls back to providing a cd command if terminal opening fails.
func (t *TerminalOpener) OpenWorktree(path string) (*OpenWorktreeResult, error) {
	// Validate the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree path does not exist: %s", path)
	}

	cdCommand := fmt.Sprintf("cd %s", shellQuote(path))

	// Try to open terminal
	terminalCmd, args := t.detectTerminal()
	if terminalCmd != "" {
		if err := t.openTerminal(terminalCmd, args, path); err == nil {
			return &OpenWorktreeResult{
				Success:   true,
				Method:    "terminal",
				Message:   fmt.Sprintf("Opened terminal at %s", path),
				CDCommand: cdCommand,
			}, nil
		}
	}

	// Fallback: return cd command for user to copy/use
	return &OpenWorktreeResult{
		Success:   false,
		Method:    "cd_command",
		Message:   fmt.Sprintf("Use this command to switch: %s", cdCommand),
		CDCommand: cdCommand,
	}, nil
}

// detectTerminal detects the available terminal emulator.
// Returns the terminal command and arguments to open a new window at a specific directory.
func (t *TerminalOpener) detectTerminal() (string, []string) {
	// If a custom terminal is set, use it
	if t.terminalCmd != "" {
		return t.terminalCmd, nil
	}

	switch runtime.GOOS {
	case "darwin":
		return t.detectMacOSTerminal()
	case "linux":
		return t.detectLinuxTerminal()
	case "windows":
		return t.detectWindowsTerminal()
	default:
		return "", nil
	}
}

// detectMacOSTerminal detects available terminal on macOS.
func (t *TerminalOpener) detectMacOSTerminal() (string, []string) {
	// Check for common terminal emulators in order of preference
	// iTerm2, Alacritty, Kitty, Terminal.app
	terminals := []struct {
		check string
		cmd   string
		args  []string
	}{
		// iTerm2
		{"/Applications/iTerm.app", "open", []string{"-a", "iTerm"}},
		// Alacritty
		{"/Applications/Alacritty.app", "open", []string{"-a", "Alacritty", "--args", "--working-directory"}},
		// Kitty
		{"/Applications/kitty.app", "open", []string{"-a", "kitty", "--args", "--directory"}},
		// WezTerm
		{"/Applications/WezTerm.app", "open", []string{"-a", "WezTerm", "--args", "start", "--cwd"}},
		// Default to Terminal.app (always available)
		{"/System/Applications/Utilities/Terminal.app", "open", []string{"-a", "Terminal"}},
	}

	for _, term := range terminals {
		if _, err := os.Stat(term.check); err == nil {
			return term.cmd, term.args
		}
	}

	// Fallback to Terminal.app
	return "open", []string{"-a", "Terminal"}
}

// detectLinuxTerminal detects available terminal on Linux.
func (t *TerminalOpener) detectLinuxTerminal() (string, []string) {
	// Check for common terminal emulators in order of preference
	terminals := []struct {
		cmd  string
		args []string
	}{
		// GNOME Terminal
		{"gnome-terminal", []string{"--working-directory"}},
		// Konsole
		{"konsole", []string{"--workdir"}},
		// xfce4-terminal
		{"xfce4-terminal", []string{"--working-directory"}},
		// Alacritty
		{"alacritty", []string{"--working-directory"}},
		// Kitty
		{"kitty", []string{"--directory"}},
		// WezTerm
		{"wezterm", []string{"start", "--cwd"}},
		// Terminator
		{"terminator", []string{"--working-directory"}},
		// xterm (fallback, uses -e cd)
		{"xterm", []string{"-e", "cd"}},
	}

	for _, term := range terminals {
		if path, err := exec.LookPath(term.cmd); err == nil && path != "" {
			return term.cmd, term.args
		}
	}

	return "", nil
}

// detectWindowsTerminal detects available terminal on Windows.
func (t *TerminalOpener) detectWindowsTerminal() (string, []string) {
	// Check for Windows Terminal first
	if path, err := exec.LookPath("wt.exe"); err == nil && path != "" {
		return "wt.exe", []string{"-d"}
	}

	// PowerShell
	if path, err := exec.LookPath("pwsh.exe"); err == nil && path != "" {
		return "pwsh.exe", []string{"-NoExit", "-Command", "Set-Location"}
	}

	// Fallback to cmd.exe
	return "cmd.exe", []string{"/K", "cd /d"}
}

// openTerminal opens a terminal window at the specified path.
func (t *TerminalOpener) openTerminal(terminalCmd string, args []string, path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = t.buildMacOSCommand(terminalCmd, args, path)
	case "linux":
		cmd = t.buildLinuxCommand(terminalCmd, args, path)
	case "windows":
		cmd = t.buildWindowsCommand(terminalCmd, args, path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// buildMacOSCommand builds the command to open a terminal on macOS.
func (t *TerminalOpener) buildMacOSCommand(terminalCmd string, args []string, path string) *exec.Cmd {
	// For iTerm and Terminal.app, we use AppleScript to open in specific directory
	if terminalCmd == "open" && len(args) >= 2 {
		app := args[1]
		switch app {
		case "iTerm":
			// iTerm2 AppleScript
			script := fmt.Sprintf(`
				tell application "iTerm"
					create window with default profile
					tell current session of current window
						write text "cd %s && clear"
					end tell
				end tell
			`, shellQuote(path))
			return exec.Command("osascript", "-e", script)
		case "Terminal":
			// Terminal.app AppleScript
			script := fmt.Sprintf(`
				tell application "Terminal"
					do script "cd %s && clear"
					activate
				end tell
			`, shellQuote(path))
			return exec.Command("osascript", "-e", script)
		default:
			// Other apps: append path to args
			fullArgs := append(args, path)
			return exec.Command(terminalCmd, fullArgs...)
		}
	}

	// Direct command with path
	fullArgs := append(args, path)
	return exec.Command(terminalCmd, fullArgs...)
}

// buildLinuxCommand builds the command to open a terminal on Linux.
func (t *TerminalOpener) buildLinuxCommand(terminalCmd string, args []string, path string) *exec.Cmd {
	// Special handling for xterm
	if terminalCmd == "xterm" {
		return exec.Command(terminalCmd, "-e", "bash", "-c", fmt.Sprintf("cd %s && bash", shellQuote(path)))
	}

	// Standard: append path to args
	fullArgs := append(args, path)
	return exec.Command(terminalCmd, fullArgs...)
}

// buildWindowsCommand builds the command to open a terminal on Windows.
func (t *TerminalOpener) buildWindowsCommand(terminalCmd string, args []string, path string) *exec.Cmd {
	fullArgs := append(args, path)
	return exec.Command(terminalCmd, fullArgs...)
}

// shellQuote quotes a string for safe use in shell commands.
func shellQuote(s string) string {
	// If the string contains single quotes, use double quotes with escaping
	if strings.Contains(s, "'") {
		escaped := strings.ReplaceAll(s, "\"", "\\\"")
		escaped = strings.ReplaceAll(escaped, "$", "\\$")
		escaped = strings.ReplaceAll(escaped, "`", "\\`")
		escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
		return "\"" + escaped + "\""
	}
	// Use single quotes for safety
	return "'" + s + "'"
}

// GetCDCommand returns the cd command to switch to the worktree.
func GetCDCommand(path string) string {
	return fmt.Sprintf("cd %s", shellQuote(path))
}
