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

TAR_FILE_PREFIX="$1"
OUTPUT_DIR="$2"
OUTPUT_BIN_DIR="$3"
TAG="$4"
BINARY_PLATFORMS="$5"
TAR_PATH="$6"
GIT_HASH="$7"

LICENSE_PATHS=($(find $OUTPUT_DIR -type d -name "LICENSES"))
ATTRIBUTION_PATHS=($(find $OUTPUT_DIR -type f -name "*ATTRIBUTION.txt"))

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

function build::simple::tarball() {
  build::common::ensure_tar
  mkdir -p "$TAR_PATH"
  SUPPORTED_PLATFORMS=(${BINARY_PLATFORMS// / })
  for platform in "${SUPPORTED_PLATFORMS[@]}"; do
    OS="$(cut -d '/' -f1 <<< ${platform})"
    ARCH="$(cut -d '/' -f2 <<< ${platform})"
    TAG_PATH="${TAG//\//-}"
    TAR_FILE="${TAR_FILE_PREFIX}-${OS}-${ARCH}-${TAG_PATH}.tar.gz"

    for path in "${LICENSE_PATHS[@]}"; do
      build::common::echo_and_run build::common::copy_if_source_destination_different $path ${OUTPUT_BIN_DIR}/${OS}-${ARCH}/
    done
    for path in "${ATTRIBUTION_PATHS[@]}"; do
      build::common::echo_and_run build::common::copy_if_source_destination_different $path ${OUTPUT_BIN_DIR}/${OS}-${ARCH}/
    done
    build::common::create_tarball ${TAR_PATH}/${TAR_FILE} ${OUTPUT_BIN_DIR}/${OS}-${ARCH} .
  done
  echo $GIT_HASH > $TAR_PATH/githash
}

build::simple::tarball
