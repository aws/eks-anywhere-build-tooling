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

EKSD_RELEASE_BRANCH="$1"
KINDNETD_IMAGE_COMPONENT="$2"
IMAGE_REPO="$3"
ARTIFACTS_BUCKET="$4"
IMAGE_TAG="$5"
LATEST="$6"
OUTPUT_FILE="$7"
IMAGE_TAG_SUFFIX="$8"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
source "${BUILD_LIB}/common.sh"

# This is used by the local-path-provisioner within the kind node
AL2_HELPER_IMAGE="public.ecr.aws/amazonlinux/amazonlinux:2"
LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE="$IMAGE_REPO/rancher/local-path-provisioner:$LATEST"
LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE="public.ecr.aws/eks-anywhere/rancher/local-path-provisioner:$(cat $MAKE_ROOT/../../rancher/local-path-provisioner/GIT_TAG)"
KIND_KINDNETD_RELEASE_OVERRIDE="public.ecr.aws/eks-anywhere/kubernetes-sigs/kind/kindnetd:$(cat $MAKE_ROOT/GIT_TAG)"
KIND_KINDNETD_IMAGE_OVERRIDE="$IMAGE_REPO/$KINDNETD_IMAGE_COMPONENT:$LATEST$IMAGE_TAG_SUFFIX"

# Preload release yaml
build::eksd_releases::load_release_yaml $EKSD_RELEASE_BRANCH

KUBE_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $EKSD_RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $EKSD_RELEASE_BRANCH)
EKSD_KUBE_VERSION="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
PAUSE_IMAGE_TAG_OVERRIDE=$(build::eksd_releases::get_eksd_kubernetes_image_url "pause-image" $EKSD_RELEASE_BRANCH)
EKSD_IMAGE_REPO=$(build::eksd_releases::get_eksd_image_repo $EKSD_RELEASE_BRANCH)
EKSD_ASSET_URL=$(build::eksd_releases::get_eksd_kubernetes_asset_base_url $EKSD_RELEASE_BRANCH)/$KUBE_VERSION

EKSD_VERSION_SUFFIX="eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
COREDNS_VERSION=$(build::eksd_releases::get_eksd_component_version "coredns" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX
ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX

# Expected versions provided by kind which are replaced in the docker build with our versions
# when updating kind check the following, they may need to be updated
# https://github.com/kubernetes-sigs/kind/blob/v0.29.0/pkg/build/nodeimage/const_cni.go#L23
KINDNETD_IMAGE_TAG="docker.io/kindest/kindnetd:v20250512-df8de77b"
# https://github.com/kubernetes-sigs/kind/blob/v0.29.0/pkg/build/nodeimage/const_storage.go#L29
LOCAL_PATH_PROVISONER_IMAGE_TAG="docker.io/kindest/local-path-provisioner:v20250214-acbabc1a"
# https://github.com/kubernetes-sigs/kind/blob/v0.29.0/pkg/build/nodeimage/const_storage.go#L29
LOCAL_PATH_HELPER_IMAGE_TAG="docker.io/kindest/local-path-helper:v20241212-8ac705d0"
# https://github.com/kubernetes-sigs/kind/blob/v0.29.0/images/base/files/etc/containerd/config.toml#L37
PAUSE_IMAGE_TAG="registry.k8s.io/pause:3.10"

mkdir -p $(dirname $OUTPUT_FILE)
cat <<EOF >> $OUTPUT_FILE
AL2_HELPER_IMAGE=$AL2_HELPER_IMAGE
LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE=$LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE
LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE=$LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE
KIND_KINDNETD_RELEASE_OVERRIDE=$KIND_KINDNETD_RELEASE_OVERRIDE
KIND_KINDNETD_IMAGE_OVERRIDE=$KIND_KINDNETD_IMAGE_OVERRIDE
KUBE_VERSION=$KUBE_VERSION
EKSD_RELEASE=$EKSD_RELEASE
EKSD_KUBE_VERSION=$EKSD_KUBE_VERSION
PAUSE_IMAGE_TAG_OVERRIDE=$PAUSE_IMAGE_TAG_OVERRIDE
EKSD_IMAGE_REPO=$EKSD_IMAGE_REPO
EKSD_ASSET_URL=$EKSD_ASSET_URL
KINDNETD_IMAGE_TAG=$KINDNETD_IMAGE_TAG
LOCAL_PATH_HELPER_IMAGE_TAG=$LOCAL_PATH_HELPER_IMAGE_TAG
LOCAL_PATH_PROVISONER_IMAGE_TAG=$LOCAL_PATH_PROVISONER_IMAGE_TAG
PAUSE_IMAGE_TAG=$PAUSE_IMAGE_TAG
NODE_IMAGE_TAG=$EKSD_KUBE_VERSION-$IMAGE_TAG
NODE_IMAGE_LATEST_TAG=$EKSD_KUBE_VERSION-$LATEST
ETCD_IMAGE_TAG=$EKSD_IMAGE_REPO/etcd-io/etcd:$ETCD_VERSION
COREDNS_IMAGE_TAG=$EKSD_IMAGE_REPO/coredns/coredns:$COREDNS_VERSION
EOF
