#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# KIND_URL needs to be provided to the script as envvar

sudo apt-get install -y tar

mkdir kind_tmp

curl --silent -L $KIND_URL  | tar xz -C ./kind_tmp
chmod +x kind_tmp/kind
sudo mv kind_tmp/kind /usr/local/bin/kind
rm -rf kind_tmp