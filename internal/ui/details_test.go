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

// TestDetailsViewWithWorktreeMetadata verifies View displays worktree metadata
func TestDetailsViewWithWorktreeMetadata(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "my-feature",
		Metadata: &WorktreeItemData{
			Path:           "/path/to/worktree",
			Branch:         "feature-branch",
			CommitHash:     "abc1234",
			ModifiedCount:  2,
			StagedCount:    1,
			UntrackedCount: 3,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should contain title
	if !strings.Contains(view, "my-feature") {
		t.Error("View() should contain item title")
	}

	// Should contain path
	if !strings.Contains(view, "/path/to/worktree") {
		t.Error("View() should contain worktree path")
	}

	// Should contain branch name
	if !strings.Contains(view, "feature-branch") {
		t.Error("View() should contain branch name")
	}
}

// TestDetailsViewShowsCleanStatus verifies View displays clean status correctly
func TestDetailsViewShowsCleanStatus(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "clean-worktree",
		Metadata: &WorktreeItemData{
			Path:           "/path/to/worktree",
			Branch:         "main",
			ModifiedCount:  0,
			StagedCount:    0,
			UntrackedCount: 0,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should show clean status
	if !strings.Contains(view, "Clean") {
		t.Error("View() should show 'Clean' for worktree with no changes")
	}
}

// TestDetailsViewShowsModifiedCount verifies View displays modified file count
func TestDetailsViewShowsModifiedCount(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "modified-worktree",
		Metadata: &WorktreeItemData{
			Path:          "/path/to/worktree",
			Branch:        "main",
			ModifiedCount: 5,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should show modified count
	if !strings.Contains(view, "5 modified") {
		t.Error("View() should show modified count")
	}
}

// TestDetailsViewShowsStagedCount verifies View displays staged file count
func TestDetailsViewShowsStagedCount(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "staged-worktree",
		Metadata: &WorktreeItemData{
			Path:        "/path/to/worktree",
			Branch:      "main",
			StagedCount: 3,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should show staged count
	if !strings.Contains(view, "3 staged") {
		t.Error("View() should show staged count")
	}
}

// TestDetailsViewShowsUntrackedCount verifies View displays untracked file count
func TestDetailsViewShowsUntrackedCount(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "untracked-worktree",
		Metadata: &WorktreeItemData{
			Path:           "/path/to/worktree",
			Branch:         "main",
			UntrackedCount: 7,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should show untracked count
	if !strings.Contains(view, "7 untracked") {
		t.Error("View() should show untracked count")
	}
}

// TestDetailsViewShowsBareRepository verifies View handles bare repository correctly
func TestDetailsViewShowsBareRepository(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/repo.git",
		Title: "repo.git",
		Metadata: &WorktreeItemData{
			Path:   "/path/to/repo.git",
			IsBare: true,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should indicate bare repository
	if !strings.Contains(view, "Bare") {
		t.Error("View() should indicate bare repository")
	}
}

// TestDetailsViewShowsDetachedHead verifies View handles detached HEAD correctly
func TestDetailsViewShowsDetachedHead(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "detached-worktree",
		Metadata: &WorktreeItemData{
			Path:       "/path/to/worktree",
			CommitHash: "abc1234",
			IsDetached: true,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should indicate detached HEAD
	if !strings.Contains(view, "Detached") {
		t.Error("View() should indicate detached HEAD")
	}

	// Should show commit hash
	if !strings.Contains(view, "abc1234") {
		t.Error("View() should show commit hash for detached HEAD")
	}
}

// TestDetailsViewFallbackToDescription verifies View falls back to description without metadata
func TestDetailsViewFallbackToDescription(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:          "1",
		Title:       "Simple Item",
		Description: "Simple description text",
		Metadata:    nil, // No metadata
	}
	details.SetItem(item)
	view := details.View()

	// Should fall back to showing description
	if !strings.Contains(view, "Simple description text") {
		t.Error("View() should show description when no metadata")
	}
}

// TestDetailsViewMultipleStatusCounts verifies View displays multiple status counts
func TestDetailsViewMultipleStatusCounts(t *testing.T) {
	details := NewDetails()
	details.SetSize(80, 20)

	item := &ListItem{
		ID:    "/path/to/worktree",
		Title: "mixed-worktree",
		Metadata: &WorktreeItemData{
			Path:           "/path/to/worktree",
			Branch:         "main",
			ModifiedCount:  2,
			StagedCount:    3,
			UntrackedCount: 1,
		},
	}
	details.SetItem(item)
	view := details.View()

	// Should show all counts
	if !strings.Contains(view, "2 modified") {
		t.Error("View() should show modified count")
	}
	if !strings.Contains(view, "3 staged") {
		t.Error("View() should show staged count")
	}
	if !strings.Contains(view, "1 untracked") {
		t.Error("View() should show untracked count")
	}
}
