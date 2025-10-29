package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
	"github.com/muesli/termenv"
)

func init() {
	// Force TrueColor output even when stdout is not a TTY
	lipgloss.SetColorProfile(termenv.TrueColor)
}

// Color definitions matching the screenshot
var (
	modelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow for model
	branchStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red for git branch
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // Green for additions
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203")) // Red for deletions
	styleColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("24"))  // Dark desaturated blue for output style
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("111")) // Blue for version
	grayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242")) // Gray for separator
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // Dim gray for empty blocks
	whiteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("255")) // White for context bar
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
	segments = append(segments, modelStyle.Render(modelSegment))

	// Git branch: "/ branch-name" or "/ no git"
	gitBranchSegment := fmt.Sprintf("/ %s", gitInfo.BranchDisplay)
	segments = append(segments, branchStyle.Render(gitBranchSegment))

	// Git changes: "(+156 -23)" or "(no git)"
	segments = append(segments, formatGitChanges(gitInfo))

	// Output style: "Style: default"
	styleSegment := fmt.Sprintf("Style: %s", hook.OutputStyle.Name)
	segments = append(segments, styleColor.Render(styleSegment))

	// Version: "v2.0.28"
	versionSegment := fmt.Sprintf("v%s", hook.Version)
	segments = append(segments, blueStyle.Render(versionSegment))

	// Context visualization: "Ctx: ████████░░ 78%"
	contextSegment := formatContextVisualization(tokenMetrics)
	segments = append(segments, contextSegment)

	// Join all segments with separator
	statusLine := strings.Join(segments, grayStyle.Render(separator))

	// Add vertical spacing (newlines before and after)
	return "\n" + statusLine + "\n"
}

// FormatStatusLineMinimal creates a minimal status line without context info (for error fallback)
func FormatStatusLineMinimal(hook *parser.StatusHook) string {
	var segments []string

	// Model info
	modelSegment := fmt.Sprintf("Model: %s", hook.Model.DisplayName)
	segments = append(segments, modelStyle.Render(modelSegment))

	// Output style
	styleSegment := fmt.Sprintf("Style: %s", hook.OutputStyle.Name)
	segments = append(segments, styleColor.Render(styleSegment))

	// Version
	versionSegment := fmt.Sprintf("v%s", hook.Version)
	segments = append(segments, blueStyle.Render(versionSegment))

	// Warning indicator (use orange/yellow for warnings)
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	segments = append(segments, warningStyle.Render("⚠ Context unavailable"))

	// Join all segments with separator
	statusLine := strings.Join(segments, grayStyle.Render(separator))

	// Add vertical spacing (newlines before and after)
	return "\n" + statusLine + "\n"
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

	// Build the block display (use white for filled blocks)
	filled := strings.Repeat(filledBlock, filledCount)
	empty := strings.Repeat(emptyBlock, emptyCount)

	blocks := whiteStyle.Render(filled) + dimStyle.Render(empty)

	// Format: "Ctx: ████████░░ 78%"
	return fmt.Sprintf("Ctx: %s %d%%", blocks, int(percentage))
}
