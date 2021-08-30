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
BIN_FILES="_output/files"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function build::kustomize-controller::fix_licenses(){
  # The mozilla-services/gopgagent does not have a license file checked into the repo, but there is a currently open PR
  # https://github.com/mozilla-services/gopgagent/pull/4 which adds it. Until this is merged, we need to fetch the license
  # file from the commit and not from main/master.
  wget https://raw.githubusercontent.com/mozilla-services/gopgagent/39936d55b621318e919509000af38573d91c42ad/LICENSE.txt
  mv LICENSE.txt ./vendor/go.mozilla.org/gopgagent/LICENSE.txt
  
  # Internal go.mod under /api directory
  cp LICENSE ./vendor/github.com/fluxcd/kustomize-controller/api/LICENSE
}

function build::kustomize-controller::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -a -ldflags "-s -w -buildid= -extldflags '-static'" -o bin/kustomize-controller .
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::kustomize-controller::binaries(){
  mkdir -p $BIN_PATH
  mkdir -p $BIN_FILES
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  cp config/kubeconfig ../$BIN_FILES/
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::kustomize-controller::create_binaries "linux/amd64"
  build::kustomize-controller::fix_licenses
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::kustomize-controller::binaries
