package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewList verifies that NewList returns a properly initialized List
func TestNewList(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1", Description: "Desc 1"},
		{ID: "2", Title: "Item 2", Description: "Desc 2"},
	}
	list := NewList(items)

	if list == nil {
		t.Fatal("NewList() returned nil")
	}
	if len(list.Items()) != 2 {
		t.Errorf("NewList() items count = %d, want 2", len(list.Items()))
	}
	if list.Selected() != 0 {
		t.Errorf("NewList() initial selection = %d, want 0", list.Selected())
	}
}

// TestNewListEmpty verifies NewList handles empty list
func TestNewListEmpty(t *testing.T) {
	list := NewList(nil)
	if list == nil {
		t.Fatal("NewList(nil) returned nil")
	}
	if len(list.Items()) != 0 {
		t.Errorf("NewList(nil) items count = %d, want 0", len(list.Items()))
	}
	if list.Selected() != 0 {
		t.Errorf("NewList(nil) selection = %d, want 0", list.Selected())
	}
}

// TestListSelected verifies Selected returns current selection index
func TestListSelected(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)

	if list.Selected() != 0 {
		t.Errorf("initial Selected() = %d, want 0", list.Selected())
	}

	list.SetSelected(2)
	if list.Selected() != 2 {
		t.Errorf("after SetSelected(2), Selected() = %d, want 2", list.Selected())
	}
}

// TestListSetSelected verifies SetSelected with bounds checking
func TestListSetSelected(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)

	tests := []struct {
		name     string
		setIndex int
		expected int
	}{
		{"valid index 0", 0, 0},
		{"valid index 1", 1, 1},
		{"valid index 2", 2, 2},
		{"negative index", -1, 0},     // Should stay at current or clamp to 0
		{"out of bounds", 10, 2},      // Should clamp to last valid index
		{"empty after clamp", 100, 2}, // Should clamp to last valid index
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list.SetSelected(0) // Reset to known state
			list.SetSelected(tt.setIndex)
			if list.Selected() != tt.expected {
				t.Errorf("SetSelected(%d) resulted in Selected() = %d, want %d",
					tt.setIndex, list.Selected(), tt.expected)
			}
		})
	}
}

// TestListMoveDown verifies moving selection down
func TestListMoveDown(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)

	list.MoveDown()
	if list.Selected() != 1 {
		t.Errorf("after MoveDown(), Selected() = %d, want 1", list.Selected())
	}

	list.MoveDown()
	if list.Selected() != 2 {
		t.Errorf("after second MoveDown(), Selected() = %d, want 2", list.Selected())
	}

	// Should stop at last item (no wrap)
	list.MoveDown()
	if list.Selected() != 2 {
		t.Errorf("after third MoveDown(), Selected() = %d, want 2 (should stop at boundary)", list.Selected())
	}
}

// TestListMoveUp verifies moving selection up
func TestListMoveUp(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)
	list.SetSelected(2) // Start at last item

	list.MoveUp()
	if list.Selected() != 1 {
		t.Errorf("after MoveUp(), Selected() = %d, want 1", list.Selected())
	}

	list.MoveUp()
	if list.Selected() != 0 {
		t.Errorf("after second MoveUp(), Selected() = %d, want 0", list.Selected())
	}

	// Should stop at first item (no wrap)
	list.MoveUp()
	if list.Selected() != 0 {
		t.Errorf("after third MoveUp(), Selected() = %d, want 0 (should stop at boundary)", list.Selected())
	}
}

// TestListMoveDownEmpty verifies MoveDown on empty list
func TestListMoveDownEmpty(t *testing.T) {
	list := NewList(nil)
	list.MoveDown() // Should not panic
	if list.Selected() != 0 {
		t.Errorf("MoveDown on empty list: Selected() = %d, want 0", list.Selected())
	}
}

// TestListMoveUpEmpty verifies MoveUp on empty list
func TestListMoveUpEmpty(t *testing.T) {
	list := NewList(nil)
	list.MoveUp() // Should not panic
	if list.Selected() != 0 {
		t.Errorf("MoveUp on empty list: Selected() = %d, want 0", list.Selected())
	}
}

