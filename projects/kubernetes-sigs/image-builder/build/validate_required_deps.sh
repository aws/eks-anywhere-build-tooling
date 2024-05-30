#!/usr/bin/env bash
# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# set -x
set -o errexit
set -o nounset
set -o pipefail

BUILD_LIB_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${BUILD_LIB_ROOT}/build/lib/common.sh"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

IMAGE_FORMAT="$1"

DEPENDENCY_YAML=$MAKE_ROOT/REQUIRED_DEPENDENCY_VERSIONS.yaml

# python 3.9+ is already checked in upstream image-builder's ensure-python.sh
# the goss binary's shasum is compared in upstream image-builder's ensure-goss.sh and if its 
#  not the expected it is redownlaoded

echo -e "\n\n************************************************"
echo "The following are key dependencies and their current versions"
echo "Please include this output in any support requests"
echo "Python $(which python3) $(python3 --version)"
echo "Packer $(which packer) $(packer --version)"
echo "Packer Plugins:"
packer plugins installed
echo "Ansible:"
ansible --version
echo "Ansible Galaxy Collection:"
ansible-galaxy collection list
echo -e "************************************************\n"

if [ "${EKSA_SKIP_VALIDATE_DEPENDENCIES:-false}" = "true" ]; then
    echo "Skipping dependency validation, proceed with caution"
    exit 0
fi

# ansible
ANSIBLE_VERSION=$(yq ".ansible" $DEPENDENCY_YAML)
if [[ "$(ansible --version | head -n 1)" != *"[core $ANSIBLE_VERSION]"* ]]; then
    echo "The version of ansible-core ($(ansible --version | head -n 1)) does not match the version ($ANSIBLE_VERSION) which has been tested by the EKS-A team."
    echo "    (Recommened) Remove this version and rerun your image build and the correct version of ansible-core will be installed."
    echo "    You can manually fix this if you would rather, or if you are sure about the version you are using you can export EKSA_SKIP_VALIDATE_DEPENDENCIES=true"
    exit 1
fi

# packer
PACKER_VERSION=$(yq ".packer.version" $DEPENDENCY_YAML)
if [[ "$(packer --version)" != "$PACKER_VERSION" ]]; then
    echo "The version of packer ($(packer --version)) does not match the version ($PACKER_VERSION) which has been tested by the EKS-A team."
    echo "    (Recommened) Remove this version and rerun your image build and the correct version of packer will be installed."
    echo "    You can manually fix this if you would rather, or if you are sure about the version you are using you can export EKSA_SKIP_VALIDATE_DEPENDENCIES=true"
    exit
fi

# ansible plugin
PACKER_PLUGIN_ANSIBLE=$(yq ".packer.plugins.ansible" $DEPENDENCY_YAML)
if [[ "$(packer plugins installed | grep plugin-ansible)" != *"v$PACKER_PLUGIN_ANSIBLE"* ]]; then
    echo "The version of packer-plugin-ansible does not match the version ($PACKER_PLUGIN_ANSIBLE) which has been tested by the EKS-A team."
    echo "Current plugin: $(packer plugins installed | grep plugin-ansible)"
    echo -e "This is unlikely to cause issues, however if you do run into problems you can install the specific version with the following: packer plugins install github.com/hashicorp/ansible $PACKER_PLUGIN_ANSIBLE\n"
fi

# nutanix plugn
PACKER_PLUGIN_NUTANIX=$(yq ".packer.plugins.nutanix" $DEPENDENCY_YAML)
if [ "${IMAGE_FORMAT}" = "nutanix" ] && [[ "$(packer plugins installed | grep plugin-nutanix)" != *"v$PACKER_PLUGIN_NUTANIX"* ]]; then
    echo "The version of packer-plugin-nutanix does not match the version ($PACKER_PLUGIN_NUTANIX) which has been tested by the EKS-A team."
    echo "Current plugin: $(packer plugins installed | grep plugin-nutanix)"
    echo -e "This is unlikely to cause issues, however if you do run into problems you can install the specific version with the following: packer plugins install github.com/nutanix-cloud-native/nutanix $PACKER_PLUGIN_NUTANIX\n"
fi
