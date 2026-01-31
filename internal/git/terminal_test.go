// Package git provides git operations for the worktree manager.
package git

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

// TestNewTerminalOpener verifies the constructor.
func TestNewTerminalOpener(t *testing.T) {
	opener := NewTerminalOpener()
	if opener == nil {
		t.Error("Expected non-nil TerminalOpener")
	}
	if opener.terminalCmd != "" {
		t.Errorf("Expected empty terminalCmd, got '%s'", opener.terminalCmd)
	}
}

// TestNewTerminalOpenerWithCmd verifies the constructor with custom command.
func TestNewTerminalOpenerWithCmd(t *testing.T) {
	opener := NewTerminalOpenerWithCmd("custom-terminal")
	if opener == nil {
		t.Error("Expected non-nil TerminalOpener")
	}
	if opener.terminalCmd != "custom-terminal" {
		t.Errorf("Expected terminalCmd 'custom-terminal', got '%s'", opener.terminalCmd)
	}
}

// TestOpenWorktreeInvalidPath tests opening a worktree with an invalid path.
func TestOpenWorktreeInvalidPath(t *testing.T) {
	opener := NewTerminalOpener()
	_, err := opener.OpenWorktree("/non/existent/path/12345")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' in error, got: %s", err.Error())
	}
}

// TestOpenWorktreeValidPath tests opening a worktree with a valid path.
func TestOpenWorktreeValidPath(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "terminaltest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	opener := NewTerminalOpener()
	result, err := opener.OpenWorktree(tmpDir)
	if err != nil {
		t.Fatalf("OpenWorktree failed: %v", err)
	}

	// The result should contain either terminal or cd_command method
	if result.Method != "terminal" && result.Method != "cd_command" {
		t.Errorf("Expected method 'terminal' or 'cd_command', got '%s'", result.Method)
	}

	// CDCommand should always be set
	if result.CDCommand == "" {
		t.Error("Expected CDCommand to be set")
	}
	if !strings.Contains(result.CDCommand, tmpDir) {
		t.Errorf("Expected CDCommand to contain path '%s', got '%s'", tmpDir, result.CDCommand)
	}
}

// TestOpenWorktreeResultFields tests the OpenWorktreeResult struct fields.
func TestOpenWorktreeResultFields(t *testing.T) {
	result := OpenWorktreeResult{
		Success:   true,
		Method:    "terminal",
		Message:   "Opened terminal at /path",
		CDCommand: "cd '/path'",
	}

	if !result.Success {
		t.Error("Expected Success true, got false")
	}
	if result.Method != "terminal" {
		t.Errorf("Expected Method 'terminal', got '%s'", result.Method)
	}
	if result.Message != "Opened terminal at /path" {
		t.Errorf("Expected Message 'Opened terminal at /path', got '%s'", result.Message)
	}
	if result.CDCommand != "cd '/path'" {
		t.Errorf("Expected CDCommand \"cd '/path'\", got '%s'", result.CDCommand)
	}
}

// TestDetectTerminal tests that terminal detection returns valid values.
func TestDetectTerminal(t *testing.T) {
	opener := NewTerminalOpener()
	cmd, args := opener.detectTerminal()

	// On all systems, we should get some terminal command
	switch runtime.GOOS {
	case "darwin":
		// macOS should return 'open' or similar
		if cmd == "" {
			t.Error("Expected terminal command on macOS, got empty")
		}
	case "linux":
		// Linux might not have a terminal, that's OK
		_ = cmd
	case "windows":
		// Windows should return cmd.exe at minimum
		if cmd == "" {
			t.Error("Expected terminal command on Windows, got empty")
		}
	}

	// args can be nil or empty, that's OK
	_ = args
}

// TestDetectTerminalWithCustomCmd tests that custom command overrides detection.
func TestDetectTerminalWithCustomCmd(t *testing.T) {
	opener := NewTerminalOpenerWithCmd("my-custom-terminal")
	cmd, args := opener.detectTerminal()

	if cmd != "my-custom-terminal" {
		t.Errorf("Expected 'my-custom-terminal', got '%s'", cmd)
	}
	if args != nil {
		t.Errorf("Expected nil args, got %v", args)
	}
}

