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

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

OUTPUT_DIR="${MAKE_ROOT}/_output"
EKS_A_DIAGNOSTIC_COLLECTOR_BINARY_DIR="${OUTPUT_DIR}/eks-a-diagnostic-collector-tools/binary"
EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR="${OUTPUT_DIR}/eks-a-diagnostic-collector-tools/licenses"

function unpack::tarballs(){
  mkdir -p $OUTPUT_DIR
  project="kubernetes"
  base=$(basename $project)
  mkdir $OUTPUT_DIR/$base
  URL=$(build::eksd_releases::get_eksd_kubernetes_asset_url "kubernetes-client-linux-amd64.tar.gz")
  curl -sSL "${URL}" -o $OUTPUT_DIR/tmp.tar.gz
  tar xzf $OUTPUT_DIR/tmp.tar.gz -C $OUTPUT_DIR/$base
}

function copy::binaries::licenses(){
  mkdir -p $EKS_A_DIAGNOSTIC_COLLECTOR_BINARY_DIR
  mkdir -p $EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR
  project="kubernetes/kubernetes"
  binary="client/bin/kubectl"
  license_prefix="KUBERNETES"
  cp ./_output/$project/$binary $EKS_A_DIAGNOSTIC_COLLECTOR_BINARY_DIR/$(basename $binary)
  cp ./_output/$project/ATTRIBUTION.txt $EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR/${license_prefix}_ATTRIBUTION.txt
  cp -r ./_output/$project/LICENSES $EKS_A_DIAGNOSTIC_COLLECTOR_TOOL_LICENSE_DIR/${license_prefix}_LICENSES
}

unpack::tarballs
copy::binaries::licenses