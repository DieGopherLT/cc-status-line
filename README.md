# Claude Code Status Line

A custom status line for Claude Code written in Go with beautiful colored output using lipgloss.

## Features

- 🎨 Colorized output with context-aware styling
- 📊 Visual context window usage indicator (block display)
- 🌿 Git branch and changes detection
- 🚀 Fast execution (< 100ms)
- 🎯 Accurate token metrics from transcript parsing

## Display Format

```
Model: Sonnet 4.5 | / main | (+156 -23) | Style: default | v2.0.28 | Ctx: ████████░░ 78%
```

### Components

- **Model**: Current Claude model (cyan)
- **Git Branch**: Current branch or "(no git)" (yellow)
- **Git Changes**: Lines added/removed or "(no git)" (green/red)
- **Output Style**: Current output style (magenta)
- **Version**: Claude Code version (blue)
- **Context**: Visual bar showing context window usage (color-coded)
  - Green: 0-60%
  - Yellow: 61-80%
  - Red: 81-100%

## Installation

### Build from Source

Use one of the provided build scripts based on your shell:

**Bash/Zsh:**
```bash
./scripts/build.sh
```

**Fish:**
```fish
./scripts/build.fish
```

**PowerShell:**
```powershell
./scripts/build.ps1
```

The binary will be created at `bin/cc-status-line`.

**Optional: Install globally**
```bash
# Linux/macOS
sudo cp bin/cc-status-line /usr/local/bin/

# Or add to PATH
export PATH="$PATH:$(pwd)/bin"
```

## Configuration

Add to your Claude Code configuration file (`~/.claude/config.json` or `.claude/config.json`):

**Using absolute path:**
```json
{
  "statusLine": {
    "type": "command",
    "command": "/path/to/cc-status-line/bin/cc-status-line",
    "padding": 0
  }
}
```

**If installed globally:**
```json
{
  "statusLine": {
    "type": "command",
    "command": "cc-status-line",
    "padding": 0
  }
}
```

## How It Works

The status line binary:

1. Reads JSON from stdin containing session information
2. Parses the transcript JSONL file to extract token usage
3. Calculates metrics:
   - **Total tokens**: Sum of ALL entries (main chain + sidechains)
   - **Context length**: Most recent MAIN CHAIN entry only
   - **Context percentage**: Context length / 200K tokens
4. Extracts git information using git commands
5. Formats and outputs a single-line colored status

### Context Calculation

Following the official Claude Code behavior:

- **Context window** = tokens loaded for the next request
- Only counts the most recent main chain entry (excludes agent/parallel requests)
- Includes input tokens + cached tokens
- Maximum context window: 200,000 tokens

## Project Structure

```
cc-status-line/
├── main.go              # Entry point
├── parser/
│   ├── status.go        # Status hook JSON parser
│   └── transcript.go    # JSONL transcript parser
├── metrics/
│   ├── tokens.go        # Token metrics calculation
│   └── git.go           # Git information extraction
└── display/
    └── formatter.go     # Status line formatter with lipgloss
```

## Testing

After building, test the binary with sample input:

```bash
cat test-input.json | ./bin/cc-status-line
```

Expected output format:
```
Model: Sonnet 4.5 | / main | (+156 -23) | Style: default | v2.0.28 | Ctx: ████████░░ 78%
```

## Requirements

- Go 1.21 or higher
- Git (optional, for git information display)

## Dependencies

- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT

## Author

Diego (@DieGopherLT)
