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
source "${MAKE_ROOT}/../../../build/lib/common.sh"

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
redhat_config_file="${MAKE_ROOT}/redhat-config.json"

function retry_image_builder() {
  local n=1
  local max=3
  local delay=30
  local failed=""
  declare -A retryable_messages=(
    ["Timeout waiting for IP."]="Failed waiting for IP"
    ["Timeout waiting for SSH"]="Wrong VM IP might be fetched"
    ["Cancelling provisioner after a timeout"]="Provisioner timed out")

  until [ $n -eq $max ]; do
    failed="false"
    timeout "1.5h" ${HOME}/image-builder "$@" && break || {
      failed="true"

      local log_file=$(find $MAKE_ROOT -name "packer.log" -type f)
      if [ ! -f "$log_file" ]; then
        >&2 echo "packer.log not found in ${MAKE_ROOT}!"
        break
      fi

      local retry="false"
      local message=""
      for key in "${!retryable_messages[@]}"; do
        if grep -q "$key" "$log_file"; then
          message="${retryable_messages[$key]}"
          retry="true"
          break
        fi
      done

      if [ "${retry}" = "true" ]; then
        ((n++))
        >&2 echo "$message. This is likely transisent, retrying. Attempt $n/$max:"
        sleep $delay
      else
        break
      fi
    }
  done

  if [ "${failed}" = "true" ]; then
    >&2 echo "The command has failed after $n attempts."
    exit 1;
  fi
}

if [[ $image_format == "ova" ]]; then
  # Setup vsphere config
  vsphere_config_file="${HOME}/vsphere_config_file"
  echo "${VSPHERE_CONNECTION_DATA}" > $vsphere_config_file

  echo "Creating VSphere image-builder config"
  if [[ $image_os == "redhat" ]]; then
    jq -s add $vsphere_config_file $redhat_config_file > $image_builder_config_file
  else
    image_builder_config_file=$vsphere_config_file
  fi

  firmware_arg=""
  if [ -n "$firmware" ] && [[ "$image_os" == "ubuntu" ]]; then
    firmware_arg="--firmware $firmware"
  fi
  cat $image_builder_config_file

  # Run image-builder CLI
  retry_image_builder build --hypervisor vsphere --os $image_os $image_os_version_arg --vsphere-config $image_builder_config_file --release-channel $release_channel $firmware_arg
elif [[ $image_format == "raw" ]]; then
  echo "Creating Bare metal image-builder config"
  if [[ $image_os == "ubuntu" ]]; then
    # Run image-builder CLI
    retry_image_builder build --hypervisor baremetal --os $image_os $image_os_version_arg --release-channel $release_channel
  elif [[ $image_os == "redhat" ]]; then
    image_builder_config_file=$redhat_config_file
    cat $image_builder_config_file

    # Run image-builder CLI
    retry_image_builder build --hypervisor baremetal --os $image_os $image_os_version_arg --release-channel $release_channel --baremetal-config $image_builder_config_file
  fi
elif [[ $image_format == "cloudstack" ]]; then
  if [[ $image_os != "redhat" ]]; then
    echo "Cloudstack builds do not support any non-redhat OS"
    exit 1
  fi

  echo "Creating Cloudstack image-builder config"
  image_builder_config_file=$redhat_config_file
  cat $image_builder_config_file

  # Run image-builder CLI
  retry_image_builder build --hypervisor cloudstack --os $image_os $image_os_version_arg --release-channel $release_channel --cloudstack-config $image_builder_config_file
elif [[ $image_format == "ami" ]]; then
  if [[ $image_os != "ubuntu" ]]; then
    echo "AMI builds do not support any non-ubuntu os"
    exit 1
  fi

  echo "Creating AMI image-builder config"
  jq --null-input \
    --arg ami_filter_owners "099720109477" \
    --arg manifest_output "$MANIFEST_OUTPUT" \
    '{"ami_filter_owners": $ami_filter_owners, "manifest_output": $manifest_output}' > $image_builder_config_file
  cat $image_builder_config_file

  # Run image-builder CLI
  retry_image_builder build --hypervisor ami --os $image_os $image_os_version_arg --release-channel $release_channel --ami-config $image_builder_config_file
elif [[ $image_format == "nutanix" ]]; then
  # Setup nutanix config
  nutanix_config_file="${HOME}/nutanix_config_file"
  echo "${NUTANIX_CONNECTION_DATA}" > $nutanix_config_file
  image_name=${image_os}-${image_os_version}-kube-v${release_channel}
  build::jq::update_in_place $nutanix_config_file '.image_name = '"\"$image_name\""'' 

  echo "Creating Nutanix image-builder config"
  if [[ $image_os == "redhat" ]]; then
    jq -s add $nutanix_config_file $redhat_config_file > $image_builder_config_file
  else
    image_builder_config_file=$nutanix_config_file
  fi
  cat $image_builder_config_file

  # Run image-builder CLI
  retry_image_builder build --hypervisor nutanix --os $image_os $image_os_version_arg --nutanix-config $image_builder_config_file --release-channel $release_channel
fi
