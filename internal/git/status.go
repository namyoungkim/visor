package git

import (
	"os/exec"
	"strings"
)

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

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func getBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		// Might be detached HEAD
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		out, err = cmd.Output()
		if err != nil {
			return ""
		}
	}
	return strings.TrimSpace(string(out))
}

func getStatusCounts() (staged, modified int) {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
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
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	parts := strings.Fields(string(out))
	if len(parts) == 2 {
		// Format: behind ahead
		if n, err := parseInt(parts[0]); err == nil {
			behind = n
		}
		if n, err := parseInt(parts[1]); err == nil {
			ahead = n
		}
	}
	return ahead, behind
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
