package formatters

import (
	"fmt"
	"strings"

	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
)

const (
	compactTotalBlocks = 20
	compactFilledBlock = "█"
	compactEmptyBlock  = "░"
)

// Unicode icons for compact style
const (
	iconModel   = "◈"
	iconStyle   = "⎔"
	iconVersion = "⌘"
	iconContext = "◐"
	iconBranch  = ""
	iconAdd     = "↑"
	iconDel     = "↓"
)

// CompactFormatter implements a compact style with Unicode icons
type CompactFormatter struct{}

// Format creates a compact status line with icons
func (f *CompactFormatter) Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var parts []string

	// Model with icon
	modelPart := fmt.Sprintf("%s %s", iconModel, hook.Model.DisplayName)
	parts = append(parts, modelStyle.Render(modelPart))

	// Output style with icon
	if hook.OutputStyle.Name != "" {
		stylePart := fmt.Sprintf("%s %s", iconStyle, hook.OutputStyle.Name)
		parts = append(parts, styleColor.Render(stylePart))
	}

	// Version with icon
	versionPart := fmt.Sprintf("%s %s", iconVersion, hook.Version)
	parts = append(parts, blueStyle.Render(versionPart))

	// Context with icon and wider bar
	if tokenMetrics != nil && tokenMetrics.ContextPercentage > 0 {
		contextPart := f.formatContextBar(tokenMetrics)
		parts = append(parts, contextPart)
	}

	// Git with icon and arrows
	if gitInfo != nil && gitInfo.IsGitRepo {
		gitPart := f.formatGitInfo(gitInfo)
		parts = append(parts, gitPart)
	}

	// Join with double space
	statusLine := strings.Join(parts, "  ")

	return "\n" + statusLine + "\n"
}

// formatContextBar creates a 20-character context bar
func (f *CompactFormatter) formatContextBar(tokenMetrics *metrics.TokenMetrics) string {
	percentage := tokenMetrics.ContextPercentage

	// Calculate filled blocks (each block = 5%)
	filledCount := int(percentage / 5)
	if filledCount > compactTotalBlocks {
		filledCount = compactTotalBlocks
	}

	// Build bar
	filledBar := whiteStyle.Render(strings.Repeat(compactFilledBlock, filledCount))
	emptyBar := dimStyle.Render(strings.Repeat(compactEmptyBlock, compactTotalBlocks-filledCount))

	return fmt.Sprintf("%s %d%% [%s%s]", iconContext, int(percentage), filledBar, emptyBar)
}

// formatGitInfo formats git with branch icon and arrows for changes
func (f *CompactFormatter) formatGitInfo(gitInfo *metrics.GitInfo) string {
	if !gitInfo.IsGitRepo {
		return ""
	}

	// Branch with icon
	branch := branchStyle.Render(fmt.Sprintf("%s %s", iconBranch, gitInfo.BranchDisplay))

	// Changes with arrows
	var changes string
	if gitInfo.Additions > 0 || gitInfo.Deletions > 0 {
		if gitInfo.Additions > 0 {
			changes += greenStyle.Render(fmt.Sprintf("%s%d", iconAdd, gitInfo.Additions))
		}
		if gitInfo.Deletions > 0 {
			if changes != "" {
				changes += " "
			}
			changes += redStyle.Render(fmt.Sprintf("%s%d", iconDel, gitInfo.Deletions))
		}
		return branch + " " + changes
	}

	return branch
}