// TestShellQuote tests the shell quoting function.
func TestShellQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/dir", "'/path/to/dir'"},
		{"/path/with spaces/dir", "'/path/with spaces/dir'"},
		{"/path/with$dollar", "'/path/with$dollar'"},
		{"/path/with`backtick", "'/path/with`backtick'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shellQuote(tt.input)
			if got != tt.expected {
				t.Errorf("shellQuote(%s) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

// TestShellQuoteWithSingleQuotes tests quoting paths with single quotes.
func TestShellQuoteWithSingleQuotes(t *testing.T) {
	input := "/path/with'quote"
	result := shellQuote(input)

	// Should use double quotes when path contains single quotes
	if !strings.HasPrefix(result, "\"") {
		t.Errorf("Expected double-quoted string for path with single quotes, got: %s", result)
	}
}

// TestGetCDCommand tests the GetCDCommand helper function.
func TestGetCDCommand(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/dir", "cd '/path/to/dir'"},
		{"/path/with spaces", "cd '/path/with spaces'"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := GetCDCommand(tt.path)
			if got != tt.expected {
				t.Errorf("GetCDCommand(%s) = %s, want %s", tt.path, got, tt.expected)
			}
		})
	}
}

// TestOpenWorktreeFallbackToCDCommand tests that we get a CD command when terminal can't open.
func TestOpenWorktreeFallbackToCDCommand(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "terminaltest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use a non-existent terminal to force fallback
	opener := NewTerminalOpenerWithCmd("")
	// The opener.terminalCmd is empty, so it will try to detect.
	// If terminal opening fails, it should fall back to cd command.

	result, err := opener.OpenWorktree(tmpDir)
	if err != nil {
		t.Fatalf("OpenWorktree failed: %v", err)
	}

	// Either terminal worked or we got cd_command fallback
	if result.Method != "terminal" && result.Method != "cd_command" {
		t.Errorf("Expected method 'terminal' or 'cd_command', got '%s'", result.Method)
	}

	// Message should be set
	if result.Message == "" {
		t.Error("Expected non-empty Message")
	}

	// CDCommand should always contain the path
	if !strings.Contains(result.CDCommand, tmpDir) {
		t.Errorf("Expected CDCommand to contain path '%s', got '%s'", tmpDir, result.CDCommand)
	}
}

// TestBuildMacOSCommand tests building macOS commands (only runs on macOS).
func TestBuildMacOSCommand(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS system")
	}

	opener := NewTerminalOpener()
	path := "/test/path"

	// Test iTerm command
	cmd := opener.buildMacOSCommand("open", []string{"-a", "iTerm"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for iTerm")
	}

	// Test Terminal command
	cmd = opener.buildMacOSCommand("open", []string{"-a", "Terminal"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for Terminal")
	}
}

// TestBuildLinuxCommand tests building Linux commands (only runs on Linux).
func TestBuildLinuxCommand(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test on non-Linux system")
	}

	opener := NewTerminalOpener()
	path := "/test/path"

	// Test gnome-terminal command
	cmd := opener.buildLinuxCommand("gnome-terminal", []string{"--working-directory"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for gnome-terminal")
	}

	// Test xterm special handling
	cmd = opener.buildLinuxCommand("xterm", []string{"-e", "cd"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for xterm")
	}
}

// TestBuildWindowsCommand tests building Windows commands (only runs on Windows).
func TestBuildWindowsCommand(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows system")
	}

	opener := NewTerminalOpener()
	path := "C:\\test\\path"

	// Test wt.exe command
	cmd := opener.buildWindowsCommand("wt.exe", []string{"-d"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for wt.exe")
	}

	// Test cmd.exe command
	cmd = opener.buildWindowsCommand("cmd.exe", []string{"/K", "cd /d"}, path)
	if cmd == nil {
		t.Error("Expected non-nil command for cmd.exe")
	}
}

// TestDetectMacOSTerminal tests macOS terminal detection (only runs on macOS).
func TestDetectMacOSTerminal(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS system")
	}

	opener := NewTerminalOpener()
	cmd, args := opener.detectMacOSTerminal()

	// Should always return something (at least Terminal.app)
	if cmd == "" {
		t.Error("Expected terminal command on macOS, got empty")
	}
	if args == nil || len(args) == 0 {
		t.Error("Expected terminal args on macOS, got empty")
	}
}

// TestDetectWindowsTerminal tests Windows terminal detection (only runs on Windows).
func TestDetectWindowsTerminal(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows system")
	}

	opener := NewTerminalOpener()
	cmd, args := opener.detectWindowsTerminal()

	// Should always return something (at least cmd.exe)
	if cmd == "" {
		t.Error("Expected terminal command on Windows, got empty")
	}
	if args == nil || len(args) == 0 {
		t.Error("Expected terminal args on Windows, got empty")
	}
}
