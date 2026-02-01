// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ilatopilskij/grove/internal/config"
)

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

// Borders defines thin (single-line) border styles for a minimal visual design.
// Uses standard box-drawing characters for clean appearance.
var Borders = struct {
	// Thin uses single-line box drawing characters (─│┌┐└┘)
	Thin lipgloss.Border
	// Rounded uses single-line with rounded corners (╭╮╰╯)
	Rounded lipgloss.Border
}{
	Thin:    lipgloss.NormalBorder(),
	Rounded: lipgloss.RoundedBorder(),
}

// Padding defines consistent padding values used throughout the UI.
// Using consistent values ensures visual harmony.
var Padding = struct {
	None   int
	Small  int
	Medium int
}{
	None:   0,
	Small:  1,
	Medium: 2,
}

// FocusIndicator defines the focus/selection indicator styling.
// Uses a subtle symbol that is visible but not overly prominent.
var FocusIndicator = struct {
	// Symbol shows for the focused/selected item
	Symbol string
	// SymbolInactive is whitespace of same width for alignment
	SymbolInactive string
}{
	Symbol:         "▸ ",
	SymbolInactive: "  ",
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
	// ListItem styles for list view items
	ListItem struct {
		Selected lipgloss.Style
		Normal   lipgloss.Style
	}
	// Box style for bordered containers
	Box lipgloss.Style
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

	ListItem: struct {
		Selected lipgloss.Style
		Normal   lipgloss.Style
	}{
		Selected: lipgloss.NewStyle().
			Foreground(Colors.Primary).
			Bold(true).
			PaddingLeft(0).
			PaddingRight(1),
		Normal: lipgloss.NewStyle().
			Foreground(Colors.Text).
			PaddingLeft(0).
			PaddingRight(1),
	},

	Box: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(0, 1),
}

// configToAdaptive converts a config.AdaptiveColor to lipgloss.AdaptiveColor.
func configToAdaptive(c config.AdaptiveColor) lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: c.Light, Dark: c.Dark}
}

// ApplyThemeConfig applies a configuration's theme colors to the global Colors
// and regenerates Styles to use those colors. This should be called during
// application initialization, before any UI rendering.
func ApplyThemeConfig(cfg config.Config) {
	// Update Colors from config
	Colors.Primary = configToAdaptive(cfg.Theme.Colors.Primary)
	Colors.OnPrimary = configToAdaptive(cfg.Theme.Colors.OnPrimary)
	Colors.Text = configToAdaptive(cfg.Theme.Colors.Text)
	Colors.TextMuted = configToAdaptive(cfg.Theme.Colors.TextMuted)
	Colors.Border = configToAdaptive(cfg.Theme.Colors.Border)
	Colors.Success = configToAdaptive(cfg.Theme.Colors.Success)
	Colors.Error = configToAdaptive(cfg.Theme.Colors.Error)
	Colors.Info = configToAdaptive(cfg.Theme.Colors.Info)
	Colors.OnSuccess = configToAdaptive(cfg.Theme.Colors.OnSuccess)
	Colors.OnError = configToAdaptive(cfg.Theme.Colors.OnError)
	Colors.OnInfo = configToAdaptive(cfg.Theme.Colors.OnInfo)

	// Regenerate Styles with new colors
	rebuildStyles()
}

// rebuildStyles recreates all styles using the current Colors values.
// This is called after Colors is updated from configuration.
func rebuildStyles() {
	Styles.Selected = lipgloss.NewStyle().
		Background(Colors.Primary).
		Foreground(Colors.OnPrimary).
		Bold(true).
		Padding(0, 1)

	Styles.Normal = lipgloss.NewStyle().
		Foreground(Colors.Text).
		Padding(0, 1)

	Styles.Muted = lipgloss.NewStyle().
		Foreground(Colors.TextMuted).
		Italic(true)

	Styles.Help = lipgloss.NewStyle().
		Foreground(Colors.TextMuted)

	Styles.ListItem.Selected = lipgloss.NewStyle().
		Foreground(Colors.Primary).
		Bold(true).
		PaddingLeft(0).
		PaddingRight(1)

	Styles.ListItem.Normal = lipgloss.NewStyle().
		Foreground(Colors.Text).
		PaddingLeft(0).
		PaddingRight(1)

	Styles.Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Colors.Border).
		Padding(0, 1)
}

// LoadAndApplyTheme loads the theme configuration from the default path
// and applies it to the global styles. Returns any error encountered
// while loading (invalid YAML), but always applies valid defaults.
func LoadAndApplyTheme() error {
	cfg, err := config.LoadConfig(config.DefaultConfigPath())
	ApplyThemeConfig(cfg)
	return err
}
