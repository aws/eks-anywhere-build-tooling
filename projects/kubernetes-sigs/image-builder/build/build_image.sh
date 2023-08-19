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

image_os="${1?Specify the first argument - image os}"
image_os_version="${2?Specify the second argument - image os version}"
release_channel="${3?Specify the third argument - release channel}"
image_format="${4?Specify the fourth argument - image format}"
artifacts_bucket=${5-$ARTIFACTS_BUCKET}
latest=${6-latest}
firmware=${7-}

image_os_version_arg=""
if [[ "$image_os" == "ubuntu" ]]; then
  image_os_version_arg="--os-version ${image_os_version:0:2}.${image_os_version:2:2}"
elif [[ "$image_os" == "redhat" ]]; then
  image_os_version_arg="--os-version $image_os_version"
fi

if [ ! -f "${HOME}/image-builder" ]; then
  ARCH="arm64"
  if [[ "$(uname -m)" == "x86_64" ]]; then
    ARCH="amd64"
  fi

  cp "$MAKE_ROOT/../../aws/image-builder/_output/bin/image-builder/linux-$ARCH/image-builder" "${HOME}" 
fi

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

  firmware_arg=""
  if [ -n "$firmware" ] && [[ "$image_os" == "ubuntu" ]]; then
    firmware_arg="--firmware $firmware"
  fi

  "${HOME}"/image-builder build --hypervisor vsphere --os $image_os $image_os_version_arg --vsphere-config $image_builder_config_file --release-channel $release_channel $firmware_arg
elif [[ $image_format == "raw" ]]; then
  # Run image-builder cli
  if [[ $image_os == "ubuntu" ]]; then
    "${HOME}"/image-builder build --hypervisor baremetal --os $image_os $image_os_version_arg --release-channel $release_channel
    echo "done with image builder"
  elif [[ $image_os == "redhat" ]]; then
    echo "Creating baremetal config"
    echo "$(jq --arg extra_rpms "https://redhat-iso-pdx.s3.us-west-2.amazonaws.com/8.4/rpms/kmod-megaraid_sas-07.719.06.00_el8.4-1.x86_64.rpm" \
      '. += {"extra_rpms": $extra_rpms}' $redhat_config_file)" > $image_builder_config_file

    "${HOME}"/image-builder build --hypervisor baremetal --os $image_os $image_os_version_arg --release-channel $release_channel --baremetal-config $image_builder_config_file
  fi
elif [[ $image_format == "cloudstack" ]]; then
  if [[ $image_os != "redhat" ]]; then
    echo "Cloudstack builds do not support any non-redhat OS"
    exit 1
  fi

  echo "Creating cloudstack config"
  image_builder_config_file=$redhat_config_file
  "${HOME}"/image-builder build --hypervisor cloudstack --os $image_os $image_os_version_arg --release-channel $release_channel --cloudstack-config $image_builder_config_file
elif [[ $image_format == "ami" ]]; then
  if [[ $image_os != "ubuntu" ]]; then
    echo "AMI builds do not support any non-ubuntu os"
    exit 1
  fi

  echo "Creating AMI config"
  jq --null-input \
    --arg ami_filter_owners "099720109477" \
    --arg manifest_output "$MANIFEST_OUTPUT" \
    '{"ami_filter_owners": $ami_filter_owners, "manifest_output": $manifest_output}' > $image_builder_config_file
  "${HOME}"/image-builder build --hypervisor ami --os $image_os $image_os_version_arg --release-channel $release_channel --ami-config $image_builder_config_file
fi
