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
set -o pipefail

REPO="$1"
REPO_OWNER="$2"
TAG="$3"
TAR_FILE_PREFIX="$4"
OUTPUT_DIR="$5"
OUTPUT_BIN_DIR="$6"
TAR_PATH="$7"
GIT_HASH="$8"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

LICENSES_PATH="$OUTPUT_DIR/LICENSES"
ATTRIBUTION_PATH="${MAKE_ROOT}/ATTRIBUTION.txt"

function build::tarball() {
  build::common::ensure_tar
  mkdir -p "$TAR_PATH"
  OS="linux"
  ARCH="amd64"
  TAR_FILE="${TAR_FILE_PREFIX}-${OS}-${ARCH}-${TAG}.tar.gz"

  cp -rf $LICENSES_PATH ${OUTPUT_BIN_DIR}/${OS}-${ARCH}/
  cp $ATTRIBUTION_PATH ${OUTPUT_BIN_DIR}/${OS}-${ARCH}/
  build::common::create_tarball ${TAR_PATH}/${TAR_FILE} ${OUTPUT_BIN_DIR}/${OS}-${ARCH} .

  echo $GIT_HASH > $TAR_PATH/githash
}

build::non-golang::copy_licenses $REPO $LICENSES_PATH/github.com/$REPO_OWNER/$REPO
build::tarball