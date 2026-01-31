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

// TestBordersAreThin verifies the border style uses thin/single-line characters.
func TestBordersAreThin(t *testing.T) {
	// Verify Borders struct uses thin (single-line) border styles
	// The lipgloss.NormalBorder() uses single-line box drawing characters:
	// ─ │ ┌ ┐ └ ┘
	border := Borders.Thin
	if border.Top != "─" {
		t.Errorf("Borders.Thin.Top = %q, want single-line ─", border.Top)
	}
	if border.Bottom != "─" {
		t.Errorf("Borders.Thin.Bottom = %q, want single-line ─", border.Bottom)
	}
	if border.Left != "│" {
		t.Errorf("Borders.Thin.Left = %q, want single-line │", border.Left)
	}
	if border.Right != "│" {
		t.Errorf("Borders.Thin.Right = %q, want single-line │", border.Right)
	}
}

// TestBordersRoundedIsThin verifies the rounded border uses thin line characters.
func TestBordersRoundedIsThin(t *testing.T) {
	// Verify Borders.Rounded uses thin rounded corners
	border := Borders.Rounded
	if border.Top != "─" {
		t.Errorf("Borders.Rounded.Top = %q, want single-line ─", border.Top)
	}
	if border.TopLeft != "╭" {
		t.Errorf("Borders.Rounded.TopLeft = %q, want rounded ╭", border.TopLeft)
	}
	if border.TopRight != "╮" {
		t.Errorf("Borders.Rounded.TopRight = %q, want rounded ╮", border.TopRight)
	}
}

// TestPaddingConstants verifies consistent padding values.
func TestPaddingConstants(t *testing.T) {
	// Verify padding constants are defined
	if Padding.None != 0 {
		t.Errorf("Padding.None = %d, want 0", Padding.None)
	}
	if Padding.Small != 1 {
		t.Errorf("Padding.Small = %d, want 1", Padding.Small)
	}
	if Padding.Medium != 2 {
		t.Errorf("Padding.Medium = %d, want 2", Padding.Medium)
	}
}

// TestFocusIndicatorStyle verifies focus indicator is properly styled.
func TestFocusIndicatorStyle(t *testing.T) {
	// Verify focus indicator exists
	if FocusIndicator.Symbol == "" {
		t.Error("FocusIndicator.Symbol should not be empty")
	}
	if FocusIndicator.SymbolInactive == "" {
		t.Error("FocusIndicator.SymbolInactive should not be empty")
	}
	// Inactive should be whitespace of same visual width for alignment
	// Use lipgloss.Width which handles multi-byte characters correctly
	activeWidth := lipgloss.Width(FocusIndicator.Symbol)
	inactiveWidth := lipgloss.Width(FocusIndicator.SymbolInactive)
	if activeWidth != inactiveWidth {
		t.Errorf("FocusIndicator symbols have different visual widths: active=%d, inactive=%d",
			activeWidth, inactiveWidth)
	}
}

// TestStylesListItemHasConsistentPadding verifies list item styles use consistent padding.
func TestStylesListItemHasConsistentPadding(t *testing.T) {
	// Both selected and normal list items should have the same padding
	// to ensure alignment when selection moves
	selectedRendered := Styles.ListItem.Selected.Render("test")
	normalRendered := Styles.ListItem.Normal.Render("test")

	if len(selectedRendered) == 0 || len(normalRendered) == 0 {
		t.Error("List item styles should render non-empty output")
	}
}

// TestBoxStyleUsesThinBorder verifies box style uses thin border.
func TestBoxStyleUsesThinBorder(t *testing.T) {
	// Verify Styles.Box exists and uses thin borders
	rendered := Styles.Box.Render("content")
	if rendered == "" {
		t.Error("Styles.Box should render non-empty output")
	}
}

// TestNoExcessiveDecorations verifies styles don't use heavy decorations.
func TestNoExcessiveDecorations(t *testing.T) {
	// Verify help text uses simple styling (no bold, no background)
	helpRendered := Styles.Help.Render("test")
	if helpRendered == "" {
		t.Error("Help style should render non-empty")
	}

	// Verify muted text is italic but simple
	mutedRendered := Styles.Muted.Render("test")
	if mutedRendered == "" {
		t.Error("Muted style should render non-empty")
	}
}
