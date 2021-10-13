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

REPO="${1?Specify first argument - repository name}"
CLONE_URL="${2?Specify second argument - git clone endpoint}"
TAG="${3?Specify third argument - git version tag}"
GOLANG_VERSION="${4?Specify fourth argument - golang version}"
BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$REPO

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function build::notification-controller::fix_licenses(){
  # The notification-controller project consumes an older version (v2.3.0) of the k0kubun/pp module 
  # which does not have a LICENSE file committed to the repo. Hence we are fetching the license file 
  # from the master branch so that it is available in the vendor directory for go-licenses to pick up
  wget https://raw.githubusercontent.com/k0kubun/pp/master/LICENSE.txt
  mv LICENSE.txt ./vendor/github.com/k0kubun/pp/LICENSE.txt
  
  # Internal go.mod under /api directory
  cp LICENSE ./vendor/github.com/fluxcd/notification-controller/api/LICENSE
}

function build::notification-controller::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -a -ldflags "-s -w -buildid= -extldflags '-static'" -o bin/notification-controller .
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::notification-controller::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  build::common::wait_for_tag $TAG
  git checkout $TAG
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::notification-controller::create_binaries "linux/amd64"
  build::notification-controller::fix_licenses
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::notification-controller::binaries
