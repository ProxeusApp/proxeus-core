Proxeus StackScript for fast DigitalOcean deployment
---

Creates a compact all-in-one instance of the Proxeus application (no code environment for smart contracts) using a bootstrapped release image for Docker. This is a good starting point for development or small installations. For more information visit https://github.com/ProxeusApp

StackScripts are private or public managed scripts which run within a DigitalOcean droplet during startup. Using a simple form, you can configure the basic details needed to quickly get a Proxeus instance up and running.

This script is maintained for the community by Proxeus Association

## Instructions

1. Create fresh DO droplet with Ubuntu v.20+ with any basic Provision configuration
2. You will need to have your API keys for Infura and Sparkpost handy - see the root README for further details.
3. add executive permission for stackscript.sh ( chmod +x stackscript.sh )
4. run $PWD/proxeus-core/deploy/digitalocean/stackscript.sh
5. It takes a few minutes for the server to boot and install, then you should be able to open `http://<your DO's IP address or domain>:1323/init`
6. A configuration screen will be shown where you can set up an admin account and check settings.

Once your server is running, visit the [User Handbook](https://github.com/ProxeusApp/community/blob/master/handbook/handbook.md) to get started.

To view the logs connect to your droplet using an SSH client program. Then paste this into the console to see the logs being updated in real time:

`cd /srv/proxeus && docker-compose logs -f`
