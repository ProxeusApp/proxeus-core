dependencies=npm yarn

.PHONY: all
all: main-hosted

.PHONY: init
init:
	@for d in $(dependencies); do (echo "Checking $$d is installed... " && which $$d ) || ( echo "Please install $$d before continuing" && exit 1 ); done
	yarn install --cache-folder .yarn

.PHONY: check-node-version
check-node-version:
	yarn run check-node-version

.PHONY: serve-main-hosted
serve-main-hosted: check-node-version
	cd core && yarn run serve

.PHONY: dev
dev: init check-node-version
	cd core && yarn run dev

.PHONY: main-hosted
main-hosted: init check-node-version
	cd core && yarn run build

.PHONY: clean
clean:
	yarn cache clean --cache-folder .yarn
