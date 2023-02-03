![](docs/_media/logo.png)

# Proxeus Core
----------------
Main repository of the Proxeus core platform.

Proxeus combines a powerful document automation tool with the wide-ranging
blockchain functionalities, enabling users to digitize and monetize their IP.
You can access the source code of this application and various extensions
on [GitHub](https://github.com/ProxeusApp).

## User guides

Get help to make the most of the platform in the **[User Handbook](https://github.com/ProxeusApp/community/blob/master/handbook/handbook.md)**.

For detailed information about payment setup, check the [XES-Payment Readme](docs/xes-payment.md).

To learn more about Smart Contracts, see the documentation in the [Solidity Readme](https://github.com/ProxeusApp/proxeus-contract).

### Quickstart

You can get all the information to run Proxeus here. In addition to the Docker installation documented below, we are making facilitated "one-click" automated build configurations available for select cloud platforms.

- [Linode StackScript](deploy/linode/README.md)

Please [contact us](https://github.com/ProxeusApp/community/discussions/3) if you are interested in additional providers on this list.

### Contributing

As an open-source project, we welcome any kind of community involvement, whether that is by contributing code, reporting issues or
engaging in insightful discussions. Especially, we are looking forward to receiving contributions for external workflow nodes.

See the **[Contributing Guide](docs/contributing.md)** for further directions.

## Development guide

Please read the [Developer Manual](https://doc.proxeus.com) to learn more about the Proxeus platform.

If you wish to build the project form the source code, follow the instructions in [Build all](docs/build_all.md)

### Installation using Docker

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

### Security Issues

If you find a vulnerability that may affect live or testnet deployments please DO NOT file a public issue - send your report privately to info@proxeus.org

## License

Licensed under the GNU GENERAL PUBLIC LICENSE. You may read a copy of the [License here](LICENSE).

## Acknowledgements

Like so many projects, this effort has roots in many places. The list can be found in [ACKNOWLEDGEMENTS](ACKNOWLEDGEMENTS).

## Supported Chains

- goerli
- mainnet
- polygon-mumbai
