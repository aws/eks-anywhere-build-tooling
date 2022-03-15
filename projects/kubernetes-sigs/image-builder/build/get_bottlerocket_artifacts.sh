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
OS="${2?Specify second argument - base os of ova built}"
BOTTLEROCKET_DOWNLOAD_PATH="${3?Specify third argument - Download path for Bottlerocket-related files}"
CARGO_HOME="${4?Specify fourth argument - Root directory for Cargo installation}"
PROJECT_PATH="${5?Specify fifth argument - Project path}"
LATEST_TAG="${6?Specify sixth argument - S3 destination folder for latest artifact upload}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
CODEBUILD_CI="${CODEBUILD_CI:-false}"

# Setting version and URL parameters for downloading the OVA
OVA_DOWNLOAD_PATH=${BOTTLEROCKET_DOWNLOAD_PATH}/ova
KUBEVERSION=$(echo $RELEASE_CHANNEL | tr '-' '.')
BOTTLEROCKET_RELEASE_VERSION=$(yq e ".${RELEASE_CHANNEL}.releaseVersion" $MAKE_ROOT/BOTTLEROCKET_OVA_RELEASES)
if [[ $BOTTLEROCKET_RELEASE_VERSION == "null" ]]; then
  echo "Bottlerocket build for ${RELEASE_CHANNEL} is not enabled. Terminating silently."
  mkdir -p $OVA_DOWNLOAD_PATH
  touch $OVA_DOWNLOAD_PATH/temp_fake.ova
  exit 0
fi
OVA="bottlerocket-vmware-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.ova"
BOTTLEROCKET_METADATA_URL="https://updates.bottlerocket.aws/2020-07-07/vmware-k8s-${KUBEVERSION}/x86_64/"
BOTTLEROCKET_TARGETS_URL="https://updates.bottlerocket.aws/targets/"

rm -rf $OVA_DOWNLOAD_PATH
# Downloading the OVA from the Bottlerocket target location using Tuftool
$CARGO_HOME/bin/tuftool download "${OVA_DOWNLOAD_PATH}" \
    --target-name "${OVA}" \
    --root "${BOTTLEROCKET_DOWNLOAD_PATH}/root.json" \
    --metadata-url "${BOTTLEROCKET_METADATA_URL}" \
    --targets-url "${BOTTLEROCKET_TARGETS_URL}"
