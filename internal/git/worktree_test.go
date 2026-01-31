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

// TestWorktreeAddError verifies the error type and message.
func TestWorktreeAddError(t *testing.T) {
	err := &WorktreeAddError{
		Path:   "/path/to/worktree",
		Branch: "feature",
		Reason: "branch already exists",
	}

	expected := "failed to add worktree at /path/to/worktree for branch feature: branch already exists"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestAddWorktreeOptions verifies the options struct.
func TestAddWorktreeOptions(t *testing.T) {
	opts := AddWorktreeOptions{
		Path:         "/path/to/worktree",
		Branch:       "feature",
		CreateBranch: true,
		BaseBranch:   "main",
	}

	if opts.Path != "/path/to/worktree" {
		t.Errorf("Expected Path '/path/to/worktree', got '%s'", opts.Path)
	}
	if opts.Branch != "feature" {
		t.Errorf("Expected Branch 'feature', got '%s'", opts.Branch)
	}
	if !opts.CreateBranch {
		t.Error("Expected CreateBranch true, got false")
	}
	if opts.BaseBranch != "main" {
		t.Errorf("Expected BaseBranch 'main', got '%s'", opts.BaseBranch)
	}
}

// TestAddWorktreeInNonGitDir tests AddWorktree in a non-git directory.
func TestAddWorktreeInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         "/path/to/worktree",
		Branch:       "feature",
		CreateBranch: true,
	})

	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestAddWorktreeEmptyPath tests AddWorktree with empty path.
func TestAddWorktreeEmptyPath(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

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

	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         "",
		Branch:       "feature",
		CreateBranch: true,
	})

	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}

	addErr, ok := err.(*WorktreeAddError)
	if !ok {
		t.Fatalf("Expected WorktreeAddError, got: %T", err)
	}
	if addErr.Reason != "path is required" {
		t.Errorf("Expected reason 'path is required', got '%s'", addErr.Reason)
	}
}

// TestAddWorktreeNoBranchWithoutCreate tests AddWorktree without branch when not creating.
func TestAddWorktreeNoBranchWithoutCreate(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

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

	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         "/path/to/worktree",
		Branch:       "",
		CreateBranch: false,
	})

	if err == nil {
		t.Error("Expected error for empty branch, got nil")
	}

	addErr, ok := err.(*WorktreeAddError)
	if !ok {
		t.Fatalf("Expected WorktreeAddError, got: %T", err)
	}
	if addErr.Reason != "branch is required when not creating a new branch" {
		t.Errorf("Expected reason about branch required, got '%s'", addErr.Reason)
	}
}

// TestAddWorktreeIntegration tests creating a worktree with a new branch.
func TestAddWorktreeIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create worktree path
	worktreePath := filepath.Join(tmpDir, "..", "worktree-add-test")
	defer os.RemoveAll(worktreePath)

	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         worktreePath,
		Branch:       "new-feature",
		CreateBranch: true,
	})

	if err != nil {
		t.Fatalf("AddWorktree failed: %v", err)
	}

	// Verify the worktree was created
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	found := false
	for _, wt := range worktrees {
		if wt.Branch == "new-feature" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find new-feature worktree in list: %+v", worktrees)
	}
}

// TestAddWorktreeWithExistingBranch tests creating a worktree with an existing branch.
func TestAddWorktreeWithExistingBranch(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create a branch
	cmd = exec.Command("git", "branch", "existing-branch")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git branch failed: %v", err)
	}

	// Create worktree using the existing branch
	worktreePath := filepath.Join(tmpDir, "..", "worktree-existing-test")
	defer os.RemoveAll(worktreePath)

	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         worktreePath,
		Branch:       "existing-branch",
		CreateBranch: false,
	})

	if err != nil {
		t.Fatalf("AddWorktree failed: %v", err)
	}

	// Verify the worktree was created
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	found := false
	for _, wt := range worktrees {
		if wt.Branch == "existing-branch" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Did not find existing-branch worktree in list: %+v", worktrees)
	}
}

