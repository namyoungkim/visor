package git

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// commandTimeout is the maximum time allowed for git commands.
const commandTimeout = 200 * time.Millisecond

// Status represents git repository status.
type Status struct {
	Branch   string
	IsRepo   bool
	IsDirty  bool
	Ahead    int
	Behind   int
	Staged   int
	Modified int
}

// GetStatus returns the current git status.
// Returns empty Status if not in a git repository.
func GetStatus() Status {
	var status Status

	// Check if we're in a git repo
	if !isGitRepo() {
		return status
	}
	status.IsRepo = true

	// Get branch name
	status.Branch = getBranch()

	// Get status counts
	status.Staged, status.Modified = getStatusCounts()
	status.IsDirty = status.Staged > 0 || status.Modified > 0

	// Get ahead/behind counts
	status.Ahead, status.Behind = getAheadBehind()

	return status
}

// gitCommand executes a git command with timeout.
func gitCommand(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	return cmd.Output()
}

// gitCommandRun executes a git command with timeout, returning only error.
func gitCommandRun(args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	return cmd.Run()
}

func isGitRepo() bool {
	err := gitCommandRun("rev-parse", "--is-inside-work-tree")
	return err == nil
}

func getBranch() string {
	out, err := gitCommand("branch", "--show-current")
	if err != nil {
		// Might be detached HEAD
		out, err = gitCommand("rev-parse", "--short", "HEAD")
		if err != nil {
			return ""
		}
	}
	return strings.TrimSpace(string(out))
}

func getStatusCounts() (staged, modified int) {
	out, err := gitCommand("status", "--porcelain")
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		index := line[0]
		worktree := line[1]

		// Staged changes (index has M, A, D, R, C)
		if index != ' ' && index != '?' {
			staged++
		}
		// Modified in worktree
		if worktree != ' ' && worktree != '?' {
			modified++
		}
	}
	return staged, modified
}

func getAheadBehind() (ahead, behind int) {
	out, err := gitCommand("rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if err != nil {
		return 0, 0
	}

	parts := strings.Fields(string(out))
	if len(parts) == 2 {
		// Format: behind ahead
		behind, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
		ahead, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}
	return ahead, behind
}
