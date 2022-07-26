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

REPO="$1"
OUTPUT_DIR="$2"
ARTIFACTS_PATH="$3"
TAG="$4"
IMAGE_REPO="$5"
IMAGE_TAG="$6"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

cd $REPO

CAPI_REGISTRY_PREFIX="${IMAGE_REPO}/kubernetes-sigs/cluster-api"
KUBE_RBAC_PROXY_LATEST_TAG=latest

make set-manifest-image \
    MANIFEST_IMG=$CAPI_REGISTRY_PREFIX/cluster-api-controller MANIFEST_TAG=$IMAGE_TAG \
    TARGET_RESOURCE="./config/default/manager_image_patch.yaml"

# Set the kubeadm bootstrap image to the production bucket.
make set-manifest-image \
    MANIFEST_IMG=$CAPI_REGISTRY_PREFIX/kubeadm-bootstrap-controller MANIFEST_TAG=$IMAGE_TAG \
    TARGET_RESOURCE="./bootstrap/kubeadm/config/default/manager_image_patch.yaml"

# Set the kubeadm control plane image to the production bucket.
make set-manifest-image \
    MANIFEST_IMG=$CAPI_REGISTRY_PREFIX/kubeadm-control-plane-controller MANIFEST_TAG=$IMAGE_TAG  \
    TARGET_RESOURCE="./controlplane/kubeadm/config/default/manager_image_patch.yaml"

make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./config/default/manager_pull_policy.yaml"
make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./bootstrap/kubeadm/config/default/manager_pull_policy.yaml"
make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./controlplane/kubeadm/config/default/manager_pull_policy.yaml"

yq eval -i -P ".spec.template.spec.containers[0].args += [\"--namespace=eksa-system\"]" config/manager/manager.yaml
yq eval -i -P ".spec.template.spec.containers[0].args += [\"--namespace=eksa-system\"]" bootstrap/kubeadm/config/manager/manager.yaml
yq eval -i -P ".spec.template.spec.containers[0].args += [\"--namespace=eksa-system\"]" controlplane/kubeadm/config/manager/manager.yaml

## Build the manifests
make release-manifests

## Build the development manifests
make set-manifest-image \
    MANIFEST_IMG=$CAPI_REGISTRY_PREFIX/capd-manager MANIFEST_TAG=$IMAGE_TAG \
    TARGET_RESOURCE="./test/infrastructure/docker/config/default/manager_image_patch.yaml"
make set-manifest-pull-policy PULL_POLICY=IfNotPresent TARGET_RESOURCE="./test/infrastructure/docker/config/default/manager_pull_policy.yaml"
make release-manifests-dev

mkdir -p $OUTPUT_DIR/manifests/{bootstrap-kubeadm,cluster-api,control-plane-kubeadm,infrastructure-docker}/$TAG
cp out/bootstrap-components.yaml "$OUTPUT_DIR/manifests/bootstrap-kubeadm/$TAG"
cp out/metadata.yaml "$OUTPUT_DIR/manifests/bootstrap-kubeadm/$TAG"

cp out/control-plane-components.yaml "$OUTPUT_DIR/manifests/control-plane-kubeadm/$TAG"
cp out/metadata.yaml "$OUTPUT_DIR/manifests/control-plane-kubeadm/$TAG"

cp out/core-components.yaml "$OUTPUT_DIR/manifests/cluster-api/$TAG"
cp out/metadata.yaml "$OUTPUT_DIR/manifests/cluster-api/$TAG"

cp out/infrastructure-components-development.yaml "$OUTPUT_DIR/manifests/infrastructure-docker/$TAG/infrastructure-components-development.yaml"
cp test/infrastructure/docker/templates/cluster-template-development.yaml "$OUTPUT_DIR/manifests/infrastructure-docker/$TAG"
cp out/metadata.yaml "$OUTPUT_DIR/manifests/infrastructure-docker/$TAG"

cp -rf $OUTPUT_DIR/manifests $ARTIFACTS_PATH
