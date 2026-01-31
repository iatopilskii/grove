// Package git provides git operations for the worktree manager.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree with its metadata.
type Worktree struct {
	// Path is the absolute path to the worktree directory.
	Path string
	// Branch is the name of the branch checked out in this worktree.
	// Empty for bare repositories or detached HEAD.
	Branch string
	// CommitHash is the short commit hash of the HEAD.
	CommitHash string
	// IsBare indicates if this is a bare repository.
	IsBare bool
	// IsDetached indicates if the worktree is in detached HEAD state.
	IsDetached bool
}

// Name returns the name of the worktree (last component of the path).
func (w *Worktree) Name() string {
	return filepath.Base(w.Path)
}

// NotGitRepoError is returned when an operation is performed outside a git repository.
type NotGitRepoError struct {
	Path string
}

func (e *NotGitRepoError) Error() string {
	return "not a git repository: " + e.Path
}

// IsNotGitRepoError checks if an error is a NotGitRepoError.
func IsNotGitRepoError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*NotGitRepoError)
	return ok
}

// IsGitRepository checks if the given directory is inside a git repository.
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	err := cmd.Run()
	return err == nil
}

// GetCurrentDirectory returns the current working directory.
func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// ListWorktrees lists all worktrees in the git repository containing the given directory.
// Returns a NotGitRepoError if the directory is not in a git repository.
func ListWorktrees(dir string) ([]Worktree, error) {
	if !IsGitRepository(dir) {
		return nil, &NotGitRepoError{Path: dir}
	}

	cmd := exec.Command("git", "worktree", "list")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return ParseWorktreeList(string(output)), nil
}

// ParseWorktreeList parses the output of "git worktree list" command.
// The format is:
//
//	/path/to/worktree  <commit> [branch]
//	/path/to/bare      (bare)
//	/path/to/detached  <commit> (detached HEAD)
func ParseWorktreeList(output string) []Worktree {
	var worktrees []Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		wt := parseWorktreeLine(line)
		if wt.Path != "" {
			worktrees = append(worktrees, wt)
		}
	}

	return worktrees
}

// parseWorktreeLine parses a single line of git worktree list output.
func parseWorktreeLine(line string) Worktree {
	var wt Worktree

	// Check for bare repository
	if strings.HasSuffix(line, "(bare)") {
		wt.IsBare = true
		wt.Path = strings.TrimSpace(strings.TrimSuffix(line, "(bare)"))
		return wt
	}

	// Check for detached HEAD
	if strings.HasSuffix(line, "(detached HEAD)") {
		wt.IsDetached = true
		remaining := strings.TrimSuffix(line, "(detached HEAD)")
		remaining = strings.TrimSpace(remaining)
		parts := splitWorktreePath(remaining)
		if len(parts) >= 1 {
			wt.Path = parts[0]
		}
		if len(parts) >= 2 {
			wt.CommitHash = parts[1]
		}
		return wt
	}

	// Regular worktree format: /path  hash [branch]
	// Find the branch in brackets
	bracketStart := strings.LastIndex(line, "[")
	bracketEnd := strings.LastIndex(line, "]")
	if bracketStart != -1 && bracketEnd != -1 && bracketEnd > bracketStart {
		wt.Branch = line[bracketStart+1 : bracketEnd]
		remaining := strings.TrimSpace(line[:bracketStart])
		parts := splitWorktreePath(remaining)
		if len(parts) >= 1 {
			wt.Path = parts[0]
		}
		if len(parts) >= 2 {
			wt.CommitHash = parts[1]
		}
	}

	return wt
}

// splitWorktreePath splits the path and hash portion of a worktree line.
// The format is: /path/to/worktree  <hash>
// Multiple spaces separate the path from the hash.
func splitWorktreePath(s string) []string {
	// Find two or more spaces which separate path from hash
	s = strings.TrimSpace(s)

	// Find the last sequence of 2+ spaces
	lastMultiSpace := -1
	for i := 0; i < len(s)-1; i++ {
		if s[i] == ' ' && s[i+1] == ' ' {
			lastMultiSpace = i
		}
	}

	if lastMultiSpace == -1 {
		// No multi-space separator found, entire string is path
		return []string{s}
	}

	path := strings.TrimSpace(s[:lastMultiSpace])
	hash := strings.TrimSpace(s[lastMultiSpace:])

	return []string{path, hash}
}

// WorktreeAddError is returned when worktree creation fails.
type WorktreeAddError struct {
	Path   string
	Branch string
	Reason string
}

func (e *WorktreeAddError) Error() string {
	return fmt.Sprintf("failed to add worktree at %s for branch %s: %s", e.Path, e.Branch, e.Reason)
}

// AddWorktreeOptions specifies options for creating a new worktree.
type AddWorktreeOptions struct {
	// Path is the absolute or relative path for the new worktree directory.
	Path string
	// Branch is the branch name to checkout. If empty and CreateBranch is true,
	// a new branch will be created with the name derived from Path.
	Branch string
	// CreateBranch indicates whether to create a new branch.
	// If true and Branch is empty, the branch name is derived from Path.
	CreateBranch bool
	// BaseBranch is the starting point for the new branch when CreateBranch is true.
	// If empty, defaults to HEAD.
	BaseBranch string
}

// AddWorktree creates a new git worktree at the specified path.
// The dir parameter is the directory of an existing git repository.
func AddWorktree(dir string, opts AddWorktreeOptions) error {
	if !IsGitRepository(dir) {
		return &NotGitRepoError{Path: dir}
	}

	if opts.Path == "" {
		return &WorktreeAddError{
			Path:   opts.Path,
			Branch: opts.Branch,
			Reason: "path is required",
		}
	}

	// Build the git worktree add command
	args := []string{"worktree", "add"}

	if opts.CreateBranch {
		// Create new branch
		branchName := opts.Branch
		if branchName == "" {
			// Derive branch name from path
			branchName = filepath.Base(opts.Path)
		}

		if opts.BaseBranch != "" {
			args = append(args, "-b", branchName, opts.Path, opts.BaseBranch)
		} else {
			args = append(args, "-b", branchName, opts.Path)
		}
	} else {
		// Use existing branch
		if opts.Branch == "" {
			return &WorktreeAddError{
				Path:   opts.Path,
				Branch: opts.Branch,
				Reason: "branch is required when not creating a new branch",
			}
		}
		args = append(args, opts.Path, opts.Branch)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		reason := strings.TrimSpace(string(output))
		if reason == "" {
			reason = err.Error()
		}
		return &WorktreeAddError{
			Path:   opts.Path,
			Branch: opts.Branch,
			Reason: reason,
		}
	}

	return nil
}

// ListBranches lists all local branches in the repository.
func ListBranches(dir string) ([]string, error) {
	if !IsGitRepository(dir) {
		return nil, &NotGitRepoError{Path: dir}
	}

	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			branches = append(branches, line)
		}
	}

	return branches, nil
}
