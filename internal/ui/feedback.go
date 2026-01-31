// Package ui provides the terminal user interface for the git worktree manager.
package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FeedbackType represents the type of feedback message.
type FeedbackType int

const (
	// FeedbackSuccess indicates a successful operation.
	FeedbackSuccess FeedbackType = iota
	// FeedbackError indicates an error occurred.
	FeedbackError
	// FeedbackInfo indicates informational message.
	FeedbackInfo
)

// Feedback displays temporary feedback messages to the user.
type Feedback struct {
	message      string
	feedbackType FeedbackType
	visible      bool
	duration     time.Duration
}

// NewFeedback creates a new feedback component.
func NewFeedback() *Feedback {
	return &Feedback{
		duration: 3 * time.Second,
	}
}

// Visible returns whether feedback is currently showing.
func (f *Feedback) Visible() bool {
	return f.visible
}

// Message returns the current feedback message.
func (f *Feedback) Message() string {
	return f.message
}

// Type returns the current feedback type.
func (f *Feedback) Type() FeedbackType {
	return f.feedbackType
}

// ClearFeedbackMsg is sent to clear the feedback message.
type ClearFeedbackMsg struct{}

// ShowSuccess displays a success message.
func (f *Feedback) ShowSuccess(message string) tea.Cmd {
	f.message = message
	f.feedbackType = FeedbackSuccess
	f.visible = true
	return f.scheduleClear()
}

// ShowError displays an error message.
func (f *Feedback) ShowError(message string) tea.Cmd {
	f.message = message
	f.feedbackType = FeedbackError
	f.visible = true
	return f.scheduleClear()
}

// ShowInfo displays an informational message.
func (f *Feedback) ShowInfo(message string) tea.Cmd {
	f.message = message
	f.feedbackType = FeedbackInfo
	f.visible = true
	return f.scheduleClear()
}

// Clear hides the feedback message.
func (f *Feedback) Clear() {
	f.visible = false
	f.message = ""
}

// SetDuration sets how long feedback messages are shown.
func (f *Feedback) SetDuration(d time.Duration) {
	f.duration = d
}

// scheduleClear returns a command that will clear the feedback after duration.
func (f *Feedback) scheduleClear() tea.Cmd {
	return tea.Tick(f.duration, func(time.Time) tea.Msg {
		return ClearFeedbackMsg{}
	})
}

// Update handles messages for the feedback component.
func (f *Feedback) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case ClearFeedbackMsg:
		f.Clear()
	}
	return nil
}

// View renders the feedback message.
func (f *Feedback) View() string {
	if !f.visible || f.message == "" {
		return ""
	}

	var style lipgloss.Style

	switch f.feedbackType {
	case FeedbackSuccess:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#4CAF50"}).
			Bold(true).
			Padding(0, 1)
	case FeedbackError:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#EF5350"}).
			Bold(true).
			Padding(0, 1)
	case FeedbackInfo:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#1565C0", Dark: "#42A5F5"}).
			Bold(true).
			Padding(0, 1)
	}

	// Add icon based on type
	var icon string
	switch f.feedbackType {
	case FeedbackSuccess:
		icon = "✓ "
	case FeedbackError:
		icon = "✗ "
	case FeedbackInfo:
		icon = "ℹ "
	}

	return style.Render(icon + f.message)
}
