#!/bin/bash
set -Eeuxo pipefail

installed () {
    which $1
}

# for linux install npm and curl
if installed apt-get; then
    apt-get install curl;
    # install go
    curl https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz > go.tar.gz
    tar -xf go.tar.gz
    rm go.tar.gz
    rm -Rf /usr/local/go
    mv go /usr/local
    export PATH=/usr/local/go/bin:$PATH
    echo "----- Please add /usr/local/go/bin to your PATH -----"
    # install node
    curl -sL https://deb.nodesource.com/setup_11.x | sudo -E bash -
    apt-get install -y nodejs
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
    apt-get update && sudo apt-get install yarn
fi


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
mkdir -p /data/hosted


# on macs use brew
if installed brew; then
    clang=$(go env GOPATH)/bin/o64-clang
    if ! [[ -e "$clang" ]]; then
        ln -s /usr/bin/clang "$clang" # symlink to be compatible with dapp/bundler_darwin.json
    fi
fi
