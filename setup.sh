#!/bin/bash

echo "Detecting operating system and architecture..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
x86_64)
  ARCH="amd64"
  ;;
i386 | i486 | i586 | i686)
  ARCH="386"
  ;;
armv6* | armv7*)
  ARCH="arm"
  ;;
aarch64 | arm64)
  ARCH="arm64"
  ;;
mips*)
  ARCH=$(echo "$ARCH" | tr '[:upper:]' '[:lower:]')
  ;;
*)
  echo "Unsupported architecture: $ARCH"
  exit 1
  ;;
esac

if [ "$(id -u)" != "0" ]; then
  if [ "$OS" == "windows" ]; then
    echo "This script must be run with administrator privileges on Windows."
    echo "Right-click on the terminal and select 'Run as administrator', then try again."
  else
    echo "This script must be run with root privileges."
    echo "Please run this script with sudo or as the root user."
  fi
  exit 1
fi

if [ "$OS" == "mingw"* ]; then
  OS="windows"
  TARGET_PATH="$SYSTEMROOT\System32"
else
  TARGET_PATH="/usr/local/bin"
fi

echo "Fetching the latest version number..."
VERSION=$(curl -s https://api.github.com/repos/yi-ge/file-sync/releases/latest | grep "tag_name" | cut -d '"' -f 4)
if [ -z "$VERSION" ]; then
  echo "Failed to fetch the latest version number"
  exit 1
fi

if [ "$OS" == "windows" ]; then
  FILENAME="file-sync-${OS}-${ARCH}.exe"
else
  FILENAME="file-sync-${OS}-${ARCH}"
fi

echo "Downloading $FILENAME version $VERSION..."

DOWNLOAD_URL="https://github.com/yi-ge/file-sync/releases/download/${VERSION}/${FILENAME}"
curl -L -O "${DOWNLOAD_URL}" || {
  echo "Download failed"
  exit 1
}

echo "Moving ${FILENAME} to ${TARGET_PATH}..."

if [ "$OS" == "windows" ]; then
  chmod +x "${FILENAME}"
  mv "${FILENAME}" "${TARGET_PATH}\\file-sync" || {
    echo "Move failed"
    exit 1
  }
else
  chmod +x "${FILENAME}"
  sudo mv "${FILENAME}" "${TARGET_PATH}/file-sync" || {
    echo "Move failed"
    exit 1
  }
fi

echo "Checking if the file-sync is working properly..."
file-sync -v || {
  echo "The file-sync command does not work as expected"
  exit 1
}

echo "Registering file-sync as a service..."
sudo file-sync service enable || {
  echo "Failed to register file-sync as a service, but continuing anyway."
  true
}

echo "Starting file-sync service..."
sudo file-sync service start || {
  echo "Failed to start file-sync service"
  exit 1
}

echo "Done! File-sync has been successfully set up, registered, and started as a service."
