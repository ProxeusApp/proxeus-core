# Docker Deployment

Proxeus docker deployment uses `docker-compose` and requires the following
dependencies:

1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)


You easily deploy Proxeus using Docker.  The repository includes several `docker-compose` YAML files that can be
used to deploy the platform in different context:

* `docker-compose.yml`
* `docker-compose.override.yml`
* `docker-compose-cloud-override.yml`

## Simple Docker Compose

This is the simplest method to experiment with Proxeus.  This will start a local Proxeus platform
using images from Docker Hub.  

Please note that with SELinux enabled: a `:z` should be added to the end of volume declarations in docker-compose.yml.

### docker-compose.yml file

```
---
version: '3.7'

networks:
  xes-platform-network:
    name: xes-platform-network

services:
  platform:
    image: proxeus/proxeus-core:latest
    container_name: xes-platform
    depends_on:
      - document-service
    networks:
      - xes-platform-network
    restart: unless-stopped
    environment:
      TZ: Europe/Zurich
      PROXEUS_PLATFORM_DOMAIN: "${PROXEUS_PLATFORM_DOMAIN:-http://xes-platform:1323}"
      PROXEUS_DOCUMENT_SERVICE_URL: "http://document-service:2115/"
      PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS: "${PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS}"
      PROXEUS_INFURA_API_KEY: "${PROXEUS_INFURA_API_KEY}"
      PROXEUS_SPARKPOST_API_KEY: "${PROXEUS_SPARKPOST_API_KEY}"
      PROXEUS_ENCRYPTION_SECRET_KEY: "${PROXEUS_ENCRYPTION_SECRET_KEY}"
      PROXEUS_EMAIL_FROM: "${PROXEUS_EMAIL_FROM:-no-reply@example.com}"
      PROXEUS_AIRDROP_WALLET_FILE: "${PROXEUS_AIRDROP_WALLET_FILE:-/root/.proxeus/settings/airdropwallet.json}"
      PROXEUS_AIRDROP_WALLET_KEY: "${PROXEUS_AIRDROP_WALLET_KEY:-/root/.proxeus/settings/airdropwallet.key}"
      PROXEUS_DATABASE_ENGINE: "${PROXEUS_DATABASE_ENGINE:-storm}"
      PROXEUS_DATABASE_URI: "${PROXEUS_DATABASE_URI:-mongodb://root:root@mongo:27017}"
      PROXEUS_TEST_MODE: "${PROXEUS_TEST_MODE:-false}"
      PROXEUS_ALLOW_HTTP: "${PROXEUS_ALLOW_HTTP:-false}"
    ports:
      - "1323:1323"
    volumes:
      - ${PROXEUS_DATA_DIR:-./data}/proxeus-platform/data:/data/hosted
      - ${PROXEUS_DATA_DIR:-./data}/proxeus-platform/settings:/root/.proxeus/settings

  document-service:
    image: proxeus/document-service:latest
    container_name: xes_document_service
    networks:
      - xes-platform-network
    restart: unless-stopped
    environment:
      TZ: Europe/Zurich
    ports:
      - "2115:2115"
    volumes:
      - ${PROXEUS_DATA_DIR:-./data}/document-service/logs:/document-service/logs
      - ${PROXEUS_DATA_DIR:-./data}/document-service/fonts:/document-service/fonts
```

### Start

```
docker-compose -f docker-compose.yml up
```

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string*|Use a hard key to ensure your database is safe.|
|PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_API_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_API_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|
|PROXEUS_PLATFORM_DOMAIN:|`http://xes-platform:1323`|The domain of the running platform.  Mainly used for display|
|PROXEUS_ALLOW_HTTP:|`false`|Allow the use of HTTP instead of HTTPS =NOT FOR PRODUCTION=|

## Development Docker Compose

