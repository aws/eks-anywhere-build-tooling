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
BASE_IMAGE="${5?Specify fifth argument - base image}"
BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$REPO

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

if [ "$(uname)" == "Darwin" ]; then
  echo "This build can only be run from a Linux builder"
  exit 1;
fi

source "${MAKE_ROOT}/build/install_deps.sh" ${BASE_IMAGE}

function build::source-controller::fix_licenses(){
  # The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
  # checked in to the repo. Hence we need to manually download licenses from Github for each of them 
  # and place them in the respective folders under vendor directory so that they is available for 
  # go-licenses to pick up
  packages=(
    "gojsonpointer"
    "gojsonreference"
    "gojsonschema"
  )
  for package in "${packages[@]}"; do
    wget https://raw.githubusercontent.com/xeipuuv/${package}/master/LICENSE-APACHE-2.0.txt
    mv LICENSE-APACHE-2.0.txt ./vendor/github.com/xeipuuv/${package}/LICENSE.txt
  done
  
  # Internal go.mod under /api directory
  cp LICENSE ./vendor/github.com/fluxcd/source-controller/api/LICENSE
}

function build::source-controller::gather_licenses(){
  build::gather_licenses $MAKE_ROOT/_output "."
  build::non-golang::gather_licenses "libgit2/libgit2" "v1.1.0"
}

function build::source-controller::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  LDFLAGS=$(pkg-config --libs libgit2)
  CGO_ENABLED=1 GOOS=$OS GOARCH=$ARCH go build -a -ldflags "'${LDFLAGS}' -s -w -buildid= -extldflags '-static'" -o bin/source-controller .
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::source-controller::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::source-controller::create_binaries "linux/amd64"
  build::source-controller::fix_licenses
  build::source-controller::gather_licenses
  cd ..
  rm -rf $REPO
  rm -rf libgit2-1.1.0
}

build::source-controller::binaries
