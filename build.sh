#/bin/bash
VERSION=0.1.0
NAME=file-sync

if ! [ -x "$(command -v upx)" ]; then
  echo 'Error: upx is not installed.'
  if [ "$(uname)" == "Darwin" ]; then
    brew install upx
  elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    apt update && apt install upx -y # osslsigncode
  fi
fi

if ! [ -x "$(command -v sed)" ]; then
  echo 'Error: sed is not installed.'
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

sed -i 's/isDev[[:space:]]=[[:space:]]os\.Getenv\("GO_ENV"\)[[:space:]]==[[:space:]]"development"/isDev[[:space:]]=[[:space:]]false/g' main.go

if [ ! -d bin ]; then
  mkdir bin
fi

if [ "$1" == "test" ]; then
  mv nac.syso nac.syso.back
  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-amd64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-amd64 server-monitor-agent
  upx ./bin/file-sync-darwin-amd64
  mv nac.syso.back nac.syso
else
  mv nac.syso nac.syso.back
  echo "env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-darwin-amd64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-darwin-amd64 server-monitor-agent
  # upx ./bin/file-sync-darwin-amd64

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-386 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-386 server-monitor-agent
  # upx ./bin/file-sync-freebsd-386

  echo "env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-freebsd-amd64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-freebsd-amd64 server-monitor-agent
  # upx ./bin/file-sync-freebsd-amd64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-386 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-linux-386 server-monitor-agent
  upx ./bin/file-sync-linux-386

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-amd64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-amd64 server-monitor-agent
  upx ./bin/file-sync-linux-amd64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-armv7l server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o ./bin/file-sync-linux-armv7l server-monitor-agent
  upx ./bin/file-sync-linux-armv7l

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-arm64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-arm64 server-monitor-agent
  # upx ./bin/file-sync-linux-arm64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64 server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64 server-monitor-agent
  # upx ./bin/file-sync-linux-mips64

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips64le server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips64le server-monitor-agent
  # upx ./bin/file-sync-linux-mips64le

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mips server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mips server-monitor-agent
  upx ./bin/file-sync-linux-mips

  echo "env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags \"-s -w\" -o ./bin/file-sync-linux-mipsle server-monitor-agent"
  env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-s -w" -o ./bin/file-sync-linux-mipsle server-monitor-agent
  upx ./bin/file-sync-linux-mipsle

  mv nac.syso.back nac.syso

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-386.exe server-monitor-agent"
  env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/file-sync-windows-386.exe server-monitor-agent
  upx ./bin/file-sync-windows-386.exe

  echo "env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags \"-s -w\" -o ./bin/file-sync-windows-amd64.exe server-monitor-agent"
  env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/file-sync-windows-amd64.exe server-monitor-agent
  upx ./bin/file-sync-windows-amd64.exe
fi

sed -i 's/isDev\s=\sfalse/isDev\s=\sos\.Getenv\("GO_ENV"\)\s==\s"development"/g' main.go
