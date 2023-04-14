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
BINARY_DEPS_DIR="$8"

LICENSES_PATH="$OUTPUT_DIR/LICENSES"
ATTRIBUTION_PATH="$OUTPUT_DIR/ATTRIBUTION.txt"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function build::simple::tarball() {
  build::common::ensure_tar
  mkdir -p "$TAR_PATH"
  SUPPORTED_PLATFORMS=(${BINARY_PLATFORMS// / })
  for platform in "${SUPPORTED_PLATFORMS[@]}"; do
    OS="$(cut -d '/' -f1 <<< ${platform})"
    ARCH="$(cut -d '/' -f2 <<< ${platform})"
    TAR_FILE="${TAR_FILE_PREFIX}-${OS}-${ARCH}-${TAG}.tar.gz"
    TAR_STAGING_DIR="${OUTPUT_BIN_DIR}/${OS}-${ARCH}/tar-staging"

    build::common::echo_and_run mkdir -p ${TAR_STAGING_DIR}/usr/local/bin ${TAR_STAGING_DIR}/sbin
		build::common::echo_and_run find ${OUTPUT_BIN_DIR}/${OS}-${ARCH} -not -type d -execdir cp "{}" ${TAR_STAGING_DIR}/usr/local/bin ";"
    build::common::echo_and_run cp $ATTRIBUTION_PATH ${TAR_STAGING_DIR}/usr/local/bin/CONTAINERD_ATTRIBUTION.txt
    build::common::echo_and_run cp ${BINARY_DEPS_DIR}/${OS}-${ARCH}/eksa/opencontainers/runc/runc ${TAR_STAGING_DIR}/sbin
    build::common::echo_and_run cp ${BINARY_DEPS_DIR}/${OS}-${ARCH}/eksa/opencontainers/runc/ATTRIBUTION.txt ${TAR_STAGING_DIR}/sbin/RUNC_ATTRIBUTION.txt

    build::common::create_tarball ${TAR_PATH}/${TAR_FILE} ${TAR_STAGING_DIR} .

    rm -rf ${TAR_STAGING_DIR}
  done
  echo $GIT_HASH > $TAR_PATH/githash
}

build::simple::tarball
