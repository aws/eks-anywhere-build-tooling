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

set -x
set -o errexit
set -o nounset
set -o pipefail


RELEASE_CHANNEL="${1?Specify first argument - EKS-D release channel}"
FORMAT="${2?Specify second argument - Image format}"
BOTTLEROCKET_DOWNLOAD_PATH="${3?Specify third argument - Download path for Bottlerocket-related files}"
PROJECT_PATH="${4?Specify fourth argument - Project path}"
LATEST_TAG="${5?Specify fifth argument - S3 destination folder for latest artifact upload}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
CODEBUILD_CI="${CODEBUILD_CI:-false}"

# Setting version and URL parameters for downloading the OVA
OS_DOWNLOAD_PATH=
VARIANT="vmware"
if [[ $FORMAT == "raw" ]]; then
  VARIANT="metal"
fi
KUBEVERSION=$(echo $RELEASE_CHANNEL | tr '-' '.')
BOTTLEROCKET_RELEASE_VERSION=$(yq e ".${RELEASE_CHANNEL}.${FORMAT}-release-version" $MAKE_ROOT/BOTTLEROCKET_RELEASES)
if [[ $BOTTLEROCKET_RELEASE_VERSION == "null" ]] ; then
  echo "Bottlerocket build for ${RELEASE_CHANNEL} is not enabled. Terminating silently."
  exit 0
fi

BOTTLEROCKET_METADATA_URL="https://updates.bottlerocket.aws/2020-07-07/${VARIANT}-k8s-${KUBEVERSION}/x86_64/"
BOTTLEROCKET_TARGETS_URL="https://updates.bottlerocket.aws/targets/"
OS_DOWNLOAD_PATH=${BOTTLEROCKET_DOWNLOAD_PATH}/${FORMAT}
TARGET=
if [[ $VARIANT == "vmware" ]]; then
  TARGET="bottlerocket-vmware-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.ova"
fi
if [[ $VARIANT == "metal" ]]; then
  TARGET="bottlerocket-metal-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.img.lz4"
fi
rm -rf $OS_DOWNLOAD_PATH
# Downloading the TARGET from the Bottlerocket target location using Tuftool
tuftool download "${OS_DOWNLOAD_PATH}" \
    --target-name "${TARGET}" \
    --root "${BOTTLEROCKET_DOWNLOAD_PATH}/root.json" \
    --metadata-url "${BOTTLEROCKET_METADATA_URL}" \
    --targets-url "${BOTTLEROCKET_TARGETS_URL}"

# Post processing for metal
if [[ $VARIANT == "metal" ]]; then
  BOTTLEROCKET_METAL_IMG="bottlerocket-metal-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.img"
  lz4 --decompress ${OS_DOWNLOAD_PATH}/${TARGET} ${OS_DOWNLOAD_PATH}/${BOTTLEROCKET_METAL_IMG}
  gzip ${OS_DOWNLOAD_PATH}/${BOTTLEROCKET_METAL_IMG}
  rm -f ${OS_DOWNLOAD_PATH}/${TARGET}
fi
