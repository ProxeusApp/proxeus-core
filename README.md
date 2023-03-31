![](docs/_media/logo.png)

# Proxeus Core
----------------

Proxeus combines a powerful automation tool with a form builder, document generator and blockchain connection - enabling you to digitize, secure, and tap into the value of data flows. Access the full source code of available modules and extensions [@ProxeusApp](https://github.com/ProxeusApp).

![Screenshot of Proxeus workflow from the handbook](https://github.com/ProxeusApp/community/raw/master/handbook/Proxeus%20-%20The%20Complete%20Handbook_html_10299e76126cc024.png)

## User Guides

Get help to make the most of the platform in the **[User Handbook](https://github.com/ProxeusApp/community/blob/master/handbook/handbook.md)**.

To learn more about Smart Contracts using Solidity, see the documentation in the [proxeus-contract](https://github.com/ProxeusApp/proxeus-contract).

For detailed information about payment setup (currently unsupported), check the [XES-Payment](docs/xes-payment.md) project.

## Installation

Proxeus is primarily a Web application, intended for access with a web browser. The [Proxeus Association](https://proxeus.org) maintains a demo instances you can use to test the product, and can recommend a service provider to help you or your business get set up. There is also a prototype [desktop application](https://github.com/ProxeusApp/storage-app/blob/master/docs/overview.md).

In addition to the developer guidelines below, several "one-click" deployment configurations are available for select cloud platforms:

- [Docker Compose](docker-compose.yml)
- [Linode StackScript](deploy/linode/README.md)
- [DigitalOcean Droplet](deploy/digitalocean/README.md)

Please [contact us](https://github.com/ProxeusApp/community/discussions/3) if you are interested in seeing additional providers on this list.

## Development

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

## Contributing

As an open-source project, we welcome any kind of community involvement, whether that is by contributing code, reporting issues or
engaging in insightful discussions. Especially, we are looking forward to receiving contributions for external workflow nodes.

See the **[Contributing Guide](docs/contributing.md)** for further directions.

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
- polygon
