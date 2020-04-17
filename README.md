![](docs/_media/logo.png)

# Proxeus
----------------
Main repository of the proxeus platform.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.

## Source Code

You can access the source code of this application on the [Proxeus GitHub repository](https://github.com/ProxeusApp).

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

### Start Proxeus

Run the following command in the project root directory (Linux and OSX):
```
export PROXEUS_INFURA_API_KEY=<Your Infura API key>
export PROXEUS_SPARKPOST_API_KEY=<Your SparkPost API Key>
export PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
export PROXEUS_ALLOW_HTTP=true
docker-compose up 
```

Please note that with SELinux enabled: a `:z` should be added to the end of volume declarations in docker-compose.yml.

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
engaging in insightful discussions. Especially, we are looking forward to receiving contributions for external workflow nodes.

See the [Contributing](docs/contributing.md) section for instructions on how to contribute.

### Security Issues

If you find a vulnerability that may affect live or testnet deployments please send your report privately to 
security@proxeus.com. Please DO NOT file a public issue.

## Misc

### XES-Payment
For more info check the [XES-Payment Readme](docs/xes-payment.md).

### Smart contracts & Solidity
For more info check the [Smart contracts & Solidity Readme](https://github.com/ProxeusApp/proxeus-contract).

## License

Licensed under the GNU GENERAL PUBLIC LICENSE. You may read a copy of the License [here](LICENSE)

## Acknowledgements

Like so many projects, this effort has roots in many places. 

The list can be found [here](ACKNOWLEDGEMENTS)
