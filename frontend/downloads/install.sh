#!/bin/bash

# DevLink CLI Installer for localhost:8080
set -e

echo "🚀 Installing DevLink CLI..."

# Download the binary from your localhost
echo "📥 Downloading DevLink CLI from localhost:8080..."
curl -L -o devlink "http://localhost:8080/downloads/devlink"

# Make it executable
chmod +x devlink

# Install globally
echo "📦 Installing globally to /usr/local/bin..."
sudo cp devlink /usr/local/bin/

# Clean up
rm devlink

echo "✅ DevLink CLI installed successfully!"
echo "🎉 You can now use 'devlink' command from anywhere!"
echo "💡 Run 'devlink --help' to see available commands"
