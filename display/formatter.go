package display

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/DieGopherLT/cc-status-line/display/formatters"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
	"github.com/muesli/termenv"
)

func init() {
	// Force TrueColor output even when stdout is not a TTY
	lipgloss.SetColorProfile(termenv.TrueColor)
}

// StatusLineFormatter defines the interface for formatting status lines
type StatusLineFormatter interface {
	Format(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string
}

// NewFormatter creates a formatter based on the style name
func NewFormatter(style string) StatusLineFormatter {
	switch style {
	case "gradient":
		return &formatters.GradientFormatter{}
	case "compact":
		return &formatters.CompactFormatter{}
	case "minimal":
		return &formatters.MinimalFormatter{}
	case "nerd":
		return &formatters.NerdFormatter{}
	default:
		return &formatters.ClassicFormatter{}
	}
}

// FormatStatusLine is a convenience function that uses the classic formatter
// Kept for backward compatibility
func FormatStatusLine(hook *parser.StatusHook, tokenMetrics *metrics.TokenMetrics, gitInfo *metrics.GitInfo) string {
	formatter := &formatters.ClassicFormatter{}
	return formatter.Format(hook, tokenMetrics, gitInfo)
}
