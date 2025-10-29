package metrics

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GitInfo contains git repository information
type GitInfo struct {
	Branch        string
	BranchDisplay string
	HasChanges    bool
	ChangesText   string
	IsGitRepo     bool
}

// GetGitInfo extracts git branch and change information
func GetGitInfo(cwd string) *GitInfo {
	info := &GitInfo{
		IsGitRepo: false,
	}

	// Try to get the current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = cwd
	output, err := cmd.Output()

	if err != nil {
		// Not in a git repository
		info.BranchDisplay = "(no git)"
		info.ChangesText = "(no git)"
		return info
	}

	info.IsGitRepo = true
	info.Branch = strings.TrimSpace(string(output))
	info.BranchDisplay = info.Branch

	// Get git changes (staged + unstaged) from git directly
	linesAdded, linesRemoved := getGitChanges(cwd)

	// Format changes
	if linesAdded > 0 || linesRemoved > 0 {
		info.HasChanges = true
		info.ChangesText = formatGitChanges(linesAdded, linesRemoved)
	} else {
		info.ChangesText = "(no changes)"
	}

	return info
}

// getGitChanges gets the number of lines added and removed from git
func getGitChanges(cwd string) (int, int) {
	// Get all changes (staged + unstaged) compared to HEAD
	cmd := exec.Command("git", "diff", "--numstat", "HEAD")
	cmd.Dir = cwd
	output, err := cmd.Output()

	if err != nil {
		return 0, 0
	}

	var totalAdded, totalRemoved int

	// Parse numstat output: "additions\tdeletions\tfilename"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// Parse additions
		if parts[0] != "-" {
			added, err := strconv.Atoi(parts[0])
			if err == nil {
				totalAdded += added
			}
		}

		// Parse deletions
		if parts[1] != "-" {
			removed, err := strconv.Atoi(parts[1])
			if err == nil {
				totalRemoved += removed
			}
		}
	}

	return totalAdded, totalRemoved
}

// formatGitChanges formats the git changes display
func formatGitChanges(added, removed int) string {
	if added > 0 && removed > 0 {
		return fmt.Sprintf("(+%d -%d)", added, removed)
	}
	if added > 0 {
		return fmt.Sprintf("(+%d)", added)
	}
	if removed > 0 {
		return fmt.Sprintf("(-%d)", removed)
	}
	return "(no changes)"
}
