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

function build::cert-manager::cherry_pick(){
  # The v1.1.0 cert-manager project repo references a 3rd-party fork of a fork of golang/crypto (https://github.com/meyskens/crypto).
  # In a later commit (https://github.com/jetstack/cert-manager/commit/b5be5a8730a7307072c1d5d5fc6f3eb57b6018f4), this crypto fork
  # was moved into the cert-manager repo and consumed from there. We are cherry-picking this commit so that we can use the changes
  # introduced by it
  git config --global user.email "eks-anywhere-bot@amazonaws.com"
  git config --global user.name "EKS Anywhere Bot"
  git cherry-pick -m 1 b5be5a8730a7307072c1d5d5fc6f3eb57b6018f4 --strategy-option theirs
}

function build::cert-manager::build_binaries(){
  platform=$1
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  export CGO_ENABLED=0
  export GOARCH=$ARCH
  export GOOS=$OS
  go build -ldflags "-s -w -extldflags -static" -o bin/cert-manager-acmesolver -v ./cmd/acmesolver
  go build -ldflags "-s -w -extldflags -static" -o bin/cert-manager-cainjector -v ./cmd/cainjector
  go build -ldflags "-s -w -extldflags -static" -o bin/cert-manager-controller -v ./cmd/controller
  go build -ldflags "-s -w -extldflags -static" -o bin/cert-manager-webhook -v ./cmd/webhook
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::cert-manager::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  build::cert-manager::cherry_pick
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::cert-manager::build_binaries "linux/amd64"
  build::gather_licenses $MAKE_ROOT/_output "./cmd/acmesolver ./cmd/cainjector ./cmd/controller ./cmd/webhook"
  
  cd ..
  rm -rf $REPO
}

build::cert-manager::binaries
