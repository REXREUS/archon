$BuildDir = "build"
$AppName = "archon"

if (!(Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir
}

Write-Host "Building ArchonCLI for multiple platforms..." -ForegroundColor Cyan

# Windows
Write-Host "Building for Windows (amd64)..."
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o "$BuildDir/$AppName-windows-amd64.exe" ./cmd/archon/main.go
Write-Host "Building for Windows (arm64)..."
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o "$BuildDir/$AppName-windows-arm64.exe" ./cmd/archon/main.go

# Linux
Write-Host "Building for Linux (amd64)..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o "$BuildDir/$AppName-linux-amd64" ./cmd/archon/main.go
Write-Host "Building for Linux (arm64)..."
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o "$BuildDir/$AppName-linux-arm64" ./cmd/archon/main.go

# macOS (Darwin)
Write-Host "Building for macOS (amd64)..."
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o "$BuildDir/$AppName-darwin-amd64" ./cmd/archon/main.go
Write-Host "Building for macOS (arm64)..."
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o "$BuildDir/$AppName-darwin-arm64" ./cmd/archon/main.go

# Reset env vars
$env:GOOS=""; $env:GOARCH=""

Write-Host "Build complete! Results available in $BuildDir/ folder" -ForegroundColor Green
