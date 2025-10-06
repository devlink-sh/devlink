#!/bin/bash

# DevLink CLI Installer with Quantum-Resistant Infrastructure
set -e

echo "🚀 Installing DevLink CLI with Quantum-Resistant Tunneling..."

# Install zrok (quantum-resistant tunneling infrastructure)
echo "🔐 Installing zrok for quantum-resistant P2P tunnels..."
if ! command -v zrok &> /dev/null; then
    curl -sSLf https://get.zrok.io | bash
    echo "✅ zrok installed successfully"
else
    echo "ℹ️  zrok already installed"
fi

# Download DevLink from GitHub releases
echo "📥 Downloading DevLink CLI..."
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Map architecture names
case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
esac

# Download the appropriate binary
curl -L -o devlink.tar.gz "https://github.com/devlink-sh/devlink/releases/latest/download/devlink_${OS}_${ARCH}.tar.gz"
tar -xzf devlink.tar.gz
chmod +x devlink

# Install globally
echo "📦 Installing globally to /usr/local/bin..."
sudo cp devlink /usr/local/bin/

# Clean up
rm devlink devlink.tar.gz

# Setup instructions
echo ""
echo "✅ DevLink CLI installed successfully!"
echo "🔐 Quantum-resistant P2P tunneling ready!"
echo ""
echo "📋 Next steps:"
echo "   1. Initialize zrok environment: zrok enable"
echo "   2. Initialize DevLink: devlink init"
echo "   3. Start sharing: devlink --help"
echo ""
echo "💡 Example: devlink git serve"
