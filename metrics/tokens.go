package metrics

import (
	"fmt"
	"time"

	"github.com/DieGopherLT/cc-status-line/parser"
)

// TokenMetrics contains calculated token usage metrics
type TokenMetrics struct {
	TotalInputTokens  int
	TotalOutputTokens int
	TotalCachedTokens int
	TotalTokens       int
	ContextLength     int
	ContextPercentage float64
	SessionDuration   string
}

const maxContextWindow = 200000 // 200K tokens

// CalculateTokenMetrics computes token usage metrics from transcript data
func CalculateTokenMetrics(data *parser.TranscriptData) *TokenMetrics {
	metrics := &TokenMetrics{}

	var mostRecentMainChainEntry *parser.TranscriptLine
	var mostRecentTimestamp time.Time

	// Iterate through all transcript lines
	for _, line := range data.Lines {
		if line.Message == nil || line.Message.Usage == nil {
			continue
		}

		usage := line.Message.Usage

		// Sum ALL entries for total tokens
		metrics.TotalInputTokens += usage.InputTokens
		metrics.TotalOutputTokens += usage.OutputTokens
		metrics.TotalCachedTokens += usage.CacheReadInputTokens + usage.CacheCreationInputTokens

		// Find most recent main chain entry for context calculation
		if !line.IsSidechain && !line.IsAPIErrorMessage && line.Timestamp != "" {
			timestamp, err := time.Parse(time.RFC3339, line.Timestamp)
			if err == nil {
				if mostRecentTimestamp.IsZero() || timestamp.After(mostRecentTimestamp) {
					mostRecentTimestamp = timestamp
					mostRecentMainChainEntry = &line
				}
			}
		}
	}

	// Calculate total tokens
	metrics.TotalTokens = metrics.TotalInputTokens + metrics.TotalOutputTokens + metrics.TotalCachedTokens

	// Calculate context length from most recent main chain entry
	if mostRecentMainChainEntry != nil && mostRecentMainChainEntry.Message != nil && mostRecentMainChainEntry.Message.Usage != nil {
		usage := mostRecentMainChainEntry.Message.Usage
		metrics.ContextLength = usage.InputTokens + usage.CacheReadInputTokens + usage.CacheCreationInputTokens
	}

	// Calculate context percentage
	if metrics.ContextLength > 0 {
		metrics.ContextPercentage = (float64(metrics.ContextLength) / float64(maxContextWindow)) * 100.0
	}

	// Calculate session duration
	metrics.SessionDuration = calculateDuration(data.FirstTimestamp, data.LastTimestamp)

	return metrics
}

// calculateDuration formats the duration between first and last timestamp
func calculateDuration(first, last time.Time) string {
	if first.IsZero() || last.IsZero() {
		return "<1m"
	}

	duration := last.Sub(first)
	minutes := int(duration.Minutes())

	if minutes < 1 {
		return "<1m"
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	if hours == 0 {
		return formatMinutes(minutes)
	}

	if remainingMinutes == 0 {
		return formatHours(hours)
	}

	return formatHoursMinutes(hours, remainingMinutes)
}

func formatMinutes(minutes int) string {
	return formatTime(minutes, "m")
}

func formatHours(hours int) string {
	return formatTime(hours, "hr")
}

func formatTime(value int, unit string) string {
	return fmt.Sprintf("%d%s", value, unit)
}

func formatHoursMinutes(hours, minutes int) string {
	return fmt.Sprintf("%dhr %dm", hours, minutes)
}
