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

GIT_TAG="${1?Specify first argument - Kind Git tag}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
SED=$(source $BUILD_LIB/common.sh && build::find::gnu_variant_on_mac sed)

declare -A image_name_source_of_truth=(
    ["docker.io/kindest/kindnetd"]="pkg/build/nodeimage/const_cni.go"
    ["docker.io/kindest/local-path-provisioner"]="pkg/build/nodeimage/const_storage.go"
    ["docker.io/kindest/local-path-helper"]="pkg/build/nodeimage/const_storage.go"
    ["registry.k8s.io/pause"]="images/base/files/etc/containerd/config.toml"
)

for image in "${!image_name_source_of_truth[@]}"; do
    upstream_image_match=$(grep -on "\"$image:.*\"" $MAKE_ROOT/kind/${image_name_source_of_truth[$image]})
    upstream_image_match_line=$(echo $upstream_image_match | cut -d: -f1)
    upstream_image_uri=$(echo $upstream_image_match | cut -d: -f2-)
    echo $upstream_image_match_line
    echo $upstream_image_uri

    $SED -i "s,\"$image:.*\",$upstream_image_uri," $MAKE_ROOT/build/node-image-build-args.sh
    $SED -i "s,https://github.com/kubernetes-sigs/kind/blob/.*/${image_name_source_of_truth[$image]}.*,https://github.com/kubernetes-sigs/kind/blob/$GIT_TAG/${image_name_source_of_truth[$image]}#L${upstream_image_match_line}," $MAKE_ROOT/build/node-image-build-args.sh
done
