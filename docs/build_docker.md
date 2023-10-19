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
BUILD_WITH_DOCKER=true make init ui server-docker build-docker
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

## Tips

Having an issue with Go? Make sure it's in your path, e.g.:

`export PATH=$PATH:/usr/local/go/bin`

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

You may also want to include custom nodes. There is a sample configuration which can be started like this:

```
docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml -f docker-compose-cnode.override.yml up
```

See `docker-compose-extra.override.yml` for examples with several other nodes.
