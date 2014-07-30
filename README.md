# gocf

A crazyflie client in Go.

## Linux Setup
Setting up with virtual box will require guest additions with full usb access. 

### Install git
```
apt-get update
apt-get install -y python-software-properties
add-apt-repository -y ppa:git-core/ppa
apt-get update
apt-get install -y git
```

### Install Go
Donwload and extract
```
apt-get install -y vim curl
curl http://golang.org/dl/go1.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```
Create workspace
```
mkdir $HOME/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### Install libusb
```
apt-get install -y libusb-dev
```

## Control 

API being developed, look into files. 
