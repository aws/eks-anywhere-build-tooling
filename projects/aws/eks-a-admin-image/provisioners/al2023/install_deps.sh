#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Install prerequisites
sudo dnf update -y
sudo dnf install -y tar wget
sudo dnf install -y docker

# Install yq
wget \
    --progress dot:giga \
    $YQ_URL
sudo mv yq_linux_amd64 /usr/local/bin/yq
chmod +x /usr/local/bin/yq
