package ui

import (
	"strings"
	"testing"
	"time"
)

// TestNewFeedback verifies the constructor creates a valid feedback component
func TestNewFeedback(t *testing.T) {
	fb := NewFeedback()
	if fb == nil {
		t.Fatal("NewFeedback() returned nil")
	}
	if fb.Visible() {
		t.Error("new Feedback should not be visible")
	}
	if fb.Message() != "" {
		t.Error("new Feedback should have empty message")
	}
}

// TestFeedbackShowSuccess verifies ShowSuccess shows a success message
func TestFeedbackShowSuccess(t *testing.T) {
	fb := NewFeedback()

	cmd := fb.ShowSuccess("Operation completed")

	if !fb.Visible() {
		t.Error("Feedback should be visible after ShowSuccess")
	}
	if fb.Message() != "Operation completed" {
		t.Errorf("Message() = %q, want %q", fb.Message(), "Operation completed")
	}
	if fb.Type() != FeedbackSuccess {
		t.Errorf("Type() = %v, want FeedbackSuccess", fb.Type())
	}
	if cmd == nil {
		t.Error("ShowSuccess should return a command to schedule clear")
	}
}

// TestFeedbackShowError verifies ShowError shows an error message
func TestFeedbackShowError(t *testing.T) {
	fb := NewFeedback()

	cmd := fb.ShowError("Something went wrong")

	if !fb.Visible() {
		t.Error("Feedback should be visible after ShowError")
	}
	if fb.Message() != "Something went wrong" {
		t.Errorf("Message() = %q, want %q", fb.Message(), "Something went wrong")
	}
	if fb.Type() != FeedbackError {
		t.Errorf("Type() = %v, want FeedbackError", fb.Type())
	}
	if cmd == nil {
		t.Error("ShowError should return a command to schedule clear")
	}
}

// TestFeedbackShowInfo verifies ShowInfo shows an info message
func TestFeedbackShowInfo(t *testing.T) {
	fb := NewFeedback()

	cmd := fb.ShowInfo("Information message")

	if !fb.Visible() {
		t.Error("Feedback should be visible after ShowInfo")
	}
	if fb.Message() != "Information message" {
		t.Errorf("Message() = %q, want %q", fb.Message(), "Information message")
	}
	if fb.Type() != FeedbackInfo {
		t.Errorf("Type() = %v, want FeedbackInfo", fb.Type())
	}
	if cmd == nil {
		t.Error("ShowInfo should return a command to schedule clear")
	}
}

// TestFeedbackClear verifies Clear hides the feedback
func TestFeedbackClear(t *testing.T) {
	fb := NewFeedback()
	fb.ShowSuccess("Test message")

	fb.Clear()

	if fb.Visible() {
		t.Error("Feedback should not be visible after Clear")
	}
	if fb.Message() != "" {
		t.Error("Message should be empty after Clear")
	}
}

// TestFeedbackUpdateClearMsg verifies ClearFeedbackMsg clears feedback
func TestFeedbackUpdateClearMsg(t *testing.T) {
	fb := NewFeedback()
	fb.ShowSuccess("Test message")

	fb.Update(ClearFeedbackMsg{})

	if fb.Visible() {
		t.Error("Feedback should not be visible after ClearFeedbackMsg")
	}
}

// TestFeedbackSetDuration verifies SetDuration changes the duration
func TestFeedbackSetDuration(t *testing.T) {
	fb := NewFeedback()
	fb.SetDuration(5 * time.Second)

	if fb.duration != 5*time.Second {
		t.Errorf("duration = %v, want 5s", fb.duration)
	}
}

// TestFeedbackViewSuccess verifies View renders success message correctly
func TestFeedbackViewSuccess(t *testing.T) {
	fb := NewFeedback()
	fb.ShowSuccess("Success!")

	view := fb.View()

	if view == "" {
		t.Error("View() returned empty for visible feedback")
	}
	if !strings.Contains(view, "Success!") {
		t.Error("View() should contain the message")
	}
	if !strings.Contains(view, "✓") {
		t.Error("View() should contain success icon")
	}
}

// TestFeedbackViewError verifies View renders error message correctly
func TestFeedbackViewError(t *testing.T) {
	fb := NewFeedback()
	fb.ShowError("Error!")

	view := fb.View()

	if view == "" {
		t.Error("View() returned empty for visible feedback")
	}
	if !strings.Contains(view, "Error!") {
		t.Error("View() should contain the message")
	}
	if !strings.Contains(view, "✗") {
		t.Error("View() should contain error icon")
	}
}

// TestFeedbackViewInfo verifies View renders info message correctly
func TestFeedbackViewInfo(t *testing.T) {
	fb := NewFeedback()
	fb.ShowInfo("Info!")

	view := fb.View()

	if view == "" {
		t.Error("View() returned empty for visible feedback")
	}
	if !strings.Contains(view, "Info!") {
		t.Error("View() should contain the message")
	}
	if !strings.Contains(view, "ℹ") {
		t.Error("View() should contain info icon")
	}
}

// TestFeedbackViewWhenHidden verifies View returns empty when hidden
func TestFeedbackViewWhenHidden(t *testing.T) {
	fb := NewFeedback()
	// Feedback is hidden by default

	view := fb.View()

	if view != "" {
		t.Errorf("View() when hidden returned %q, want empty", view)
	}
}

// TestFeedbackViewAfterClear verifies View returns empty after clear
func TestFeedbackViewAfterClear(t *testing.T) {
	fb := NewFeedback()
	fb.ShowSuccess("Test")
	fb.Clear()

	view := fb.View()

	if view != "" {
		t.Errorf("View() after Clear returned %q, want empty", view)
	}
}

// TestFeedbackTypeConstants verifies FeedbackType constants are distinct
func TestFeedbackTypeConstants(t *testing.T) {
	if FeedbackSuccess == FeedbackError {
		t.Error("FeedbackSuccess should not equal FeedbackError")
	}
	if FeedbackSuccess == FeedbackInfo {
		t.Error("FeedbackSuccess should not equal FeedbackInfo")
	}
	if FeedbackError == FeedbackInfo {
		t.Error("FeedbackError should not equal FeedbackInfo")
	}
}
