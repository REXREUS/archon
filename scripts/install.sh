#!/bin/bash

REPO="rexreus/archon" # Replace with the correct GitHub repository
APP_NAME="archon"

echo "Installing $APP_NAME from GitHub Release..."

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
    ARCH="arm64"
fi

# Determine binary filename in release
BINARY_NAME="${APP_NAME}-${OS}-${ARCH}"
if [ "$OS" == "darwin" ]; then
    BINARY_NAME="${APP_NAME}-darwin-${ARCH}"
elif [ "$OS" == "linux" ]; then
    BINARY_NAME="${APP_NAME}-linux-${ARCH}"
fi

# Fetch latest version from GitHub API
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Failed to fetch latest version from GitHub. Ensure REPO '$REPO' is correct and has a release."
    exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$BINARY_NAME"

echo "Downloading $BINARY_NAME version $LATEST_TAG..."
curl -L -o "$APP_NAME" "$URL"

if [ $? -ne 0 ]; then
    echo "Download failed. Ensure release URL is correct: $URL"
    exit 1
fi

chmod +x "$APP_NAME"

# Move to /usr/local/bin
echo "Moving $APP_NAME to /usr/local/bin (requires sudo)..."
sudo mv "$APP_NAME" /usr/local/bin/

echo "$APP_NAME successfully installed!"
archon version
