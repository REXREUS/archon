#!/bin/bash

APP_NAME="archon"
INSTALL_PATH="/usr/local/bin/$APP_NAME"

echo "Removing $APP_NAME..."

if [ -f "$INSTALL_PATH" ]; then
    sudo rm "$INSTALL_PATH"
    echo "$APP_NAME has been removed from $INSTALL_PATH"
else
    echo "$APP_NAME not found in $INSTALL_PATH"
fi

# Optional: delete configuration
read -p "Delete configuration folder (~/.archon.yaml and vector database)? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf ~/.archon.yaml
    rm -rf ./chromem_db
    echo "Configuration and database have been deleted."
fi

echo "Uninstall complete."
