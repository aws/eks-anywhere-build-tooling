#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Install packages
sudo apt-get update -y
sudo apt-get install -y python python3 python3-pip
sudo apt-get install -y awscli
sudo snap install jq