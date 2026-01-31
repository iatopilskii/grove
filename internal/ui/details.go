// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"fmt"
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
		content = d.renderItemDetails()
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

// renderItemDetails renders the detailed view for the selected item.
func (d *Details) renderItemDetails() string {
	// Title with primary color for emphasis
	titleStyle := lipgloss.NewStyle().
		Foreground(Colors.Text).
		Bold(true)
	title := titleStyle.Render(d.item.Title)

	// Label style for field names
	labelStyle := lipgloss.NewStyle().
		Foreground(Colors.TextMuted).
		Bold(true)

	// Value style for field values
	valueStyle := lipgloss.NewStyle().
		Foreground(Colors.Text)

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	// Check if we have worktree metadata
	if wtData, ok := d.item.Metadata.(*WorktreeItemData); ok && wtData != nil {
		// Show full path
		lines = append(lines, labelStyle.Render("Path"))
		lines = append(lines, valueStyle.Render(wtData.Path))
		lines = append(lines, "")

		// Show branch name
		if wtData.IsBare {
			lines = append(lines, labelStyle.Render("Type"))
			lines = append(lines, valueStyle.Render("Bare repository"))
		} else if wtData.IsDetached {
			lines = append(lines, labelStyle.Render("State"))
			lines = append(lines, valueStyle.Render("Detached HEAD"))
			if wtData.CommitHash != "" {
				lines = append(lines, "")
				lines = append(lines, labelStyle.Render("Commit"))
				lines = append(lines, valueStyle.Render(wtData.CommitHash))
			}
		} else {
			lines = append(lines, labelStyle.Render("Branch"))
			lines = append(lines, valueStyle.Render(wtData.Branch))
		}
		lines = append(lines, "")

		// Show status with modified/staged file counts
		if !wtData.IsBare {
			lines = append(lines, labelStyle.Render("Status"))
			statusLine := d.renderStatusLine(wtData)
			lines = append(lines, statusLine)
		}
	} else if d.item.Description != "" {
		// Fallback to simple description
		descStyle := lipgloss.NewStyle().
			Foreground(Colors.TextMuted)
		lines = append(lines, descStyle.Render(d.item.Description))
	}

	return strings.Join(lines, "\n")
}

// renderStatusLine renders the status line showing modified/staged/untracked counts.
func (d *Details) renderStatusLine(wtData *WorktreeItemData) string {
	// Style for clean status
	cleanStyle := lipgloss.NewStyle().
		Foreground(Colors.Success)

	// Style for modified files (yellow/warning)
	modifiedStyle := lipgloss.NewStyle().
		Foreground(Colors.Error)

	// Style for staged files (green/success)
	stagedStyle := lipgloss.NewStyle().
		Foreground(Colors.Success)

	// Style for untracked files (muted)
	untrackedStyle := lipgloss.NewStyle().
		Foreground(Colors.TextMuted)

	totalChanges := wtData.ModifiedCount + wtData.StagedCount + wtData.UntrackedCount
	if totalChanges == 0 {
		return cleanStyle.Render("✓ Clean")
	}

	var parts []string

	if wtData.StagedCount > 0 {
		parts = append(parts, stagedStyle.Render(fmt.Sprintf("%d staged", wtData.StagedCount)))
	}

	if wtData.ModifiedCount > 0 {
		parts = append(parts, modifiedStyle.Render(fmt.Sprintf("%d modified", wtData.ModifiedCount)))
	}

	if wtData.UntrackedCount > 0 {
		parts = append(parts, untrackedStyle.Render(fmt.Sprintf("%d untracked", wtData.UntrackedCount)))
	}

	return strings.Join(parts, ", ")
}
