#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Install prerequisites
sudo yum update -y
sudo amazon-linux-extras enable docker
sudo yum install -y docker tar

# Install yq
wget \
    --progress dot:giga \
    $YQ_URL
sudo mv yq_linux_amd64 /usr/local/bin/yq
chmod +x /usr/local/bin/yq
