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
CLONE_URL="${2?Specify second argument - git clone endpoint}"
TAG="${3?Specify third argument - git version tag}"
GOLANG_VERSION="${4?Specify fourth argument - golang version}"
IMAGE_REPO="${5?Specify fifth argument - ecr image repo}"
IMAGE_TAG="${6?Specify sixth argument - ecr image tag}"
BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$REPO

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function build::cluster-api-provider-vsphere::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags "${LDFLAGS} -s -w -buildid= -extldflags '-static'" -o bin/manager .
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::cluster-api-provider-vsphere::manifests(){
  KUBE_RBAC_PROXY_IMAGE_OVERRIDE=${IMAGE_REPO}/brancz/kube-rbac-proxy:latest

  sed -i 's,image: .*,image: '"${KUBE_RBAC_PROXY_IMAGE_OVERRIDE}"',' ./config/manager/manager_auth_proxy_patch.yaml
  make manifests STAGE="release" \
    MANIFEST_DIR="out" \
    PULL_POLICY="IfNotPresent" \
    IMAGE="${IMAGE_REPO}/kubernetes-sigs/cluster-api-provider-vsphere/release/manager:$IMAGE_TAG"

  mkdir -p ../_output/manifests/infrastructure-vsphere/$TAG
  cp out/cluster-template.yaml "../_output/manifests/infrastructure-vsphere/$TAG"
  cp out/infrastructure-components.yaml "../_output/manifests/infrastructure-vsphere/$TAG"
  cp metadata.yaml "../_output/manifests/infrastructure-vsphere/$TAG"
}

function build::cluster-api-provider-vsphere::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  source ./hack/version.sh;
  LDFLAGS=$(version::ldflags)
  git checkout $TAG
  git apply --verbose $MAKE_ROOT/patches/*
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::cluster-api-provider-vsphere::create_binaries "linux/amd64"
  build::cluster-api-provider-vsphere::manifests
  # Pattern source: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/Makefile/#L139-L147
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::cluster-api-provider-vsphere::binaries
