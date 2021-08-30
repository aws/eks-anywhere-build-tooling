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

function build::etcdadm-bootstrap-provider::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GO111MODULE=on GOOS=$OS GOARCH=$ARCH go build -v -o bin/manager -ldflags "-s -w -buildid=''" $(pwd)
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::etcdadm-bootstrap-provider::manifests(){
  MANIFEST_IMAGE_OVERRIDE="${IMAGE_REPO}/mrajashree/etcdadm-bootstrap-provider:${IMAGE_TAG}"
  if [[ -v CODEBUILD_CI ]]; then
    KUBE_RBAC_PROXY_LATEST_TAG=$(aws ecr-public describe-images --region us-east-1 --output text --repository-name brancz/kube-rbac-proxy --query 'sort_by(imageDetails,& imagePushedAt)[-1].imageTags[0]')
  else
    KUBE_RBAC_PROXY_LATEST_TAG=latest
  fi
  KUBE_RBAC_PROXY_IMAGE_OVERRIDE=${IMAGE_REPO}/brancz/kube-rbac-proxy:${KUBE_RBAC_PROXY_LATEST_TAG}

  sed -i "s,\${ETCDADM_BOOTSTRAP_IMAGE},${MANIFEST_IMAGE_OVERRIDE}," ./config/manager/manager.yaml
  sed -i 's,image: .*,image: '"${KUBE_RBAC_PROXY_IMAGE_OVERRIDE}"',' ./config/default/manager_auth_proxy_patch.yaml

  mkdir -p ../_output/manifests/bootstrap-etcdadm-bootstrap/${TAG}
  kustomize build config/default > bootstrap-components.yaml
  cp bootstrap-components.yaml "../_output/manifests/bootstrap-etcdadm-bootstrap/${TAG}"
  cp ../manifests/metadata.yaml "../_output/manifests/bootstrap-etcdadm-bootstrap/${TAG}"
}

function build::etcdadm-bootstrap-provider::binaries(){
  mkdir -p $BIN_PATH
  mkdir $KUSTOMIZE_BIN
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  build::install::kustomize
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::etcdadm-bootstrap-provider::create_binaries "linux/amd64"
  build::etcdadm-bootstrap-provider::manifests
  build::gather_licenses $MAKE_ROOT/_output "."
  cd ..
  rm -rf $REPO
}

build::etcdadm-bootstrap-provider::binaries
