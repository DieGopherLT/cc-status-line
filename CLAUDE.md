# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go CLI tool that generates a colorized status line for Claude Code, reading JSON from stdin and outputting formatted terminal text with lipgloss.

**Pipeline:**

```
stdin (JSON) → parser → metrics → display → stdout (styled text)
```

## Build and Test Commands

**Build:**

```bash
./scripts/build.sh      # Bash/Zsh (output: bin/cc-status-line)
./scripts/build.fish    # Fish
./scripts/build.ps1     # PowerShell
```

**Test:**

```bash
cat test-input.json | ./bin/cc-status-line
```

**Install globally (optional):**

```bash
sudo cp bin/cc-status-line /usr/local/bin/
```

## Architecture

### Data Flow (main.go:16-44)

1. `parser.ParseStatusHook(os.Stdin)` - Parse Claude Code's status hook JSON
2. `parser.ParseTranscript(hook.TranscriptPath)` - Read JSONL transcript file
3. `metrics.CalculateTokenMetrics(transcriptData)` - Calculate context window usage
4. `metrics.GetGitInfo(hook.Workspace.CurrentDir)` - Execute git commands
5. `display.FormatStatusLine(...)` - Produce styled output

### Package Responsibilities

| Package | Purpose | Key Details |
|---------|---------|-------------|
| **parser/** | JSON/JSONL parsing | 10MB buffer for large thinking blocks; skips malformed lines silently |
| **metrics/** | Token & git calculations | Context from **most recent main chain entry** (non-sidechain, non-error); git uses `diff --numstat HEAD` |
| **display/** | Terminal styling | Forces TrueColor mode; 10-block visual context indicator (█/░) |

### Key Implementation Details

**Context Window Calculation (metrics/tokens.go:24-73)**

- Uses most recent non-sidechain, non-error entry's tokens
- Formula: `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`
- Percentage calculated against 200K token window (hardcoded constant)

**Git Integration (metrics/git.go:20-53)**

- Branch: `git rev-parse --abbrev-ref HEAD`
- Changes: `git diff --numstat HEAD` (staged + unstaged)
- Gracefully returns "(no git)" when not in a repository
- Handles binary files (marked with "-" in numstat)

**Color Scheme (display/formatter.go:19-29)**

- Yellow (226): model name
- Red (196/203): branch, deletions
- Green (76): additions
- Blue (24): output style
- Blue (111): version
- Gray (242): separators

**Error Handling Pattern**

- **Fatal** (stdin parse fails): stderr error + generic status
- **Partial** (transcript parse fails): stderr warning + `FormatStatusLineMinimal()` without context
- **Graceful**: shows "(no git)" or "(no changes)" when git commands fail

## Code Organization

```text
.
├── main.go                 # Entry point and orchestration
├── parser/
│   ├── status.go          # Status hook JSON parsing
│   └── transcript.go      # Transcript JSONL parsing (10MB buffer)
├── metrics/
│   ├── tokens.go          # Token usage calculations
│   └── git.go             # Git information extraction
└── display/
    └── formatter.go       # Terminal output formatting (lipgloss)
```

## Coding Conventions

- **Structs with JSON tags** in parser/status.go define Claude Code's API contract
- **Lipgloss styles** defined as package-level vars in display/formatter.go
- **Token metrics always come from transcript parsing**, never from stdin hook's `Cost` field
- **Git commands execute in** `hook.Workspace.CurrentDir`
- **Transcript parsing** uses 10MB scanner buffer to handle large thinking blocks

## Dependencies

- `github.com/charmbracelet/lipgloss` - Terminal styling framework
- `github.com/muesli/termenv` - Terminal environment detection (used by lipgloss)

## Important Constants

- `maxContextWindow = 200000` (metrics/tokens.go:21) - Claude's context window size
- `totalBlocks = 10` (display/formatter.go:32) - Number of blocks in context visualization
- `maxCapacity = 10 * 1024 * 1024` (parser/transcript.go:52) - Scanner buffer size
