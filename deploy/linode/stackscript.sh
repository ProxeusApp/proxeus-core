#!/bin/bash

# <UDF name="fqdn" Label="Fully Qualified Domain Name" example="web.example.com" />
# <UDF name="infura" Label="Infura.io API key" example="a0e728c9fd444a123456789000b9370f" />
# <UDF name="sparkpost" Label="Sparkpost.com API key" example="27ed8e1234567890000014863f9e2cf553a7bd87" />

# Logs: tail -f /var/log/stackscript.log
# Logs: cat /var/log/stackscript.log

# Log to /var/log/stackscript.log for future troubleshooting

# Logging set up
exec 1> >(tee -a "/var/log/stackscript.log") 2>&1
function log {
  echo "### $1 -- `date '+%D %T'`"
}

# Common bash functions
source <ssinclude StackScriptID=1>
log "Common lib loaded"

# Apply harden script
source <ssinclude StackScriptID=394223>
log "Hardening activated"

log "Configuring System Updates"
apt-get -o Acquire::ForceIPv4=true update -y
DEBIAN_FRONTEND=noninteractive apt-get -y -o DPkg::options::="--force-confdef" -o DPkg::options::="--force-confold" install grub-pc
apt-get -o Acquire::ForceIPv4=true update -y

## Set hostname, configure apt and perform update/upgrade
log "Setting hostname"
IP=`hostname -I | awk '{print$1}'`
hostnamectl set-hostname $fqdn
echo $IP $fqdn  >> /etc/hosts

log "Updating .."
export DEBIAN_FRONTEND=noninteractive
apt-get update -y

## Remove older installations and get set for Docker install
log "Getting ready to install Docker"
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    make \
    gnupg-agent \
    software-properties-common \
    apache2-utils

log "Installing Docker Engine for $lsb_dist"
lsb_dist="$(. /etc/os-release && echo "$ID")"
lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"

## Add Docker’s official GPG key
curl -fsSL "https://download.docker.com/linux/$lsb_dist/gpg" | sudo apt-key add -

## Install stable docker as daemon
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/$lsb_dist \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce docker-ce-cli docker-compose containerd.io
systemctl enable docker

## Set up fail2ban
log "Installing fail2ban"
apt-get install fail2ban -y
cd /etc/fail2ban
cp fail2ban.conf fail2ban.local
cp jail.conf jail.local
systemctl start fail2ban
systemctl enable fail2ban

## Set up firewall with port 1323 open to default Proxeus platform
# Set up nginx separately to proxy to 443
log "Configuring firewall"
apt-get install ufw -y
ufw default allow outgoing
ufw default deny incoming

ufw allow ssh
ufw allow https
ufw allow http
ufw allow 1323

ufw enable

systemctl enable ufw
ufw logging off

## ----------------------------------------------
## Install & configure proxeus

log "Installing Proxeus"
mkdir -p /srv
cd /srv

wget https://raw.githubusercontent.com/loleg/proxeus-core/release/bootstrap.sh;
bash bootstrap.sh

cd proxeus
cat <<END >.env.prod
PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS="0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71"
PROXEUS_ALLOW_HTTP=true
PROXEUS_DATA_DIR=./data
PROXEUS_INFURA_API_KEY="$infura"
PROXEUS_SPARKPOST_API_KEY="$sparkpost"
PROXEUS_PLATFORM_DOMAIN="http://$fqdn:1323"

END

log "Starting Proxeus Core"
docker-compose --env-file .env.prod -f docker-compose.yml -f docker-compose-cloud.override.yml up -d


# Open http://$fqdn:1323/init to configure your server
log "Open http://$fqdn:1323/init to finish install"

## ----------------------------------------------

echo "Installation complete!"
