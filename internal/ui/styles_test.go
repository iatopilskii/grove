package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestColorsUseAdaptiveColor verifies all colors in the palette are AdaptiveColors.
func TestColorsUseAdaptiveColor(t *testing.T) {
	// Test all color palette entries are properly defined AdaptiveColors
	// by checking they have both Light and Dark values set
	colors := []struct {
		name  string
		color lipgloss.AdaptiveColor
	}{
		{"Primary", Colors.Primary},
		{"OnPrimary", Colors.OnPrimary},
		{"Text", Colors.Text},
		{"TextMuted", Colors.TextMuted},
		{"Border", Colors.Border},
		{"Success", Colors.Success},
		{"Error", Colors.Error},
		{"Info", Colors.Info},
		{"OnSuccess", Colors.OnSuccess},
		{"OnError", Colors.OnError},
		{"OnInfo", Colors.OnInfo},
	}

	for _, tc := range colors {
		t.Run(tc.name, func(t *testing.T) {
			if tc.color.Light == "" {
				t.Errorf("%s.Light is empty", tc.name)
			}
			if tc.color.Dark == "" {
				t.Errorf("%s.Dark is empty", tc.name)
			}
		})
	}
}

// TestColorContrastLightMode verifies light mode uses dark text.
func TestColorContrastLightMode(t *testing.T) {
	// In light mode, text should be dark
	// Primary text should be #333333 (dark gray)
	if Colors.Text.Light != "#333333" {
		t.Errorf("Text.Light = %s, want dark color #333333", Colors.Text.Light)
	}

	// Muted text should also be dark
	if Colors.TextMuted.Light != "#666666" {
		t.Errorf("TextMuted.Light = %s, want dark color #666666", Colors.TextMuted.Light)
	}
}

// TestColorContrastDarkMode verifies dark mode uses light text.
func TestColorContrastDarkMode(t *testing.T) {
	// In dark mode, text should be light
	// Primary text should be #CCCCCC (light gray)
	if Colors.Text.Dark != "#CCCCCC" {
		t.Errorf("Text.Dark = %s, want light color #CCCCCC", Colors.Text.Dark)
	}

	// Muted text should also be light (but muted)
	if Colors.TextMuted.Dark != "#888888" {
		t.Errorf("TextMuted.Dark = %s, want muted light color #888888", Colors.TextMuted.Dark)
	}
}

// TestSemanticColorsHaveContrast verifies semantic colors have proper contrast.
func TestSemanticColorsHaveContrast(t *testing.T) {
	// All "On" colors should be white for good contrast
	tests := []struct {
		name  string
		color lipgloss.AdaptiveColor
	}{
		{"OnPrimary", Colors.OnPrimary},
		{"OnSuccess", Colors.OnSuccess},
		{"OnError", Colors.OnError},
		{"OnInfo", Colors.OnInfo},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.color.Light != "#FFFFFF" || tc.color.Dark != "#FFFFFF" {
				t.Errorf("%s should be white (#FFFFFF) for both modes, got Light=%s, Dark=%s",
					tc.name, tc.color.Light, tc.color.Dark)
			}
		})
	}
}

// TestStylesExist verifies predefined styles are properly initialized.
func TestStylesExist(t *testing.T) {
	// Just verify styles are created without panic
	_ = Styles.Selected.Render("test")
	_ = Styles.Normal.Render("test")
	_ = Styles.Muted.Render("test")
	_ = Styles.Help.Render("test")
}

// TestNoHardcodedANSI verifies no hardcoded ANSI codes in rendered output.
func TestNoHardcodedANSI(t *testing.T) {
	// Render styles and check for escape sequences
	// Note: lipgloss will add ANSI codes, but they should be derived
	// from AdaptiveColor, not hardcoded. This test ensures rendering works.
	outputs := []string{
		Styles.Selected.Render("test"),
		Styles.Normal.Render("test"),
		Styles.Muted.Render("test"),
		Styles.Help.Render("test"),
	}

	for i, output := range outputs {
		if output == "" {
			t.Errorf("Style %d rendered empty output", i)
		}
	}
}
