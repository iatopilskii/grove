package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTabString(t *testing.T) {
	tests := []struct {
		tab      Tab
		expected string
	}{
		{TabWorktrees, "Worktrees"},
		{TabBranches, "Branches"},
		{TabSettings, "Settings"},
		{Tab(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.tab.String(); got != tt.expected {
				t.Errorf("Tab.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestNewTabs(t *testing.T) {
	tabs := NewTabs()
	if tabs == nil {
		t.Fatal("NewTabs() returned nil")
	}
	if tabs.Active() != TabWorktrees {
		t.Errorf("NewTabs() active = %v, want TabWorktrees", tabs.Active())
	}
}

func TestTabsActive(t *testing.T) {
	tabs := NewTabs()
	if tabs.Active() != TabWorktrees {
		t.Errorf("initial Active() = %v, want TabWorktrees", tabs.Active())
	}
}

func TestTabsSetActive(t *testing.T) {
	tests := []struct {
		name     string
		setTab   Tab
		expected Tab
	}{
		{"set to worktrees", TabWorktrees, TabWorktrees},
		{"set to branches", TabBranches, TabBranches},
		{"set to settings", TabSettings, TabSettings},
		{"negative tab ignored", Tab(-1), TabWorktrees},
		{"out of bounds tab ignored", Tab(99), TabWorktrees},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := NewTabs()
			tabs.SetActive(tt.setTab)
			if tabs.Active() != tt.expected {
				t.Errorf("SetActive(%v) resulted in Active() = %v, want %v",
					tt.setTab, tabs.Active(), tt.expected)
			}
		})
	}
}

func TestTabsNext(t *testing.T) {
	tabs := NewTabs()

	// Start at Worktrees (0)
	if tabs.Active() != TabWorktrees {
		t.Fatalf("expected initial tab to be TabWorktrees")
	}

	// Next should go to Branches (1)
	tabs.Next()
	if tabs.Active() != TabBranches {
		t.Errorf("after first Next(), Active() = %v, want TabBranches", tabs.Active())
	}

	// Next should go to Settings (2)
	tabs.Next()
	if tabs.Active() != TabSettings {
		t.Errorf("after second Next(), Active() = %v, want TabSettings", tabs.Active())
	}

	// Next should wrap to Worktrees (0)
	tabs.Next()
	if tabs.Active() != TabWorktrees {
		t.Errorf("after third Next(), Active() = %v, want TabWorktrees (wrap)", tabs.Active())
	}
}

func TestTabsPrev(t *testing.T) {
	tabs := NewTabs()

	// Start at Worktrees (0)
	if tabs.Active() != TabWorktrees {
		t.Fatalf("expected initial tab to be TabWorktrees")
	}

	// Prev should wrap to Settings (2)
	tabs.Prev()
	if tabs.Active() != TabSettings {
		t.Errorf("after first Prev(), Active() = %v, want TabSettings (wrap)", tabs.Active())
	}

	// Prev should go to Branches (1)
	tabs.Prev()
	if tabs.Active() != TabBranches {
		t.Errorf("after second Prev(), Active() = %v, want TabBranches", tabs.Active())
	}

	// Prev should go to Worktrees (0)
	tabs.Prev()
	if tabs.Active() != TabWorktrees {
		t.Errorf("after third Prev(), Active() = %v, want TabWorktrees", tabs.Active())
	}
}

func TestTabsUpdateTabKey(t *testing.T) {
	tabs := NewTabs()

	// Tab key should move to next
	tabs.Update(tea.KeyMsg{Type: tea.KeyTab})
	if tabs.Active() != TabBranches {
		t.Errorf("after Tab key, Active() = %v, want TabBranches", tabs.Active())
	}
}

func TestTabsUpdateShiftTabKey(t *testing.T) {
	tabs := NewTabs()

	// Shift+Tab should move to previous (wrap to Settings)
	tabs.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if tabs.Active() != TabSettings {
		t.Errorf("after Shift+Tab key, Active() = %v, want TabSettings", tabs.Active())
	}
}

func TestTabsUpdateOtherKeys(t *testing.T) {
	tabs := NewTabs()

	// Other keys should not change active tab
	tabs.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if tabs.Active() != TabWorktrees {
		t.Errorf("after Enter key, Active() = %v, want TabWorktrees (unchanged)", tabs.Active())
	}

	tabs.Update(tea.KeyMsg{Type: tea.KeyUp})
	if tabs.Active() != TabWorktrees {
		t.Errorf("after Up key, Active() = %v, want TabWorktrees (unchanged)", tabs.Active())
	}
}

func TestTabsView(t *testing.T) {
	tabs := NewTabs()
	view := tabs.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

	// Check that all tab names appear in the view
	for i := Tab(0); i < TabCount; i++ {
		if !strings.Contains(view, i.String()) {
			t.Errorf("View() does not contain tab name %q", i.String())
		}
	}
}

func TestTabsViewContainsBorder(t *testing.T) {
	tabs := NewTabs()
	view := tabs.View()

	// Should contain a horizontal border line
	if !strings.Contains(view, "─") {
		t.Error("View() does not contain border character")
	}
}

func TestTabsViewUpdatesWithActiveTab(t *testing.T) {
	// This test ensures the view changes based on active tab
	// The active tab should be visually different (we can't easily test colors,
	// but we can verify the render doesn't crash and produces output)
	tabs := NewTabs()

	for i := Tab(0); i < TabCount; i++ {
		tabs.SetActive(i)
		view := tabs.View()
		if view == "" {
			t.Errorf("View() returned empty string for active tab %v", i)
		}
		// Ensure the active tab name is present
		if !strings.Contains(view, i.String()) {
			t.Errorf("View() does not contain active tab name %q", i.String())
		}
	}
}

func TestTabsSetWidth(t *testing.T) {
	tabs := NewTabs()
	tabs.SetWidth(80)

	view := tabs.View()
	if view == "" {
		t.Error("View() returned empty string after SetWidth")
	}
}

func TestTabCount(t *testing.T) {
	// Ensure TabCount matches the expected number of tabs
	if TabCount != 3 {
		t.Errorf("TabCount = %d, want 3", TabCount)
	}
}
