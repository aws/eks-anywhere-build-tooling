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
latest=${5-latest}

# Download and setup latest image-builder cli
image_build::common::download_latest_dev_image_builder_cli "${HOME}" $artifacts_bucket 'amd64' $latest

image_builder_config_file="${HOME}/image_builder_config_file"
redhat_config_file="${HOME}/redhat_config_file"
if [[ $image_os == "redhat" ]]; then
  jq --null-input \
    --arg rhel_username $RHSM_USERNAME \
    --arg rhel_password $RHSM_PASSWORD \
    --arg iso_url "https://redhat-iso-pdx.s3.us-west-2.amazonaws.com/8.4/rhel-8.4-x86_64-dvd.iso" \
    --arg iso_checksum_type "sha256" \
    --arg iso_checksum "ea5f349d492fed819e5086d351de47261c470fc794f7124805d176d69ddf1fcd" \
    '{"rhel_username": $rhel_username, "rhel_password": $rhel_password, "iso_url": $iso_url, "iso_checksum_type": $iso_checksum_type, "iso_checksum": $iso_checksum}' > $redhat_config_file
fi

if [[ $image_format == "ova" ]]; then
  # Setup vsphere config
  vsphere_config_file="${HOME}/vsphere_config_file"
  echo "${VSPHERE_CONNECTION_DATA}" > $vsphere_config_file

  # Run image-builder cli
  if [[ $image_os == "redhat" ]]; then
    jq -s add $vsphere_config_file $redhat_config_file > $image_builder_config_file
  else
    image_builder_config_file=$vsphere_config_file
  fi

  "${HOME}"/image-builder build --hypervisor vsphere --os $image_os --vsphere-config $image_builder_config_file --release-channel $release_channel
elif [[ $image_format == "raw" ]]; then
  # Run image-builder cli
  if [[ $image_os == "ubuntu" ]]; then
    "${HOME}"/image-builder build --hypervisor baremetal --os $image_os --release-channel $release_channel
  elif [[ $image_os == "redhat" ]]; then
    echo "Creating baremetal config"
    echo "$(jq --arg extra_rpms "https://redhat-iso-pdx.s3.us-west-2.amazonaws.com/8.4/rpms/kmod-megaraid_sas-07.719.06.00_el8.4-1.x86_64.rpm" \
      '. += {"extra_rpms": $extra_rpms}' $redhat_config_file)" > $image_builder_config_file

    "${HOME}"/image-builder build --hypervisor baremetal --os $image_os --release-channel $release_channel --baremetal-config $image_builder_config_file
  fi
elif [[ $image_format == "cloudstack" ]]; then
  echo "Creating cloudstack config"
  if [[ $image_os != "redhat" ]]; then
    echo "Cloudstack builds does not support any non-redhat os"
    exit 1
  fi
  image_builder_config_file=$redhat_config_file
  "${HOME}"/image-builder build --hypervisor cloudstack --os $image_os --release-channel $release_channel --cloudstack-config $image_builder_config_file
fi
