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

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

EKSD_RELEASE_BRANCH="${1?Specify first argument - release branch}"
ARCH="${2?Specify second argument - Targetarch}"

DEP_FOLDER="$MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies/linux-$ARCH"
OUTPUT_FOLDER="$DEP_FOLDER/files/rootfs"
mkdir -p $OUTPUT_FOLDER/LICENSES

for FOLDER in eksd/kubernetes eksa/kubernetes-sigs/etcdadm eksd/cni-plugins eksa/kubernetes-sigs/cri-tools; do
    cp -rf $DEP_FOLDER/$FOLDER/LICENSES "$OUTPUT_FOLDER/LICENSES/$(echo $(basename $FOLDER) | tr a-z A-Z  | tr -d '-'  )_LICENSES"
    cp $DEP_FOLDER/$FOLDER/ATTRIBUTION.txt "$OUTPUT_FOLDER/LICENSES/$(echo $(basename $FOLDER) | tr a-z A-Z  | tr -d '-'  )_ATTRIBUTION.txt"
done

mkdir -p $OUTPUT_FOLDER/usr/bin/ $OUTPUT_FOLDER/usr/local/bin/
cp $DEP_FOLDER/eksa/kubernetes-sigs/etcdadm/etcdadm $OUTPUT_FOLDER/usr/bin/
cp $DEP_FOLDER/eksa/kubernetes-sigs/cri-tools/{crictl,critest} $OUTPUT_FOLDER/usr/local/bin/

# Place etcd tarball etcdadm cache directory to avoid downloading at runtime   
ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH $ARCH)
FOLDER="$OUTPUT_FOLDER/var/cache/etcdadm/etcd/$ETCD_VERSION"
mkdir -p $FOLDER
cp $DEP_FOLDER/eksd/etcd/etcd.tar.gz $FOLDER/etcd-$ETCD_VERSION-linux-$ARCH.tar.gz
