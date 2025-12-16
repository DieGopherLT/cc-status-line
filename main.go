package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DieGopherLT/cc-status-line/display"
	"github.com/DieGopherLT/cc-status-line/metrics"
	"github.com/DieGopherLT/cc-status-line/parser"
)

func main() {
	style := flag.String("style", "classic", "Status line style: classic, gradient, compact")
	flag.Parse()

	run(*style)
}

func run(style string) {
	// Parse status hook JSON from stdin
	hook, err := parser.ParseStatusHook(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cc-status-line error: %v\n", err)
		fmt.Println("Status: Error parsing input")
		return
	}

	// Calculate token metrics (handles nil gracefully)
	tokenMetrics := metrics.CalculateTokenMetrics(hook.ContextWindow)

	// Get git information
	gitInfo := metrics.GetGitInfo(hook.Workspace.CurrentDir)

	// Format and output status line using selected formatter
	formatter := display.NewFormatter(style)
	statusLine := formatter.Format(hook, tokenMetrics, gitInfo)
	fmt.Println(statusLine)
}
