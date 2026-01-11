package formatters

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Fractional blocks for horizontal fill (1/8 increments, empty to full)
var HorizontalBlocks = []string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}

// Fractional blocks for vertical fill (1/8 increments, empty to full)
var VerticalBlocks = []string{"", "▁", "▂", "▃", "▄", "▅", "▆", "▇"}

const (
	fullBlock  = "█"
	emptyBlock = "░"
)

// RenderProgressBar creates a progress bar with 1/8 character precision.
// Parameters:
//   - percentage: 0-100 value representing fill level
//   - totalBlocks: number of character positions for the bar
//   - fractionalBlocks: slice of partial block characters (HorizontalBlocks or VerticalBlocks)
//   - filledStyle: lipgloss style for filled portion
//   - emptyStyle: lipgloss style for empty portion
//
// Returns a styled string with the progress bar.
func RenderProgressBar(
	percentage float64,
	totalBlocks int,
	fractionalBlocks []string,
	filledStyle, emptyStyle lipgloss.Style,
) string {
	// Clamp percentage to valid range
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	// Calculate total segments (8 per block for 1/8 precision)
	totalSegments := totalBlocks * 8
	filledSegments := int(percentage * float64(totalSegments) / 100.0)

	// Calculate full blocks and fractional remainder
	fullBlocks := filledSegments / 8
	remainder := filledSegments % 8

	// Build filled portion
	var filledPart strings.Builder

	// Add full blocks
	filledPart.WriteString(strings.Repeat(fullBlock, fullBlocks))

	// Track how many character positions we've used
	usedBlocks := fullBlocks

	// Add fractional block if there's a remainder and space available
	if remainder > 0 && usedBlocks < totalBlocks {
		filledPart.WriteString(fractionalBlocks[remainder])
		usedBlocks++
	}

	// Render filled portion with style
	result := filledStyle.Render(filledPart.String())

	// Add empty blocks
	emptyCount := totalBlocks - usedBlocks
	if emptyCount > 0 {
		result += emptyStyle.Render(strings.Repeat(emptyBlock, emptyCount))
	}

	return result
}
