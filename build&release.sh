#/bin/bash
VERSION=0.1.0
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
  alias sed='gsed'
fi

sed -i 's/isDev[[:space:]]=[[:space:]]os\.Getenv("GO_ENV")[[:space:]]==[[:space:]]"development"/isDev = false/g' main.go

if [ ! -d bin ]; then
  mkdir bin
fi

if [ "$1" == "test" ]; then
  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-arm64 ."
  env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-arm64 .
else
  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-arm64 ."
  env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-arm64 .

  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-amd64 ."
  env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-amd64 .

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-386 ."
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-386 .

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-amd64 ."
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-amd64 .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-386 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-linux-386 .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-amd64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-amd64 .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-armv7l ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o ./bin/file-sync-linux-armv7l .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-arm64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-arm64 .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64 .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64le ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64le .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips .

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mipsle ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mipsle .

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-386.exe ."
  env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-windows-386.exe .

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-amd64.exe ."
  env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-windows-amd64.exe .
fi

sed -i 's/isDev[[:space:]]=[[:space:]]false/isDev = os.Getenv("GO_ENV") == "development"/g' main.go

# Ensure GitHub CLI is installed
if ! [ -x "$(command -v gh)" ]; then
  echo 'Error: GitHub CLI is not installed.' >&2
  echo 'Please visit https://github.com/cli/cli#installation for installation instructions.' >&2
  exit 1
fi

# Set your GitHub repository
GITHUB_USER="file-sync"
GITHUB_REPO="yi-ge"

# Check if the git repository is clean
if ! git diff-index --quiet HEAD --; then
  echo "Error: Your git repository contains uncommitted changes. Please commit or stash them before proceeding."
  exit 1
fi

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
