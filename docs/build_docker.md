# Docker


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
docker-compose --env-file .env -f docker-compose.yml -f docker-compose.override.yml build
```

## Using Docker for deployment

For deployment, a `docker-compose-cloud.override.yml` file is provide and must be used
instead of the default `docker-compose.override.yml`:

```
docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml
```

## Docker Light version

There is also a Docker Compose configure in one file with a 'minimal' Proxeus installation. The only extra nodes are 'mail-sender' and 'json-sender':

```
docker-compose -f docker-compose-light.yml up
```
