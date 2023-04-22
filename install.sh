#!/bin/bash

set -e

REPO="https://github.com/rocketblend/rocketblend"
APP_NAME="rktb"
RELEASES_API="https://api.github.com/repos/$REPO/releases/latest"

# Determine the platform and architecture
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm*)
        ARCH="arm"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Download the appropriate binary for the platform and architecture
DOWNLOAD_URL="$(curl -s $RELEASES_API | grep "browser_download_url.*${OS}-${ARCH}" | cut -d : -f 2,3 | tr -d \" | tr -d ' ')"
DESTINATION="/usr/local/bin/$APP_NAME"

echo "Downloading $APP_NAME for $OS-$ARCH..."
curl -L -o "$DESTINATION" "$DOWNLOAD_URL"

# Set the executable permissions
chmod +x "$DESTINATION"

echo "Installation complete. $APP_NAME is now available in $DESTINATION"