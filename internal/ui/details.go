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
	// Border style
	borderColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	titleColor := lipgloss.AdaptiveColor{Light: "#333333", Dark: "#EEEEEE"}
	descColor := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#AAAAAA"}
	placeholderColor := lipgloss.AdaptiveColor{Light: "#888888", Dark: "#666666"}

	// Calculate inner dimensions
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
		placeholderStyle := lipgloss.NewStyle().
			Foreground(placeholderColor).
			Italic(true)
		content = placeholderStyle.Render("Select an item to view details")
	} else {
		// Title
		titleStyle := lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true)
		title := titleStyle.Render(d.item.Title)

		// Description
		descStyle := lipgloss.NewStyle().
			Foreground(descColor)

		var descLines []string
		descLines = append(descLines, title)
		descLines = append(descLines, "")

		if d.item.Description != "" {
			descLines = append(descLines, descStyle.Render(d.item.Description))
		}

		content = strings.Join(descLines, "\n")
	}

	// Border style using rounded corners
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	if innerWidth > 0 {
		boxStyle = boxStyle.Width(innerWidth)
	}
	if innerHeight > 0 {
		boxStyle = boxStyle.Height(innerHeight)
	}

	return boxStyle.Render(content)
}
