#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo dnf clean all
sudo dnf update -y
sudo dnf upgrade -y

sudo reboot
