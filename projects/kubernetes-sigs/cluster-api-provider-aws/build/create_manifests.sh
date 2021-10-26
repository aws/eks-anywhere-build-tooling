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

MANIFEST_IMAGE="public.ecr.aws/l0g8r8j6/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller:v0.6.4"
MANIFEST_IMAGE_OVERRIDE="${IMAGE_REPO}/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller:${IMAGE_TAG}"
KUBE_RBAC_PROXY_MANIFEST_IMAGE="gcr.io/kubebuilder/kube-rbac-proxy:v0.4.1"
KUBE_RBAC_PROXY_MANIFEST_IMAGE_OVERRIDE=${IMAGE_REPO}/brancz/kube-rbac-proxy:latest

mkdir -p $OUTPUT_DIR/manifests/infrastructure-aws/$TAG
sed -i "s,${MANIFEST_IMAGE},${MANIFEST_IMAGE_OVERRIDE}," ../manifests/infrastructure-components.yaml
sed -i "s,${KUBE_RBAC_PROXY_MANIFEST_IMAGE},${KUBE_RBAC_PROXY_MANIFEST_IMAGE_OVERRIDE}," ../manifests/infrastructure-components.yaml
cp ../manifests/infrastructure-components.yaml "$OUTPUT_DIR/manifests/infrastructure-aws/$TAG"
cp templates/cluster-template.yaml "$OUTPUT_DIR/manifests/infrastructure-aws/$TAG"
cp metadata.yaml "$OUTPUT_DIR/manifests/infrastructure-aws/$TAG"
