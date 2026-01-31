// Package ui provides the terminal user interface for the git worktree manager.
package ui

import "github.com/charmbracelet/lipgloss"

// Colors defines the adaptive color palette for the application.
// All colors use lipgloss.AdaptiveColor to automatically adjust
// for light and dark terminal themes.
var Colors = struct {
	// Primary colors (purple accent)
	Primary   lipgloss.AdaptiveColor
	OnPrimary lipgloss.AdaptiveColor

	// Text colors
	Text      lipgloss.AdaptiveColor
	TextMuted lipgloss.AdaptiveColor

	// Border colors
	Border lipgloss.AdaptiveColor

	// Semantic colors
	Success   lipgloss.AdaptiveColor
	Error     lipgloss.AdaptiveColor
	Info      lipgloss.AdaptiveColor
	OnSuccess lipgloss.AdaptiveColor
	OnError   lipgloss.AdaptiveColor
	OnInfo    lipgloss.AdaptiveColor
}{
	// Primary colors - purple accent for active/selected states
	Primary:   lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
	OnPrimary: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

	// Text colors
	Text:      lipgloss.AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"},
	TextMuted: lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"},

	// Border colors
	Border: lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},

	// Success (green)
	Success:   lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#4CAF50"},
	OnSuccess: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

	// Error (red)
	Error:   lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#EF5350"},
	OnError: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},

	// Info (blue)
	Info:   lipgloss.AdaptiveColor{Light: "#1565C0", Dark: "#42A5F5"},
	OnInfo: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
}

// Styles defines reusable lipgloss styles for the application.
var Styles = struct {
	// Selected item style
	Selected lipgloss.Style
	// Normal item style
	Normal lipgloss.Style
	// Muted/placeholder text style
	Muted lipgloss.Style
	// Help text style
	Help lipgloss.Style
}{
	Selected: lipgloss.NewStyle().
		Background(Colors.Primary).
		Foreground(Colors.OnPrimary).
		Bold(true).
		Padding(0, 1),

	Normal: lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1),

	Muted: lipgloss.NewStyle().
		Foreground(Colors.TextMuted).
		Italic(true),

	Help: lipgloss.NewStyle().
		Foreground(Colors.TextMuted),
}
