![](docs/_media/logo.png)

# Proxeus
----------------
Main repository of the proxeus platform.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.

## Quick Start

The quickest way to try Proxeus is to use `docker-compose`.

### Install docker and docker-compose
1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

### Get API Keys for Infura and SparkPost
The Proxeus platform depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/) 
for Ethereum and email integration respectively.

Please create an account on those platform and get an API Keys.

### Proxeus Demo Etherum Smart Contract

For your convenience, a demo smart contract is deployed on the Ropsten network at the following address:

```
0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
```

### Start Proxeus

Run the following command (Linux and OSX):
```
export PROXEUS_INFURA_KEY=<Your Infura API key>
export PROXEUS_SPARKPOST_KEY=<Your SparkPost API Key>
export PROXEUS_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
docker-compose up 
```

Proxeus should be available at http://localhost:1323

The next step is to [configure](configure.md) your instance for the first time.

## User manual

The user manual is available here: [User Manual](https://docs.google.com/document/d/1SP0ZimG7uemfZ2cF2JkY5enUZnBJLDyfcJGZnyWOejQ)

## Developer manual

Please read the [Developer Manual](docs/_sidebar.md) to learn more about the 
Proxeus platform. *TODO: link to the github pages documentation site when ready*


## 3 Misc

### 3.1 XES-Payment
For more info check the [XES-Payment Readme](docs/xes-payment.md).

### 3.2 Smart contracts & Solidity
For more info check the [Smart contracts & Solidity Readme](https://github.com/ProxeusApp/proxeus-contract).
