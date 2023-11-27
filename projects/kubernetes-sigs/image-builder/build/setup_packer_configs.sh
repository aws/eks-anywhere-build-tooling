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

RELEASE_BRANCH="${1?Specify first argument - release branch}"
IMAGE_FORMAT="${2?Specify second argument - image format}"
IMAGE_OS="${3?Specify third argument - image OS}"
ARTIFACTS_BUCKET="${4?Specify fourth argument - artifact bucket}"
OVA_PATH="${5? Specify fifth argument - ova output path}"
ADDITIONAL_PAUSE_IMAGE_FROM="${6? Specify sixth argument - additional pause image}"
LATEST_TAG="${7? Specify seventh argument - latest tag}"
IMAGE_BUILDER_DIR="${8? Specify eighth argument - image-builder directory}"
IMAGE_OS_DIR="${9? Specify ninth argument - image-os-dir}"
IMAGE_OS_VERSION="${10? Specify tenth argument - image-os-version}"

CI="${CI:-false}"
CODEBUILD_CI="${CODEBUILD_CI:-false}"
EKSA_USE_DEV_RELEASE="${EKSA_USE_DEV_RELEASE:-false}"
DEV_RELEASE=false
if [[ $CI == true || $CODEBUILD_CI == true || $EKSA_USE_DEV_RELEASE == true ]]; then
  DEV_RELEASE=true
fi

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"
source "${MAKE_ROOT}/../../../build/lib/eksa_releases.sh"

