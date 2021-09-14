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

function build::troubleshoot::build_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  echo $GOPATH
  export PATH=$PATH:$(go env GOPATH)/bin
  go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0
  make support-bundle
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::troubleshoot::fix_licenses(){
  # The tj/go-spin dependency github repo has a license file however it does not have a go.mod file
  # checked in to the repo. Hence we need to manually download the license from Github and place
  # it in the respective folder under vendor directory so that they is available for go-licenses
  # to pick up
  wget https://raw.githubusercontent.com/tj/go-spin/master/LICENSE -O vendor/github.com/tj/go-spin/LICENSE
}

function build::troubleshoot::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  git apply --verbose $MAKE_ROOT/patches/*
  build::common::use_go_version $GOLANG_VERSION
  go mod tidy
  go mod vendor
  build::troubleshoot::build_binaries "linux/amd64"
  build::troubleshoot::gather_licenses
  cd ..
  rm -rf $REPO
}

function build::troubleshoot::gather_licenses(){
  (go mod vendor && build::troubleshoot::fix_licenses && build::gather_licenses $MAKE_ROOT/_output "./cmd/troubleshoot")
}

build::troubleshoot::binaries
