#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo amazon-linux-extras install kernel-5.15 -y
sudo yum clean all
sudo yum update -y
sudo yum upgrade -y

sudo reboot
