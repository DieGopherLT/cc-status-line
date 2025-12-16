package formatters

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
)

// MinimalFormatter implements an ultra-compact style with no decorations
type MinimalFormatter struct{}

// Format creates a compact single-line status line
func (f *MinimalFormatter) Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var parts []string

	// Model name (always present)
	parts = append(parts, modelStyle.Render(hook.Model.DisplayName))

	// Git branch and changes
	if gitInfo != nil && gitInfo.IsGitRepo {
		parts = append(parts, branchStyle.Render(gitInfo.BranchDisplay))

		// Git changes in compact format: +156-23
		if gitInfo.Additions > 0 || gitInfo.Deletions > 0 {
			coloredChanges := greenStyle.Render(fmt.Sprintf("+%d", gitInfo.Additions)) + redStyle.Render(fmt.Sprintf("-%d", gitInfo.Deletions))
			parts = append(parts, coloredChanges)
		}
	} else {
		parts = append(parts, grayStyle.Render("(no git)"))
	}

	// Output style (only if present)
	if hook.OutputStyle.Name != "" {
		parts = append(parts, styleColor.Render(hook.OutputStyle.Name))
	}

	// Version (always present)
	parts = append(parts, blueStyle.Render(hook.Version))

	// Context percentage (only if available)
	if tokenMetrics != nil && tokenMetrics.ContextPercentage > 0 {
		ctxPart := fmt.Sprintf("%d%%", int(tokenMetrics.ContextPercentage))
		parts = append(parts, ctxPart)
	}

	// Join with single space
	statusLine := strings.Join(parts, " ")

	// Add subtle horizontal lines for visual breathing room
	width := lipgloss.Width(statusLine)
	line := lineStyle.Render(strings.Repeat("─", width))

	return line + "\n" + statusLine + "\n" + line
}
