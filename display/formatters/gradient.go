package formatters

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
)

const (
	gradientTotalBlocks = 10
	gradientSeparator   = " │ "
)

// Height-variable block characters (from lowest to highest)
var gradientBlocks = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// Color thresholds for gradient bar
var (
	gradientGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))  // Bright green
	gradientYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow
	gradientRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
)

// GradientFormatter implements a style with height-variable context bar and dynamic colors
type GradientFormatter struct{}

// Format creates the status line with gradient-style context visualization
func (f *GradientFormatter) Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var segments []string

	// Model name (compact, no "Model:" prefix)
	segments = append(segments, modelStyle.Render(hook.Model.DisplayName))

	// Context visualization with gradient bar
	if tokenMetrics != nil && tokenMetrics.ContextPercentage > 0 {
		contextSegment := f.formatGradientBar(tokenMetrics)
		segments = append(segments, contextSegment)
	}

	// Git info
	if gitInfo != nil && gitInfo.IsGitRepo {
		gitSegment := f.formatGitInfo(gitInfo)
		segments = append(segments, gitSegment)
	}

	// Output style (compact)
	if hook.OutputStyle.Name != "" {
		segments = append(segments, styleColor.Render(hook.OutputStyle.Name))
	}

	// Version
	versionSegment := fmt.Sprintf("v%s", hook.Version)
	segments = append(segments, blueStyle.Render(versionSegment))

	// Join with vertical bar separator
	statusLine := strings.Join(segments, grayStyle.Render(gradientSeparator))

	return "\n" + statusLine + "\n"
}

// formatGradientBar creates a height-variable bar with color based on usage
func (f *GradientFormatter) formatGradientBar(tokenMetrics *metrics.TokenMetrics) string {
	percentage := tokenMetrics.ContextPercentage

	// Determine color based on percentage
	var barStyle lipgloss.Style
	switch {
	case percentage >= 75:
		barStyle = gradientRed
	case percentage >= 50:
		barStyle = gradientYellow
	default:
		barStyle = gradientGreen
	}

	// Calculate filled blocks
	filledCount := int(percentage / 10)
	remainder := percentage - float64(filledCount*10)

	// Build the bar with height-variable blocks
	var bar strings.Builder

	// Full height blocks for filled portion
	for i := 0; i < filledCount && i < gradientTotalBlocks; i++ {
		bar.WriteString(gradientBlocks[7]) // Full block █
	}

	// Partial block based on remainder
	if filledCount < gradientTotalBlocks && remainder > 0 {
		// Map remainder (0-10) to block index (0-7)
		idx := int(remainder * 8 / 10)
		if idx > 7 {
			idx = 7
		}
		bar.WriteString(gradientBlocks[idx])
		filledCount++
	}

	// Empty blocks
	emptyCount := gradientTotalBlocks - filledCount
	if emptyCount > 0 {
		bar.WriteString(dimStyle.Render(strings.Repeat("░", emptyCount)))
	}

	// Render filled portion with appropriate color
	filledBar := barStyle.Render(bar.String()[:len(bar.String())-emptyCount*len("░")])
	if emptyCount > 0 {
		filledBar += dimStyle.Render(strings.Repeat("░", emptyCount))
	}

	return fmt.Sprintf("%s %d%%", filledBar, int(percentage))
}

// formatGitInfo formats git branch and changes
func (f *GradientFormatter) formatGitInfo(gitInfo *metrics.GitInfo) string {
	if !gitInfo.IsGitRepo {
		return ""
	}

	branch := branchStyle.Render(gitInfo.BranchDisplay)

	// Format changes if present
	if gitInfo.Additions > 0 || gitInfo.Deletions > 0 {
		changes := fmt.Sprintf("(+%d/-%d)", gitInfo.Additions, gitInfo.Deletions)
		return branch + " " + grayStyle.Render(changes)
	}

	return branch
}
