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

set -o errexit
set -o nounset
set -o pipefail

EKSD_RELEASE_BRANCH="$1"
ARTIFACTS_BUCKET="$2"
OUTPUT_FILE="$3"
LATEST_TAG="$4"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
source "${BUILD_LIB}/common.sh"

# Preload release yaml
build::eksd_releases::load_release_yaml $EKSD_RELEASE_BRANCH false

KUBE_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $EKSD_RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $EKSD_RELEASE_BRANCH)
EKSD_IMAGE_REPO=$(build::eksd_releases::get_eksd_image_repo $EKSD_RELEASE_BRANCH)
EKSD_VERSION_SUFFIX="eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
EKSD_KUBE_VERSION="$KUBE_VERSION-$EKSD_VERSION_SUFFIX"
COREDNS_VERSION=$(build::eksd_releases::get_eksd_component_version "coredns" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX
ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH)-$EKSD_VERSION_SUFFIX

CNI_PLUGINS_AMD64_URL=$(build::eksd_releases::get_eksd_component_url "cni-plugins" $EKSD_RELEASE_BRANCH amd64)
CNI_PLUGINS_ARM64_URL=$(build::eksd_releases::get_eksd_component_url "cni-plugins" $EKSD_RELEASE_BRANCH arm64)
CNI_PLUGINS_AMD64_SHA256SUM=$(build::eksd_releases::get_eksd_component_sha "cni-plugins" $EKSD_RELEASE_BRANCH amd64)
CNI_PLUGINS_ARM64_SHA256SUM=$(build::eksd_releases::get_eksd_component_sha "cni-plugins" $EKSD_RELEASE_BRANCH arm64)

mkdir -p $(dirname $OUTPUT_FILE)
cat <<EOF >> $OUTPUT_FILE
EKSD_KUBE_VERSION=$EKSD_KUBE_VERSION
EKSD_IMAGE_REPO=$EKSD_IMAGE_REPO
COREDNS_VERSION=$COREDNS_VERSION
ETCD_VERSION=$ETCD_VERSION
CNI_PLUGINS_AMD64_URL=$CNI_PLUGINS_AMD64_URL
CNI_PLUGINS_ARM64_URL=$CNI_PLUGINS_ARM64_URL
CNI_PLUGINS_AMD64_SHA256SUM=$CNI_PLUGINS_AMD64_SHA256SUM
CNI_PLUGINS_ARM64_SHA256SUM=$CNI_PLUGINS_ARM64_SHA256SUM
EOF
