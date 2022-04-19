#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo apt-get update
sudo apt-get upgrade -y

sudo reboot
