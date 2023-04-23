#!/bin/bash
set -e

EXECUTABLE_NAME="rocketblend"
API_URL="https://api.github.com/repos/rocketblend/rocketblend/releases/latest"

printf "Fetching the latest release info...\n"
LATEST_RELEASE=$(curl -sL "${API_URL}")

printf "Determining the OS and Architecture...\n"
OS="$(uname -s)"
ARCH="$(uname -m)"
case "${OS}" in
  Darwin*)  OS="Darwin" ;;
  Linux*)   OS="Linux" ;;
  *)        printf "Your OS is not supported by this script\n"; exit 1 ;;
esac
case "${ARCH}" in
  x86_64) ARCH="x86_64" ;;
  *)      printf "Your architecture is not supported by this script\n"; exit 1 ;;
esac

printf "Finding the desired asset based on the OS and Architecture...\n"
DOWNLOAD_URL=$(echo "${LATEST_RELEASE}" | grep -o "https://.*${EXECUTABLE_NAME}_${OS}_${ARCH}\.tar\.gz" || true)

# Fail if the download URL could not be found or grep failed
if [ -z "${DOWNLOAD_URL}" ]; then
  printf "Failed to find the download URL or parse the release information for the %s version of %s\n" "${OS}" "${EXECUTABLE_NAME}"
  exit 1
fi

printf "Downloading and extracting the tarball...\n"
TEMP_DIR=$(mktemp -d)
if ! curl -sL "${DOWNLOAD_URL}" | tar xz -C "${TEMP_DIR}"; then
  printf "Failed to download and extract the release\n"
  rm -rf "${TEMP_DIR}"
  exit 1
fi

printf "Moving the binary to a directory in PATH...\n"
INSTALL_PATH="/usr/local/bin"
sudo mv "${TEMP_DIR}/${EXECUTABLE_NAME}" "${INSTALL_PATH}"
rm -rf "${TEMP_DIR}"

printf "%s has been installed successfully!\n" "${EXECUTABLE_NAME}"
printf "You may need to restart your terminal session to refresh your PATH\n"