#!/bin/bash

# Default version
DEFAULT_VERSION="v0.0.5"

# Get the version from the command line argument or use default
VERSION=${1:-$DEFAULT_VERSION}

# Define the base URL for the release artifacts
BASE_URL="https://github.com/struckchure/gv/releases/download/${VERSION}"

# Define the file names (adjust these as needed)
LINUX_AMD64="gv_Linux_x86_64.tar.gz"
LINUX_ARM64="gv_Linux_arm64.tar.gz"
MACOS_AMD64="gv_Darwin_x86_64.tar.gz"
MACOS_ARM64="gv_Darwin_arm64.tar.gz"

# Determine the OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Set the file to download based on OS and architecture
case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                FILE="$LINUX_AMD64"
                ;;
            aarch64)
                FILE="$LINUX_ARM64"
                ;;
            *)
                echo "Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64)
                FILE="$MACOS_AMD64"
                ;;
            arm64)
                FILE="$MACOS_ARM64"
                ;;
            *)
                echo "Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Define the destination directory
DEST_DIR="$HOME/.gv/bin"

# Create the destination directory if it does not exist
mkdir -p "$DEST_DIR"

# Download the file
echo "Downloading $FILE..."
curl -fsSL "${BASE_URL}/${FILE}" -o "${FILE}"

# Extract the downloaded file to the .gv directory
echo "Extracting $FILE to $DEST_DIR..."
tar -xzf "$FILE" -C "$DEST_DIR"

# Remove the downloaded file
rm "$FILE"

# Add .gv directory to PATH if not already present
PATH_ENTRY="export PATH=\"$DEST_DIR:\$PATH\""
PROFILE_FILE="$HOME/.bashrc"

# Check if we are running on Alpine Linux
if [ -f /etc/alpine-release ]; then
    PROFILE_FILE="$HOME/.profile"
fi

if ! grep -Fxq "$PATH_ENTRY" "$PROFILE_FILE"; then
    echo "$PATH_ENTRY" >> "$PROFILE_FILE"
    echo "Updated $PROFILE_FILE to include $DEST_DIR in PATH"
else
    echo "$DEST_DIR is already in PATH in $PROFILE_FILE"
fi
