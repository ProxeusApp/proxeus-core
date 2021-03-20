#!/usr/bin/env bash
set -eo pipefail
[[ $TRACE ]] && set -x

# A script to bootstrap proxeus-core.

# Based on https://github.com/dokku/dokku/blob/master/bootstrap.sh
# We encourage you to also add Dokku for managing your instance,
# however this is not done by this script.

# It expects to be run on Debian, Ubuntu, or CentOS 7 via 'sudo'

# It checks out the proxeus source code from Github into ~/proxeus and then runs 'make install'.


log-fail() {
  declare desc="log fail formatter"
  echo "$@" 1>&2
  exit 1
}

ensure-environment() {
  local FREE_MEMORY
  if [[ -z "$GIT_TAG" ]]; then
    echo "Preparing to install $GIT_REPO..."
  else
    echo "Preparing to install $GIT_TAG from $GIT_REPO..."
  fi

  hostname -f >/dev/null 2>&1 || {
    log-fail "This installation script requires that you have a hostname set for the instance. Please set a hostname for 127.0.0.1 in your /etc/hosts"
  }

  if ! command -v apt-get &>/dev/null; then
    log-fail "This installation script supports Debian-based systems and expects apt-get."
  fi

  if ! command -v docker &> /dev/null; then
    log-fail "Docker needs to be installed."
  fi

  if ! command -v docker-compose &> /dev/null; then
    log-fail "Docker Compose needs to be installed."
  fi

  FREE_MEMORY=$(grep MemTotal /proc/meminfo | awk '{print $2}')
  if [[ "$FREE_MEMORY" -lt 1003600 ]]; then
    echo "To build containers, it is strongly suggested that you have 1024 megabytes or more of free memory"
  fi
}

install-requirements() {
  echo "--> Ensuring we have the proper dependencies"

  case "$SRV_DISTRO" in
    debian)
      if ! dpkg -l | grep -q software-properties-common; then
        apt-get update -qq >/dev/null
        apt-get -qq -y --no-install-recommends install software-properties-common
      fi
      ;;
    ubuntu)
      if ! dpkg -l | grep -q software-properties-common; then
        apt-get update -qq >/dev/null
        apt-get -qq -y --no-install-recommends install software-properties-common
      fi

      add-apt-repository universe >/dev/null
      apt-get update -qq >/dev/null
      ;;
  esac

  apt-get -qq -y --no-install-recommends install sudo git make software-properties-common
}

install-proxeus() {
  echo "--> Starting Proxeus install"

  if [[ -n $GIT_BRANCH ]]; then
    install-proxeus-from-source "origin/$GIT_BRANCH"
  elif [[ -n $GIT_TAG ]]; then
    local GIT_SEMVER="${GIT_TAG//v/}"
    major=$(echo "$GIT_SEMVER" | awk '{split($0,a,"."); print a[1]}')
    minor=$(echo "$GIT_SEMVER" | awk '{split($0,a,"."); print a[2]}')
    patch=$(echo "$GIT_SEMVER" | awk '{split($0,a,"."); print a[3]}')

    install-proxeus-from-source "$GIT_TAG"
  else
    install-proxeus-from-source
  fi
}

install-proxeus-from-source() {
  local GIT_CHECKOUT="$1"

  if [[ ! -d ./proxeus ]]; then
    git clone "$GIT_REPO" ./proxeus
  fi

  cd ./proxeus
  touch ../.env
  cp ../.env .
  git fetch origin
  [[ -n $GIT_CHECKOUT ]] && git checkout "$GIT_CHECKOUT"
  make
}

main() {
  export SRV_DISTRO SRV_DISTRO_VERSION
  # shellcheck disable=SC1091
  SRV_DISTRO=$(. /etc/os-release && echo "$ID")
  # shellcheck disable=SC1091
  SRV_DISTRO_VERSION=$(. /etc/os-release && echo "$VERSION_ID")

  export DEBIAN_FRONTEND=noninteractive
  export GIT_REPO=${GIT_REPO:-"https://github.com/ProxeusApp/proxeus-core.git"}

  ensure-environment
  install-requirements
  install-proxeus
}

main "$@"
