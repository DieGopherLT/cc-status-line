#!/usr/bin/env pwsh

Write-Host "Building cc-status-line..." -ForegroundColor Cyan

# Create bin directory if it doesn't exist
if (-not (Test-Path -Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build the binary
$output = "bin/cc-status-line"
if ($IsWindows) {
    $output = "bin\cc-status-line.exe"
}

go build -ldflags="-s -w" -o $output

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful: $output" -ForegroundColor Green
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}
