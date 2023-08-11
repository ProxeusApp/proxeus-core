# Docker

Please make sure that you always pull Docker images from the official `proxeus` DockerHub repository and that you are using the latest version.

## Start Proxeus

Run the following command in the directory containing your `docker-compose.yml` file (Linux and OSX):
```
export PROXEUS_EMAIL_FROM=<Your valid Sender Email Address>
export PROXEUS_INFURA_API_KEY=<Your Infura project ID>
export PROXEUS_SPARKPOST_API_KEY=<Your SparkPost API Key>
export PROXEUS_ENCRYPTION_SECRET_KEY=<A 32-character random string>
export PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
export PROXEUS_ALLOW_HTTP=true
docker-compose up
```

You can also put these settings into an `.env` file in the same folder as the application.

See [Configuration settings](configure.md) for more details.

## Tweak your Docker setup

Besides the basic Docker Compose configuration, you can extend your installation as follows:

- `docker-compose-cloud.override` for cloud installations, which includes Nginx and Let's Encrypt
- `docker-compose-example.override` shows how to add another Proxeus Node to your installation
- `docker-compose-extra.override` includes all officially supported Proxeus Nodes
- `docker-compose-local.override` if you want to use your local Docker image (details below)

To use one or more of these overrides, start Proxeus as follows:

`docker-compose -f docker-compose.yml -f docker-compose-example.override.yml up`

(you always have to first include `docker-compose.yml`)

## Build a Docker image

The project includes a `Dockerfile` that is used to create the `proxeus-core` docker image.

First you need to have the [Docker Engine](https://docs.docker.com/install/) installed.

Then you can build a docker image:

```
make server-docker
docker build .
```

Please refer to the `docker-compose.yml` file to learn how to configure a Proxeus container.

## Using Docker during development

If you want to use a Docker image to build your server (for example, when you are not on a Linux machine or have a newer version of GLIBC than is supported in production), compile as follows:

```
BUILD_WITH_DOCKER=true make init server-docker
docker build .
```

You can of course use Docker to run the platform during development:

```
docker-compose build
docker-compose up
```

and

```
docker-compose build
docker-compose restart
```

This will build the proxeus-core image based on your current project and use a deployed image
for the document service.

## Using Docker for the build

If you're having trouble, try a clean full Docker build, specifying each of the configuration files:

```
make clean
BUILD_WITH_DOCKER=true make init server-docker
docker-compose --env-file .env -f docker-compose.yml -f docker-compose-local.yml build
```

## Using Docker for production

For deployment, a `docker-compose-cloud.override.yml` file is provided which includes Nginx, Let's Encrypt and other services used in larger deployments:

```
docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml up
```

## Docker Light version

There is also a Docker Compose configure in one file with a 'minimal' Proxeus installation. The only extra nodes are 'mail-sender' and 'json-sender':

```
docker-compose -f docker-compose-light.yml up
```
