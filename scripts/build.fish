#!/usr/bin/env fish

echo "Building cc-status-line..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the binary
go build -ldflags="-s -w" -o bin/cc-status-line

if test $status -eq 0
    echo "✓ Build successful: bin/cc-status-line"
else
    echo "✗ Build failed"
    exit 1
end
