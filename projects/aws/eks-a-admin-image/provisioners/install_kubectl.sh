#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# KUBECTL_URL needs to be provided to the script

curl --silent -L $KUBECTL_URL --output kubectl
chmod +x kubectl
sudo mv kubectl /usr/local/bin/kubectl
