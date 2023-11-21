Proxeus Script for fast DigitalOcean deployment
---

Creates a compact all-in-one instance of the Proxeus application (no code environment for smart contracts) using a bootstrapped release image for Docker. This is a good starting point for development or small installations. For more information visit https://github.com/ProxeusApp

Here you can find a script which run within a DigitalOcean droplet during startup. Using a simple form, you can configure the basic details needed to quickly get a Proxeus instance up and running.

This script is maintained for the community by Proxeus Association

## Instructions

1. Create fresh DO droplet with Ubuntu v.20+ with any basic Provision configuration. [Here](https://docs.digitalocean.com/products/droplets/getting-started/recommended-droplet-setup) you can find some help how to set up a production-ready droplet.
2. Make sure your API keys for Infura and Sparkpost.
3. Get the deployment script and make it executable:
```bash
wget https://raw.githubusercontent.com/ProxeusApp/proxeus-core/main/deploy/digitalocean/deploy.sh && chmod +x deploy.sh
```
4. Enter all the necessary variables, where **FQDN** - domain name for the future server, **INFURA** - your Infura API key, **SPARKPOST** - your SparkPost API key, **ADMINEMAIL** - admin email (in format like admin@proxeus.org). It can be done by editing the script file directly or by shell "export" command:
```bash
export FQDN=[value]
export INFURA=[value]
export SPARKPOST=[value]
export ADMINEMAIL=[value]
```
5. Run deployment process (go through it carefully, it may ask you any additional confirmation):
```bash
./deploy.sh
```
7. It takes a few minutes for the server to boot and install, then you should be able to open `http://<your DO's IP address or domain>:1323/init`
8. A configuration screen will be shown where you can set up an admin account and check settings.

Once your server is running, visit the [User Handbook](https://doc.proxeus.org/#/handbook) to get started.

To view the logs connect to your droplet using an SSH client program. Then paste this into the console to see the logs being updated in real time:

`cd /srv/proxeus && docker-compose logs -f`
