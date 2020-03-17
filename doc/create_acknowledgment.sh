#/bin/bash

# Comand line tool to generate the ACKNOWLEDGEMENT file
#
# ./create_acknowledgment.sh <path to dependency_decisions.yml file> ACKNOWLEDGEMENT

set -o nounset -o errexit -o pipefail 

which -s hub || echo Please install the hub command line tool from https://hub.github.com
which -s yq || echo Please install the yq YAML command line tool from https://github.com/kislyuk/yq
which -s jq || echo Please install the jq JSON command line tool from https://stedolan.github.io/jq

dependency_file=${1:-dependency_decisions.yml}
workdir=$(mktemp -d /tmp/create-acknowledgement.XXXXX)
trap "rm -fr ${workdir}; exit" INT TERM EXIT

# module to git repo map
key(){
    echo $(md5 -rq -s $1)
}

map(){
    export map_$(key $1)=$2
}

get(){
    local k=map_$(key $1)
    echo ${!k}
}

map go.etcd.io/bbolt github.com/etcd-io/bbolt
map golang.org/x/crypto github.com/golang/crypto
map golang.org/x/image github.com/golang/image
map google.golang.org/appengine github.com/golang/appengine
map gopkg.in/gavv/httpexpect.v2 github.com/gavv/httpexpect
map go.mongodb.org/mongo-driver github.com/mongodb/mongo-go-driver
map go.starlark.net github.com/google/starlark-go
map golang.org/x/arch github.com/golang/arch
map golang.org/x/net github.com/golang/net
map golang.org/x/sys github.com/golang/sys
map golang.org/x/tools github.com/golang/tools
map gopkg.in/yaml.v2 github.com/go-yaml/yaml

#####################

license() {
    local repo=$1
    [[ "x$(get ${repo})" != "x" ]] && repo=$(get ${repo})

    local dir=${workdir}/$(basename ${repo})
    local url=https://${repo}

    [[ -d ${dir} ]] && return # Exit if already processed

    git clone --bare --depth=1 ${url} ${dir}
    cd ${dir}

    local project_name=$(hub api repos/{owner}/{repo} | jq -r .full_name)
    local license_url=$(hub api repos/{owner}/{repo}/license | jq .download_url | xargs echo)

    echo "############################################"
    echo Project name : ${project_name} 
    echo Project repository: ${url}

    if [[ "x${license_url}" != "xnull" ]]; then
        echo License file URL: ${license_url} 
        echo License:
        echo
        curl ${license_url}
    fi
    echo
    echo "############################################"
    echo
    echo
}


for repo in `cat "${dependency_file}" | yq -r .[][1]`
do
    license ${repo}
done
