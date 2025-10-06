#!/usr/bin/env bash
set -euo pipefail

VERSION="0.1.0"  # Change this when you cut a new release
DIST_DIR="dist"

echo "🚀 Building DevLink CLI v$VERSION for Homebrew..."

# Clean and recreate dist directory
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# Build for macOS Intel (amd64)
echo "🔨 Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/devlink" ./cmd/devlink
tar -czf "$DIST_DIR/devlink_${VERSION}_darwin_amd64.tar.gz" -C "$DIST_DIR" devlink

# Build for macOS ARM (arm64)
echo "🔨 Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/devlink" ./cmd/devlink
tar -czf "$DIST_DIR/devlink_${VERSION}_darwin_arm64.tar.gz" -C "$DIST_DIR" devlink

# Compute SHA256 checksums
echo "🔑 Computing SHA256 checksums..."
shasum -a 256 "$DIST_DIR"/*.tar.gz
echo
echo "✅ Done! Upload the two tar.gz files in $DIST_DIR to your GitHub Release."
