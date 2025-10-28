#!/bin/bash

set -e

echo "Building cc-status-line..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the binary
go build -ldflags="-s -w" -o bin/cc-status-line

echo "✓ Build successful: bin/cc-status-line"
