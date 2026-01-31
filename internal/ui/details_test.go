package ui

import (
	"strings"
	"testing"
)

// TestNewDetails verifies that NewDetails returns a properly initialized Details pane
func TestNewDetails(t *testing.T) {
	details := NewDetails()
	if details == nil {
		t.Fatal("NewDetails() returned nil")
	}
}

// TestDetailsSetItem verifies SetItem updates the displayed item
func TestDetailsSetItem(t *testing.T) {
	details := NewDetails()
	item := &ListItem{ID: "1", Title: "Test Item", Description: "Test Description"}
	details.SetItem(item)

	if details.item != item {
		t.Error("SetItem did not set the item correctly")
	}
}

// TestDetailsSetItemNil verifies SetItem handles nil
func TestDetailsSetItemNil(t *testing.T) {
	details := NewDetails()
	details.SetItem(nil) // Should not panic

	if details.item != nil {
		t.Error("SetItem(nil) should set item to nil")
	}
}

// TestDetailsSetSize verifies SetSize updates dimensions
func TestDetailsSetSize(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)

	if details.width != 60 {
		t.Errorf("SetSize width = %d, want 60", details.width)
	}
	if details.height != 15 {
		t.Errorf("SetSize height = %d, want 15", details.height)
	}
}

// TestDetailsViewEmpty verifies View with no item shows placeholder
func TestDetailsViewEmpty(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)
	view := details.View()

	if view == "" {
		t.Error("View() with no item should not return empty string")
	}
}

// TestDetailsViewShowsTitle verifies View displays item title
func TestDetailsViewShowsTitle(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)
	details.SetItem(&ListItem{ID: "1", Title: "My Title", Description: "Desc"})
	view := details.View()

	if !strings.Contains(view, "My Title") {
		t.Errorf("View() does not contain item title 'My Title'")
	}
}

// TestDetailsViewShowsDescription verifies View displays item description
func TestDetailsViewShowsDescription(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)
	details.SetItem(&ListItem{ID: "1", Title: "Title", Description: "My Description"})
	view := details.View()

	if !strings.Contains(view, "My Description") {
		t.Errorf("View() does not contain item description 'My Description'")
	}
}

// TestDetailsViewUpdatesWithItem verifies View changes when item changes
func TestDetailsViewUpdatesWithItem(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)

	details.SetItem(&ListItem{ID: "1", Title: "First"})
	view1 := details.View()

	details.SetItem(&ListItem{ID: "2", Title: "Second"})
	view2 := details.View()

	if view1 == view2 {
		t.Error("View() should change when item changes")
	}
}

// TestDetailsViewWithBorder verifies View includes a border
func TestDetailsViewWithBorder(t *testing.T) {
	details := NewDetails()
	details.SetSize(60, 15)
	details.SetItem(&ListItem{ID: "1", Title: "Test"})
	view := details.View()

	// Border characters should be present (box drawing characters)
	hasBorder := strings.Contains(view, "─") || strings.Contains(view, "│") ||
		strings.Contains(view, "┌") || strings.Contains(view, "└")
	if !hasBorder {
		t.Error("View() should include border characters")
	}
}

// TestDetailsItem verifies Item getter returns current item
func TestDetailsItem(t *testing.T) {
	details := NewDetails()

	if details.Item() != nil {
		t.Error("Item() on new Details should return nil")
	}

	item := &ListItem{ID: "1", Title: "Test"}
	details.SetItem(item)

	if details.Item() != item {
		t.Error("Item() should return the set item")
	}
}
