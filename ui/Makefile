
.PHONY: check-node-version init serve-main-hosted all main-hosted

init:
	yarn install --cache-folder .yarn

check-node-version:
	yarn run check-node-version

serve-main-hosted: check-node-version
	cd core && yarn run serve

all: main-hosted

main-hosted: init check-node-version
	cd core && yarn run build
