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
CLONE_URL="${2?Specify second argument - git clone url}"
TAG="${3?Specify third argument - git tag to checkout}"
GOLANG_VERSION="${4?Specify fourth argument - golang version}"
BINARY_NAME="${5?Specify fifth argument - binary name}"
IMAGE_REPO="${6?Specify sixth argument - ecr image repo}"
IMAGE_TAG="${7?Specify seventh argument - ecr image tag}"

BIN_ROOT="_output/bin"
BIN_PATH=$BIN_ROOT/$BINARY_NAME

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

KUSTOMIZE_BIN="${MAKE_ROOT}/_output/kustomize-bin"

function build::install::kustomize(){
  curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
  mv kustomize $KUSTOMIZE_BIN
  export PATH=$KUSTOMIZE_BIN:$PATH
}

function build::eks-anywhere-cluster-controller::create_binaries(){
  platform=${1}
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH make build-cluster-controller
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
}

function build::eks-anywhere-cluster-controller::manifests(){
  MANIFEST_IMAGE="public.ecr.aws/l0g8r8j6/eks-anywhere-cluster-controller:.*"
  MANIFEST_IMAGE_OVERRIDE="${IMAGE_REPO}/eks-anywhere-cluster-controller:${IMAGE_TAG}"

  if [[ -v CODEBUILD_CI ]]; then
    KUBE_RBAC_PROXY_LATEST_TAG=$(aws ecr-public describe-images --region us-east-1 --output text --repository-name brancz/kube-rbac-proxy --query 'sort_by(imageDetails,& imagePushedAt)[-1].imageTags[0]')
  else
    KUBE_RBAC_PROXY_LATEST_TAG=latest
  fi
  KUBE_RBAC_PROXY_IMAGE_OVERRIDE=${IMAGE_REPO}/brancz/kube-rbac-proxy:${KUBE_RBAC_PROXY_LATEST_TAG}

  sed -i 's,image: .*,image: '"${MANIFEST_IMAGE_OVERRIDE}"',' ./config/manager/manager_image_patch.yaml
  sed -i 's,image: .*,image: '"${KUBE_RBAC_PROXY_IMAGE_OVERRIDE}"',' ./config/default/manager_auth_proxy_patch.yaml

  mkdir -p ../_output/manifests/cluster-controller
  make release-manifests RELEASE_DIR=.
  cp eksa-components.yaml "../_output/manifests/cluster-controller/"
}

function build::eks-anywhere-cluster-controller::binaries(){
  mkdir -p $BIN_PATH
  mkdir -p $KUSTOMIZE_BIN
  if [ "$CI" = "true" ]; then
      cp -r /home/prow/go/src/github.com/aws/$REPO ./
  else
      git clone $CLONE_URL $REPO
  fi
  cd $REPO
  git checkout -f $TAG
  build::install::kustomize
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::eks-anywhere-cluster-controller::create_binaries "linux/amd64"
  build::eks-anywhere-cluster-controller::manifests
  build::gather_licenses $MAKE_ROOT/_output "./controllers"
  cd ..
  rm -rf $REPO
}

build::eks-anywhere-cluster-controller::binaries
