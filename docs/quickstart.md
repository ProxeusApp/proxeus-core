# Quick Start

The quickest way to try Proxeus is to use `docker-compose`.

## Install docker and docker-compose
1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

## Get API Keys for Infura and SparkPost
The Proxeus platform depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/) 
for Ethereum and email integration respectively.

Please create an account on those platform and get an API Keys.

## Proxeus Demo Ethereum Smart Contract

For your convenience, a demo smart contract is deployed on the Ropsten network at the following address:

```
0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
```

## Create a docker-compose.yml file 

**Note: Please make sure that you always pull Docker images from the official `proxeus` DockerHub repository and that you are using the latest version.**

User the example below as your `docker-compose.yml` file:

```
version: '3.7'

networks:
  xes-platform-network:
    name: xes-platform-network

services:
  platform:
    image: proxeus/proxeus-core:latest
    container_name: xes_platform
    depends_on:
      - document-service
    networks:
      - xes-platform-network
    restart: unless-stopped
    environment:
      TZ: Europe/Zurich
      DataDir: "/data/hosted"
      DocumentServiceUrl: "http://document-service:2115/"
      InfuraApiKey: "${PROXEUS_INFURA_KEY}"
      SparkpostApiKey: "${PROXEUS_SPARKPOST_KEY}"
      BlockchainContractAddress: "${PROXEUS_CONTRACT_ADDRESS}"
      EmailFrom: "${PROXEUS_EMAIL_FROM:-no-reply@proxeus.com}"
      AirdropWalletfile: "${PROXEUS_AIRDROP_WALLET_FILE:-./data/proxeus-platform/settings/airdropwallet.json}"
      AirdropWalletkey: "${PROXEUS_AIRDROP_WALLET_KEY:-./data/proxeus-platform/settings/airdropwallet.key}"
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

## Start Proxeus

Run the following command in the directory containing your `docker-compose.yml` file (Linux and OSX):
```
export PROXEUS_INFURA_KEY=<Your Infura API key>
export PROXEUS_SPARKPOST_KEY=<Your SparkPost API Key>
export PROXEUS_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
docker-compose up 
```

Proxeus should be available at http://localhost:1323

The next step is to [configure](configure.md) your instance for the first time.
