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


RELEASE_CHANNEL="${1?Specify first argument - EKS-D release channel}"
FORMAT="${2?Specify second argument - Image format}"
ARTIFACTS_PATH="${3?Specify third argument - Artifacts path}"
ARTIFACTS_UPLOAD_PATH="${4?Specify fourth argument - Artifacts upload path}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

release_availability=$(build::bottlerocket::check_release_availablilty $MAKE_ROOT/BOTTLEROCKET_RELEASES $RELEASE_CHANNEL $FORMAT)
if [ $release_availability -ne 0 ]; then
  echo "No Bottlerocket release found for release branch. Terminating silently..."
  exit 0
fi

build::common::echo_and_run make -C $MAKE_ROOT upload-artifacts ARTIFACTS_PATH=$ARTIFACTS_PATH ARTIFACTS_UPLOAD_PATH=$ARTIFACTS_UPLOAD_PATH IMAGE_FORMAT=$FORMAT IMAGE_OS=bottlerocket
