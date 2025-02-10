# Quick Start

This document is a brief installation guide, with links to further details to get started with Proxeus. We hope that you, too, will enjoy working with our product! Please remember to  star us on GitHub drop us any questions or feedback in the [Discussion Forum](https://github.com/ProxeusApp/community/discussions).


_"Proxeus is a platform for quick and convenient document digitalization, signing, processing, and distribution. It allows users to keep their important documents secure and registered on the blockchain. Proxeus empowers anyone to create blockchain applications and is available for free as an open-source project." --[S-PRO](https://s-pro.io/)_

➡️ [Demo server](https://proxeus-demo.s-pro-services.com/)

Kindly hosted by S-PRO, the demo server allows you to try a full-featured Proxeus instance. Please note that all content will be erased every 24 hours.

## One-Click Installation

Proxeus is primarily a Web application, intended for access with a web browser. The [Proxeus Association](https://proxeus.org) maintains a demo instances you can use to test the product, and can recommend a service provider to help you or your business get set up. There is also a prototype [desktop application](https://github.com/ProxeusApp/storage-app/blob/master/docs/overview.md).

In addition to the developer guidelines below, several "one-click" deployment configurations are available for select cloud platforms. Join the [Discussions](https://github.com/ProxeusApp/community/discussions/3) on GitHub if you are interested in seeing additional providers on this list:

- [Docker Compose](docs/docker.md)
- [Linode StackScript](deploy/linode/README.md)
- [DigitalOcean Droplet](deploy/digitalocean/README.md)

You will still need API keys for Infura and Sparkpost, as [detailed below](#get-keys).

Please read the [Developer Manual](https://doc.proxeus.com) to learn more about developing for the Proxeus platform.

**By installing and using the Proxeus software you agree with the terms of the [Proxeus License Agreement](https://github.com/ProxeusApp/proxeus-core/blob/main/LICENSE)** (GPL-3.0)

## Installing from Source

You can access the source code of this application on the [ProxeusApp](https://github.com/ProxeusApp) repository on GitHub, or the [Proxeus mirror](https://codeberg.org/proxeus/) on Codeberg.

If you wish to build the project from source, follow the instructions in [Build all](build_all.md).

The quickest way to try Proxeus is to use `docker compose`:

1. Install [Docker Engine](https://docs.docker.com/install/) or [Docker Desktop](https://docs.docker.com/desktop/)
2. Install [Docker Compose](https://docs.docker.com/compose/install/) if you are not using Docker Desktop
3. Get API keys for Infura and Sparkpost as [detailed below](#get-keys)
4. Configure your environment variables by copying `.env.example` to `.env` and modifying it
5. Start the application using `docker compose up`

Please make sure that you always pull Docker images from the official `proxeus` DockerHub repository and that you are using the latest version.

A `docker-compose.yml` file and some other possible configuraitons are in the root folder of the proxeus-core repository. See [Configuration settings](configure.md) and [Building with Docker](build_docker.md) for more details.

<a name="get-keys"></a>

## Get API Keys for Infura and SparkPost

The Proxeus platform depends on [Infura](https://infura.io/) and [SparkPost](https://www.sparkpost.com/). Please create an account on those platforms to get API Keys. These keys need to be added to corresponding environment variables, or entered when deploying a "one-click" instance.

We use [Infura](https://infura.io/) for blockchain services. Supported chains include:

- Ethereum: sepolia, goerli, mainnet
- Polygon: mumbai, mainnet

Create an account on [SparkPost](https://www.sparkpost.com/) for email integration.

- Please note that the domain you set up on SparkPost MUST match the **reply-to** e-mail address that you configure in the next step in order to create accounts and receive e-mails on your instance.

## Next steps

Proxeus should at this point be available at http://localhost:1323 or the IP address of your new cloud machine.

The next step is to [configure](configure.md) your instance for the first time.

To learn more about Smart Contracts using Solidity, see the documentation in the [proxeus-contract](https://github.com/ProxeusApp/proxeus-contract). For detailed information about token setup, check the [XES-Payment](xes-payment.md) project. Software architecture and other details are in the [Developer Manual](develop.md). For creating a custom node, visit [External workflow nodes](external_workflow_nodes.md).
