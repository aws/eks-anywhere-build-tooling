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

function build::kube-rbac-proxy::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GO111MODULE=auto GOOS=$OS GOARCH=$ARCH go build -trimpath -v -o bin/kube-rbac-proxy-linux-amd64 -ldflags "-s -w -buildid=''" $(pwd)
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::kube-rbac-proxy::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  build::common::wait_for_tag $TAG
  git checkout $TAG
  build::common::use_go_version $GOLANG_VERSION
  build::common::set_go_cache kube-rbac-proxy $TAG
  go mod vendor
  build::kube-rbac-proxy::create_binaries "linux/amd64"
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::kube-rbac-proxy::binaries
