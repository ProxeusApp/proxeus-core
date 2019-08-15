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

if [ ! -d ./artifacts/dist ]; then
  make ui
fi

go-bindata ${DEBUG_FLAG} -pkg assets -o ./main/handlers/assets/bindata.go -prefix ./artifacts/dist ./artifacts/dist/...

go build -ldflags="${LDFLAGS}" -tags nocgo -o ./artifacts/server ./main
