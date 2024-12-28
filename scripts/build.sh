#!/bin/bash

# Default values
BUILD_DIR=${1:-"build"}
BINARY_NAME=${2:-"dashboard"}
OS=${GOOS:-"linux"}
ARCH=${GOARCH:-"amd64"}

# Add .exe extension for Windows
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

# Clean previous build
echo "Cleaning previous build..."
rm -rf "$BUILD_DIR/${OS}_${ARCH}"

# Create release structure
RELEASE_PATH="$BUILD_DIR/${OS}_${ARCH}"
BIN_PATH="$RELEASE_PATH/bin"
DOC_PATH="$RELEASE_PATH/docs"

echo "Creating release structure..."
mkdir -p "$BIN_PATH"
mkdir -p "$DOC_PATH"

# Build binary
echo "Building ${OS} binary..."
if [ "$OS" = "windows" ]; then
    echo "Using MinGW for Windows build..."
    export CC=x86_64-w64-mingw32-gcc
    export CXX=x86_64-w64-mingw32-g++
fi

go build -o "$BIN_PATH/$BINARY_NAME" ./cmd/dashboard

# Copy release documentation
echo "Copying release documentation..."

# Core files (stay in root)
cp README.md LICENSE "$RELEASE_PATH/"

# Documentation (goes in docs/)
if [ -d "docs" ]; then
    echo "  Copying documentation"
    cp docs/*.md "$DOC_PATH/" 2>/dev/null || true
    
    # Documentation assets
    if [ -d "docs/images" ]; then
        echo "  Copying documentation images"
        cp -r docs/images "$DOC_PATH/"
    fi
fi

echo "Release structure created at $RELEASE_PATH"
echo "  bin/  - Binary files"
echo "  docs/ - Documentation" 