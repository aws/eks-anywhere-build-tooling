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
RELEASE_BRANCH="${2?Specify second argument - release branch}"
CLONE_URL="${3?Specify third argument - git clone endpoint}"
TAG="${4?Specify fourth argument - git version tag}"
GOLANG_VERSION="${5?Specify fifth argument - golang version}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

OUTPUT_DIR="$MAKE_ROOT/_output/$RELEASE_BRANCH"
BIN_ROOT="$OUTPUT_DIR/bin"
BIN_PATH=$BIN_ROOT/$REPO

function build::cloud-provider-vsphere::fix_licenses(){
  # The vsphere-automation-sdk-go dependency github repo has a license however
  # it is broken up into three separate go modules. Since the license file does not live in the same
  # folder as the go.mod files it is not being included in the downloaded package. Manually
  # downloading from github and placing in each of the packages from vsphere-automation-sdk-go
  # under vendor to make go-licenses happy.  The license needs to be copied into each package
  # folder, otherwise go-licenses will group them all together as vsphere-automation-sdk-go
  # which would be wrong
  wget https://raw.githubusercontent.com/vmware/vsphere-automation-sdk-go/master/LICENSE.txt
  packages=(
    "lib"
    "runtime"
    "services"
  )
  for package in "${packages[@]}"; do
    cp LICENSE.txt ./vendor/github.com/vmware/vsphere-automation-sdk-go/${package}/LICENSE.txt
  done
}

function build::cloud-provider-vsphere::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -a -ldflags='-s -w -extldflags=static' -o bin/vsphere-cloud-controller-manager ./cmd/vsphere-cloud-controller-manager
  mkdir -p ${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ${BIN_PATH}/${OS}-${ARCH}/
}

function build::cloud-provider-vsphere::binaries(){
  mkdir -p $BIN_PATH
  rm -rf $REPO
  git clone $CLONE_URL $REPO
  cd $REPO
  build::common::wait_for_tag $TAG
  git checkout $TAG
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::cloud-provider-vsphere::create_binaries "linux/amd64"
  build::cloud-provider-vsphere::fix_licenses
  build::gather_licenses $OUTPUT_DIR "./cmd/vsphere-cloud-controller-manager"
  cd ..
  rm -rf $REPO
}

build::cloud-provider-vsphere::binaries
