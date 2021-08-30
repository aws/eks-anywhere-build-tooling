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

function build::cluster-api::gather_licenses(){
  # Pattern source: https://github.com/kubernetes-sigs/cluster-api/blob/master/Makefile/#L156-L176
  build::gather_licenses $MAKE_ROOT/_output "./cmd/clusterctl ./bootstrap/kubeadm ./controlplane/kubeadm"
  # since capd is a sperate module, go-licenses has to be run seperately
  (cd ./test/infrastructure/docker && go mod vendor && build::gather_licenses $MAKE_ROOT/_output/capd ".")
}

function build::cluster-api::build_binaries(){
  platform=$1
  OS="$(cut -d '/' -f1 <<< ${platform})"
  ARCH="$(cut -d '/' -f2 <<< ${platform})"
  export CGO_ENABLED=0
  export GOARCH=$ARCH
  export GOOS=$OS
  make manager-core
  make manager-kubeadm-bootstrap
  make manager-kubeadm-control-plane
  make clusterctl
  make manager -C test/infrastructure/docker
  mkdir -p ../${BIN_PATH}/${OS}-${ARCH}/
  mv bin/* ../${BIN_PATH}/${OS}-${ARCH}/
  mv test/infrastructure/docker/bin/manager ../${BIN_PATH}/${OS}-${ARCH}/cluster-api-provider-docker-manager
  make clean

  unset CGO_ENABLED GOARCH GOOS
}

function build::cluster-api::manifests(){
  PROD_REGISTRY="${IMAGE_REPO}/kubernetes-sigs/cluster-api"
  if [[ -v CODEBUILD_CI ]]; then
    KUBE_RBAC_PROXY_LATEST_TAG=$(aws ecr-public describe-images --region us-east-1 --output text --repository-name brancz/kube-rbac-proxy --query 'sort_by(imageDetails,& imagePushedAt)[-1].imageTags[0]')
  else
    KUBE_RBAC_PROXY_LATEST_TAG=latest
  fi

  make set-manifest-image \
    MANIFEST_IMG=$PROD_REGISTRY/cluster-api-controller MANIFEST_TAG=$IMAGE_TAG \
    TARGET_RESOURCE="./config/manager/manager_image_patch.yaml"
  # Set the kubeadm bootstrap image to the production bucket.
  make set-manifest-image \
    MANIFEST_IMG=$PROD_REGISTRY/kubeadm-bootstrap-controller MANIFEST_TAG=$IMAGE_TAG \
    TARGET_RESOURCE="./bootstrap/kubeadm/config/manager/manager_image_patch.yaml"
  # Set the kubeadm control plane image to the production bucket.
  make set-manifest-image \
    MANIFEST_IMG=$PROD_REGISTRY/kubeadm-control-plane-controller MANIFEST_TAG=$IMAGE_TAG  \
    TARGET_RESOURCE="./controlplane/kubeadm/config/manager/manager_image_patch.yaml"
  make set-manifest-image \
    MANIFEST_IMG=${IMAGE_REPO}/brancz/kube-rbac-proxy MANIFEST_TAG=$KUBE_RBAC_PROXY_LATEST_TAG \
    TARGET_RESOURCE="./config/manager/manager_auth_proxy_patch.yaml"
  make set-manifest-image \
    MANIFEST_IMG=${IMAGE_REPO}/brancz/kube-rbac-proxy MANIFEST_TAG=$KUBE_RBAC_PROXY_LATEST_TAG \
    TARGET_RESOURCE="./bootstrap/kubeadm/config/manager/manager_auth_proxy_patch.yaml"
  make set-manifest-image \
    MANIFEST_IMG=${IMAGE_REPO}/brancz/kube-rbac-proxy MANIFEST_TAG=$KUBE_RBAC_PROXY_LATEST_TAG \
    TARGET_RESOURCE="./controlplane/kubeadm/config/manager/manager_auth_proxy_patch.yaml"
  make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./config/manager/manager_pull_policy.yaml"
  make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./bootstrap/kubeadm/config/manager/manager_pull_policy.yaml"
  make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./controlplane/kubeadm/config/manager/manager_pull_policy.yaml"

  ## Build the manifests
  make release-manifests
  ## Build the development manifests
  make -C test/infrastructure/docker set-manifest-image \
    MANIFEST_IMG=$PROD_REGISTRY/capd-manager MANIFEST_TAG=$IMAGE_TAG
  sed -i 's,image: .*,image: '"${IMAGE_REPO}/brancz/kube-rbac-proxy:${KUBE_RBAC_PROXY_LATEST_TAG}"',' ./test/infrastructure/docker/config/manager/manager_auth_proxy_patch.yaml
  make -C test/infrastructure/docker set-manifest-pull-policy PULL_POLICY=IfNotPresent
  PATH="$(pwd)/hack/tools/bin:$PATH" make -C test/infrastructure/docker release-manifests 

  mkdir -p ../_output/manifests/{bootstrap-kubeadm,cluster-api,control-plane-kubeadm,infrastructure-docker}/$TAG
  cp out/bootstrap-components.yaml "../_output/manifests/bootstrap-kubeadm/$TAG"
  cp out/metadata.yaml "../_output/manifests/bootstrap-kubeadm/$TAG"

  cp out/control-plane-components.yaml "../_output/manifests/control-plane-kubeadm/$TAG"
  cp out/metadata.yaml "../_output/manifests/control-plane-kubeadm/$TAG"

  cp out/core-components.yaml "../_output/manifests/cluster-api/$TAG"
  cp out/metadata.yaml "../_output/manifests/cluster-api/$TAG"

  cp test/infrastructure/docker/out/infrastructure-components.yaml "../_output/manifests/infrastructure-docker/$TAG/infrastructure-components-development.yaml"
  cp test/infrastructure/docker/templates/cluster-template-development.yaml "../_output/manifests/infrastructure-docker/$TAG"
  cp out/metadata.yaml "../_output/manifests/infrastructure-docker/$TAG"
}

function build::cluster-api::binaries(){
  mkdir -p $BIN_PATH
  git clone $CLONE_URL $REPO
  cd $REPO
  git checkout $TAG
  git apply --verbose $MAKE_ROOT/patches/*
  build::common::use_go_version $GOLANG_VERSION
  go mod vendor
  build::cluster-api::build_binaries "linux/amd64"
  build::cluster-api::manifests
  build::cluster-api::gather_licenses
  cd ..
  rm -rf $REPO
}

build::cluster-api::binaries
