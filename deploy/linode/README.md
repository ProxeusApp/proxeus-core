Proxeus StackScript for fast Linode deployment
---

Creates a compact all-in-one instance of the Proxeus application (no code environment for smart contracts) using a bootstrapped release image for Docker. This is a good starting point for development or small installations. For more information visit https://github.com/ProxeusApp

StackScripts are private or public managed scripts which run within a Linode instance during startup. Using a simple form, you can configure the basic details needed to quickly get a Proxeus instance up and running.

This script is maintained for the community by Proxeus Association

## Instructions

1. Search for "proxeus" when deploying a new Linode, or log in and navigate to https://cloud.linode.com/stackscripts/758453
1. Additional documentation is available from [Linode Guides](https://www.linode.com/docs/guides/platform/stackscripts/)
1. You will need to have your API keys for Infura and Sparkpost handy - see the root README for further details.
1. It takes a few minutes for the server to boot and install, then you should be able to open `http://<your Linode's IP address or domain>:1323/init`
1. A configuration screen will be shown where you can set up an admin account and check settings.

Once your server is running, visit the [User Handbook](https://docs.google.com/document/d/e/2PACX-1vTchv7PotoQeH2cBA2VIHcqV0I0N_IQpFnbESR-8C19cgBikek3HAMVdPtfJJcYkANzPWbfy_S3bf8X/pub) to get started.

## References

The basic set-up of a Debian or Ubuntu server is based roughly on Linode's [Basic OCA Helper One-Click](https://cloud.linode.com/stackscripts/401712).

We suggest [Securing Public Shadowsocks Server](https://github.com/shadowsocks/shadowsocks/wiki/Securing-Public-Shadowsocks-Server) as one example guide to follow for further 'buttoning down' your instance.

[Linode](https://linode.com) is a privately-owned American cloud hosting company that provides virtual private and managed servers around the world.
