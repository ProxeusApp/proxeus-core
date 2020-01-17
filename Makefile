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
mocks=main/handlers/blockchain/mock/adapter_mock.go
bindata=main/handlers/assets/bindata.go test/bindata.go

.PHONY: all
all: ui server

.PHONY: init
init:
	@for d in $(dependencies); do (echo "Checking $$d is installed... " && which $$d ) || ( echo "Please install $$d before continuing" && exit 1 ); done
	go install golang.org/x/tools/cmd/goimports
	go install github.com/asticode/go-bindata/go-bindata
	go install github.com/golang/mock/mockgen
	go install github.com/wadey/gocovmerge

.PHONY: ui
ui:
	$(MAKE) -C ui

.PHONY: ui-dev
ui-dev:
	$(MAKE) -C ui serve-main-hosted

.PHONY: generate
generate: $(bindata) $(mocks) storage/database/mock/mocks.go

.PHONY: server
server: generate
	go build $(GO_OPTS) -tags nocgo -o ./artifacts/server ./main 

.PHONY: server-docker
server-docker: generate
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
test: generate 
	go test  ./main/... ./sys/... ./storage/...

.PHONY: test-api
test-api: generate
	go clean -testcache && go test ./test

.PHONY: coverage
comma:=,
space:= $() $()
coverpkg=$(subst $(space),$(comma), $(filter-out %/mock %/assets, $(shell go list ./main/... ./sys/... ./storage/...)))
coverage: generate
	go test -coverprofile artifacts/cover_unittests.out ./main/... ./sys/... ./storage/...
	go test -v -tags coverage -coverprofile artifacts/cover_integration.out -coverpkg="$(coverpkg)" ./main

.PHONY: print-coverage
print-coverage:
	gocovmerge artifacts/cover_unittests.out artifacts/cover_integration.out > artifacts/cover_merged.out
	go tool cover -func artifacts/cover_merged.out
	go tool cover -html artifacts/cover_merged.out

.PHONY: clean
clean:
	cd artifacts && rm -rf `ls . | grep -v 'cache'`
	cd ui && yarn cache clean && cd ..

.PHONY: run
run:
	artifacts/server -DataDir ./data/proxeus-platform/data/  -DocumentServiceUrl=http://document-service:2115 \
		-BlockchainContractAddress=${PROXEUS_CONTRACT_ADDRESS} -InfuraApiKey=${PROXEUS_INFURA_KEY} \
		-SparkpostApiKey=${PROXEUS_SPARKPOST_KEY} -EmailFrom=${PROXEUS_EMAIL_FROM} -TestMode=${PROXEUS_TEST_MODE}

main/handlers/assets/bindata.go: $(wildcard ./ui/core/dist/**)
	go-bindata ${BINDATA_OPTS} -pkg assets -o ./main/handlers/assets/bindata.go -prefix ./ui/core/dist ./ui/core/dist/...
	goimports -w $@

test/bindata.go: $(wildcard ./test/assets/**)
	go-bindata ${BINDATA_OPTS} -pkg test -o ./test/bindata.go ./test/assets
	goimports -w $@

.SECONDEXPANSION: # See https://www.gnu.org/software/make/manual/make.html#Secondary-Expansion
$(mocks): $$(patsubst %_mock.go, %.go, $$(subst /mock,, $$@))
	mockgen -package mock  -source $<  -destination $@  -self_package github.com/ProxeusApp/proxeus-core/$(shell dirname $@)
	goimports -w $@

storage/database/mock/mocks.go: storage/interfaces.go
	mockgen -package mock  -source storage/interfaces.go -destination $@  -self_package github.com/ProxeusApp/proxeus-core/$(shell dirname $@)
	goimports -w $@