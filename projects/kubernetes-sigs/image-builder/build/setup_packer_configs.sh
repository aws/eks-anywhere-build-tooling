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

RELEASE_BRANCH="${1?Specify first argument - release branch}"
IMAGE_FORMAT="${2?Specify second argument - image format}"
ARTIFACTS_BUCKET="${3?Specify third argument - artifact bucket}"
OVA_PATH="${4? Specify fourth argument - ova output path}"
ADDITIONAL_PAUSE_IMAGE_FROM="${5? Specify fifth argument - additional pause image}"
LATEST_TAG="${6? Specify sixth argument - latest tag}"

CI="${CI:-false}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

# Preload release yaml
build::eksd_releases::load_release_yaml $RELEASE_BRANCH

OUTPUT_CONFIGS="$MAKE_ROOT/_output/$RELEASE_BRANCH/$IMAGE_FORMAT/config"
mkdir -p $OUTPUT_CONFIGS

export CNI_SHA="sha256:$(build::eksd_releases::get_eksd_component_sha 'cni-plugins' $RELEASE_BRANCH)"
export PLUGINS_ASSET_BASE_URL=$(build::eksd_releases::get_eksd_cni_asset_base_url $RELEASE_BRANCH)
export CNI_VERSION=$(build::eksd_releases::get_eksd_component_version 'cni-plugins' $RELEASE_BRANCH)

# Get CNI sha256 to validate cni plugins in image are from eks-d
# Use sha from manfiest to validate tar and then generate sha from specific binary
TMP_CNI="/tmp/eks-image-builder-cni"
mkdir -p $TMP_CNI
curl -o $TMP_CNI/cni-plugins.tar.gz "$PLUGINS_ASSET_BASE_URL/$CNI_VERSION/cni-plugins-linux-amd64-$CNI_VERSION.tar.gz"
echo "$(echo $CNI_SHA | sed -E 's/.*sha256:(.*)$/\1/') $TMP_CNI/cni-plugins.tar.gz" > $TMP_CNI/cni.sha256
sha256sum -c $TMP_CNI/cni.sha256
tar -zxvf $TMP_CNI/cni-plugins.tar.gz -C $TMP_CNI ./host-device 
export CNI_HOST_DEVICE_SHA256="$(sha256sum /tmp/cni/host-device | awk -F ' ' '{print $1}')"
rm -rf $TMP_CNI

envsubst '$CNI_SHA:$PLUGINS_ASSET_BASE_URL:$CNI_VERSION:$CNI_HOST_DEVICE_SHA256' \
    < "$MAKE_ROOT/packer/config/cni.json.tmpl" \
    > "$OUTPUT_CONFIGS/cni.json"

export PAUSE_IMAGE=$(build::eksd_releases::get_eksd_kubernetes_image_url 'pause-image' $RELEASE_BRANCH)
envsubst '$PAUSE_IMAGE' \
    < "$MAKE_ROOT/packer/config/common.json.tmpl" \
    > "$OUTPUT_CONFIGS/common.json"


export IMAGE_REPO=$(build::eksd_releases::get_eksd_image_repo $RELEASE_BRANCH)
export KUBERNETES_ASSET_BASE_URL=$(build::eksd_releases::get_eksd_kubernetes_asset_base_url $RELEASE_BRANCH)
export KUBERNETES_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $RELEASE_BRANCH)
export KUBERNETES_SERIES="v${RELEASE_BRANCH/-/.}"
export EKSD_NAME=$(build::eksd_releases::get_eksd_release_name $RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $RELEASE_BRANCH)
export KUBERNETES_FULL_VERSION="$KUBERNETES_VERSION-eks-$RELEASE_BRANCH-$EKSD_RELEASE"
export ETCD_HTTP_SOURCE=$(build::eksd_releases::get_eksd_component_url "etcd" $RELEASE_BRANCH)
export ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $RELEASE_BRANCH)
export ETCD_SHA256=$(build::eksd_releases::get_eksd_component_sha "etcd" $RELEASE_BRANCH)
export ETCDADM_HTTP_SOURCE=${ETCDADM_HTTP_SOURCE:-$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/etcdadm' 'amd64' $LATEST_TAG)}
export ETCDADM_VERSION='v0.1.5'
export CRICTL_URL=${CRICTL_URL:-$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/cri-tools' 'amd64' $LATEST_TAG)}
export CRICTL_SHA256="$CRICTL_URL.sha256"

envsubst '$IMAGE_REPO:$KUBERNETES_ASSET_BASE_URL:$KUBERNETES_VERSION:$KUBERNETES_SERIES:$CRICTL_URL:$CRICTL_SHA256:$ETCD_HTTP_SOURCE:$ETCD_VERSION:$ETCDADM_HTTP_SOURCE:$ETCD_SHA256:$ETCDADM_VERSION:$KUBERNETES_FULL_VERSION' \
    < "$MAKE_ROOT/packer/config/kubernetes.json.tmpl" \
    > "$OUTPUT_CONFIGS/kubernetes.json"

ADDITIONAL_PAUSE_IMAGE_VERSION_BASE_URL=$(build::eksd_releases::get_eksd_kubernetes_asset_base_url $ADDITIONAL_PAUSE_IMAGE_FROM)
ADDITIONAL_PAUSE_KUBERNETES_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $ADDITIONAL_PAUSE_IMAGE_FROM)
export ADDITIONAL_PAUSE_IMAGE=$ADDITIONAL_PAUSE_IMAGE_VERSION_BASE_URL/$ADDITIONAL_PAUSE_KUBERNETES_VERSION/bin/linux/amd64/pause.tar
envsubst '$ADDITIONAL_PAUSE_IMAGE' \
    < "$MAKE_ROOT/packer/config/additional_components.json.tmpl" \
    > "$OUTPUT_CONFIGS/additional_components.json"

# Write kubernetes version and the eksd manifest url consumed to output directory
mkdir -p $OVA_PATH
echo "$KUBERNETES_VERSION" > "$OVA_PATH"/KUBERNETES_VERSION
export EKSD_MANIFEST_URL=$(build::eksd_releases::get_release_yaml_url $RELEASE_BRANCH)
echo "$EKSD_MANIFEST_URL" > "$OVA_PATH"/EKSD_MANIFEST_URL

envsubst '$CNI_VERSION:$ETCD_VERSION:$ETCD_SHA256:$ETCDADM_VERSION:$PAUSE_IMAGE:$CNI_HOST_DEVICE_SHA256' \
    < "$MAKE_ROOT/packer/config/validate_goss_inline_vars.json.tmpl" \
    > "$OUTPUT_CONFIGS/validate_goss_inline_vars.json"

envsubst '$EKSD_NAME' \
    < "$MAKE_ROOT/packer/config/ovf_custom_properties.json.tmpl" \
    > "$OUTPUT_CONFIGS/ovf_custom_properties.json"
