#!/bin/bash
set -Eeuo pipefail

# is dep in sync
# dep check

# imports

installed () {
    which $1
}
require () {
    if ! installed $1; then echo "Please install $1"; exit 1; fi
}
require goimports
require gofmt


if [[ "$(goimports -l -local git.proxeus.com main sys\
 | grep -v bindata.go \
 | wc -l)" -ne "0" ]]
then
    goimports -l -local git.proxeus.com main sys
    echo "code not formatted, run make fmt to fix"; exit 1;
fi

# formatting

#gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
#[ -z "$gofiles" ] && exit 0

gofiles=$(find . -regex "^.*\.go$" \
    | grep -v "^./vendor/" \
    | grep -v "^./artifacts/" \
    | grep -v "/bindata.go$" \
)

unformatted=$(gofmt -l ${gofiles})
[[ -z "${unformatted}" ]] && exit 0

echo >&2 "Go files must be formatted with gofmt. Please run:"
for fn in ${unformatted}; do
    echo >&2 "  gofmt -w $PWD/$fn"
done

exit 1
