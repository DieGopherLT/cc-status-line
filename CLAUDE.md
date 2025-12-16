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

### Data Flow (main.go:16-35)

1. `parser.ParseStatusHook(os.Stdin)` - Parse Claude Code's status hook JSON
2. `metrics.CalculateTokenMetrics(hook.ContextWindow)` - Calculate context window usage from status hook
3. `metrics.GetGitInfo(hook.Workspace.CurrentDir)` - Execute git commands
4. `display.FormatStatusLine(...)` - Produce styled output

### Package Responsibilities

| Package | Purpose | Key Details |
|---------|---------|-------------|
| **parser/** | JSON parsing | Parses status hook JSON including `context_window` with `current_usage` |
| **metrics/** | Token & git calculations | Context from `current_usage` field; git uses `diff --numstat HEAD` |
| **display/** | Terminal styling | Forces TrueColor mode; 10-block visual context indicator (█/░) |

### Key Implementation Details

**Context Window Calculation (metrics/tokens.go:25-32)**

- Uses `current_usage` from status hook's `context_window` field
- Formula: `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`
- Window size from `context_window.context_window_size` (dynamic, not hardcoded)
- Fallback to `total_input_tokens + total_output_tokens` if `current_usage` is nil

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
│   └── status.go          # Status hook JSON parsing (includes ContextWindow, CurrentUsage)
├── metrics/
│   ├── tokens.go          # Token usage calculations
│   └── git.go             # Git information extraction
└── display/
    └── formatter.go       # Terminal output formatting (lipgloss)
```

## Coding Conventions

- **Structs with JSON tags** in parser/status.go define Claude Code's API contract
- **Lipgloss styles** defined as package-level vars in display/formatter.go
- **Token metrics come from** `context_window.current_usage` in the status hook JSON
- **Git commands execute in** `hook.Workspace.CurrentDir`

## Dependencies

- `github.com/charmbracelet/lipgloss` - Terminal styling framework
- `github.com/muesli/termenv` - Terminal environment detection (used by lipgloss)

## Important Constants

- `totalBlocks = 10` (display/formatter.go:32) - Number of blocks in context visualization