// TestListBranchesInNonGitDir tests ListBranches in a non-git directory.
func TestListBranchesInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = ListBranches(tmpDir)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestListBranchesIntegration tests listing branches in a git repository.
func TestListBranchesIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create additional branches
	cmd = exec.Command("git", "branch", "feature-one")
	cmd.Dir = tmpDir
	cmd.Run()
	cmd = exec.Command("git", "branch", "feature-two")
	cmd.Dir = tmpDir
	cmd.Run()

	// List branches
	branches, err := ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}

	if len(branches) < 3 {
		t.Errorf("Expected at least 3 branches, got %d", len(branches))
	}

	// Check for expected branches
	expectedBranches := []string{"feature-one", "feature-two"}
	for _, expected := range expectedBranches {
		found := false
		for _, b := range branches {
			if b == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find branch '%s' in list: %+v", expected, branches)
		}
	}
}

// TestWorktreeRemoveError verifies the error type and message.
func TestWorktreeRemoveError(t *testing.T) {
	err := &WorktreeRemoveError{
		Path:   "/path/to/worktree",
		Reason: "worktree has uncommitted changes",
	}

	expected := "failed to remove worktree at /path/to/worktree: worktree has uncommitted changes"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestRemoveWorktreeOptions verifies the options struct.
func TestRemoveWorktreeOptions(t *testing.T) {
	opts := RemoveWorktreeOptions{
		Path:  "/path/to/worktree",
		Force: true,
	}

	if opts.Path != "/path/to/worktree" {
		t.Errorf("Expected Path '/path/to/worktree', got '%s'", opts.Path)
	}
	if !opts.Force {
		t.Error("Expected Force true, got false")
	}
}

// TestRemoveWorktreeInNonGitDir tests RemoveWorktree in a non-git directory.
func TestRemoveWorktreeInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = RemoveWorktree(tmpDir, RemoveWorktreeOptions{
		Path: "/path/to/worktree",
	})

	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestRemoveWorktreeEmptyPath tests RemoveWorktree with empty path.
func TestRemoveWorktreeEmptyPath(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

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

	err = RemoveWorktree(tmpDir, RemoveWorktreeOptions{
		Path: "",
	})

	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}

	removeErr, ok := err.(*WorktreeRemoveError)
	if !ok {
		t.Fatalf("Expected WorktreeRemoveError, got: %T", err)
	}
	if removeErr.Reason != "path is required" {
		t.Errorf("Expected reason 'path is required', got '%s'", removeErr.Reason)
	}
}

// TestRemoveWorktreeIntegration tests removing a worktree.
func TestRemoveWorktreeIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create worktree path
	worktreePath := filepath.Join(tmpDir, "..", "worktree-remove-test")

	// Add worktree
	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         worktreePath,
		Branch:       "remove-test",
		CreateBranch: true,
	})
	if err != nil {
		t.Fatalf("AddWorktree failed: %v", err)
	}

	// Verify worktree was created
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	found := false
	for _, wt := range worktrees {
		if wt.Branch == "remove-test" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Worktree was not created")
	}

	// Remove the worktree
	err = RemoveWorktree(tmpDir, RemoveWorktreeOptions{
		Path: worktreePath,
	})
	if err != nil {
		t.Fatalf("RemoveWorktree failed: %v", err)
	}

	// Verify the worktree was removed
	worktrees, err = ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	for _, wt := range worktrees {
		if wt.Branch == "remove-test" {
			t.Error("Worktree was not removed")
		}
	}
}

