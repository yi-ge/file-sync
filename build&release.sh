#!/bin/bash
VERSION=0.1.11
NAME=file-sync

if [ "$(uname)" == "Darwin" ]; then
  if ! [ -x "$(command -v gsed)" ]; then
    echo 'Notice: sed is not installed, missing dependencies being installed.'
    brew install gnu-sed
  fi
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  if ! [ -x "$(command -v sed)" ]; then
    echo 'Notice: sed is not installed, missing dependencies being installed.'
    apt install sed -y # osslsigncode
  fi
fi

echo "// Do not modify this file.
package config

const (
  name = \"${NAME}\"
  version = \"${VERSION}\"
)
" >config/auto_config.go

if [ "$(uname)" == "Darwin" ]; then
  gsed -i 's/isDev[[:space:]]=[[:space:]]true/isDev = false/g' main.go
else
  sed -i 's/isDev[[:space:]]=[[:space:]]true/isDev = false/g' main.go
fi

if [ ! -d bin ]; then
  mkdir bin
fi

build() {
  echo "env CGO_ENABLED=0 GOOS=$1 GOARCH=$2 $3 go build -ldflags \"-s -w\" -o ./bin/file-sync-$1-$2$4 ."
  env CGO_ENABLED=0 GOOS=$1 GOARCH=$2 $3 go build -ldflags "-s -w" -o "./bin/file-sync-$1-$2$4" .
}

if [ "$1" == "--build" ]; then
  current_os="$(uname | awk '{print tolower($0)}')"
  current_arch="$(uname -m)"
  build "$current_os" "$current_arch"
elif [ "$1" == "--release" ]; then
  build darwin arm64
  build darwin amd64
  build freebsd 386
  build freebsd amd64
  build linux 386
  build linux amd64
  build linux arm GOARM=7
  build linux arm64
  build linux mips64
  build linux mips64le
  build linux mips GOMIPS=softfloat
  build linux mipsle GOMIPS=softfloat
  build windows 386 "" .exe
  build windows amd64 "" .exe
else
  echo "Invalid argument. Use '--build' to build only, or '--release' to create a release."
  exit 1
fi

if [ "$(uname)" == "Darwin" ]; then
  gsed -i 's/isDev[[:space:]]=[[:space:]]false/isDev = true/g' main.go
else
  sed -i 's/isDev[[:space:]]=[[:space:]]false/isDev = true/g' main.go
fi

if [ "$1" == "--release" ]; then
  # Ensure GitHub CLI is installed
  if ! [ -x "$(command -v gh)" ]; then
    echo 'Error: GitHub CLI is not installed.' >&2
    echo 'Please visit https://github.com/cli/cli#installation for installation instructions.' >&2
    exit 1
  fi

  # Set your GitHub repository
  GITHUB_USER="file-sync"
  GITHUB_REPO="yi-ge"

  # Create a release tag
  TAG_NAME="v${VERSION}"
  TAG_MESSAGE="Release ${NAME} ${TAG_NAME}"

  # Check if tag already exists
  if git rev-parse "$TAG_NAME" >/dev/null 2>&1; then
    echo "Error: Tag ${TAG_NAME} already exists."
    exit 1
  fi

  git tag -a "$TAG_NAME" -m "$TAG_MESSAGE"
  git push origin "$TAG_NAME"

  # Create a new GitHub release using the tag
  gh release create "$TAG_NAME" --title "$TAG_MESSAGE" --notes "Release notes for ${NAME} ${TAG_NAME}"

  # Upload the compiled files to the release
  for file in ./bin/*; do
    gh release upload "$TAG_NAME" "$file"
  done

  echo "Release ${TAG_NAME} has been created and the compiled files have been uploaded."
fi
