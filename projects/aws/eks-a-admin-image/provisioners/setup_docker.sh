#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Enabling docker
sudo systemctl start docker
sudo systemctl enable docker

USER="${USER:-ubuntu}"
# Setting permisions to avoid sudo
sudo usermod -a -G docker $USER

# Rebooting to get permissions
sudo reboot
