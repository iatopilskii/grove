// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Details is the details pane component that shows information about the selected item.
type Details struct {
	item   *ListItem
	width  int
	height int
}

// NewDetails creates a new details pane.
func NewDetails() *Details {
	return &Details{}
}

// Item returns the currently displayed item.
func (d *Details) Item() *ListItem {
	return d.item
}

// SetItem sets the item to display.
func (d *Details) SetItem(item *ListItem) {
	d.item = item
}

// SetSize sets the details pane dimensions.
func (d *Details) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// View renders the details pane.
func (d *Details) View() string {
	// Calculate inner dimensions (accounting for border)
	innerWidth := d.width - 2
	if innerWidth < 0 {
		innerWidth = 0
	}
	innerHeight := d.height - 2
	if innerHeight < 0 {
		innerHeight = 0
	}

	var content string
	if d.item == nil {
		// Show placeholder when no item selected
		content = Styles.Muted.Render("Select an item to view details")
	} else {
		// Title with primary color for emphasis
		titleStyle := lipgloss.NewStyle().
			Foreground(Colors.Text).
			Bold(true)
		title := titleStyle.Render(d.item.Title)

		// Description with muted color
		descStyle := lipgloss.NewStyle().
			Foreground(Colors.TextMuted)

		var descLines []string
		descLines = append(descLines, title)
		descLines = append(descLines, "")

		if d.item.Description != "" {
			descLines = append(descLines, descStyle.Render(d.item.Description))
		}

		content = strings.Join(descLines, "\n")
	}

	// Use centralized box style with thin rounded border
	boxStyle := Styles.Box

	if innerWidth > 0 {
		boxStyle = boxStyle.Width(innerWidth)
	}
	if innerHeight > 0 {
		boxStyle = boxStyle.Height(innerHeight)
	}

	return boxStyle.Render(content)
}
