SHELL:= /bin/bash
DEBUG_FLAG?=false
GO_VERSION=1.13

ifeq ($(DEBUG), "true")
	BINDATA_OPTS="-debug"
endif

ifdef BUILD_ID
    GO_OPTS=-ldflags="-X main.ServerVersion=build-$(BUILD_ID)"
endif

DOCKER_GATEWAY=172.17.0.1
ifeq ($(shell uname), Darwin)
	DOCKER_LINUX=docker run --rm -v "$(PWD):/usr/src" -w /usr/src golang:$(GO_VERSION)
	DOCKER_GATEWAY=host.docker.internal
endif

# Default proxeus environment
export PROXEUS_TEST_MODE?=false
export PROXEUS_ALLOW_HTTP?=true
export PROXEUS_SETTINGS_FILE?=~/.proxeus/settings/main.json
export PROXEUS_PLATFORM_DOMAIN?=http://localhost:1323
export PROXEUS_DOCUMENT_SERVICE_URL?=http://localhost:2115
export PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS?=${PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS}
export PROXEUS_INFURA_API_KEY?=${PROXEUS_INFURA_API_KEY}
export PROXEUS_SPARKPOST_API_KEY?=${PROXEUS_SPARKPOST_API_KEY}
export PROXEUS_EMAIL_FROM=no-reply@proxeus.com
export PROXEUS_DATA_DIR?=./data
export PROXEUS_DATABASE_ENGINE?=storm
export PROXEUS_DATABASE_URI?=mongodb://localhost:27017
#Only used for the mongo engine

#########################################################

#########################################################

# Coverage
coverage?=false
comma:=,
space:= $() $()
coverpkg=$(subst $(space),$(comma), $(filter-out %/mock %/assets, $(shell go list ./main/... ./sys/... ./storage/... ./service/...)))

startproxeus=PROXEUS_DATA_DIR=$(1)/data PROXEUS_SETTINGS_FILE=$(1)/settings/main.json PROXEUS_TEST_MODE=true artifacts/proxeus &
stopproxeus=pkill proxeus
startds=curl -s http://localhost:2115 > /dev/null || ( PROXEUS_DATA_DIR=$(1) docker-compose up -d document-service && touch $(1)/ds-started )
startnodes=curl -s http://localhost:8011 > /dev/null || (PROXEUS_PLATFORM_DOMAIN=http://$(DOCKER_GATEWAY):1323 NODE_CRYPTO_RATES_URL=http://localhost:8011 REGISTER_RETRY_INTERVAL=1 docker-compose up -d node-crypto-forex-rates && touch $(1)/nodes-started )
startmongo=nc -z localhost 27017 2> /dev/null || (docker run -d -p 27017:27017 -p 27018:27018 -p 27019:27019 proxeus/mongo-dev-cluster && sleep 10 && touch $(1)/mongo-started)

ifeq ($(coverage),true)
	COVERAGE_OPTS=-coverprofile artifacts/$@.coverage -coverpkg="$(coverpkg)"
	startproxeus=PROXEUS_DATA_DIR=$(1)/data PROXEUS_SETTINGS_FILE=$(1)/settings/main.json PROXEUS_TEST_MODE=true go test -v -tags coverage -coverprofile artifacts/$@-$(PROXEUS_DATABASE_ENGINE).coverage -coverpkg="$(coverpkg)" ./main &
	stopproxeus=pkill main.test
endif

dependencies=go curl
mocks=main/handlers/blockchain/mock/adapter_mock.go storage/mock/interfaces.go
bindata=main/handlers/assets/bindata.go test/assets/bindata.go
golocalimport=github.com/ProxeusApp/proxeus-core

.PHONY: all
all: ui server

.PHONY: init
init:
	@for d in $(dependencies); do (echo "Checking $$d is installed... " && which $$d ) || ( echo "Please install $$d before continuing" && exit 1 ); done
	go install golang.org/x/tools/cmd/goimports
	go install github.com/asticode/go-bindata/go-bindata
	go install github.com/golang/mock/mockgen
	go install github.com/wadey/gocovmerge
	go install golang.org/x/tools/cmd/godoc

.PHONY: ui
ui:
	$(MAKE) -C ui

.PHONY: ui-dev
ui-dev:
	$(MAKE) -C ui serve-main-hosted

.PHONY: generate
generate: $(bindata) $(mocks)

.PHONY: server
server: generate
	go build $(GO_OPTS) -tags nocgo -o ./artifacts/proxeus ./main

.PHONY: server-docker
server-docker: generate
	$(DOCKER_LINUX) go build $(GO_OPTS) -tags nocgo -o ./artifacts/proxeus-docker ./main

.PHONY: validate
validate: init
	@if [[ "$$(goimports -l -local $(golocalimport) main sys storage service | grep -v bindata.go | tee /dev/tty | wc -l | xargs)" != "0" ]]; then \
		echo "Format validation error.  Please run make fmt"; exit 1; \
	fi
	@echo "Format validated"

.PHONY: license
license:
	# https://github.com/pivotal/LicenseFinder
	license_finder

