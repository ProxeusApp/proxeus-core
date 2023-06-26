#!/bin/bash

# <UDF name="FQDN" Label="Fully Qualified Domain Name" example="web.example.com" />
# <UDF name="INFURA" Label="Infura.io API key" example="a0e728c9fd444a123456789000b9370f" />
# <UDF name="SPARKPOST" Label="Sparkpost.com API key" example="27ed8e1234567890000014863f9e2cf553a7bd87" />
# <UDF name="ADMINEMAIL" Label="Admin e-mail address" example="admin@example.com" />

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

log "Installing Proxeus with demo smart contract"
mkdir -p /srv
cd /srv

cat <<END >.env
PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS=0x66FF4FBF80D4a3C85a54974446309a2858221689
PROXEUS_ENCRYPTION_SECRET_KEY=`hexdump -n 16 -e '4/4 "%08X"' /dev/random`
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
