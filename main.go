package main

import (
	"fmt"
	"os"

	"github.com/diegopher/cc-status-line/display"
	"github.com/diegopher/cc-status-line/metrics"
	"github.com/diegopher/cc-status-line/parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse status hook JSON from stdin
	hook, err := parser.ParseStatusHook(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to parse status hook: %w", err)
	}

	// Parse transcript file
	transcriptData, err := parser.ParseTranscript(hook.TranscriptPath)
	if err != nil {
		return fmt.Errorf("failed to parse transcript: %w", err)
	}

	// Calculate token metrics
	tokenMetrics := metrics.CalculateTokenMetrics(transcriptData)

	// Get git information
	gitInfo := metrics.GetGitInfo(hook.Workspace.CurrentDir, hook.Cost.TotalLinesAdded, hook.Cost.TotalLinesRemoved)

	// Format and output status line
	statusLine := display.FormatStatusLine(hook, tokenMetrics, gitInfo)
	fmt.Println(statusLine)

	return nil
}
