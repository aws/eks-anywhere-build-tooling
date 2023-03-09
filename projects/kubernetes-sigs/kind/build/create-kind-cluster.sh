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

IMAGE="$1"
ARCH="$2"
EKSD_RELEASE_BRANCH="$3"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
source "${BUILD_LIB}/common.sh"

# Preload release yaml
build::eksd_releases::load_release_yaml $EKSD_RELEASE_BRANCH false

KUBE_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $EKSD_RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $EKSD_RELEASE_BRANCH)
EKSD_VERSION_SUFFIX="eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
EKSD_KUBE_VERSION="$KUBE_VERSION-$EKSD_VERSION_SUFFIX"
COREDNS_VERSION=$(build::eksd_releases::get_eksd_component_version "coredns" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX
ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX
EKSD_IMAGE_REPO=$(build::eksd_releases::get_eksd_image_repo $EKSD_RELEASE_BRANCH)

# Make sure the correct arch image is pulled and tagged
build::docker::retry_pull --platform linux/$ARCH $IMAGE

KIND_PATH="$MAKE_ROOT/_output/bin/kind/$(uname | tr '[:upper:]' '[:lower:]')-$(go env GOHOSTARCH)/kind"
cat << EOF \
  | $KIND_PATH create cluster \
    --name "eks-a-kind-test-$ARCH" \
    --image="$IMAGE" \
    --wait="5m" -v9 --retain \
    --config=/dev/stdin \
    || true
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
kubeadmConfigPatches:
- |
    kind: ClusterConfiguration
    dns:
        type: CoreDNS
        imageRepository: $EKSD_IMAGE_REPO/coredns
        imageTag: $COREDNS_VERSION
    etcd:
        local:
            imageRepository: $EKSD_IMAGE_REPO/etcd-io
            imageTag: $ETCD_VERSION
    imageRepository: $EKSD_IMAGE_REPO/kubernetes
    kubernetesVersion: $EKSD_KUBE_VERSION
EOF
