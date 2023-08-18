# Docker Deployment

Proxeus docker deployment uses `docker-compose` and requires the following
dependencies:

1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

## Start with Docker Compose

You easily deploy Proxeus using Docker.  The repository includes several `docker-compose` YAML files that can be used to deploy the platform in different context:

* `docker-compose.yml`
* `docker-compose-extra.override.yml`
* `docker-compose-cloud.override.yml`
* `docker-compose-local.override.yml`

In each case you can then use the **logs** command to see the system status (with the useful `-f` parameter):

`docker-compose logs`

For more usage instructions, visit the [Docker Compose CLI reference](https://docs.docker.com/compose/reference/).


In each case you can then use the **logs** command to see the system status (with the useful `-f` parameter):

`docker-compose logs`

For more usage instructions, visit the [Docker Compose CLI reference](https://docs.docker.com/compose/reference/).


## Simple Docker Compose

`docker-compose -f docker-compose.yml -f docker-compose-extra.override.yml up`

In each case you can then use the **logs** command to see the system status (with the useful `-f` parameter):

`docker-compose logs`

For more usage instructions, visit the [Docker Compose CLI reference](https://docs.docker.com/compose/reference/).
  
Please note that with **SELinux enabled**: a `:z` should be added to the end of volume declarations in docker-compose.yml.

### Environment variables

Add these to a `.env` file in the same folder as `docker-compose.yml`:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string of 32 characters*|Use a hard key to ensure your database is safe.|
|PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS|*0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71*|The address of the Proxeus contract.|
|PROXEUS_INFURA_API_KEY|*Your Infura API Key*|An Infura API Key for Ethereum integration.|
|PROXEUS_SPARKPOST_API_KEY|*Your SpartPost Key*|A SparkPost API Key for email integration.|
|PROXEUS_EMAIL_FROM|`no-reply@example.com`|The email address used as sender when Proxeus sends an email.|
|PROXEUS_PLATFORM_DOMAIN:|`http://xes-platform:1323`|The domain of the running platform.  Mainly used for display|
|PROXEUS_ALLOW_HTTP:|`false`|Allow the use of HTTP instead of HTTPS =NOT FOR PRODUCTION=|

## Development Docker Compose

This file will start the document service available from DockerHub but will start
the local Platform built from your local files.  This method is preferred during development.
This is the default configuration when using `docker-compose up` from the project root directory.
See [Share Compose configurations between files and projects](https://docs.docker.com/compose/extends/)

```
make server-docker
docker-compose build
docker-compose -f docker-compose.yml -f docker-compose-local.override.yml up --remove-orphans
```

Please not than in this case, you do not need to specify the docker-compose YAML files as reading the
`docker-compose.yml` and `docker-compose.override.yml` is the default behaviour.

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string of 32 characters*|Use a hard key to ensure your database is safe.|
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
* [letsencrypt](https://hub.docker.com/r/jrcs/letsencrypt-nginx-proxy-companion/) HTTPS provider and
* [watchtower](https://hub.docker.com/r/v2tec/watchtower/) automatic container update.

Please refer to [Use Compose in production](https://docs.docker.com/compose/production/) for more information about
running docker compose in production.

This is the method that we use to deploy the Proxeus Demo site. Start with:

```
docker-compose -f docker-compose.yml -f docker-compose-cloud.override.yml -d up
```

To simplify your deployment, you can rename `docker-compose-cloud.override.yml` to
`docker-compose.override.yml` and avoid having to specify the file names in the command.

Environment:

|Name           | Default Value | Description |
|---------------|-----------------------|------------------------------|
|PROXEUS_DATA_DIR| `./data` | Path to the directory to use a data store.|
|PROXEUS_ENCRYPTION_SECRET_KEY|*A random string of 32 characters*|Use a hard key to ensure your database is safe.|
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
