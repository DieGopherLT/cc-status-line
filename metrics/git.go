package metrics

import (
	"fmt"
	"os/exec"
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
func GetGitInfo(cwd string, linesAdded, linesRemoved int) *GitInfo {
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

	// Format changes
	if linesAdded > 0 || linesRemoved > 0 {
		info.HasChanges = true
		info.ChangesText = formatGitChanges(linesAdded, linesRemoved)
	} else {
		info.ChangesText = "(no changes)"
	}

	return info
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
