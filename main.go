package main

import (
	"fmt"
	"os"

	"github.com/DieGopherLT/cc-status-line/display"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
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

	// Calculate token metrics (handles nil gracefully)
	tokenMetrics := metrics.CalculateTokenMetrics(hook.ContextWindow)

	// Get git information
	gitInfo := metrics.GetGitInfo(hook.Workspace.CurrentDir)

	// Format and output status line
	statusLine := display.FormatStatusLine(hook, tokenMetrics, gitInfo)
	fmt.Println(statusLine)
}
