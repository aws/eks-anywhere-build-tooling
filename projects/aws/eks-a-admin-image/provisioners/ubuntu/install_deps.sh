#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Install prerequisites
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common tar

# Install docker gpg key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Adding docker repo
LSB_RELEASE=$(lsb_release -cs)
ARCH=$(dpkg --print-architecture)
SOURCE="deb [arch=$ARCH] https://download.docker.com/linux/ubuntu $LSB_RELEASE stable"
sudo add-apt-repository "$SOURCE"

# Install docker and containerd
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install yq
sudo snap install yq
