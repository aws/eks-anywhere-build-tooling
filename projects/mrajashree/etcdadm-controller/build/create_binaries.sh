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

KUSTOMIZE_BIN="${MAKE_ROOT}/_output/kustomize-bin"

function build::install::kustomize(){
  curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
  mv kustomize $KUSTOMIZE_BIN
  export PATH=$KUSTOMIZE_BIN:$PATH
}

function build::etcdadm-controller::fix_licenses(){
  cp LICENSE ./vendor/github.com/mrajashree/etcdadm-bootstrap-provider/api/v1alpha3/LICENSE
}

function build::etcdadm-controller::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GO111MODULE=on GOOS=$OS GOARCH=$ARCH go build -v -o bin/manager -ldflags "-s -w -buildid=''" $(pwd)
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::etcdadm-controller::manifests(){
  MANIFEST_IMAGE_OVERRIDE="${IMAGE_REPO}/mrajashree/etcdadm-controller:${IMAGE_TAG}"
  KUBE_RBAC_PROXY_IMAGE_OVERRIDE=${IMAGE_REPO}/brancz/kube-rbac-proxy:latest

  sed -i "s,\${ETCDADM_CONTROLLER_IMAGE},${MANIFEST_IMAGE_OVERRIDE}," ./config/manager/manager.yaml
  sed -i 's,image: .*,image: '"${KUBE_RBAC_PROXY_IMAGE_OVERRIDE}"',' ./config/default/manager_auth_proxy_patch.yaml
  mkdir -p ../_output/manifests/bootstrap-etcdadm-controller/${TAG}
  kustomize build config/default > bootstrap-components.yaml
  sed -i "s,\${ETCDADM_CONTROLLER_IMAGE},$MANIFEST_IMAGE_OVERRIDE," bootstrap-components.yaml
  cp bootstrap-components.yaml "../_output/manifests/bootstrap-etcdadm-controller/${TAG}"
  cp ../manifests/metadata.yaml "../_output/manifests/bootstrap-etcdadm-controller/${TAG}"
}

function build::etcdadm-controller::binaries(){
  mkdir -p $BIN_PATH
  mkdir $KUSTOMIZE_BIN
  git clone $CLONE_URL $REPO
  cd $REPO
  build::common::wait_for_tag $TAG
  git checkout $TAG
  build::install::kustomize
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::etcdadm-controller::create_binaries "linux/amd64"
  build::etcdadm-controller::manifests
  build::etcdadm-controller::fix_licenses
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::etcdadm-controller::binaries
