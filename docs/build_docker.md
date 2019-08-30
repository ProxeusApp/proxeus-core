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

You can use Docker during development by using the `docker-compose-dev.yml` file:

```
docker-compose -f docker-compose-dev.yml build 
docker-compose -f docker-compose-dev.yml up
```

and

```
docker-compose -f docker-compose-dev.yml build
docker-compose -f docker-compose-dev.yml restart
```

This will build the proxeus-core image based on your current project and use a deployed image 
for the document service.


