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

RELEASE_BRANCH="${1?Specify first argument - release branch}"
OUTPUT_DIR="_output/$RELEASE_BRANCH/kubernetes"
mkdir -p $OUTPUT_DIR
TARBALLS=(
    "kubernetes-client-linux-amd64.tar.gz"
    "kubernetes-server-linux-amd64.tar.gz"
)
for TARBALL in "${TARBALLS[@]}"; do
    URL=$(build::eksd_releases::get_eksd_kubernetes_asset_url $TARBALL $RELEASE_BRANCH)
    curl -sSL $URL -o $OUTPUT_DIR/$TARBALL
    tar xzf $OUTPUT_DIR/$TARBALL -C $OUTPUT_DIR
done

OUTPUT_DIR="_output/$RELEASE_BRANCH/etcdadm"
mkdir -p $OUTPUT_DIR
TARBALL="etcdadm-linux-amd64.tar.gz"
URL=$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/etcdadm')
curl -sSL $URL -o $OUTPUT_DIR/$TARBALL
tar xzf $OUTPUT_DIR/$TARBALL -C $OUTPUT_DIR