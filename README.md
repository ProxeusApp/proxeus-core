![](docs/_media/logo.png)

# Proxeus
----------------
Main repository of the proxeus platform.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.

## Quick Start with docker
The quickest way to try Proxeus is to use `docker-compose`.

### Install docker and docker-compose
1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

### Get API Keys for Infura and SparkPost
The Proxeus platform depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/) 
for Ethereum and email integration respectively.

Please create an account on those platform and get an API Keys.

## Proxeus Demo Ethereum Smart Contract

For your convenience, a demo smart contract is deployed on the Ropsten network at the following address:

```
0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
```

### Create a docker-compose.yml file 

**Note: Please make sure that you always pull Docker images from the official `proxeus` DockerHub repository and that you are using the latest version.**

User the example below as your `docker-compose.yml` file:

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
      PROXEUS_EMAIL_FROM: "${PROXEUS_EMAIL_FROM:-no-reply@example.com}"
      PROXEUS_AIRDROP_WALLET_FILE: "${PROXEUS_AIRDROP_WALLET_FILE:-/root/.proxeus/settings/airdropwallet.json}"
      PROXEUS_AIRDROP_WALLET_KEY: "${PROXEUS_AIRDROP_WALLET_KEY:-/root/.proxeus/settings/airdropwallet.key}"
      PROXEUS_DATABASE_ENGINE: "${PROXEUS_DATABASE_ENGINE:-storm}"
      PROXEUS_DATABASE_URI: "${PROXEUS_DATABASE_URI:-mongodb://root:root@mongo:27017}"
      PROXEUS_TEST_MODE: "${PROXEUS_TEST_MODE:-false}"
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

### Start Proxeus

Run the following command in the directory containing your `docker-compose.yml` file (Linux and OSX):
```
export PROXEUS_INFURA_KEY=<Your Infura API key>
export PROXEUS_SPARKPOST_KEY=<Your SparkPost API Key>
export PROXEUS_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
docker-compose up 
```

Proxeus should be available at http://localhost:1323

The next step is to [configure](docs/configure.md) your instance for the first time.

## Build Proxeus Platform from the source code

If you are a developer and want to build the project form the source code follow the instructions in [Build all](docs/build_all.md)

## Developer manual

Please read the [Developer Manual](https://doc.proxeus.com) to learn more about the 
Proxeus platform. 

## User manual

The user manual is available here: [User Manual](https://docs.google.com/document/d/1SP0ZimG7uemfZ2cF2JkY5enUZnBJLDyfcJGZnyWOejQ)

## Contributing

As an open-source project, we welcome any kind of community involvement, whether that is by contributing code, reporting issues or 
engaging in insightful discussions. Especially, we are looking forward to receiving contributions for custom workflow nodes.

See the [coding style](coding_style.md) section for instructions on the coding style we use.

Please use our support@proxeus.com email to open improvement and bug tickets.

### Security Issues

If you find a vulnerability that may affect live or testnet deployments please send your report privately to 
security@proxeus.com. Please DO NOT file a public issue.

## Misc

### XES-Payment
For more info check the [XES-Payment Readme](docs/xes-payment.md).

### Smart contracts & Solidity
For more info check the [Smart contracts & Solidity Readme](https://github.com/ProxeusApp/proxeus-contract).


