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

set -o errexit
set -o nounset
set -o pipefail

BUILD_LIB_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${BUILD_LIB_ROOT}/build/lib/common.sh"

SED=$(build::find::gnu_variant_on_mac sed)

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

IMAGE_BUILDER_DIR=$MAKE_ROOT/image-builder/images/capi

DEPENDENCY_YAML=$MAKE_ROOT/REQUIRED_DEPENDENCY_VERSIONS.yaml

MINIMUM_PYTHON_VERSION="$($SED -n "s/^minimum_python_version=\(\S*\)/\1/p" $IMAGE_BUILDER_DIR/hack/ensure-python.sh)"
ANSIBLE_VERSION="$($SED -n "s/^_version_ansible_core=\"\(\S*\)\"/\1/p" $IMAGE_BUILDER_DIR/hack/utils.sh)"

PACKER_VERSION="$($SED -n "s/^_version=\"\(\S*\)\"/\1/p" $IMAGE_BUILDER_DIR/hack/ensure-packer.sh)"

PACKER_PLUGIN_ANSIBLE_VERSION="$($SED -n "s/^\s*version = \">= \(\S*\)\"/\1/p" $IMAGE_BUILDER_DIR/packer/config.pkr.hcl | sed -n 1p)"
PACKER_PLUGIN_GOSS_VERSION="$($SED -n "s/^\s*version = \">= \(\S*\)\"/\1/p" $IMAGE_BUILDER_DIR/packer/config.pkr.hcl | sed -n 2p)"
PACKER_PLUGIN_NUTANIX_VERSION="$($SED -n "s/^\s*version = \">= \(\S*\)\"/\1/p" $IMAGE_BUILDER_DIR/packer/nutanix/config.pkr.hcl)"

rm $DEPENDENCY_YAML

yq e ".packer={\"version\":\"$PACKER_VERSION\",\"plugins\":{\"ansible\":\"$PACKER_PLUGIN_ANSIBLE_VERSION\",\"goss\":\"$PACKER_PLUGIN_GOSS_VERSION\",\"nutanix\":\"$PACKER_PLUGIN_NUTANIX_VERSION\"}}" $DEPENDENCY_YAML > $DEPENDENCY_YAML
echo "ansible: $ANSIBLE_VERSION" >> $DEPENDENCY_YAML
echo "python: $MINIMUM_PYTHON_VERSION" >> $DEPENDENCY_YAML

yq -i 'sort_keys(..)' $DEPENDENCY_YAML 

yq $DEPENDENCY_YAML