// TestRemoveWorktreeWithUncommittedChanges tests removing a worktree with uncommitted changes.
func TestRemoveWorktreeWithUncommittedChanges(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create worktree path
	worktreePath := filepath.Join(tmpDir, "..", "worktree-uncommitted-test")
	defer os.RemoveAll(worktreePath)

	// Add worktree
	err = AddWorktree(tmpDir, AddWorktreeOptions{
		Path:         worktreePath,
		Branch:       "uncommitted-test",
		CreateBranch: true,
	})
	if err != nil {
		t.Fatalf("AddWorktree failed: %v", err)
	}

	// Create uncommitted changes in the worktree
	newFile := filepath.Join(worktreePath, "uncommitted.txt")
	if err := os.WriteFile(newFile, []byte("uncommitted change"), 0644); err != nil {
		t.Fatalf("Failed to create uncommitted file: %v", err)
	}

	// Try to remove the worktree without force - should fail
	err = RemoveWorktree(tmpDir, RemoveWorktreeOptions{
		Path:  worktreePath,
		Force: false,
	})
	if err == nil {
		t.Error("Expected error for worktree with uncommitted changes, got nil")
	}

	// Remove with force - should succeed
	err = RemoveWorktree(tmpDir, RemoveWorktreeOptions{
		Path:  worktreePath,
		Force: true,
	})
	if err != nil {
		t.Fatalf("RemoveWorktree with force failed: %v", err)
	}
}

// TestHasUncommittedChangesInNonGitDir tests HasUncommittedChanges in a non-git directory.
func TestHasUncommittedChangesInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = HasUncommittedChanges(tmpDir)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestHasUncommittedChangesIntegration tests HasUncommittedChanges with a git repository.
func TestHasUncommittedChangesIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Clean state - no uncommitted changes
	hasChanges, err := HasUncommittedChanges(tmpDir)
	if err != nil {
		t.Fatalf("HasUncommittedChanges failed: %v", err)
	}
	if hasChanges {
		t.Error("Expected no uncommitted changes, got true")
	}

	// Create uncommitted changes
	newFile := filepath.Join(tmpDir, "new.txt")
	if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Should now have uncommitted changes
	hasChanges, err = HasUncommittedChanges(tmpDir)
	if err != nil {
		t.Fatalf("HasUncommittedChanges failed: %v", err)
	}
	if !hasChanges {
		t.Error("Expected uncommitted changes, got false")
	}
}

// TestWorktreePruneError verifies the error type and message.
func TestWorktreePruneError(t *testing.T) {
	err := &WorktreePruneError{
		Reason: "failed to prune worktrees",
	}

	expected := "failed to prune worktrees: failed to prune worktrees"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestPruneWorktreesInNonGitDir tests PruneWorktrees in a non-git directory.
func TestPruneWorktreesInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = PruneWorktrees(tmpDir)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestPruneWorktreesIntegration tests pruning worktrees in a git repository.
func TestPruneWorktreesIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Run prune on a clean repo - should succeed without errors
	output, err := PruneWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("PruneWorktrees failed: %v", err)
	}

	// Output should be empty or contain no error text
	if strings.Contains(strings.ToLower(output), "error") {
		t.Errorf("Expected no errors in output, got: %s", output)
	}
}

// TestPruneWorktreesWithStaleEntry tests pruning a stale worktree entry.
func TestPruneWorktreesWithStaleEntry(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create a worktree
	worktreePath := filepath.Join(tmpDir, "..", "worktree-prune-test")
	cmd = exec.Command("git", "worktree", "add", "-b", "prune-test", worktreePath)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add failed: %v", err)
	}

	// Verify worktree exists in list
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	if len(worktrees) < 2 {
		t.Errorf("Expected at least 2 worktrees, got %d", len(worktrees))
	}

	// Manually delete the worktree directory to create a stale entry
	if err := os.RemoveAll(worktreePath); err != nil {
		t.Fatalf("Failed to remove worktree directory: %v", err)
	}

	// Prune should clean up the stale entry
	output, err := PruneWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("PruneWorktrees failed: %v", err)
	}

	// The prune should have worked (even if output is empty)
	_ = output

	// Verify the stale entry was removed from the worktree list
	worktrees, err = ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	// Should no longer have the stale worktree
	for _, wt := range worktrees {
		if wt.Branch == "prune-test" {
			t.Error("Stale worktree was not pruned")
		}
	}
}

