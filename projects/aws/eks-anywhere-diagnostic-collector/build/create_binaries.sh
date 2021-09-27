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

REPO="${1?Specify first argument - repository name}"
CLONE_URL="${2?Specify second argument - git clone url}"
TAG="${3?Specify third argument - git tag to checkout}"
GOLANG_VERSION="${4?Specify fourth argument - golang version}"
BINARY_NAME="${5?Specify fifth argument - binary name}"
BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$BINARY_NAME

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function build::eks-anywhere-diagnostic-collector::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH make build
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::eks-anywhere-diagnostic-collector::binaries(){
  mkdir -p $BIN_PATH
  if [ "$CI" = "true" ]; then
      cp -r /home/prow/go/src/github.com/aws/$REPO ./
  else
      git clone $CLONE_URL $REPO
  fi
  cd $REPO
  git checkout -f $TAG
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::eks-anywhere-diagnostic-collector::create_binaries "linux/amd64"
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::eks-anywhere-diagnostic-collector::binaries
