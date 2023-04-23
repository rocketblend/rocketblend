#!/bin/sh
set -e

EXECUTABLE_NAME="rktb"

API_URL="https://api.github.com/repos/rocketblend/rocketblend/releases/latest"

# Fetch the latest release info
LATEST_RELEASE=$(curl -sL "$API_URL")

# Determine the OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"
case "$OS" in
  Darwin*)  OS="Darwin" ;;
  Linux*)   OS="Linux" ;;
  *)        printf "Your OS is not supported by this script\n"; exit 1 ;;
esac
case "$ARCH" in
  x86_64) ARCH="x86_64" ;;
  *)      printf "Your architecture is not supported by this script\n"; exit 1 ;;
esac

# Find the desired asset based on the OS and Architecture
DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep -o "https://.*${EXECUTABLE_NAME}_${OS}_${ARCH}\.tar\.gz")

if [ -z "$DOWNLOAD_URL" ]; then
  printf "Failed to find the download URL for the %s version of %s\n" "$OS" "$EXECUTABLE_NAME"
  exit 1
fi

# Download and extract the tarball
TEMP_DIR=$(mktemp -d)
curl -sL "$DOWNLOAD_URL" | tar xz -C "$TEMP_DIR"

# Move the binary to a directory in PATH
INSTALL_PATH="/usr/local/bin"
sudo mv "$TEMP_DIR/$EXECUTABLE_NAME" "$INSTALL_PATH"
rm -rf "$TEMP_DIR"

printf "%s has been installed successfully!\n" "$EXECUTABLE_NAME"