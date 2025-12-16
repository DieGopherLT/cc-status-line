#!/usr/bin/env fish

# Preview all available status line styles

set styles classic gradient compact minimal nerd
set input_file (test -n "$argv[1]" && echo "$argv[1]" || echo "status-line.json")

# Check if input file exists
if not test -f "$input_file"
    echo "Error: Input file '$input_file' not found"
    echo "Usage: $0 [input-file]"
    exit 1
end

echo "Previewing all status line styles with: $input_file"
echo ""

for style in $styles
    echo "=== $style ==="
    cat $input_file | cc-status-line --style $style
    echo ""
end
