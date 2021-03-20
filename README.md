![](docs/_media/logo.png)

# Proxeus
----------------
Main repository of the Proxeus platform core.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.

You can access the source code of this application and various extensions
on [GitHub](https://github.com/ProxeusApp).

# Quick Starts

## Automated builds

We are making facilitated "one click" build configurations available for select cloud platforms:

- [Linode StackScript](deploy/linode/README.md)

## Using Docker

The quickest way to set up Proxeus for development is to use Docker, and the `docker-compose` tool.

### Install docker and docker-compose

1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

### Get API Keys for Infura and SparkPost

The Proxeus platform depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/)
for Ethereum and email integration respectively. Create an account on those platforms
to get API Keys. These keys need to be added to corresponding environment variables.

See further deployment instructions in [docs/docker](docs/docker.md) to set up your server using Docker.

If all goes well, Proxeus should be available at http://localhost:1323

The next step is to [configure](docs/configure.md) your instance for the first time.

## Build Proxeus Platform from the source code

If you are a developer and want to build the project form the source code follow the instructions in [Build all](docs/build_all.md)

## Developer manual

Please read the [Developer Manual](https://doc.proxeus.com) to learn more about the
Proxeus platform.

## User manual

Get help to make the most of the platform in the [User Handbook](https://docs.google.com/document/d/e/2PACX-1vTchv7PotoQeH2cBA2VIHcqV0I0N_IQpFnbESR-8C19cgBikek3HAMVdPtfJJcYkANzPWbfy_S3bf8X/pub).

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
