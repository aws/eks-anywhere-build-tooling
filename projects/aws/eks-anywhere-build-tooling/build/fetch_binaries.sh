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

ARTIFACTS_BUCKET="$1"

declare -A repos=(["fluxcd/flux2"]="true"
                  ["kubernetes-sigs/cluster-api"]="true"
                  ["kubernetes-sigs/cluster-api-provider-aws"]="true"
                  ["kubernetes-sigs/kind"]="true"
                  ["kubernetes"]="false"
                  ["replicatedhq/troubleshoot"]="true"
                  ["vmware/govmomi"]="true")

declare -A project_bin_licenses=(["flux2"]="flux FLUX2"
                                 ["cluster-api"]="clusterctl CAPI"
                                 ["cluster-api-provider-aws"]="clusterawsadm CAPA"
                                 ["kind"]="kind KIND"
                                 ["kubernetes/kubernetes"]="client/bin/kubectl KUBERNETES"
                                 ["troubleshoot"]="support-bundle TROUBLESHOOT"
                                 ["govmomi"]="govc GOVMOMI")

function unpack::tarballs(){
  local -r output_dir="$1"
  local -r arch="$2"

  mkdir -p $output_dir
  for repo in "${!repos[@]}"
  do
    project=$repo
    base=$(basename $project)
    mkdir $output_dir/$base
    private=${repos[$repo]}
    if [ $private = "true" ]; then
      URL=$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET $project $arch)
    else      
      URL=$(build::eksd_releases::get_eksd_kubernetes_asset_url "kubernetes-client-linux-$arch.tar.gz" $(build::eksd_releases::get_release_branch) $arch)
    fi
    curl -sSL "${URL}" -o $output_dir/tmp.tar.gz
    tar xzf $output_dir/tmp.tar.gz -C $output_dir/$base
  done
}

function copy::binaries::licenses(){
  local -r output_dir="$1"
  local -r eks_a_tool_binary_dir="$2"
  local -r eks_a_tool_license_dir="$3"

  mkdir -p $eks_a_tool_binary_dir
  mkdir -p $eks_a_tool_license_dir
  for project in "${!project_bin_licenses[@]}"
  do
    bin_license_map=(${project_bin_licenses[$project]})
    binary=${bin_license_map[0]}
    license_prefix=${bin_license_map[1]}
    cp $output_dir/$project/$binary $eks_a_tool_binary_dir/$(basename $binary)
    cp $output_dir/$project/ATTRIBUTION.txt $eks_a_tool_license_dir/${license_prefix}_ATTRIBUTION.txt
    cp -r $output_dir/$project/LICENSES $eks_a_tool_license_dir/${license_prefix}_LICENSES
  done
}
for arch in amd64 arm64; do
  OUTPUT_DIR="${MAKE_ROOT}/_output/linux-$arch"
  EKS_A_TOOL_BINARY_DIR="${OUTPUT_DIR}/eks-a-tools/binary"
  EKS_A_TOOL_LICENSE_DIR="${OUTPUT_DIR}/eks-a-tools/licenses"
  
  unpack::tarballs $OUTPUT_DIR $arch
  copy::binaries::licenses $OUTPUT_DIR $EKS_A_TOOL_BINARY_DIR $EKS_A_TOOL_LICENSE_DIR
done
