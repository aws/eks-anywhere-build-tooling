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

set -xe
set -o errexit
set -o nounset

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/build/setup_image_builder_cli.sh"

image_os="${1?Specify the first argument - image os}"
release_channel="${2?Specify the second argument - release channel}"
image_format="${3?Specify the third argument - image format}"
artifacts_bucket=${4-$ARTIFACTS_BUCKET}

# Download and setup latest image-builder cli
image_build::common::download_latest_dev_image_builder_cli "${HOME}" $artifacts_bucket

if [[ $image_format == "ova" ]]; then
  # Setup vsphere config
  vsphere_config_file="${HOME}/vsphere_config_file"
  echo "${VSPHERE_CONNECTION_DATA}" > $vsphere_config_file

  # Run image-builder cli
  if [[ $image_os == "redhat" ]]; then
    echo "$(jq --arg rhel_username $RHSM_USERNAME \
               --arg rhel_password $RHSM_PASSWORD \
               --arg iso_url "https://redhat-iso-pdx.s3.us-west-2.amazonaws.com/7.9/rhel-server-7.9-x86_64-dvd.iso" \
               --arg iso_checksum_type "sha256" \
               --arg iso_checksum "2cb36122a74be084c551bc7173d2d38a1cfb75c8ffbc1489c630c916d1b31b25" \
               '. += {"rhel_username": $rhel_username, "rhel_password": $rhel_password, "iso_url": $iso_url, "iso_checksum_type": $iso_checksum_type, "iso_checksum": $iso_checksum}' $vsphere_config_file)" > $vsphere_config_file
  fi
  "${HOME}"/image-builder build --hypervisor vsphere --os $image_os --vsphere-config $vsphere_config_file --release-channel $release_channel
elif [[ $image_format == "raw" ]]; then
  # Run image-builder cli
  if [[ $image_os == "ubuntu" ]]; then
    "${HOME}"/image-builder build --hypervisor baremetal --os $image_os --release-channel $release_channel
  elif [[ $image_os == "redhat" ]]; then
    echo "Redhat image building is not yet supported"
    exit 1
  fi
fi
