package formatters

import (
	"fmt"
	"strings"

	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
	"github.com/charmbracelet/lipgloss"
)

const (
	nerdTotalBlocks = 10
	nerdFilledBlock = "█"
	nerdEmptyBlock  = "░"
)

// NerdFormatter implements a technical panel style with borders and absolute token counts
type NerdFormatter struct{}

// Format creates a bordered panel with detailed token metrics
func (f *NerdFormatter) Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var segments []string

	// Model name
	segments = append(segments, modelStyle.Render(hook.Model.DisplayName))

	// Git branch and changes
	if gitInfo != nil && gitInfo.IsGitRepo {
		gitSegment := fmt.Sprintf("%s %s%d %s%d",
			branchStyle.Render(gitInfo.BranchDisplay),
			greenStyle.Render("⇡"),
			gitInfo.Additions,
			redStyle.Render("⇣"),
			gitInfo.Deletions)
		segments = append(segments, gitSegment)
	}

	// Output style (if present)
	if hook.OutputStyle.Name != "" {
		segments = append(segments, styleColor.Render(hook.OutputStyle.Name))
	}

	// Version
	versionSegment := blueStyle.Render("v" + hook.Version)
	segments = append(segments, versionSegment)

	// Context with absolute tokens
	if tokenMetrics != nil && tokenMetrics.ContextPercentage > 0 {
		currentTokens := tokenMetrics.ContextLength
		maxTokens := tokenMetrics.ContextWindowSize
		bar := f.formatContextBar(tokenMetrics.ContextPercentage)

		ctxSegment := fmt.Sprintf("CTX: %s/%s (%d%%) %s",
			f.formatTokens(currentTokens),
			f.formatTokens(maxTokens),
			int(tokenMetrics.ContextPercentage),
			bar)
		segments = append(segments, ctxSegment)
	}

	// Join segments with box separator
	content := strings.Join(segments, grayStyle.Render(" │ "))

	// Calculate width using lipgloss (strips ANSI codes)
	contentWidth := lipgloss.Width(content)
	totalWidth := contentWidth + 4 // +4 for "│ " prefix and " │" suffix

	// Build bordered output
	topBorder := grayStyle.Render("┌" + strings.Repeat("─", totalWidth) + "┐")
	middle := grayStyle.Render("│ ") + content + grayStyle.Render(" │")
	bottomBorder := grayStyle.Render("└" + strings.Repeat("─", totalWidth) + "┘")

	return topBorder + "\n" + middle + "\n" + bottomBorder
}

// formatContextBar creates a 10-block context visualization
func (f *NerdFormatter) formatContextBar(percentage float64) string {
	filled := int(percentage / 10)
	if filled > nerdTotalBlocks {
		filled = nerdTotalBlocks
	}

	filledBar := whiteStyle.Render(strings.Repeat(nerdFilledBlock, filled))
	emptyBar := dimStyle.Render(strings.Repeat(nerdEmptyBlock, nerdTotalBlocks-filled))

	return filledBar + emptyBar
}

// formatTokens formats token count in human-readable format (k, M, etc.)
func (f *NerdFormatter) formatTokens(tokens int) string {
	if tokens >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
	}
	if tokens >= 1000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}
