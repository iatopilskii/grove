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

// TestListPageDown verifies PageDown moves selection by page size
func TestListPageDown(t *testing.T) {
	// Create a list with 20 items
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	list := NewList(items)
	list.SetSize(80, 5) // height=5 means page size of 5

	// PageDown should move by page size
	list.PageDown()
	if list.Selected() != 5 {
		t.Errorf("after PageDown with height=5, Selected() = %d, want 5", list.Selected())
	}

	// Another PageDown
	list.PageDown()
	if list.Selected() != 10 {
		t.Errorf("after second PageDown, Selected() = %d, want 10", list.Selected())
	}

	// PageDown near end should clamp to last item
	list.SetSelected(18)
	list.PageDown()
	if list.Selected() != 19 {
		t.Errorf("PageDown from 18 should clamp to 19, got %d", list.Selected())
	}
}

// TestListPageUp verifies PageUp moves selection by page size
func TestListPageUp(t *testing.T) {
	// Create a list with 20 items
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	list := NewList(items)
	list.SetSize(80, 5) // height=5 means page size of 5
	list.SetSelected(15)

	// PageUp should move by page size
	list.PageUp()
	if list.Selected() != 10 {
		t.Errorf("after PageUp with height=5 from 15, Selected() = %d, want 10", list.Selected())
	}

	// Another PageUp
	list.PageUp()
	if list.Selected() != 5 {
		t.Errorf("after second PageUp, Selected() = %d, want 5", list.Selected())
	}

	// PageUp near beginning should clamp to first item
	list.SetSelected(2)
	list.PageUp()
	if list.Selected() != 0 {
		t.Errorf("PageUp from 2 should clamp to 0, got %d", list.Selected())
	}
}

// TestListPageDownEmpty verifies PageDown on empty list doesn't panic
func TestListPageDownEmpty(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 10)
	list.PageDown() // Should not panic
	if list.Selected() != 0 {
		t.Errorf("PageDown on empty list: Selected() = %d, want 0", list.Selected())
	}
}

// TestListPageUpEmpty verifies PageUp on empty list doesn't panic
func TestListPageUpEmpty(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 10)
	list.PageUp() // Should not panic
	if list.Selected() != 0 {
		t.Errorf("PageUp on empty list: Selected() = %d, want 0", list.Selected())
	}
}

// TestListPageDownZeroHeight verifies PageDown with zero height uses fallback
func TestListPageDownZeroHeight(t *testing.T) {
	items := make([]ListItem, 10)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	list := NewList(items)
	// Don't set size (height=0)

	list.PageDown()
	// With zero height, should use fallback page size of 1 (moves at least one)
	if list.Selected() != 1 {
		t.Errorf("PageDown with zero height, Selected() = %d, want 1", list.Selected())
	}
}

// TestListPageUpZeroHeight verifies PageUp with zero height uses fallback
func TestListPageUpZeroHeight(t *testing.T) {
	items := make([]ListItem, 10)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	list := NewList(items)
	list.SetSelected(5)
	// Don't set size (height=0)

	list.PageUp()
	// With zero height, should use fallback page size of 1 (moves at least one)
	if list.Selected() != 4 {
		t.Errorf("PageUp with zero height from 5, Selected() = %d, want 4", list.Selected())
	}
}

// TestListUpdatePageKeys verifies PageUp/PageDown key handling
func TestListUpdatePageKeys(t *testing.T) {
	items := make([]ListItem, 20)
	for i := range items {
		items[i] = ListItem{ID: string(rune('a' + i)), Title: "Item"}
	}
	list := NewList(items)
	list.SetSize(80, 5)

	// PageDown key
	list.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	if list.Selected() != 5 {
		t.Errorf("after KeyPgDown, Selected() = %d, want 5", list.Selected())
	}

	// PageUp key
	list.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	if list.Selected() != 0 {
		t.Errorf("after KeyPgUp, Selected() = %d, want 0", list.Selected())
	}
}

// TestListMouseClickSelect verifies clicking on a list item selects it
func TestListMouseClickSelect(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)
	list.SetSize(80, 10)
	list.SetOffset(0, 2) // List starts at Y=2

	// Click on second item (each item is 1 line, so Y=3 is item index 1)
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      5,
		Y:      3, // Offset 2 + item index 1
	})
	if list.Selected() != 1 {
		t.Errorf("after mouse click on second item, Selected() = %d, want 1", list.Selected())
	}

	// Click on third item
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      5,
		Y:      4, // Offset 2 + item index 2
	})
	if list.Selected() != 2 {
		t.Errorf("after mouse click on third item, Selected() = %d, want 2", list.Selected())
	}
}

