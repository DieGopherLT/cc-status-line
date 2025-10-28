package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/diegopher/cc-status-line/metrics"
	"github.com/diegopher/cc-status-line/parser"
)

// Color definitions matching the screenshot
var (
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))  // Cyan for model
	yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("220")) // Yellow for git branch
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // Green for additions
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203")) // Red for deletions
	magentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")) // Magenta for style
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("111")) // Blue for version
	grayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242")) // Gray for separator
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // Dim gray for empty blocks
)

const (
	totalBlocks  = 10
	filledBlock  = "█"
	emptyBlock   = "░"
	separator    = " | "
)

// FormatStatusLine creates the formatted status line output
func FormatStatusLine(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var segments []string

	// Model info: "Model: Sonnet 4.5"
	modelSegment := fmt.Sprintf("Model: %s", hook.Model.DisplayName)
	segments = append(segments, cyanStyle.Render(modelSegment))

	// Git branch: "/ branch-name" or "/ no git"
	gitBranchSegment := fmt.Sprintf("/ %s", gitInfo.BranchDisplay)
	segments = append(segments, yellowStyle.Render(gitBranchSegment))

	// Git changes: "(+156 -23)" or "(no git)"
	segments = append(segments, formatGitChanges(gitInfo))

	// Output style: "Style: default"
	styleSegment := fmt.Sprintf("Style: %s", hook.OutputStyle.Name)
	segments = append(segments, magentaStyle.Render(styleSegment))

	// Version: "v2.0.28"
	versionSegment := fmt.Sprintf("v%s", hook.Version)
	segments = append(segments, blueStyle.Render(versionSegment))

	// Context visualization: "Ctx: ████████░░ 78%"
	contextSegment := formatContextVisualization(tokenMetrics)
	segments = append(segments, contextSegment)

	// Join all segments with separator
	return strings.Join(segments, grayStyle.Render(separator))
}

// formatGitChanges formats git changes with appropriate colors
func formatGitChanges(gitInfo *metrics.GitInfo) string {
	if !gitInfo.IsGitRepo {
		return grayStyle.Render("(no git)")
	}

	text := gitInfo.ChangesText
	if strings.Contains(text, "+") && strings.Contains(text, "-") {
		// Mixed changes: "(+156 -23)"
		// Strip outer parentheses first
		inner := strings.Trim(text, "()")
		parts := strings.Split(inner, " ")
		if len(parts) == 2 {
			addPart := greenStyle.Render(parts[0])
			remPart := redStyle.Render(parts[1])
			return grayStyle.Render("(") + addPart + grayStyle.Render(" ") + remPart + grayStyle.Render(")")
		}
	} else if strings.Contains(text, "+") {
		// Only additions
		return grayStyle.Render("(") + greenStyle.Render(text[1:len(text)-1]) + grayStyle.Render(")")
	} else if strings.Contains(text, "-") {
		// Only deletions
		return grayStyle.Render("(") + redStyle.Render(text[1:len(text)-1]) + grayStyle.Render(")")
	}

	return grayStyle.Render(text)
}

// formatContextVisualization creates the context window block display
func formatContextVisualization(tokenMetrics *metrics.TokenMetrics) string {
	percentage := tokenMetrics.ContextPercentage
	filledCount := min(int((percentage / 100.0) * float64(totalBlocks)), totalBlocks)

	emptyCount := totalBlocks - filledCount

	// Determine color based on percentage
	var blockStyle lipgloss.Style
	switch {
	case percentage >= 81:
		blockStyle = redStyle
	case percentage >= 61:
		blockStyle = yellowStyle
	default:
		blockStyle = greenStyle
	}

	// Build the block display
	filled := strings.Repeat(filledBlock, filledCount)
	empty := strings.Repeat(emptyBlock, emptyCount)

	blocks := blockStyle.Render(filled) + dimStyle.Render(empty)

	// Format: "Ctx: ████████░░ 78%"
	return fmt.Sprintf("Ctx: %s %d%%", blocks, int(percentage))
}
