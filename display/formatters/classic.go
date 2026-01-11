package formatters

import (
	"fmt"
	"strings"

	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
	"github.com/charmbracelet/lipgloss"
)

const (
	classicTotalBlocks = 10
	classicSeparator   = " | "
)

// ClassicFormatter implements the original status line style
type ClassicFormatter struct{}

// Format creates the formatted status line output in the classic style
func (f *ClassicFormatter) Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	var segments []string

	// Model info (always present - required field)
	modelSegment := fmt.Sprintf("Model: %s", hook.Model.DisplayName)
	segments = append(segments, modelStyle.Render(modelSegment))

	// Git branch (only if git repo detected)
	if gitInfo != nil && gitInfo.IsGitRepo {
		gitBranchSegment := fmt.Sprintf("/ %s", gitInfo.BranchDisplay)
		segments = append(segments, branchStyle.Render(gitBranchSegment))

		// Git changes (only if git repo detected)
		segments = append(segments, f.formatGitChanges(gitInfo))
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
		contextSegment := f.formatContextVisualization(tokenMetrics)
		segments = append(segments, contextSegment)
	}

	// Join all segments with separator
	statusLine := strings.Join(segments, grayStyle.Render(classicSeparator))

	// Add subtle horizontal lines for visual breathing room
	width := lipgloss.Width(statusLine)
	line := lineStyle.Render(strings.Repeat("─", width))

	return line + "\n" + statusLine + "\n" + line
}

// formatGitChanges formats git changes with appropriate colors
func (f *ClassicFormatter) formatGitChanges(gitInfo *metrics.GitInfo) string {
	if !gitInfo.IsGitRepo {
		return grayStyle.Render("(no git)")
	}

	text := gitInfo.ChangesText
	if strings.Contains(text, "+") && strings.Contains(text, "-") {
		// Mixed changes: "(+156 -23)"
		inner := strings.Trim(text, "()")
		parts := strings.Split(inner, " ")
		if len(parts) == 2 {
			addPart := greenStyle.Render(parts[0])
			remPart := redStyle.Render(parts[1])
			return grayStyle.Render("(") + addPart + grayStyle.Render(" ") + remPart + grayStyle.Render(")")
		}
	} else if strings.Contains(text, "+") {
		return grayStyle.Render("(") + greenStyle.Render(text[1:len(text)-1]) + grayStyle.Render(")")
	} else if strings.Contains(text, "-") {
		return grayStyle.Render("(") + redStyle.Render(text[1:len(text)-1]) + grayStyle.Render(")")
	}

	return grayStyle.Render(text)
}

// formatContextVisualization creates the context window block display with 1/8 precision
func (f *ClassicFormatter) formatContextVisualization(tokenMetrics *metrics.TokenMetrics) string {
	percentage := tokenMetrics.ContextPercentage

	bar := RenderProgressBar(percentage, classicTotalBlocks, HorizontalBlocks, whiteStyle, dimStyle)

	return fmt.Sprintf("Ctx: %s %d%%", bar, int(percentage))
}