function env_and_envsubst()
{
    echo -e "\nRunning envsubst with input file $2 saving to $3 with ENV: "
    local -r env_vars=(${1//:/ })
    for var in "${env_vars[@]}"; do
        local var_name=${var:1}
        echo "$var_name: ${!var_name:-}"
    done
    envsubst "$1" < "$2" > "$3"
}

# Preload release yaml
echo "Loading EKSD manifest from: $(build::eksd_releases::get_release_yaml_url $RELEASE_BRANCH)"
build::eksd_releases::load_release_yaml $RELEASE_BRANCH > /dev/null

echo "Loading EKSA manifest from: $(build::eksa_releases::get_eksa_release_manifest_url $DEV_RELEASE $LATEST_TAG)"
build::eksa_releases::load_bundle_manifest $DEV_RELEASE $LATEST_TAG > /dev/null

OUTPUT_CONFIGS="$MAKE_ROOT/_output/$RELEASE_BRANCH/$IMAGE_FORMAT/$IMAGE_OS_DIR/config"
mkdir -p $OUTPUT_CONFIGS

export CNI_SHA="sha256:$(build::eksd_releases::get_eksd_component_sha 'cni-plugins' $RELEASE_BRANCH)"
export PLUGINS_ASSET_BASE_URL=$(build::eksd_releases::get_eksd_cni_asset_base_url $RELEASE_BRANCH)
export CNI_VERSION=$(build::eksd_releases::get_eksd_component_version 'cni-plugins' $RELEASE_BRANCH)

# Get CNI sha256 to validate cni plugins in image are from eks-d
# Use sha from manfiest to validate tar and then generate sha from specific binary
echo "Downloading cni plugins to get shasum of host-device plugin for validation during image building"
TMP_CNI="$HOME/tmp/eks-image-builder-cni"
mkdir -p $TMP_CNI
curl -o $TMP_CNI/cni-plugins.tar.gz "$PLUGINS_ASSET_BASE_URL/$CNI_VERSION/cni-plugins-linux-amd64-$CNI_VERSION.tar.gz"
echo "$(echo $CNI_SHA | sed -E 's/.*sha256:(.*)$/\1/') $TMP_CNI/cni-plugins.tar.gz" > $TMP_CNI/cni.sha256
sha256sum -c $TMP_CNI/cni.sha256
tar -zxvf $TMP_CNI/cni-plugins.tar.gz -C $TMP_CNI ./host-device 
export CNI_HOST_DEVICE_SHA256="$(sha256sum $TMP_CNI/host-device | awk -F ' ' '{print $1}')"
rm -rf $TMP_CNI

env_and_envsubst '$CNI_SHA:$PLUGINS_ASSET_BASE_URL:$CNI_VERSION:$CNI_HOST_DEVICE_SHA256' \
    "$MAKE_ROOT/packer/config/cni.json.tmpl" \
    "$OUTPUT_CONFIGS/cni.json"

export PAUSE_IMAGE=$(build::eksd_releases::get_eksd_kubernetes_image_url 'pause-image' $RELEASE_BRANCH)
env_and_envsubst '$PAUSE_IMAGE:$HTTP_PROXY:$HTTPS_PROXY:$NO_PROXY' \
    "$MAKE_ROOT/packer/config/common.json.tmpl" \
    "$OUTPUT_CONFIGS/common.json"

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
export ETCDADM_HTTP_SOURCE=${ETCDADM_HTTP_SOURCE:-$(build::eksa_releases::get_eksa_component_asset_url 'eksD' 'etcdadm' $RELEASE_BRANCH $DEV_RELEASE $LATEST_TAG)}
export ETCDADM_VERSION='v0.1.5'
export CRICTL_URL=${CRICTL_URL:-$(build::eksa_releases::get_eksa_component_asset_url 'eksD' 'crictl' $RELEASE_BRANCH $DEV_RELEASE $LATEST_TAG)}
export CRICTL_SHA256="${CRICTL_SHA256:-$(build::eksa_releases::get_eksa_component_asset_artifact_checksum 'eksD' 'crictl' 'sha256' $RELEASE_BRANCH $DEV_RELEASE $LATEST_TAG)}"
export CONTAINERD_URL=${CONTAINERD_URL:-$(build::eksa_releases::get_eksa_component_asset_url 'eksD' 'containerd' $RELEASE_BRANCH $DEV_RELEASE $LATEST_TAG)}
export CONTAINERD_SHA256="${CONTAINERD_SHA256:-$(build::eksa_releases::get_eksa_component_asset_artifact_checksum 'eksD' 'containerd' 'sha256' $RELEASE_BRANCH $DEV_RELEASE $LATEST_TAG)}"

EMPTY_URL=false
URL_RETURNED_NOT_OK_STATUS=false
EKSA_BINARY_URLS=(CRICTL_URL CONTAINERD_URL ETCDADM_HTTP_SOURCE)
for url in ${EKSA_BINARY_URLS[@]}; do
    if [ -z "${!url}" ]; then
        EMPTY_URL=true
        echo "$url is empty"
    else
        http_code=$(curl -I -L -s -o /dev/null -w "%{http_code}" ${!url})
        if [[ "$http_code" != "200" ]]; then 
            URL_RETURNED_NOT_OK_STATUS=true
            echo "$url returned non-200 HTTP status"
        fi
    fi
done
if [[ $EMPTY_URL == "true" || $URL_RETURNED_NOT_OK_STATUS == "true" ]]; then
    exit 1
fi

env_and_envsubst '$IMAGE_REPO:$KUBERNETES_ASSET_BASE_URL:$KUBERNETES_VERSION:$KUBERNETES_SERIES:$CRICTL_URL:$CRICTL_SHA256:$CONTAINERD_URL:$CONTAINERD_SHA256:$ETCD_HTTP_SOURCE:$ETCD_VERSION:$ETCDADM_HTTP_SOURCE:$ETCD_SHA256:$ETCDADM_VERSION:$KUBERNETES_FULL_VERSION' \
    "$MAKE_ROOT/packer/config/kubernetes.json.tmpl" \
    "$OUTPUT_CONFIGS/kubernetes.json"

ADDITIONAL_PAUSE_IMAGE_VERSION_BASE_URL=$(build::eksd_releases::get_eksd_kubernetes_asset_base_url $ADDITIONAL_PAUSE_IMAGE_FROM)
ADDITIONAL_PAUSE_KUBERNETES_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $ADDITIONAL_PAUSE_IMAGE_FROM)
export ADDITIONAL_PAUSE_IMAGE=$ADDITIONAL_PAUSE_IMAGE_VERSION_BASE_URL/$ADDITIONAL_PAUSE_KUBERNETES_VERSION/bin/linux/amd64/pause.tar
env_and_envsubst '$ADDITIONAL_PAUSE_IMAGE' \
    "$MAKE_ROOT/packer/config/additional_components.json.tmpl" \
    "$OUTPUT_CONFIGS/additional_components.json"

# Write kubernetes version and the eksd manifest url consumed to output directory
mkdir -p $OVA_PATH
echo "$KUBERNETES_VERSION" > "$OVA_PATH"/KUBERNETES_VERSION
export EKSD_MANIFEST_URL=$(build::eksd_releases::get_release_yaml_url $RELEASE_BRANCH)
echo "$EKSD_MANIFEST_URL" > "$OVA_PATH"/EKSD_MANIFEST_URL

env_and_envsubst '$CNI_VERSION:$ETCD_VERSION:$ETCD_SHA256:$ETCDADM_VERSION:$PAUSE_IMAGE:$CNI_HOST_DEVICE_SHA256' \
    "$MAKE_ROOT/packer/config/validate_goss_inline_vars.json.tmpl" \
    "$OUTPUT_CONFIGS/validate_goss_inline_vars.json"

env_and_envsubst '$EKSD_NAME' \
    "$MAKE_ROOT/packer/config/ovf_custom_properties.json.tmpl" \
    "$OUTPUT_CONFIGS/ovf_custom_properties.json"

# This is the IP address that Packer will create the server on to host the local
# directory containing the kickstart config
if [ "$IMAGE_FORMAT" = "ova" ] && \
    ( [ "$IMAGE_OS" = "redhat" ] || ( [ "$IMAGE_OS" = "ubuntu" ] && [ "$IMAGE_OS_VERSION" != "2004" ] )  ); then
    
    # Get Packer HTTP server IP from env var if set
    PACKER_HTTP_SERVER_IP=${PACKER_HTTP_SERVER_IP:-}
    if [ -z $PACKER_HTTP_SERVER_IP ]; then
        # Get active network interface from env var if set
        PACKER_ACTIVE_INTERFACE="${PACKER_ACTIVE_INTERFACE:-}"
        if [ -z $PACKER_ACTIVE_INTERFACE ]; then
            echo "Finding interface for Packer temporary HTTP server"
            if [ "$(uname -s)" = "Linux" ]; then
                if command -v route &> /dev/null; then
                    PACKER_ACTIVE_INTERFACE=$(route |awk '/^default/{print $NF}')
                else
                    INTERFACES=($(ls /sys/class/net))
                    for interface in "${INTERFACES[@]}"; do
                        if [ "$interface" = "eth0" ] || [ "$interface" = "en0" ] || [ "$interface" = "eno1" ]; then
                            PACKER_ACTIVE_INTERFACE=$interface
                            echo "Found interface: $interface"
                            break
                        fi
                    done
                fi
            elif [ "$(uname -s)" = "Darwin" ]; then
                PACKER_ACTIVE_INTERFACE=$(route -n get default | grep 'interface:' | grep -o '[^ ]*$')
            fi
        fi

        if [ -z $PACKER_ACTIVE_INTERFACE ]; then
            echo "PACKER_ACTIVE_INTERFACE cannot be automatically determined. Please export PACKER_ACTIVE_INTERFACE=<Current Host's primary network interface>"
            exit 1
        fi

        if [ "$(uname -s)" = "Linux" ]; then
            PACKER_HTTP_SERVER_IP=$(ip a l $PACKER_ACTIVE_INTERFACE | awk '/inet / {print $2}' | cut -d/ -f1)
        elif
            [ "$(uname -s)" = "Darwin" ]; then
            PACKER_HTTP_SERVER_IP=$(ifconfig $PACKER_ACTIVE_INTERFACE | awk '/inet / {print $2}' | cut -d/ -f1)
        fi
    fi

    if [ -z $PACKER_HTTP_SERVER_IP ]; then
        echo "PACKER_HTTP_SERVER_IP cannot be automatically determined.  Please export PACKER_HTTP_SERVER_IP=<Current Host's IP>"
        exit 1
    fi

    SED=sed
    if [ "$(uname -s)" = "Darwin" ]; then
      SED=gsed
    fi
    find ${MAKE_ROOT}/${IMAGE_BUILDER_DIR}/packer/ova -type f -name '*.json' | xargs $SED -i "s/{{ .HTTPIP }}/$PACKER_HTTP_SERVER_IP/g"
fi

# If the image format is AMI and our Packer config specifies a non-gp3 volume, then
# remove the throughput field in the upstream Packer config since only gp3 supports
# specifying throughput
if [ "$IMAGE_FORMAT" = "ami" ]; then
    volume_type="$(cat $MAKE_ROOT/packer/ami/ami.json | jq -r '.volume_type')"
    if [[ $volume_type != "null" && $volume_type != "gp3" ]]; then
        build::jq::update_in_place ${MAKE_ROOT}/${IMAGE_BUILDER_DIR}/packer/ami/packer.json "del(.builders[0].launch_block_device_mappings[].throughput)"
    else
        build::jq::update_in_place ${MAKE_ROOT}/${IMAGE_BUILDER_DIR}/packer/ami/packer.json '.builders[0].launch_block_device_mappings[].throughput = "{{ user `throughput` }}"'
    fi
fi

EKSA_ANSIBLE_VERBOSITY="${EKSA_ANSIBLE_VERBOSITY:-}"
if [ -n "$EKSA_ANSIBLE_VERBOSITY" ]; then
    PACKER_FILE_NAME="packer.json"
    if [ "$IMAGE_FORMAT" = "ova" ]; then
        PACKER_FILE_NAME="packer-node.json"
    fi

    HYPERVISOR_DIRNAME=$IMAGE_FORMAT
    if [ "$IMAGE_FORMAT" = "cloudstack" ]; then
        HYPERVISOR_DIRNAME="qemu"
    fi

    build::jq::update_in_place ${MAKE_ROOT}/${IMAGE_BUILDER_DIR}/packer/$HYPERVISOR_DIRNAME/$PACKER_FILE_NAME "(.provisioners[] | select(.type == \"ansible\").extra_arguments) += [\"$EKSA_ANSIBLE_VERBOSITY\"]"
fi