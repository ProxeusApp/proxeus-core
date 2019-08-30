# Docker Deployment

Proxeus docker deployment uses `docker-compose` and requires the following 
dependencies:

1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)


You easily deploy Proxeus using Docker.  The repository includes several `docker-compose` YAML files that can be 
used to deploy the platform in different context:

* `docker-compose.yml`
* `docker-compose-dev.yml`
* `docker-compose-cloud-override.yml`

## Simple Docker Compose
This is the simplest method to experiment with Proxeus.  This will start a local Proxeus platform 
using images from Docker Hub.  

```
docker-compose up 
```

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|

## Development Docker Compose
This file will start the document service available from Docker Hub but will start
the local Platform built from your local files.  This method is preferred during development.

```
make server-docker
docker-compose -f docker-compose-dev.yml build
docker-compose -f docker-compose-dev.yml up
```

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|


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


```
docker-compose -f docker-compose.yml -f docker-compose-cloud-override.yml -d up
```

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|
|DOCKER_SOCK|`/var/run/docker.sock`| The path to the Docker socket file. Used to allow `Letsencrypt` and `watchtower` to find other containers. |
|DOCKER_CONFIG|`/root/.docker/config.json`|The path to the docker config file.  Used to access the image repository authentication parameters |



## Custom Docker Deployment

The first method to adapt Proxeus to your infrastructure need is to define the environment variables corresponding to your situation.
The next level will be to customize a `docker-compose.yml` file. 


