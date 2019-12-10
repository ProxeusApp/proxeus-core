PROXEUS_TEST_MODE?=false

# required tooling
init:
	./build/deps.sh

main: ui server

ui:
	cd ui && yarn cache clean && cd .. && make -C ui main-hosted

ui-dev:
	make -C ui serve-main-hosted

server:
	./build/server.sh

server-docker:
	GOOS=linux GOARCH=amd64 ./build/server.sh --docker

validate:
	./build/validate.sh

validate-ui:
	make -C ui validate

fmt:
	goimports -w -local git.proxeus.com main sys

test:
	./build/test.sh

test-payment:
	go test ./main/handlers/payment  ./main/handlers/blockchain

test/bindata.go: $(wildcard ./test/assets/**)
	go-bindata ${DEBUG_FLAG} -pkg test -o ./test/bindata.go ./test/assets

test-api: test/bindata.go
	go clean -testcache && go test ./test

coverage:
	build/coverage.sh

print-coverage:
	grep -v "_mock.go" artifacts/cover.out > artifacts/cover2.out
	gocovmerge artifacts/cover_unittests.out artifacts/cover2.out > artifacts/cover_merged.out
	go tool cover -func artifacts/cover_merged.out
	go tool cover -html artifacts/cover_merged.out

clean:
	cd artifacts && rm -rf `ls . | grep -v 'cache'`

all: ui server

run:
	artifacts/server -DataDir ./data/proxeus-platform/data/  -DocumentServiceUrl=http://document-service:2115 \
		-BlockchainContractAddress=${PROXEUS_CONTRACT_ADDRESS} -InfuraApiKey=${PROXEUS_INFURA_KEY} \
		-SparkpostApiKey=${PROXEUS_SPARKPOST_KEY} -EmailFrom=${PROXEUS_EMAIL_FROM} -TestMode=${PROXEUS_TEST_MODE}

# used by CI
link-repo:
	mkdir -p /go/src/git.proxeus.com/core
	ln -s ${PROJECT_ROOT_FOLDER} /go/src/git.proxeus.com/core/central

.PHONY: init main all all-debug generate test clean fmt validate link-repo
.PHONY: ui server
.PHONY: coverage print-coverage