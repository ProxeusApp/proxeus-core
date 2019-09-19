#!/bin/bash
set -Eeuxo pipefail

installed () {
    which $1
}

require () {
    if ! installed $1; then echo "Please manually install $1"; exit 1; fi
}
require go
require curl
require npm
require yarn

# install golang's dep
mkdir -p $(go env GOPATH)/bin
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

go get golang.org/x/tools/cmd/goimports

echo "Installing dependencies. This may take a while."
dep ensure

if [[ "$OSTYPE" == "darwin"* ]]; then
	sudo mkdir -p /data/hosted
	sudo chown $(id -u):$(id -g) /data/hosted
else
	mkdir -p /data/hosted
fi


# on macs use brew
if installed brew; then
    clang=$(go env GOPATH)/bin/o64-clang
    if ! [[ -e "$clang" ]]; then
        ln -s /usr/bin/clang "$clang" # symlink to be compatible with dapp/bundler_darwin.json
    fi
fi
