#!/bin/bash

INSTALL_DIR="/usr/local/bin"
APP_NAME="clir"

# --- Script Start ---

set -e # Exit immediately if a command exits with a non-zero status.

# Function to print messages
info() {
    echo -e "\033[36mINFO:\033[0m $1"
}

error() {
    echo -e "\033[31mERROR:\033[0m $1" >&2
    exit 1
}

success() {
    echo -e "\033[32mSUCCESS:\033[0m $1"
}

# Check for required tools
command -v curl >/dev/null 2>&1 || error "curl is required but not installed. Please install curl."
command -v uname >/dev/null 2>&1 || error "uname is required but not installed."

# Determine architecture
ARCH=$(uname -m)
OS_ARCH=""

if [[ "$ARCH" == "arm64" ]]; then
    OS_ARCH="darwin-arm64"
elif [[ "$ARCH" == "x86_64" ]]; then
    OS_ARCH="darwin-amd64"
else
    error "Unsupported architecture: $ARCH. Only arm64 (Apple Silicon) and x86_64 (Intel) are supported on macOS."
fi

info "Detected architecture: $OS_ARCH"

BINARY_NAME="${APP_NAME}"

info "Preparing to download $APP_NAME..."

DOWNLOAD_URL="https://clir.pashuk.info/builds/${OS_ARCH}/${APP_NAME}"

# Temporary download path
TEMP_DOWNLOAD_PATH="/tmp/${BINARY_NAME}"

info "Downloading $APP_NAME ($OS_ARCH) from $DOWNLOAD_URL..."
curl -SL "$DOWNLOAD_URL" -o "$TEMP_DOWNLOAD_PATH"

if [[ ! -f "$TEMP_DOWNLOAD_PATH" ]] || [[ ! -s "$TEMP_DOWNLOAD_PATH" ]]; then # Check if file exists and is not empty
    error "Download failed or downloaded file is empty. Please check the URL and asset name ($APP_NAME) on clir.pashuk.info."
fi

info "Making the binary executable..."
chmod +x "$TEMP_DOWNLOAD_PATH"

INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"
info "Attempting to install $APP_NAME to $INSTALL_PATH..."

# Check if INSTALL_DIR is writable, if not, try with sudo
if [[ -w "$INSTALL_DIR" ]]; then
    mv "$TEMP_DOWNLOAD_PATH" "$INSTALL_PATH"
    info "$APP_NAME moved to $INSTALL_PATH."
else
    info "$INSTALL_DIR is not writable. Attempting with sudo..."
    if sudo mv "$TEMP_DOWNLOAD_PATH" "$INSTALL_PATH"; then
        info "$APP_NAME moved to $INSTALL_PATH using sudo."
    else
        error "Failed to move $APP_NAME to $INSTALL_PATH even with sudo. Please check permissions or try installing manually."
    fi
fi

success "$APP_NAME installed successfully to $INSTALL_PATH"
info "You can now run '$BINARY_NAME' from your terminal."

exit 0