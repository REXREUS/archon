#!/bin/bash

# Output folder
BUILD_DIR="build"
APP_NAME="archon"

# Create build folder if it doesn't exist
mkdir -p $BUILD_DIR

echo "Building ArchonCLI for multiple platforms..."

# Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/$APP_NAME-windows-amd64.exe ./cmd/archon/main.go
echo "Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -o $BUILD_DIR/$APP_NAME-windows-arm64.exe ./cmd/archon/main.go

# Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/$APP_NAME-linux-amd64 ./cmd/archon/main.go
echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o $BUILD_DIR/$APP_NAME-linux-arm64 ./cmd/archon/main.go

# macOS (Darwin)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/$APP_NAME-darwin-amd64 ./cmd/archon/main.go
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/$APP_NAME-darwin-arm64 ./cmd/archon/main.go

echo "Build complete! Results available in $BUILD_DIR/ folder"
