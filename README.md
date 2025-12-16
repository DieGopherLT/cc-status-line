# Claude Code Status Line

A custom status line for Claude Code written in Go with beautiful colored output using lipgloss.

## Features

- 🎨 Colorized output with context-aware styling
- 📊 Visual context window usage indicator (block display)
- 🌿 Git branch and changes detection
- 🚀 Fast execution (< 100ms)
- 🎯 Accurate token metrics from transcript parsing

## Styles

The status line supports multiple visual styles, selectable via the `--style` flag.

### Classic (default)

```
Model: Sonnet 4.5 | / main | (+156 -23) | Style: default | v2.0.28 | Ctx: ████████░░ 78%
```

The original style with labeled sections and pipe separators.

### Gradient

```
Sonnet 4.5 │ ▁▂▃▄▅▆▇█░░ 78% │ main (+156/-23) │ default │ v2.0.28
```

- Height-variable blocks for context visualization
- Dynamic color: green (<50%), yellow (50-75%), red (>75%)
- Compact layout without labels

### Compact

```
◈ Sonnet 4.5  ⎔ default  ⌘ 2.0.28  ◐ 78% [████████████████░░░░]   main ↑156 ↓23
```

- Unicode icons as prefixes
- 20-character progress bar for precision
- Arrow indicators for git changes (↑ adds, ↓ dels)

### Minimal

```
Sonnet 4.5 78% main +156-23 default 2.0.28
```

- Ultra-compact, space-separated only
- No labels or decorations
- Git format: `+156-23` (no parentheses)
- Single line, no vertical spacing
- Perfect for: tmux users, tight terminal layouts, minimalists

### Nerd

```
┌──────────────────────────────────────────────────────────────────┐
│ Sonnet │ CTX: 15.5k/200k (78%) ████████░░ │  main ⇡156 ⇣23 │ 2.0.28 │
└──────────────────────────────────────────────────────────────────┘
```

- Box-drawing borders with UTF-8 characters
- **Absolute token counts**: "15.5k/200k" format
- Arrow symbols for git: ⇡ (additions) ⇣ (deletions)
- Technical panel aesthetic (like htop/btop)
- Dynamic width to fit content

## Usage

```bash
# Classic style (default)
cc-status-line

# Gradient style
cc-status-line --style gradient

# Compact style
cc-status-line --style compact

# Minimal style
cc-status-line --style minimal

# Nerd style
cc-status-line --style nerd
```

### Components

- **Model**: Current Claude model (yellow)
- **Git Branch**: Current branch or "(no git)" (red)
- **Git Changes**: Lines added/removed or "(no git)" (green for additions, red for deletions)
- **Output Style**: Current output style (dark blue)
- **Version**: Claude Code version (light blue)
- **Context**: Visual bar showing context window usage

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

Add to your Claude Code configuration file (`~/.claude/settings.json`):

**Classic style (default):**
```json
{
  "statusLine": {
    "type": "command",
    "command": "cc-status-line",
    "padding": 0
  }
}
```

**Using a different style:**
```json
{
  "statusLine": {
    "type": "command",
    "command": "cc-status-line --style compact",
    "padding": 0
  }
}
```

Available styles: `classic`, `gradient`, `compact`, `minimal`, `nerd`

## Testing

After building, test the binary with sample input:

```bash
# Test classic style (default)
cat status-line.json | cc-status-line

# Test gradient style
cat status-line.json | cc-status-line --style gradient

# Test compact style
cat status-line.json | cc-status-line --style compact

# Test minimal style
cat status-line.json | cc-status-line --style minimal

# Test nerd style
cat status-line.json | cc-status-line --style nerd
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
