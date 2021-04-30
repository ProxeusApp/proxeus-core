![](docs/_media/logo.png)

# Proxeus
----------------
Main repository of the Proxeus platform core.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.

You can access the source code of this application and various extensions
on [GitHub](https://github.com/ProxeusApp).

## Quickstart

We are making facilitated "one-click" automated build configurations available for select cloud platforms. Please contact us if you are interested in additional providers on this list:

- [Linode StackScript](deploy/linode/README.md)

## Installation using Docker

The quickest way to set up Proxeus for development is to use Docker, and the `docker-compose` tool.

1. [Install Docker Engine](https://docs.docker.com/install/)
2. [Install docker-compose](https://docs.docker.com/compose/install/)

See further deployment instructions in [docs/docker](docs/docker.md) to set up your server using Docker.

### Infura and SparkPost

The Proxeus platform currently depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/)
for Ethereum and email integration respectively. Create an account on those platforms
to get API Keys. These keys need to be added to corresponding environment variables, or
entered when deploying a "one-click" instance.

Please note that the domain you set up on SparkPost MUST match the **reply-to** e-mail address that you configure in the next step in order to create accounts and receive e-mails on your instance.

If all goes well, Proxeus should be available at http://localhost:1323. The next step will be to [configure](docs/configure.md) your instance for the first time.

## User manual

Get help to make the most of the platform in the [User Handbook](https://docs.google.com/document/d/e/2PACX-1vTchv7PotoQeH2cBA2VIHcqV0I0N_IQpFnbESR-8C19cgBikek3HAMVdPtfJJcYkANzPWbfy_S3bf8X/pub).

### XES Payments
For detailed information about payments setup, check the [XES-Payment Readme](docs/xes-payment.md).

### Smart contracts
Check the [Smart contracts & Solidity Readme](https://github.com/ProxeusApp/proxeus-contract) to learn more.

## Development guide

Please read the [Developer Manual](https://doc.proxeus.com) to learn more about the Proxeus platform.

If you wish to build the project form the source code, follow the instructions in [Build all](docs/build_all.md)

## Contributing

As an open-source project, we welcome any kind of community involvement, whether that is by contributing code, reporting issues or
engaging in insightful discussions. Especially, we are looking forward to receiving contributions for external workflow nodes.

See the [Contributing](docs/contributing.md) section for instructions on how to contribute.

## Security Issues

If you find a vulnerability that may affect live or testnet deployments please DO NOT file a public issue - send your report privately to info@proxeus.org

## License

Licensed under the GNU GENERAL PUBLIC LICENSE. You may read a copy of the [License here](LICENSE).

## Acknowledgements

Like so many projects, this effort has roots in many places. The list can be found in [ACKNOWLEDGEMENTS](ACKNOWLEDGEMENTS).