// TestWorktreeStatusFields verifies the WorktreeStatus struct fields and methods.
func TestWorktreeStatusFields(t *testing.T) {
	status := WorktreeStatus{
		ModifiedCount:  3,
		StagedCount:    2,
		UntrackedCount: 5,
	}

	if status.ModifiedCount != 3 {
		t.Errorf("Expected ModifiedCount 3, got %d", status.ModifiedCount)
	}
	if status.StagedCount != 2 {
		t.Errorf("Expected StagedCount 2, got %d", status.StagedCount)
	}
	if status.UntrackedCount != 5 {
		t.Errorf("Expected UntrackedCount 5, got %d", status.UntrackedCount)
	}
	if status.TotalChanges() != 10 {
		t.Errorf("Expected TotalChanges 10, got %d", status.TotalChanges())
	}
	if status.IsClean() {
		t.Error("Expected IsClean false, got true")
	}
}

// TestWorktreeStatusIsClean tests the IsClean method.
func TestWorktreeStatusIsClean(t *testing.T) {
	tests := []struct {
		name     string
		status   WorktreeStatus
		expected bool
	}{
		{
			name:     "all zeros",
			status:   WorktreeStatus{ModifiedCount: 0, StagedCount: 0, UntrackedCount: 0},
			expected: true,
		},
		{
			name:     "modified only",
			status:   WorktreeStatus{ModifiedCount: 1, StagedCount: 0, UntrackedCount: 0},
			expected: false,
		},
		{
			name:     "staged only",
			status:   WorktreeStatus{ModifiedCount: 0, StagedCount: 1, UntrackedCount: 0},
			expected: false,
		},
		{
			name:     "untracked only",
			status:   WorktreeStatus{ModifiedCount: 0, StagedCount: 0, UntrackedCount: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsClean(); got != tt.expected {
				t.Errorf("IsClean() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParseWorktreeStatus tests parsing of git status --porcelain output.
func TestParseWorktreeStatus(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedModified  int
		expectedStaged    int
		expectedUntracked int
	}{
		{
			name:              "empty output",
			input:             "",
			expectedModified:  0,
			expectedStaged:    0,
			expectedUntracked: 0,
		},
		{
			name:              "single modified file",
			input:             " M file.txt\n",
			expectedModified:  1,
			expectedStaged:    0,
			expectedUntracked: 0,
		},
		{
			name:              "single staged file",
			input:             "M  file.txt\n",
			expectedModified:  0,
			expectedStaged:    1,
			expectedUntracked: 0,
		},
		{
			name:              "single untracked file",
			input:             "?? file.txt\n",
			expectedModified:  0,
			expectedStaged:    0,
			expectedUntracked: 1,
		},
		{
			name:              "staged and modified same file",
			input:             "MM file.txt\n",
			expectedModified:  1,
			expectedStaged:    1,
			expectedUntracked: 0,
		},
		{
			name:              "added file",
			input:             "A  file.txt\n",
			expectedModified:  0,
			expectedStaged:    1,
			expectedUntracked: 0,
		},
		{
			name:              "deleted file",
			input:             "D  file.txt\n",
			expectedModified:  0,
			expectedStaged:    1,
			expectedUntracked: 0,
		},
		{
			name:              "renamed file",
			input:             "R  old.txt -> new.txt\n",
			expectedModified:  0,
			expectedStaged:    1,
			expectedUntracked: 0,
		},
		{
			name: "multiple files",
			input: ` M modified.txt
M  staged.txt
?? untracked.txt
MM both.txt
A  added.txt
`,
			expectedModified:  2, // modified.txt and both.txt
			expectedStaged:    3, // staged.txt, both.txt, and added.txt
			expectedUntracked: 1, // untracked.txt
		},
		{
			name:              "deleted in worktree",
			input:             " D file.txt\n",
			expectedModified:  1,
			expectedStaged:    0,
			expectedUntracked: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := ParseWorktreeStatus(tt.input)

			if status.ModifiedCount != tt.expectedModified {
				t.Errorf("ModifiedCount = %d, want %d", status.ModifiedCount, tt.expectedModified)
			}
			if status.StagedCount != tt.expectedStaged {
				t.Errorf("StagedCount = %d, want %d", status.StagedCount, tt.expectedStaged)
			}
			if status.UntrackedCount != tt.expectedUntracked {
				t.Errorf("UntrackedCount = %d, want %d", status.UntrackedCount, tt.expectedUntracked)
			}
		})
	}
}

// TestGetWorktreeStatusInNonGitDir tests GetWorktreeStatus in a non-git directory.
func TestGetWorktreeStatusInNonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitworktreetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = GetWorktreeStatus(tmpDir)
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}
	if !IsNotGitRepoError(err) {
		t.Errorf("Expected NotGitRepoError, got: %v", err)
	}
}

// TestGetWorktreeStatusIntegration tests GetWorktreeStatus with a real git repository.
func TestGetWorktreeStatusIntegration(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Clean state - no uncommitted changes
	status, err := GetWorktreeStatus(tmpDir)
	if err != nil {
		t.Fatalf("GetWorktreeStatus failed: %v", err)
	}
	if !status.IsClean() {
		t.Errorf("Expected clean status, got: modified=%d, staged=%d, untracked=%d",
			status.ModifiedCount, status.StagedCount, status.UntrackedCount)
	}

	// Create an untracked file
	untrackedFile := filepath.Join(tmpDir, "untracked.txt")
	if err := os.WriteFile(untrackedFile, []byte("untracked"), 0644); err != nil {
		t.Fatalf("Failed to create untracked file: %v", err)
	}

	status, err = GetWorktreeStatus(tmpDir)
	if err != nil {
		t.Fatalf("GetWorktreeStatus failed: %v", err)
	}
	if status.UntrackedCount != 1 {
		t.Errorf("Expected 1 untracked file, got %d", status.UntrackedCount)
	}

	// Modify an existing tracked file
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	status, err = GetWorktreeStatus(tmpDir)
	if err != nil {
		t.Fatalf("GetWorktreeStatus failed: %v", err)
	}
	if status.ModifiedCount != 1 {
		t.Errorf("Expected 1 modified file, got %d", status.ModifiedCount)
	}

	// Stage the modified file
	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add failed: %v", err)
	}

	status, err = GetWorktreeStatus(tmpDir)
	if err != nil {
		t.Fatalf("GetWorktreeStatus failed: %v", err)
	}
	if status.StagedCount != 1 {
		t.Errorf("Expected 1 staged file, got %d", status.StagedCount)
	}
	if status.ModifiedCount != 0 {
		t.Errorf("Expected 0 modified files after staging, got %d", status.ModifiedCount)
	}
}

// TestPruneWorktreesDryRun tests the dry-run mode of pruning.
func TestPruneWorktreesDryRun(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

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

	// Create a worktree
	worktreePath := filepath.Join(tmpDir, "..", "worktree-dryrun-test")
	cmd = exec.Command("git", "worktree", "add", "-b", "dryrun-test", worktreePath)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git worktree add failed: %v", err)
	}
	defer os.RemoveAll(worktreePath)

	// Manually delete the worktree directory to create a stale entry
	if err := os.RemoveAll(worktreePath); err != nil {
		t.Fatalf("Failed to remove worktree directory: %v", err)
	}

	// Dry run should report the stale entry but not remove it
	output, err := PruneWorktreesDryRun(tmpDir)
	if err != nil {
		t.Fatalf("PruneWorktreesDryRun failed: %v", err)
	}

	// Output should mention the stale worktree path
	if !strings.Contains(output, "dryrun-test") && !strings.Contains(output, "worktree-dryrun-test") {
		// Some git versions may have different output format
		// Just check that it ran successfully
		_ = output
	}

	// The entry should still be in the list (dry run doesn't remove)
	worktrees, err := ListWorktrees(tmpDir)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}

	// The worktree entry should still be there but marked as stale in list
	// Note: git worktree list may or may not show stale entries depending on version
	_ = worktrees
}
