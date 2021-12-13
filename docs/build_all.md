# Building Backend and Frontend

This section describes how to build and run the platform for developers.

### Dependencies

Here is the list of dependencies:

+ make
+ go (1.10+, 64bit for Windows)
+ GOBIN added to your PATH (to check your GOBIN: `echo $(go env GOPATH)/bin`)
+ curl
+ yarn (1.12+)
+ node (14)
+ vue-cli
+ git
+ docker-compose (18.06.0+)
+ wget (for godoc generation)
+ (OSX) Brew Package Manager
+ (OSX) Command Line Tools (Xcode)


#### Linux (Debian)

```
sudo apt-get install make golang curl npm git
```

#### OSX
If you currently do not have the OSX `Command Line Tools` installed run
```
xcode-select --install
```

```
brew install make golang curl npm git
```

### Export PATH variables
```
PATH=$PATH:$(go env GOPATH)/bin
```

### Clone repository
The project uses go modules.

Clone the repository outside your GOPATH:
```
cd <your workspace>
git clone https://github.com/ProxeusApp/proxeus-core.git
cd proxeus-core
```

### Build
All the build projects are stated in `./Makefile`.

Before building Proxeus, make sure to set all [required environment variables](../README.md#quick-start-with-docker).
```
export PROXEUS_DATA_DIR=./data
export PROXEUS_INFURA_API_KEY=<Your Infura project ID>
export PROXEUS_SPARKPOST_API_KEY=<Your SparkPost API Key>
export PROXEUS_ENCRYPTION_SECRET_KEY=<A long random string>
export PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
export PROXEUS_ALLOW_HTTP=true
export PROXEUS_PLATFORM_DOMAIN=localhost
export PROXEUS_EMAIL_FROM=youremail@example.com

```
Make sure that the email domain for `PROXEUS_EMAIL_FROM` is a configured sending domain in sparkpost

To initialize dependencies run:
```
make init
```

To update go dependencies run:
```
make update
```

The make all command build the `server` and `ui`:
```
make all
```

### Start

Run `server`.
```
./artifacts/proxeus
```

The platform will be available at the following URL: http://localhost:1323

It is now time to [configure your platform](configure.md)

### Tips

Update command failing? Try:

- doing a quick verification and config-file check
`go mod verify` / `go mod tidy`
- clearing your module cache:
`go clean -cache -modcache -i -r`
- checking your dependency graph to isolate the issue:
`go mod graph`
