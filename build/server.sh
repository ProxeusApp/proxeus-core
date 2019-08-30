#!/bin/bash
set -Eeuxo pipefail

DEBUG_FLAG=""
if [[ "$(printenv DEBUG || true)" == "1" ]]
then
    echo "bindata set to DEBUG mode"
    DEBUG_FLAG="-debug"
fi

LDFLAGS=""
if [[ -n "${BUILD_ID:-}" ]]
then
    LDFLAGS="-X main.ServerVersion=build-${BUILD_ID}"
fi

# make sure go-bindata is up to date
go install ./vendor/github.com/asticode/go-bindata/go-bindata

if [[ ! -d ./artifacts/dist ]]; then
  make ui
fi

go-bindata ${DEBUG_FLAG} -pkg assets -o ./main/handlers/assets/bindata.go -prefix ./artifacts/dist ./artifacts/dist/...

if  [[ "$#" > 0 && -n "$1" && "$1" == "--docker" ]]; then
    echo "Building server for docker"
    if [[ "$OSTYPE" == "linux-gnu" ]]; then
        go build -ldflags="${LDFLAGS}" -tags nocgo -o ./artifacts/server-docker ./main #build for host os
    else
        # cross compile for docker container
        echo "cross compile"
        go get github.com/karalabe/xgo
        xgo -ldflags="${LDFLAGS}" -tags nocgo -dest ./artifacts -targets=linux/amd64 ./main && mv ./artifacts/main-linux-amd64 ./artifacts/server-docker #build for docker
    fi
else
    echo "Building server"
    go build -ldflags="${LDFLAGS}" -tags nocgo -o ./artifacts/server ./main #build for host os
fi

