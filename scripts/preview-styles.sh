#!/bin/bash

# Preview all available status line styles

STYLES=("classic" "gradient" "compact" "minimal" "nerd")
INPUT_FILE="${1:-status-line.json}"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found"
    echo "Usage: $0 [input-file]"
    exit 1
fi

echo "Previewing all status line styles with: $INPUT_FILE"
echo ""

for style in "${STYLES[@]}"; do
    echo "=== $style ==="
    cat "$INPUT_FILE" | cc-status-line --style "$style"
    echo ""
done
