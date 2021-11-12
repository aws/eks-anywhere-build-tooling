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

RELEASE_BRANCH="$1"
ARTIFACTS_BUCKET="$2"

for arch in amd64 arm64; do
    OUTPUT_DIR="_output/linux-$arch/kubernetes/$RELEASE_BRANCH"
    mkdir -p $OUTPUT_DIR

    TARBALLS=(
        "kubernetes-client-linux-$arch.tar.gz"
        "kubernetes-server-linux-$arch.tar.gz"
    )
    for TARBALL in "${TARBALLS[@]}"; do
        URL=$(build::eksd_releases::get_eksd_kubernetes_asset_url $TARBALL $RELEASE_BRANCH $arch)
        curl -sSL $URL -o $OUTPUT_DIR/$TARBALL
        tar xzf $OUTPUT_DIR/$TARBALL -C $OUTPUT_DIR
    done

    OUTPUT_DIR="_output/linux-$arch/etcdadm"
    mkdir -p $OUTPUT_DIR
    TARBALL="etcdadm-linux-$arch.tar.gz"
    URL=$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/etcdadm' $arch)
    curl -sSL $URL -o $OUTPUT_DIR/$TARBALL
    tar xzf $OUTPUT_DIR/$TARBALL -C $OUTPUT_DIR
done
