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
IMAGE_REPO="${5?Specify fifth argument - ecr image repo}"
BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$REPO


MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

KUSTOMIZE_BIN="${MAKE_ROOT}/_output/kustomize-bin"

function build::install::kustomize(){
  curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
  mv kustomize $KUSTOMIZE_BIN
  export PATH=$KUSTOMIZE_BIN:$PATH
}

function build::flux::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  make build
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::flux::create_manifests(){
  mkdir -p ../_output/manifests/gotk/$TAG
  FLUX_PATH="$MAKE_ROOT/$BIN_PATH/linux-amd64/flux"
  $FLUX_PATH install --export > gotk-components.yaml
  images=(
    helm-controller
    kustomize-controller
    notification-controller
    source-controller
  )
  for image in "${images[@]}"; do
    sed -i "s,ghcr.io/fluxcd/$image:.*,$IMAGE_REPO/fluxcd/$image:latest," gotk-components.yaml
  done

  cp gotk-components.yaml "../_output/manifests/gotk/$TAG"
}

function build::flux::binaries(){
  mkdir -p $BIN_PATH
  mkdir $KUSTOMIZE_BIN
  git clone $CLONE_URL $REPO
  cd $REPO
  build::common::wait_for_tag $TAG
  git checkout $TAG
  git apply --verbose $MAKE_ROOT/patches/*
  build::install::kustomize
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::flux::create_binaries "linux/amd64"
  build::flux::create_manifests
  build::gather_licenses $MAKE_ROOT/_output "./cmd/flux"
  cd ..
  rm -rf $REPO
}

build::flux::binaries
