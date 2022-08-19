#!/bin/bash
#specify values for local variables

#FQDN=<value>
#INFURA=<value>
#SPARKPOST=<value>
#ADMINEMAIL=<value>


log "Configuring System Updates"
apt-get -o Acquire::ForceIPv4=true update -y
DEBIAN_FRONTEND=noninteractive apt-get -y -o DPkg::options::="--force-confdef" -o DPkg::options::="--force-confold" install grub-pc
apt-get -o Acquire::ForceIPv4=true update -y

## Set hostname, configure apt and perform update/upgrade
log "Setting hostname"
IP=`hostname -I | awk '{print$1}'`
hostnamectl set-hostname $FQDN
echo $IP $FQDN  >> /etc/hosts

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

## Set up firewall with defaults ports of Proxeus platform
log "Configuring firewall"
apt-get install ufw -y
ufw default allow outgoing
ufw default deny incoming

ufw allow ssh
ufw allow https
ufw allow http
ufw allow 1323
ufw allow 2115

ufw enable

systemctl enable ufw
ufw logging off

## ----------------------------------------------
## Install & configure proxeus

log "Installing Proxeus"
mkdir -p /srv
cd /srv

cat <<END >.env
PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS=0x1d3e5c81bf4bc60d41a8fbbb3d1bae6f03a75f71
PROXEUS_ENCRYPTION_SECRET_KEY=`dd bs=32 count=1 if="/dev/urandom" | base64 | tr +/ _.`
PROXEUS_DATA_DIR=./data
PROXEUS_ALLOW_HTTP=true
PROXEUS_INFURA_API_KEY=$INFURA
PROXEUS_SPARKPOST_API_KEY=$SPARKPOST
PROXEUS_PLATFORM_DOMAIN=http://$FQDN:1323
PROXEUS_VIRTUAL_HOST=$FQDN
LETSENCRYPT_EMAIL=$ADMINEMAIL

END

log "Warning: you should disable port 80 in production by removing the PROXEUS_ALLOW_HTTP line in your .env"

# Commence Proxeus installation
wget https://raw.githubusercontent.com/ProxeusApp/proxeus-core/master/bootstrap.sh;
bash bootstrap.sh

cd /srv/proxeus

# Compilation should not be necessary for a cloud install
# make init server-docker

log "Starting cloud deployment via docker-compose"
docker-compose --env-file .env -f docker-compose.yml -f docker-compose-cloud.override.yml up -d &

# Open http://$FQDN:1323/init to configure your server
log "After a minute, open: http://$FQDN:1323/init"

## ----------------------------------------------

echo "Installation complete!"
