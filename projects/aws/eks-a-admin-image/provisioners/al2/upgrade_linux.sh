#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo yum clean all
sudo yum update -y
sudo yum upgrade -y

sudo reboot
