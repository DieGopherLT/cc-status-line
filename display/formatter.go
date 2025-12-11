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
	totalBlocks = 10
	filledBlock = "█"
	emptyBlock  = "░"
	separator   = " | "
)

// Partial block characters for sub-block precision (87.5% down to 12.5%)
var partialBlocks = []string{"▉", "▊", "▋", "▌", "▍", "▎", "▏"}

// FormatStatusLine creates the formatted status line output dynamically
func FormatStatusLine(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var segments []string

	// Model info (always present - required field)
	modelSegment := fmt.Sprintf("Model: %s", hook.Model.DisplayName)
	segments = append(segments, modelStyle.Render(modelSegment))

	// Git branch (only if git repo detected)
	if gitInfo != nil && gitInfo.IsGitRepo {
		gitBranchSegment := fmt.Sprintf("/ %s", gitInfo.BranchDisplay)
		segments = append(segments, branchStyle.Render(gitBranchSegment))

		// Git changes (only if git repo detected)
		segments = append(segments, formatGitChanges(gitInfo))
	}

	// Output style (only if present)
	if hook.OutputStyle.Name != "" {
		styleSegment := fmt.Sprintf("Style: %s", hook.OutputStyle.Name)
		segments = append(segments, styleColor.Render(styleSegment))
	}

	// Version (always present - required field)
	versionSegment := fmt.Sprintf("v%s", hook.Version)
	segments = append(segments, blueStyle.Render(versionSegment))

	// Context visualization (only if context data available)
	if tokenMetrics != nil && tokenMetrics.ContextPercentage > 0 {
		contextSegment := formatContextVisualization(tokenMetrics)
		segments = append(segments, contextSegment)
	}

	// Join all segments with separator
	statusLine := strings.Join(segments, grayStyle.Render(separator))

	// Add vertical spacing
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

// formatContextVisualization creates the context window block display with partial blocks
func formatContextVisualization(tokenMetrics *metrics.TokenMetrics) string {
	percentage := tokenMetrics.ContextPercentage

	// Calculate full blocks (each block = 10%)
	fullBlocks := int(percentage / 10)
	remainder := percentage - float64(fullBlocks*10) // 0-10

	// Build filled portion
	filledStr := strings.Repeat(filledBlock, fullBlocks)

	// Add partial block if remainder is significant (>1.25%)
	usedBlocks := fullBlocks
	if remainder > 1.25 && fullBlocks < totalBlocks {
		// Map remainder (0-10) to partial block index (0-6)
		// Each partial block represents ~1.25% (10% / 8 states)
		idx := int(remainder / 1.25)
		if idx > 6 {
			idx = 6 // Cap at the last partial block
		}
		// Reverse the array: larger remainder = earlier character (fuller)
		partialChar := partialBlocks[6-idx]
		filledStr += partialChar
		usedBlocks++
	}

	// Render filled blocks with white color
	filledBar := whiteStyle.Render(filledStr)

	// Add empty blocks with dim color
	emptyCount := totalBlocks - usedBlocks
	if emptyCount > 0 {
		emptyBar := dimStyle.Render(strings.Repeat(emptyBlock, emptyCount))
		filledBar += emptyBar
	}

	// Format: "Ctx: ████▋░░░░░ 47%"
	return fmt.Sprintf("Ctx: %s %d%%", filledBar, int(percentage))
}
