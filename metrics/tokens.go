package metrics

import (
	"github.com/DieGopherLT/cc-status-line/parser"
)

// TokenMetrics contains calculated token usage metrics
type TokenMetrics struct {
	ContextLength     int     // Input + Output tokens
	ContextPercentage float64 // Percentage of context window used
	ContextWindowSize int     // Maximum context window size
}

// CalculateTokenMetrics computes token usage metrics from context window data
func CalculateTokenMetrics(contextWindow *parser.ContextWindow) *TokenMetrics {
	metrics := &TokenMetrics{}

	// Graceful degradation: return empty metrics if context_window is nil
	if contextWindow == nil {
		return metrics
	}

	// Calculate context length: total_input + total_output
	metrics.ContextLength = contextWindow.TotalInputTokens + contextWindow.TotalOutputTokens
	metrics.ContextWindowSize = contextWindow.ContextWindowSize

	// Calculate context percentage
	if metrics.ContextWindowSize > 0 && metrics.ContextLength > 0 {
		metrics.ContextPercentage = (float64(metrics.ContextLength) / float64(metrics.ContextWindowSize)) * 100.0
	}

	return metrics
}
