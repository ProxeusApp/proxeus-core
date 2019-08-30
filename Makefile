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
	go test ./main/handlers/api ./main/handlers/workflow ./main/handlers/blockchain

clean:
	cd artifacts && rm -rf `ls . | grep -v 'cache'`

all: ui server

# used by CI
link-repo:
	mkdir -p /go/src/git.proxeus.com/core
	ln -s /builds/core/central /go/src/git.proxeus.com/core/central

.PHONY: init main all all-debug generate test clean fmt validate link-repo
.PHONY: ui server
