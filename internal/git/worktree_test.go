// Package git provides git operations for the worktree manager.
package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestWorktreeFields verifies the Worktree struct has required fields.
func TestWorktreeFields(t *testing.T) {
	wt := Worktree{
		Path:   "/path/to/worktree",
		Branch: "main",
		IsBare: false,
	}

	if wt.Path != "/path/to/worktree" {
		t.Errorf("Expected Path '/path/to/worktree', got '%s'", wt.Path)
	}
	if wt.Branch != "main" {
		t.Errorf("Expected Branch 'main', got '%s'", wt.Branch)
	}
	if wt.IsBare {
		t.Errorf("Expected IsBare false, got true")
	}
}

// TestWorktreeName verifies the Name() method returns the correct name.
func TestWorktreeName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/myrepo", "myrepo"},
		{"/path/to/feature-branch", "feature-branch"},
		{"/home/user/projects/main", "main"},
		{"simple", "simple"},
		{"/", "/"},
	}

	for _, tt := range tests {
		wt := Worktree{Path: tt.path}
		if got := wt.Name(); got != tt.expected {
			t.Errorf("Name() for path '%s': expected '%s', got '%s'", tt.path, tt.expected, got)
		}
	}
}

// TestParseWorktreeList tests parsing of git worktree list output.
func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Worktree
	}{
		{
			name: "single worktree",
			input: `/path/to/main  abc1234 [main]
`,
			expected: []Worktree{
				{Path: "/path/to/main", Branch: "main", CommitHash: "abc1234", IsBare: false},
			},
		},
		{
			name: "multiple worktrees",
			input: `/path/to/main  abc1234 [main]
/path/to/feature  def5678 [feature-branch]
`,
			expected: []Worktree{
				{Path: "/path/to/main", Branch: "main", CommitHash: "abc1234", IsBare: false},
				{Path: "/path/to/feature", Branch: "feature-branch", CommitHash: "def5678", IsBare: false},
			},
		},
		{
			name: "bare repository",
			input: `/path/to/repo.git  (bare)
/path/to/worktree  abc1234 [main]
`,
			expected: []Worktree{
				{Path: "/path/to/repo.git", Branch: "", CommitHash: "", IsBare: true},
				{Path: "/path/to/worktree", Branch: "main", CommitHash: "abc1234", IsBare: false},
			},
		},
		{
			name: "detached HEAD",
			input: `/path/to/detached  abc1234 (detached HEAD)
`,
			expected: []Worktree{
				{Path: "/path/to/detached", Branch: "", CommitHash: "abc1234", IsBare: false, IsDetached: true},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Worktree{},
		},
		{
			name:     "whitespace only",
			input:    "   \n\t\n   ",
			expected: []Worktree{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseWorktreeList(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d worktrees, got %d", len(tt.expected), len(result))
			}
			for i, wt := range result {
				if wt.Path != tt.expected[i].Path {
					t.Errorf("Worktree %d: expected Path '%s', got '%s'", i, tt.expected[i].Path, wt.Path)
				}
				if wt.Branch != tt.expected[i].Branch {
					t.Errorf("Worktree %d: expected Branch '%s', got '%s'", i, tt.expected[i].Branch, wt.Branch)
				}
				if wt.CommitHash != tt.expected[i].CommitHash {
					t.Errorf("Worktree %d: expected CommitHash '%s', got '%s'", i, tt.expected[i].CommitHash, wt.CommitHash)
				}
				if wt.IsBare != tt.expected[i].IsBare {
					t.Errorf("Worktree %d: expected IsBare %v, got %v", i, tt.expected[i].IsBare, wt.IsBare)
				}
				if wt.IsDetached != tt.expected[i].IsDetached {
					t.Errorf("Worktree %d: expected IsDetached %v, got %v", i, tt.expected[i].IsDetached, wt.IsDetached)
				}
			}
		})
	}
}

// TestIsGitRepository tests the IsGitRepository function.
func TestIsGitRepository(t *testing.T) {
	// Create a temporary directory that is NOT a git repo
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test non-git directory
	if IsGitRepository(tmpDir) {
		t.Errorf("Expected IsGitRepository to return false for non-git directory")
	}

	// Initialize a git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Skipping test: git init failed: %v", err)
	}

	// Test git directory
	if !IsGitRepository(tmpDir) {
		t.Errorf("Expected IsGitRepository to return true for git directory")
	}
}

// TestListWorktreesInNonGitDir tests that ListWorktrees returns error for non-git directory.
func TestListWorktreesInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = ListWorktrees(tmpDir)
	if err == nil {
		t.Errorf("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestListWorktreesIntegration tests ListWorktrees with an actual git repository.
func TestListWorktreesIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

	// Create a temporary directory and initialize a git repo
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Configure git user (required for commit)
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config email failed: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config name failed: %v", err)
	}

	// Create an initial commit (required for worktrees)
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// List worktrees (should return at least the main worktree)
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	if len(worktrees) < 1 {
		t.Errorf("Expected at least 1 worktree, got %d", len(worktrees))
	}

	// Verify the main worktree
	found := false
	for _, wt := range worktrees {
		if strings.HasSuffix(wt.Path, filepath.Base(tmpDir)) || wt.Path == tmpDir {
			found = true
			if wt.Branch == "" && !wt.IsBare && !wt.IsDetached {
				// New repos might have 'master' or 'main' or no branch set yet
				// This is acceptable
			}
			break
		}
	}
	if !found {
		t.Errorf("Did not find main worktree in list: %+v", worktrees)
	}
}

// TestListWorktreesWithMultipleWorktrees tests listing with multiple worktrees.
func TestListWorktreesWithMultipleWorktrees(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

	// Create a temporary directory and initialize a git repo
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	cmd.Run()
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	cmd.Run()

	// Create an initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// Create a new branch for the worktree
	worktreePath := filepath.Join(tmpDir, "..", "worktree-test-feature")
	cmd = exec.Command("git", "worktree", "add", "-b", "feature-test", worktreePath)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add failed: %v", err)
	}
	defer os.RemoveAll(worktreePath)

	// List worktrees
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	if len(worktrees) < 2 {
		t.Errorf("Expected at least 2 worktrees, got %d", len(worktrees))
	}

	// Verify we have a feature branch worktree
	foundFeature := false
	for _, wt := range worktrees {
		if wt.Branch == "feature-test" {
			foundFeature = true
			break
		}
	}
	if !foundFeature {
		t.Errorf("Did not find feature-test worktree in list: %+v", worktrees)
	}
}

// TestNotGitRepoError verifies the error type.
func TestNotGitRepoError(t *testing.T) {
	err := &NotGitRepoError{Path: "/some/path"}
	if !IsNotGitRepoError(err) {
		t.Error("Expected IsNotGitRepoError to return true for NotGitRepoError")
	}

	if IsNotGitRepoError(nil) {
		t.Error("Expected IsNotGitRepoError to return false for nil")
	}

	otherErr := os.ErrNotExist
	if IsNotGitRepoError(otherErr) {
		t.Error("Expected IsNotGitRepoError to return false for other errors")
	}

	// Test error message
	expected := "not a git repository: /some/path"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestGetCurrentDirectory tests the GetCurrentDirectory helper.
func TestGetCurrentDirectory(t *testing.T) {
	dir, err := GetCurrentDirectory()
	if err != nil {
		t.Fatalf("GetCurrentDirectory failed: %v", err)
	}
	if dir == "" {
		t.Error("GetCurrentDirectory returned empty string")
	}

	// Should be an absolute path
	if !filepath.IsAbs(dir) {
		t.Errorf("Expected absolute path, got: %s", dir)
	}
}
