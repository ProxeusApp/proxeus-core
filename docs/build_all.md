# Building Backend and Frontend

This section describes how to build and run the platform for developers.

If you intend to deploy Proxeus to a production server, please follow the [Quick start](README) guide.

### Dependencies

Here is the list of dependencies:

+ make
+ go (1.10+, 64bit for Windows)
+ GOBIN added to your PATH (to check your GOBIN: `echo $(go env GOPATH)/bin`)
+ curl
+ git
+ yarn (1.12+)
+ node (16+)
+ docker-compose (18+)
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

Before building Proxeus, make sure to set all [required environment variables](build_docker.md) (i.e. with an `.env` file).

Make sure that the email domain for `PROXEUS_EMAIL_FROM` is a configured sending domain in Sparkpost.

To initialize dependencies run:
```
make init
```

The make all command build the `server` and `ui`:
```
make all
```

### Start

To run the server (`artifacts/proxeus`):

```
make run
```

The platform should in a few moments be available at the following URL: http://localhost:1323

It is now time to [configure your platform](configure.md).

### Tips

Update command failing? Try:

- doing a quick verification and config-file check
`go mod verify` / `go mod tidy`
- clearing your module cache:
`go clean -cache -modcache -i -r`
- checking your dependency graph to isolate the issue:
`go mod graph`
- updating Go & JS dependencies:
`make update`
