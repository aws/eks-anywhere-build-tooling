#!/usr/bin/env bash
# Copyright 2020 Amazon.com Inc. or its affiliates. All Rights Reserved.
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

EKSD_RELEASE_BRANCH="${1?Specify first argument - release branch}"
KIND_BASE_IMAGE_NAME="${2?Specify second argument - kind base tag}"
ARTIFACTS_BUCKET="${3?Specify third argument - artifact bucket}"
BASE_IMAGE="${4?Specify fourth argument - base image}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
source "${BUILD_LIB}/common.sh"

# Preload release yaml
build::eksd_releases::load_release_yaml $EKSD_RELEASE_BRANCH

KUBE_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $EKSD_RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $EKSD_RELEASE_BRANCH)
EKSD_KUBE_VERSION="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
COREDNS_VERSION=$(build::eksd_releases::get_eksd_component_version "coredns" $EKSD_RELEASE_BRANCH)
ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH)
CNI_PLUGINS_URL=$(build::eksd_releases::get_eksd_component_url "cni-plugins" $EKSD_RELEASE_BRANCH)
CNI_PLUGINS_AMD64_SHA256SUM=$(build::eksd_releases::get_eksd_component_sha "cni-plugins" $EKSD_RELEASE_BRANCH)

CRICTL_URL=$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/cri-tools')
# TODO: need to publish sha256 files to artifact bucket during build so it can be pulled here
CRICTL_AMD64_SHA256SUM=$(curl $CRICTL_URL | sha256sum | cut -d ' ' -f1)

# Tweak the kind/base image to have a hardcode kubeadm config
# so that during the image pull phase it pulls eks-d images
# vs upstream images
# kubeadm-override and config are copied into kind/images/base/files/etc
# so they are automatically added into the image by the dockerfile
function build::kind::override_kubeadm(){
    export EKSD_KUBE_VERSION="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
    export COREDNS_VERSION="$COREDNS_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
    export ETCD_VERSION="$ETCD_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"

    envsubst '$COREDNS_VERSION:$ETCD_VERSION:$EKSD_KUBE_VERSION' \
        < "$MAKE_ROOT/images/base/kubeadm.config.tmpl" \
        > "$MAKE_ROOT/kind/images/base/files/etc/kubeadm.config"

    cp "$MAKE_ROOT/images/base/kubeadm-override.sh" "$MAKE_ROOT/kind/images/base/files/usr/local/bin/kubeadm"
}

function build::kind::base() {
    mkdir -p "$MAKE_ROOT/_output/images/$EKSD_RELEASE_BRANCH"
    $BUILD_LIB/buildkit.sh build \
        --frontend dockerfile.v0 \
        --opt platform=linux/amd64 \
        --opt build-arg:BASE_IMAGE=$BASE_IMAGE \
        --opt build-arg:CNI_PLUGINS_URL=$CNI_PLUGINS_URL \
        --opt build-arg:CNI_PLUGINS_AMD64_SHA256SUM=$CNI_PLUGINS_AMD64_SHA256SUM \
        --opt build-arg:CRICTL_AMD64_SHA256SUM=$CRICTL_AMD64_SHA256SUM \
        --opt build-arg:CRICTL_URL=$CRICTL_URL \
        --local dockerfile="$MAKE_ROOT/kind/images/base" \
        --local context="$MAKE_ROOT/kind/images/base" \
        --opt filename=Dockerfile \
        --progress plain \
        --output type=oci,oci-mediatypes=true,name=$KIND_BASE_IMAGE_NAME:$EKSD_KUBE_VERSION,dest="$MAKE_ROOT/_output/images/$EKSD_RELEASE_BRANCH/base.tar" 
}

build::kind::override_kubeadm
build::kind::base
