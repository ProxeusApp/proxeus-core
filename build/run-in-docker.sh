#!/bin/bash
set -Eeuxo pipefail

# usage: run-in-docker.sh image_tag shell_command

if [[ "$(printenv OSXCROSS_REPO || true)" != "" ]]; then
    # this script is executed both from CI and manually (through make),
    # on CI it runs already inside prepared docker image
    echo "Already in builder docker so skip wrapping..."
    /bin/sh -c "${2}"
    exit 0
fi

# docker chowns some files as root, bring them back
# TODO: ideally mapped volume should be used with non-root privileges to remove chown fix
function chown-artifacts {
    sudo chown -R $(whoami) ./artifacts
}
trap chown-artifacts EXIT

sudo docker run --rm \
    --workdir /go/src/git.proxeus.com/core/central \
    -v $(go env GOPATH)/src/git.proxeus.com/core/central:/go/src/git.proxeus.com/core/central \
    "${1}" /bin/sh -c \
    "${2}"