.PHONY: doc
doc: init server
	$(eval serverurl=localhost:6061)
	# godoc only reads documentation from downloaded packages, therefore we simulate a src structure. godoc doesn't follow symlinks
	mkdir -p /tmp/tmpgoroot/doc
	mkdir -p /tmp/tmpgopath/src/git.proxeus.com/
	cp -R * /tmp/tmpgopath/src/git.proxeus.com/
	GOPATH=/tmp/tmpgopath GO111MODULE=on godoc -goroot=/tmp/tmpgoroot -http=$(serverurl) &
	sleep 3
	# Download css & js first
	wget -P artifacts/$(serverurl)/lib/godoc http://$(serverurl)/lib/godoc/style.css
	wget -P artifacts/$(serverurl)/lib/godoc http://$(serverurl)/lib/godoc/jquery.js
	wget -P artifacts/$(serverurl)/lib/godoc http://$(serverurl)/lib/godoc/godocs.js
	# Now, only the package we're interested into. not the whole standard library
	wget -r -P artifacts -np -e robots=off "http://$(serverurl)/pkg/git.proxeus.com/"
	mkdir -p artifacts/godoc/lib/godoc
	cp -r artifacts/$(serverurl)/pkg/git.proxeus.com/* artifacts/godoc
	cp -r artifacts/$(serverurl)/lib/godoc/* artifacts/godoc/lib/godoc/
	rm -R artifacts/$(serverurl)
	pkill godoc
	tar -zcvf artifacts/godoc.tar.gz artifacts/godoc
	rm -R artifacts/godoc

.PHONY: fmt
fmt:
	goimports -w -local $(golocalimport) main sys service

.PHONY: test
test: generate
	go test $(COVERAGE_OPTS)  ./main/... ./sys/... ./storage/... ./service/...

.PHONY: test-integration
test-integration:
	$(eval testdir := $(shell mktemp -d /tmp/proxeus-test-integration.XXXXX ))
	mkdir -p $(testdir)
	$(eval cid=$(shell $(call startmongo,$(testdir))))
	go test $(COVERAGE_OPTS) -count=1 -tags integration ./storage/database/db/...; ret=$$?; \
		$(if $(cid), docker rm -f $(cid);) \
		rm -fr $(testdir); \
		exit $$ret

.PHONY: test-api
test-api: server
	$(eval testdir := $(shell mktemp -d /tmp/proxeus-test-api.XXXXX ))
	mkdir -p $(testdir)
	$(eval cid=$(shell [[ "x$(PROXEUS_DATABASE_ENGINE)" == "xmongo" ]] && $(call startmongo,$(testdir))))
	$(call startds,$(testdir))
	$(call startnodes,$(testdir))
	$(call startproxeus,$(testdir))
	PROXEUS_URL=http://localhost:1323  go test -count=1 ./test; ret=$$?; \
		$(stopproxeus); \
		[ -e  $(testdir)/ds-started ] && docker-compose down; \
		$(if $(cid), docker rm -f $(cid);) \
		rm -fr $(testdir); \
		exit $$ret

.PHONY: test-ui
test-ui: server ui
	$(eval testdir := $(shell mktemp -d /tmp/proxeus-test-ui.XXXXX ))
	mkdir -p $(testdir)
	$(eval cid=$(shell [[ "x$(PROXEUS_DATABASE_ENGINE)" == "xmongo" ]] && $(call startmongo,$(testdir))))
	$(call startds,$(testdir))
	$(call startnodes,$(testdir))
	$(call startproxeus,$(testdir))
	$(MAKE) -C test/e2e test; ret=$$?; \
		$(stopproxeus); \
		[ -e  $(testdir)/ds-started ] && docker-compose down; \
		$(if $(cid), docker rm -f $(cid);) \
		rm -fr $(testdir); \
		exit $$ret


.PHONY: coverage
coverage:
	gocovmerge artifacts/*.coverage > artifacts/coverage
	go tool cover -func artifacts/coverage > artifacts/coverage.txt
	go tool cover -html artifacts/coverage -o artifacts/coverage.html

.PHONY: clean
clean:
	cd artifacts && rm -rf `ls . | grep -v 'cache'`
	cd ui && yarn cache clean && cd ..

.PHONY: run
run: server
	artifacts/proxeus -DataDir $(PROXEUS_DATA_DIR)/proxeus-platform/data

main/handlers/assets/bindata.go: $(wildcard ./ui/core/dist/**)
	go-bindata ${BINDATA_OPTS} -pkg assets -o ./main/handlers/assets/bindata.go -prefix ./ui/core/dist ./ui/core/dist/...
	goimports -w $@

test/assets/bindata.go: $(filter-out bindata.go,$(shell find ./test/assets/))
	go-bindata ${BINDATA_OPTS} -pkg assets -o ./test/assets/bindata.go ./test/assets/...
	goimports -w $@

.SECONDEXPANSION: # See https://www.gnu.org/software/make/manual/make.html#Secondary-Expansion
$(mocks): $$(patsubst %_mock.go, %.go, $$(subst /mock,, $$@))
	mockgen -package mock  -source $<  -destination $@  -self_package github.com/ProxeusApp/proxeus-core/$(shell dirname $@)
	goimports -w $@