// TestListSelectedItem verifies getting currently selected item
func TestListSelectedItem(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "First", Description: "First item"},
		{ID: "2", Title: "Second", Description: "Second item"},
	}
	list := NewList(items)

	item := list.SelectedItem()
	if item == nil {
		t.Fatal("SelectedItem() returned nil")
	}
	if item.Title != "First" {
		t.Errorf("SelectedItem().Title = %q, want %q", item.Title, "First")
	}

	list.SetSelected(1)
	item = list.SelectedItem()
	if item.Title != "Second" {
		t.Errorf("after SetSelected(1), SelectedItem().Title = %q, want %q", item.Title, "Second")
	}
}

// TestListSelectedItemEmpty verifies SelectedItem on empty list
func TestListSelectedItemEmpty(t *testing.T) {
	list := NewList(nil)
	item := list.SelectedItem()
	if item != nil {
		t.Errorf("SelectedItem() on empty list = %v, want nil", item)
	}
}

// TestListUpdateArrowKeys verifies arrow key handling
func TestListUpdateArrowKeys(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)

	// Down arrow
	list.Update(tea.KeyMsg{Type: tea.KeyDown})
	if list.Selected() != 1 {
		t.Errorf("after KeyDown, Selected() = %d, want 1", list.Selected())
	}

	// Up arrow
	list.Update(tea.KeyMsg{Type: tea.KeyUp})
	if list.Selected() != 0 {
		t.Errorf("after KeyUp, Selected() = %d, want 0", list.Selected())
	}
}

// TestListUpdateJKKeys verifies j/k vim-style navigation
func TestListUpdateJKKeys(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)

	// j key (down)
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if list.Selected() != 1 {
		t.Errorf("after 'j' key, Selected() = %d, want 1", list.Selected())
	}

	// k key (up)
	list.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if list.Selected() != 0 {
		t.Errorf("after 'k' key, Selected() = %d, want 0", list.Selected())
	}
}

// TestListSetSize verifies setting list dimensions
func TestListSetSize(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 20)

	if list.width != 80 {
		t.Errorf("after SetSize, width = %d, want 80", list.width)
	}
	if list.height != 20 {
		t.Errorf("after SetSize, height = %d, want 20", list.height)
	}
}

// TestListViewShowsItems verifies View renders items
func TestListViewShowsItems(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "First Item"},
		{ID: "2", Title: "Second Item"},
	}
	list := NewList(items)
	list.SetSize(80, 20)
	view := list.View()

	if !strings.Contains(view, "First Item") {
		t.Error("View() does not contain 'First Item'")
	}
	if !strings.Contains(view, "Second Item") {
		t.Error("View() does not contain 'Second Item'")
	}
}

// TestListViewHighlightsSelected verifies selected item is visually distinct
func TestListViewHighlightsSelected(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
	}
	list := NewList(items)
	list.SetSize(80, 20)

	// Change selection and verify view changes
	view1 := list.View()
	list.SetSelected(1)
	view2 := list.View()

	if view1 == view2 {
		t.Error("View() should change when selection changes")
	}
}

// TestListViewEmptyMessage verifies empty list shows message
func TestListViewEmptyMessage(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 20)
	view := list.View()

	if view == "" {
		t.Error("View() on empty list should not be empty")
	}
}

// TestListItems verifies Items returns all items
func TestListItems(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
	}
	list := NewList(items)

	got := list.Items()
	if len(got) != 2 {
		t.Errorf("Items() count = %d, want 2", len(got))
	}
}

// TestListSetItems verifies SetItems updates the list
func TestListSetItems(t *testing.T) {
	list := NewList(nil)

	newItems := []ListItem{
		{ID: "a", Title: "New Item A"},
		{ID: "b", Title: "New Item B"},
	}
	list.SetItems(newItems)

	if len(list.Items()) != 2 {
		t.Errorf("after SetItems, Items() count = %d, want 2", len(list.Items()))
	}
	if list.Items()[0].Title != "New Item A" {
		t.Errorf("after SetItems, first item title = %q, want %q", list.Items()[0].Title, "New Item A")
	}
}

// TestListSetItemsResetsSelection verifies selection is reset when items change
func TestListSetItemsResetsSelection(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)
	list.SetSelected(2)

	newItems := []ListItem{
		{ID: "a", Title: "New A"},
	}
	list.SetItems(newItems)

	// Selection should be clamped to valid range
	if list.Selected() != 0 {
		t.Errorf("after SetItems with shorter list, Selected() = %d, want 0", list.Selected())
	}
}
