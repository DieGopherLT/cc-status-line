package formatters

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	// Force TrueColor output even when stdout is not a TTY
	lipgloss.SetColorProfile(termenv.TrueColor)
}

// Shared color definitions for all formatters
var (
	modelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208")) // Claude orange for model
	branchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red for git branch
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // Green for additions
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("203")) // Red for deletions
	styleColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("24"))  // Dark desaturated blue for output style
	blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("111")) // Blue for version
	grayStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("242")) // Gray for separator
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // Dim gray for empty blocks
	whiteStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("255")) // White for context bar
	lineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("232")) // Almost black for border lines
)