// TestListMouseClickOutOfBounds verifies clicking outside list bounds doesn't crash
func TestListMouseClickOutOfBounds(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
	}
	list := NewList(items)
	list.SetSize(80, 10)
	list.SetOffset(0, 2)

	initial := list.Selected()

	// Click above the list
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      5,
		Y:      0, // Above offset
	})
	if list.Selected() != initial {
		t.Errorf("click above list should not change selection, got %d want %d", list.Selected(), initial)
	}

	// Click below all items
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      5,
		Y:      100, // Way below list
	})
	if list.Selected() != initial {
		t.Errorf("click below list should not change selection, got %d want %d", list.Selected(), initial)
	}

	// Click to the left of list
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      -5,
		Y:      3,
	})
	// Should not crash
}

// TestListMouseClickEmpty verifies clicking on empty list doesn't crash
func TestListMouseClickEmpty(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 10)
	list.SetOffset(0, 2)

	// Click should not panic
	list.Update(tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      5,
		Y:      3,
	})
	if list.Selected() != 0 {
		t.Errorf("click on empty list should keep selection at 0, got %d", list.Selected())
	}
}

// TestListMouseWheelDown verifies mouse wheel scrolls down
func TestListMouseWheelDown(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)
	list.SetSize(80, 10)
	list.SetOffset(0, 2)

	list.Update(tea.MouseMsg{
		Type:   tea.MouseWheelDown,
		Button: tea.MouseButtonWheelDown,
		X:      5,
		Y:      3,
	})
	if list.Selected() != 1 {
		t.Errorf("after mouse wheel down, Selected() = %d, want 1", list.Selected())
	}
}

// TestListMouseWheelUp verifies mouse wheel scrolls up
func TestListMouseWheelUp(t *testing.T) {
	items := []ListItem{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
	list := NewList(items)
	list.SetSize(80, 10)
	list.SetOffset(0, 2)
	list.SetSelected(2) // Start at last item

	list.Update(tea.MouseMsg{
		Type:   tea.MouseWheelUp,
		Button: tea.MouseButtonWheelUp,
		X:      5,
		Y:      3,
	})
	if list.Selected() != 1 {
		t.Errorf("after mouse wheel up, Selected() = %d, want 1", list.Selected())
	}
}

// TestListMouseWheelEmpty verifies mouse wheel on empty list doesn't crash
func TestListMouseWheelEmpty(t *testing.T) {
	list := NewList(nil)
	list.SetSize(80, 10)
	list.SetOffset(0, 2)

	// Should not panic
	list.Update(tea.MouseMsg{
		Type:   tea.MouseWheelDown,
		Button: tea.MouseButtonWheelDown,
		X:      5,
		Y:      3,
	})
	list.Update(tea.MouseMsg{
		Type:   tea.MouseWheelUp,
		Button: tea.MouseButtonWheelUp,
		X:      5,
		Y:      3,
	})
}

// TestListSetOffset verifies SetOffset sets position correctly
func TestListSetOffset(t *testing.T) {
	list := NewList(nil)
	list.SetOffset(10, 20)

	if list.offsetX != 10 {
		t.Errorf("after SetOffset, offsetX = %d, want 10", list.offsetX)
	}
	if list.offsetY != 20 {
		t.Errorf("after SetOffset, offsetY = %d, want 20", list.offsetY)
	}
}

// TestListIsInBounds verifies IsInBounds correctly detects mouse position
func TestListIsInBounds(t *testing.T) {
	list := NewList(nil)
	list.SetSize(40, 10)
	list.SetOffset(5, 5)

	tests := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{"inside", 10, 8, true},
		{"top-left corner", 5, 5, true},
		{"bottom-right corner", 44, 14, true},
		{"above", 10, 4, false},
		{"below", 10, 16, false},
		{"left of", 4, 8, false},
		{"right of", 46, 8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := list.IsInBounds(tt.x, tt.y); got != tt.expected {
				t.Errorf("IsInBounds(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}
