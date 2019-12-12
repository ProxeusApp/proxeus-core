SHELL:= /bin/bash
DEBUG_FLAG?=false
PROXEUS_TEST_MODE?=false
GO_VERSION=1.13

ifeq ($(DEBUG), "true")
	BINDATA_OPTS="-debug"
endif

ifdef BUILD_ID
    GO_OPTS=-ldflags="-X main.ServerVersion=build-$(BUILD_ID)"
endif

ifeq ($(shell uname), Darwin)
	DOCKER_LINUX=docker run --rm -v "$(PWD):/usr/src" -w /usr/src golang:$(GO_VERSION)
endif
#########################################################
dependencies=go curl
mocks=main/handlers/blockchain/mock/adapter_mock.go sys/db/storm/mock/workflow_payments_mock.go sys/db/storm/mock/user_mock.go sys/db/storm/mock/workflow_mock.go 

.PHONY: all
all: ui server

.PHONY: init
init:
	@for d in $(dependencies); do (echo "Checking $$d is installed... " && which $$d ) || ( echo "Please install $$d before continuing" && exit 1 ); done
	go install golang.org/x/tools/cmd/goimports
	go install github.com/asticode/go-bindata/go-bindata
	go install github.com/golang/mock/mockgen

.PHONY: ui
ui:
	make -C ui

.PHONY: ui-dev
ui-dev:
	make -C ui serve-main-hosted

.PHONY: bindata
bindata: main/handlers/assets/bindata.go  

main/handlers/assets/bindata.go: $(wildcard ./ui/core/dist/**)
	go-bindata ${BINDATA_OPTS} -pkg assets -o ./main/handlers/assets/bindata.go -prefix ./ui/core/dist ./ui/core/dist/...

.PHONY: mock
mock: $(mocks)

.PHONY: server
server: bindata mock
	go build $(GO_OPTS) -tags nocgo -o ./artifacts/server ./main 

.PHONY: server-docker
server-docker: bindata mock
	$(DOCKER_LINUX) go build $(GO_OPTS) -tags nocgo -o ./artifacts/server-docker ./main

.PHONY: validate
validate:
	@if [[ "$$(goimports -l -local git.proxeus.com main sys | grep -v bindata.go | tee /dev/tty | wc -l | xargs)" != "0" ]]; then \
		echo "Format validation error.  Please run make fmt"; exit 1; \
	fi
	@echo "Format validated"

.PHONY: fmt
fmt:
	goimports -w -local git.proxeus.com main sys

.PHONY: test
test: main/handlers/assets/bindata.go
	go test  ./main/... ./sys/... 

.PHONY:test-payment
test-payment:
	go test ./main/handlers/payment  ./main/handlers/blockchain

test/bindata.go: $(wildcard ./test/assets/**)
	go-bindata ${BINDATA_OPTS} -pkg test -o ./test/bindata.go ./test/assets

.PHONY: test-api
test-api: test/bindata.go
	go clean -testcache && go test ./test

.PHONY: coverage
coverpkg=$(filter-out %/assets %/mock, $(shell go list ./main/... ./sys/...))
coverage: bindata mock
	go test -v -tags coverage -coverprofile artifacts/cover.out -coverpkg="$(coverpkg)" ./main

.PHONY: print-coverage
print-coverage:
	go tool cover -html artifacts/cover.out

.PHONY: clean
clean:
	cd artifacts && rm -rf `ls . | grep -v 'cache'`
	cd ui && yarn cache clean && cd ..

.PHONY: run
run:
	artifacts/server -DataDir ./data/proxeus-platform/data/  -DocumentServiceUrl=http://document-service:2115 \
		-BlockchainContractAddress=${PROXEUS_CONTRACT_ADDRESS} -InfuraApiKey=${PROXEUS_INFURA_KEY} \
		-SparkpostApiKey=${PROXEUS_SPARKPOST_KEY} -EmailFrom=${PROXEUS_EMAIL_FROM} -TestMode=${PROXEUS_TEST_MODE}

.SECONDEXPANSION: # See https://www.gnu.org/software/make/manual/make.html#Secondary-Expansion
$(mocks): $$(patsubst %_mock.go, %.go, $$(subst /mock,, $$@))
	mockgen -package mock  -source $<  -destination $@  -self_package github.com/ProxeusApp/proxeus-core/$(shell dirname $@)
	goimports -w $@
