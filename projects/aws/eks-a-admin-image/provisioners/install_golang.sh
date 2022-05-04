#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

if [[ "$GO_VERSION" == "latest" ]]; then
    GO_VERSION="$(curl --silent https://go.dev/VERSION?m=text)";
    GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
else
    GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
fi

sudo apt-get install -y tar
sudo apt-get install -y git

curl -L --silent $GO_URL | sudo tar xz -C /usr/local

export GOROOT=/usr/local/go
export GOPATH=/usr/local/go/packages

sudo sh -c "echo 'PATH=\$PATH:$GOROOT/bin:$GOPATH/bin' >> /etc/profile.d/golang.sh"

source /etc/profile.d/golang.sh

go version