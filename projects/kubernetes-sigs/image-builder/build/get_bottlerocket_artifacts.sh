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

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

RELEASE_CHANNEL="${1?Specify first argument - EKS-D release channel}"
FORMAT="${2?Specify second argument - Image format}"
BOTTLEROCKET_DOWNLOAD_PATH="${3?Specify third argument - Download path for Bottlerocket-related files}"
PROJECT_PATH="${4?Specify fourth argument - Project path}"
LATEST_TAG="${5?Specify fifth argument - S3 destination folder for latest artifact upload}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
CODEBUILD_CI="${CODEBUILD_CI:-false}"

# Setting version and URL parameters for downloading the OVA
OS_DOWNLOAD_PATH=
KUBEVERSION=$(echo $RELEASE_CHANNEL | tr '-' '.')
BOTTLEROCKET_RELEASE_VERSION=$(yq e ".${RELEASE_CHANNEL}.${FORMAT}-release-version" $MAKE_ROOT/BOTTLEROCKET_RELEASES)
if [[ $BOTTLEROCKET_RELEASE_VERSION == "null" ]] ; then
  echo "Bottlerocket build for ${RELEASE_CHANNEL} is not enabled. Terminating silently."
  exit 0
fi
case $FORMAT in
ami)
  VARIANT="aws"
  TARGET="bottlerocket-${VARIANT}-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.img.lz4"
  ;;
raw)
  VARIANT="metal"
  TARGET="bottlerocket-${VARIANT}-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.img.lz4"
  ;;
ova)
  VARIANT="vmware"
  TARGET="bottlerocket-${VARIANT}-k8s-${KUBEVERSION}-x86_64-${BOTTLEROCKET_RELEASE_VERSION}.ova"
  ;;
*)
  echo "Invalid image format: $FORMAT"
  exit 1
  ;;
esac

BOTTLEROCKET_METADATA_URL="https://updates.bottlerocket.aws/2020-07-07/${VARIANT}-k8s-${KUBEVERSION}/x86_64/"
BOTTLEROCKET_TARGETS_URL="https://updates.bottlerocket.aws/targets/"
OS_DOWNLOAD_PATH=${BOTTLEROCKET_DOWNLOAD_PATH}/${FORMAT}
rm -rf $OS_DOWNLOAD_PATH

# If TARGET image is available at CloudFront URL, download it from CloudFront
if [[ "$(build::common::echo_and_run curl -I -L -s -o /dev/null -w "%{http_code}" https://$BOTTLEROCKET_CLOUDFRONT_ENDPOINT/$TARGET)" == "200" ]]; then
  mkdir $OS_DOWNLOAD_PATH
  build::common::echo_and_run curl -L -s https://$BOTTLEROCKET_CLOUDFRONT_ENDPOINT/$TARGET -o $OS_DOWNLOAD_PATH/$TARGET
else
  # Otherwise download the TARGET from the Bottlerocket target location using Tuftool
  build::common::echo_and_run tuftool download "${OS_DOWNLOAD_PATH}" \
      --target-name "${TARGET}" \
      --root "${BOTTLEROCKET_DOWNLOAD_PATH}/root.json" \
      --metadata-url "${BOTTLEROCKET_METADATA_URL}" \
      --targets-url "${BOTTLEROCKET_TARGETS_URL}"
fi

# Post processing for aws and metal variants
if [[ $VARIANT == "aws" ]] || [[ $VARIANT == "metal" ]]; then
  BOTTLEROCKET_IMG="${TARGET%.lz4}"
  build::common::echo_and_run lz4 --decompress ${OS_DOWNLOAD_PATH}/${TARGET} ${OS_DOWNLOAD_PATH}/${BOTTLEROCKET_IMG}
  build::common::echo_and_run gzip ${OS_DOWNLOAD_PATH}/${BOTTLEROCKET_IMG}
  rm -f ${OS_DOWNLOAD_PATH}/${TARGET}
fi
