# Copilot Instructions for cc-status-line

A Go CLI tool that generates a colorized status line for Claude Code, reading JSON from stdin and outputting formatted terminal text with lipgloss.

## Architecture Overview

```
stdin (JSON) → parser → metrics → display → stdout (styled text)
```

**Data Flow in [main.go](../main.go):**
1. `parser.ParseStatusHook(os.Stdin)` - Parses Claude Code's status hook JSON
2. `parser.ParseTranscript(hook.TranscriptPath)` - Reads JSONL transcript file
3. `metrics.CalculateTokenMetrics()` - Calculates context window usage
4. `metrics.GetGitInfo()` - Executes git commands for branch/changes
5. `display.FormatStatusLine()` - Produces styled output with lipgloss

## Package Responsibilities

| Package | Purpose | Key Details |
|---------|---------|-------------|
| `parser/` | JSON/JSONL parsing | 10MB buffer for large thinking blocks; skips malformed lines |
| `metrics/` | Token & git calculations | Context from **most recent main chain entry** (non-sidechain, non-error) |
| `display/` | Terminal styling | Forces TrueColor mode; 10-block visual context indicator |

## Key Implementation Details

**Context Window Calculation** ([metrics/tokens.go](../metrics/tokens.go)):
- Uses most recent non-sidechain, non-error entry's tokens
- Formula: `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`
- Hardcoded 200K token max window

**Git Integration** ([metrics/git.go](../metrics/git.go)):
- `git rev-parse --abbrev-ref HEAD` for branch
- `git diff --numstat HEAD` for changes (staged + unstaged)
- Returns "(no git)" gracefully when not in a repo

**Color Scheme** ([display/formatter.go](../display/formatter.go)):
- Yellow (226): model name
- Red (196/203): branch, deletions  
- Green (76): additions
- Blue (111): version
- Gray (242): separators

## Build & Test

```bash
./scripts/build.sh          # Output: bin/cc-status-line
cat test-input.json | ./bin/cc-status-line   # Manual test
```

## Error Handling Pattern

The tool uses graceful degradation:
- **Fatal** (stdin parse fails): stderr error + generic status
- **Partial** (transcript parse fails): stderr warning + `FormatStatusLineMinimal()` without context

## Coding Conventions

- Structs with JSON tags define Claude Code's API contract in `parser/status.go`
- Use lipgloss styles defined as package-level vars in `display/formatter.go`
- Token metrics always come from transcript parsing, never from stdin hook's `Cost` field
- Git commands execute in `hook.Workspace.CurrentDir`
