#!/usr/bin/env bash

# Update, get python-software-properties in order to get add-apt-repository,
# then update (for latest git version):
apt-get update
apt-get install -y python-software-properties
add-apt-repository -y ppa:git-core/ppa
apt-get update
apt-get install -y git

# Vim & Curl:
apt-get install -y vim curl

# Libusb
apt-get install -y libusb-dev

# Install Go
curl http://golang.org/dl/go1.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin