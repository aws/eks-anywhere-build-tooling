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

function build::kind::gather_licenses(){
  # Pattern source: https://github.com/kubernetes-sigs/kind/blob/main/Makefile/#L57-L64
  build::gather_licenses $MAKE_ROOT/_output "./cmd/kind"
  (cd ./images/kindnetd && go mod vendor && build::gather_licenses $MAKE_ROOT/_output/kindnetd "./cmd/kindnetd")
}

function build::kind::build_binaries(){
  platform=$1
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -v -o bin/kind -trimpath -ldflags="-s -buildid= -w -X=sigs.k8s.io/kind/pkg/cmd/kind/version.GitCommit=${COMMIT}"
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::kindnetd::build_binaries(){
  platform=$1
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  cd ./images/kindnetd
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -v -o bin/kindnetd -ldflags="-s -buildid= -w" ./cmd/kindnetd
  mv bin/* ../../../${BIN_PATH}/${OS}-${ARCH}/
  cd ../../
}

function build::kind::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  COMMIT=$(git rev-parse HEAD 2>/dev/null)
  git checkout $TAG    
  git apply --verbose $MAKE_ROOT/patches/*
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::kind::build_binaries "linux/amd64"
  build::kind::build_binaries "darwin/amd64"
  build::kindnetd::build_binaries "linux/amd64"
  build::kind::gather_licenses
  cd ..
}

build::kind::binaries
