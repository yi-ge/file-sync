#/bin/bash
VERSION=0.1.0
NAME=file-sync

if ! [ -x "$(command -v upx)" ]; then
  echo 'Notice: upx is not installed, missing dependencies being installed.'
  if [ "$(uname)" == "Darwin" ]; then
    brew install upx
  elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    apt update && apt install upx -y # osslsigncode
  fi
fi

if ! [ -x "$(command -v sed)" ]; then
  echo 'Notice: sed is not installed, missing dependencies being installed.'
  if [ "$(uname)" == "Darwin" ]; then
    brew install gnu-sed
  elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
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
  upx ./bin/file-sync-darwin-arm64
else
  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-amd64 ."
  env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-amd64 .

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-386 ."
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-386 .
  # upx ./bin/file-sync-freebsd-386

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-amd64 ."
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-amd64 .
  # upx ./bin/file-sync-freebsd-amd64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-386 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-linux-386 .
  upx ./bin/file-sync-linux-386

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-amd64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-amd64 .
  upx ./bin/file-sync-linux-amd64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-armv7l ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o ./bin/file-sync-linux-armv7l .
  upx ./bin/file-sync-linux-armv7l

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-arm64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-arm64 .
  # upx ./bin/file-sync-linux-arm64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64 ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64 .
  # upx ./bin/file-sync-linux-mips64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64le ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64le .
  # upx ./bin/file-sync-linux-mips64le

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips .
  upx ./bin/file-sync-linux-mips

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mipsle ."
  env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mipsle .
  upx ./bin/file-sync-linux-mipsle

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-386.exe ."
  env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-windows-386.exe .
  upx ./bin/file-sync-windows-386.exe

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-amd64.exe ."
  env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-windows-amd64.exe .
  upx ./bin/file-sync-windows-amd64.exe
fi

sed -i 's/isDev[[:space:]]=[[:space:]]false/isDev = os.Getenv("GO_ENV") == "development"/g' main.go
