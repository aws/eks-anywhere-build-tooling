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
GOLANG_VERSION="$7"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

build::common::use_go_version $GOLANG_VERSION

cd $REPO

yq eval -i -P ".spec.template.spec.containers[0].args += [\"--namespace=eksa-system\"]" config/manager/manager.yaml

make manifest-modification STAGE="release" \
  REGISTRY="${IMAGE_REPO}" \
  IMAGE_NAME="kubernetes-sigs/cluster-api-provider-vsphere/release/manager" \
  RELEASE_TAG="${IMAGE_TAG}" \
  PULL_POLICY="IfNotPresent"

make release-manifests STAGE="release" \
  MANIFEST_DIR="out"

mkdir -p $OUTPUT_DIR/manifests/infrastructure-vsphere/$TAG
cp out/cluster-template.yaml "$OUTPUT_DIR/manifests/infrastructure-vsphere/$TAG"
cp out/infrastructure-components.yaml "$OUTPUT_DIR/manifests/infrastructure-vsphere/$TAG"
cp metadata.yaml "$OUTPUT_DIR/manifests/infrastructure-vsphere/$TAG"

cp -rf $OUTPUT_DIR/manifests $ARTIFACTS_PATH
