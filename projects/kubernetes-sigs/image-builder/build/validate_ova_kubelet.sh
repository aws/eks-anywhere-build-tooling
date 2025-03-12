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

OVA_PATH="$1"
RELEASE_BRANCH="$2"
KUBERNETES_PACKER_CONFIG="$3"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

if [ ! -f "$OVA_PATH" ]; then
    echo "ERROR: $OVA_PATH does not exist!"
    exit 1
fi

TMP_FOLDER="$MAKE_ROOT/_output/tmp"
BIN_FOLDER="$MAKE_ROOT/_output/.bin"
ZIP="$BIN_FOLDER/7zz"
TAR=$(build::find::gnu_variant_on_mac tar)

mkdir -p $BIN_FOLDER $TMP_FOLDER

if [ ! -f "$ZIP" ]; then
    build::common::echo_and_run curl -L https://www.7-zip.org/a/7z2406-$([ "Darwin" = "$(uname -s)" ] && echo mac || echo linux)$([ "Linux" = "$(uname -s)" ] && ([ "x86_64" = "$(uname -m)" ] && echo -x64 || echo -arm64)).tar.xz  | $TAR -xJ -C $BIN_FOLDER 7zz
fi

VMDK="$($TAR --wildcards -tf $OVA_PATH '*.vmdk')"
build::common::echo_and_run $TAR -C $TMP_FOLDER -xf $OVA_PATH $VMDK

if $ZIP l $TMP_FOLDER/$VMDK | grep "usr/bin/kubelet"; then
    build::common::echo_and_run $ZIP -y -o$TMP_FOLDER e $TMP_FOLDER/$VMDK usr/bin/kubelet > /dev/null
    ACTUAL_VERSION="$($TMP_FOLDER/kubelet --version)"
else
    IMG=$($ZIP l $TMP_FOLDER/$VMDK | grep -oP "\d+[A-Za-z .]*\.img$" | sort | tail -n 1)
    if [ -n "$IMG" ]; then
        build::common::echo_and_run $ZIP -y -o$TMP_FOLDER e $TMP_FOLDER/$VMDK $IMG
        if [ -f "$TMP_FOLDER/$IMG" ]; then
            ACTUAL_VERSION=$(strings "$TMP_FOLDER/$IMG" | grep -P "^v\d+\.\d+\.\d+-eks-.*$" | head -1)
        fi
    fi
fi
if [ -z "$ACTUAL_VERSION" ]; then
    echo "ERROR: unable to determine Kubelet version on image!"
    exit 1
fi

EXPECTED_VERSION="$(jq -r '.kubernetes_semver' $KUBERNETES_PACKER_CONFIG)"
if [[ $ACTUAL_VERSION != *"$EXPECTED_VERSION"* ]]; then
    echo "ERROR: kubelet version unexpected!"
    echo "expected: $EXPECTED_VERSION, actual: $ACTUAL_VERSION"
    exit 1
fi

echo "kubelet version matches, expected: $EXPECTED_VERSION, actual: $ACTUAL_VERSION"

rm -rf $TMP_FOLDER
