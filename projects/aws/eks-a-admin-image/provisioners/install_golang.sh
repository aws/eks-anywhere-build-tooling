#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

if [[ "$GO_VERSION" == "latest" ]]; then
    GO_VERSION="$(curl --silent https://go.dev/VERSION?m=text | head -n 1)";
    GO_URL="https://go.dev/dl/${GO_VERSION}.linux-amd64.tar.gz"
else
    GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
fi

sudo apt-get install -y tar
sudo apt-get install -y git

curl -L --silent $GO_URL | sudo tar xz -C /usr/local

export GOROOT=/usr/local/go
export GOPATH=/usr/local/go/packages
export GOPROXY=direct

sudo sh -c "echo 'PATH=\$PATH:$GOROOT/bin:$GOPATH/bin' >> /etc/profile.d/golang.sh"

source /etc/profile.d/golang.sh

GO=/usr/local/go/bin/go

$GO version

# required for e2e tests
sudo GOPROXY=direct GOPATH=/usr/local/go/packages $GO install gotest.tools/gotestsum@latest
sudo cp $GOPATH/bin/gotestsum /usr/local/bin/gotestsum

sudo env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $GO build -o test2json -ldflags="-s -w" cmd/test2json
sudo mv test2json /usr/local/bin/test2json