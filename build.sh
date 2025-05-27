#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

APP_NAME="clir"

echo "🚀 Starting build for $APP_NAME..."

# Define architectures
ARCHS=("darwin/arm64" "darwin/amd64")

# Output directory
OUTPUT_DIR_BASE="site/builds"

for ARCH_PAIR in "${ARCHS[@]}"; do
  GOOS=$(echo $ARCH_PAIR | cut -d'/' -f1)
  GOARCH=$(echo $ARCH_PAIR | cut -d'/' -f2)
  ARCH_DIR_SUFFIX=$(echo $ARCH_PAIR | tr '/' '-') # e.g., darwin-arm64

  OUTPUT_DIR="${OUTPUT_DIR_BASE}/${ARCH_DIR_SUFFIX}"
  OUTPUT_PATH="${OUTPUT_DIR}/${APP_NAME}" # Renamed binary

  echo "🛠️ Building for $GOOS/$GOARCH..."
  mkdir -p "$OUTPUT_DIR"
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT_PATH" main.go
  echo "✅ Built $APP_NAME for $GOOS/$GOARCH at $OUTPUT_PATH"
done

echo "🎉 All builds completed for $APP_NAME!"
echo "Find your binaries in the '$OUTPUT_DIR_BASE' directory, organized by architecture."