package main

import (
	"fmt"
	"os"

	"github.com/diegopher/cc-status-line/display"
	"github.com/diegopher/cc-status-line/metrics"
	"github.com/diegopher/cc-status-line/parser"
)

func main() {
	run()
}

func run() {
	// Parse status hook JSON from stdin
	hook, err := parser.ParseStatusHook(os.Stdin)
	if err != nil {
		// Fallback: show error status
		fmt.Fprintf(os.Stderr, "cc-status-line error: %v\n", err)
		fmt.Println("Status: Error parsing input")
		return
	}

	// Parse transcript file
	transcriptData, err := parser.ParseTranscript(hook.TranscriptPath)
	if err != nil {
		// Partial fallback: show basic info without context
		fmt.Fprintf(os.Stderr, "cc-status-line warning: failed to parse transcript: %v\n", err)
		statusLine := display.FormatStatusLineMinimal(hook)
		fmt.Println(statusLine)
		return
	}

	// Calculate token metrics
	tokenMetrics := metrics.CalculateTokenMetrics(transcriptData)

	// Get git information (changes obtained directly from git)
	gitInfo := metrics.GetGitInfo(hook.Workspace.CurrentDir)

	// Format and output status line
	statusLine := display.FormatStatusLine(hook, tokenMetrics, gitInfo)
	fmt.Println(statusLine)
}
