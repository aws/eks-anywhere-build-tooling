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

PROJECT_ROOT="$1"
TARGET_FILE="$2"
REPO="$3"
GOLANG_VERSION="$4"
TAG="$5"
BINARY_PLATFORMS="$6"
SOURCE_PATTERN="$7"
GOBUILD_COMMAND="$8" 
EXTRA_GOBUILD_FLAGS="$9"
GO_LDFLAGS="${10}"
REPO_SUBPATH="${11:-}"


SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

function build::simple::binaries(){
  mkdir -p $(dirname $TARGET_FILE)
  cd "$PROJECT_ROOT/$REPO/$REPO_SUBPATH"
  local -r cache_key=$(echo $PROJECT_ROOT | sed 's/\(.*\)\//\1-/' | xargs basename)
  build::common::use_go_version $GOLANG_VERSION
  build::common::set_go_cache $cache_key $TAG
  go mod vendor
  SUPPORTED_PLATFORMS=(${BINARY_PLATFORMS// / })
  for platform in "${SUPPORTED_PLATFORMS[@]}"; do
    OS="$(cut -d '/' -f1 <<< ${platform})"
    ARCH="$(cut -d '/' -f2 <<< ${platform})"
    CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH \
      go $GOBUILD_COMMAND -trimpath -a -ldflags "$GO_LDFLAGS" $EXTRA_GOBUILD_FLAGS -o $TARGET_FILE $SOURCE_PATTERN
  done
}

build::simple::binaries
