#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Install packages
sudo apt-get update -y
sudo apt-get install -y ipmitool