This file will start the document service available from Docker Hub but will start
the local Platform built from your local files.  This method is preferred during development.
This is the default configuration when using `docker-compose up` from the project root directory.
See [Share Compose configurations between files and projects](https://docs.docker.com/compose/extends/)

### docker-compose.override.yml

```
---
version: '3.7'

services:
  platform:
    build:
      context: .
    environment:
      PROXEUS_ALLOW_HTTP: "${PROXEUS_ALLOW_HTTP:-true}"
```


### Start
```
make server-docker
docker-compose build
docker-compose up
```

Please not than in this case, you do not need to specify the docker-compose YAML files as reading the
`docker-compose.yml` and `docker-compose.override.yml` is the default behaviour.

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string*|Use a hard key to ensure your database is safe.|
|PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_API_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_API_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|
|PROXEUS_PLATFORM_DOMAIN:|`http://xes-platform:1323`|The domain of the running platform.  Mainly used for display|
|PROXEUS_ALLOW_HTTP:|`false`|Allow the use of HTTP instead of HTTPS =NOT FOR PRODUCTION=|


## Cloud Docker Compose

This is a docker compose override file, i.e. it must be used in conjunction
with the `docker-compose.yml` file as described in [Multiple Compose files](https://docs.docker.com/compose/extends/#multiple-compose-files).
It will add the required configuration to deploy Proxeus on a hosted VM for example on Google Cloud or AWS,
including
* [nginx](https://hub.docker.com/r/jwilder/nginx-proxy/) reverse proxy,
* [letsencrypt ](https://hub.docker.com/r/jrcs/letsencrypt-nginx-proxy-companion/) HTTPS provider and
* [watchtower](https://hub.docker.com/r/v2tec/watchtower/) automatic container update.

Please refer to [Use Compose in production](https://docs.docker.com/compose/production/) for more information about
running docker compose in production.

This is the method that we use to deploy the Proxeus Demo site.

### docker-compose-cloud.override.yml

```
---
# This file is an override and needs to be used in combination with
# docker-compose.yml using one of the following methods:
#
# 1. using the -f parameter:
#  docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml up
#
# 2. renaming the file docker-compose.override.yml (default override):
#  docker-compose up

version: '3.7'

networks:
  # Add Network for reverse-proxy
  reverse-proxy:
    name: reverse-proxy
    driver: bridge

volumes:
  # Add volume for nginx-proxy and letsencrypt
  nginx-share:

services:
  # Add Nginx reverse-proxy
  # https://hub.docker.com/r/jwilder/nginx-proxy/
  # Automated Nginx reverse proxy for docker containers
  nginx-proxy:
    container_name: nginx-proxy
    image: jwilder/nginx-proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - nginx-share:/etc/nginx/vhost.d
      - nginx-share:/usr/share/nginx/html
      - ${PROXEUS_DATA_DIR:-./data}/certs:/etc/nginx/certs:ro
      - ${DOCKER_SOCK:-/var/run/docker.sock}:/tmp/docker.sock:ro
    networks:
      - reverse-proxy
    labels:
      com.github.jrcs.letsencrypt_nginx_proxy_companion.nginx_proxy: "true"
    restart: unless-stopped

  # Add Letsencrypt
  # https://hub.docker.com/r/jrcs/letsencrypt-nginx-proxy-companion/
  # LetsEncrypt container to use with nginx as proxy
  letsencrypt:
    image: jrcs/letsencrypt-nginx-proxy-companion
    depends_on:
      - nginx-proxy
    networks:
      - reverse-proxy
    volumes:
      - nginx-share:/etc/nginx/vhost.d
      - nginx-share:/usr/share/nginx/html
      - ${PROXEUS_DATA_DIR:-./data}/certs:/etc/nginx/certs:rw
      - ${DOCKER_SOCK:-/var/run/docker.sock}:/var/run/docker.sock:ro
    restart: unless-stopped

  platform:
    networks:
      - reverse-proxy
    labels:
      com.centurylinklabs.watchtower.enable: "true"
    environment:
      # Replace values for reverse-proxy
      VIRTUAL_HOST: ${PROXEUS_VIRTUAL_HOST:-proxeus.example.com}
      VIRTUAL_PORT: 1323
      # Replace values for letsencrypt
      LETSENCRYPT_HOST: ${PROXEUS_VIRTUAL_HOST:-proxeus.example.com}
      LETSENCRYPT_EMAIL: ${LETSENCRYPT_EMAIL:-admin@example.com}

  document-service:
    networks:
      - reverse-proxy
    labels:
      com.centurylinklabs.watchtower.enable: "true"
    environment:
      # Replace values for reverse-proxy
      VIRTUAL_HOST: ${DOCUMENT_SERVICE_VIRTUAL_HOST:-proxeus.example.com}
      VIRTUAL_PORT: 2115
      # Replace values for letsencrypt
      LETSENCRYPT_HOST: ${DOCUMENT_SERVICE_VIRTUAL_HOST:-proxeus.example.com}
      LETSENCRYPT_EMAIL: ${LETSENCRYPT_EMAIL:-admin@example.com}

  # Add Watchtower
  # https://hub.docker.com/r/v2tec/watchtower/
  # Watches your containers and automatically restarts them whenever their
  # image is refreshed
  watchtower:
    image: v2tec/watchtower
    container_name: watchtower
    restart: always
    volumes:
      - ${DOCKER_SOCK:-/var/run/docker.sock}:/var/run/docker.sock:ro
      - ${DOCKER_CONFIG_FILE:-/root/.docker/config.json}:/config.json
    command: --interval 60 --label-enable
```

### Start

```
docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml -d up
```

To simplify your deployment, you can rename `docker-compose-cloud.override.yml` to
`docker-compose.override.yml` and avoid having to specify the file names in the command.

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string*|Use a hard key to ensure your database is safe.|
|PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_API_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_API_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|
|PROXEUS_PLATFORM_DOMAIN:|`http://xes-platform:1323`|The domain of the running platform.  Mainly used for display|
|PROXEUS_ALLOW_HTTP:|`false`|Allow the use of HTTP instead of HTTPS =NOT FOR PRODUCTION=|
|DOCKER_SOCK|`/var/run/docker.sock`| The path to the Docker socket file. Used to allow `Letsencrypt` and `watchtower` to find other containers. |
|DOCKER_CONFIG|`/root/.docker/config.json`|The path to the docker config file.  Used to access the image repository authentication parameters |


## Custom Docker Deployment

The first method to adapt Proxeus to your infrastructure need is to define the environment variables corresponding to your situation.
The next level will be to customize a `docker-compose.yml` file.

For more information see the [Configuration docs](./configure.md).
