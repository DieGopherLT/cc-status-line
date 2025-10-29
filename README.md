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

- **Model**: Current Claude model (yellow)
- **Git Branch**: Current branch or "(no git)" (red)
- **Git Changes**: Lines added/removed or "(no git)" (green for additions, red for deletions)
- **Output Style**: Current output style (dark blue)
- **Version**: Claude Code version (light blue)
- **Context**: Visual bar showing context window usage (white for filled blocks, dim gray for empty blocks)

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

## Testing

After building, test the binary with sample input:

```bash
cat test-input.json | ./bin/cc-status-line
```

**Normal output:**
```
Model: Sonnet 4.5 | / main | (+156 -23) | Style: default | v2.0.28 | Ctx: ████████░░ 78%
```

**Fallback output** (when context information is unavailable):
```
Model: Sonnet 4.5 | Style: default | v2.0.28 | ⚠ Context unavailable
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
