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
KUBERNETES_PACKER_CONFIG="$2"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

TMP_FOLDER="$MAKE_ROOT/_output/tmp"
BIN_FOLDER="$MAKE_ROOT/_output/.bin"
ZIP="$BIN_FOLDER/7zz"

mkdir -p $BIN_FOLDER $TMP_FOLDER

if [ ! -f "$ZIP" ]; then
    build::common::echo_and_run curl -L https://www.7-zip.org/a/7z2406-linux-$([ "x86_64" = "$(uname -m)" ] && echo x64 || echo arm64).tar.xz  | tar -xJ -C $BIN_FOLDER 7zz
fi

EXPECTED_VERSION="$(jq -r '.kubernetes_semver' $KUBERNETES_PACKER_CONFIG)"

FAKE_KUBELET=$TMP_FOLDER/fake-ova/vmdk/usr/bin/kubelet
mkdir -p $(dirname $FAKE_KUBELET)

cat <<EOF > $FAKE_KUBELET
#!/usr/bin/env bash
echo "Kuberentes $EXPECTED_VERSION"
EOF

chmod +x $FAKE_KUBELET

build::common::echo_and_run $ZIP a $TMP_FOLDER/fake-ova/disk-1.7z $TMP_FOLDER/fake-ova/vmdk/*
build::common::echo_and_run mv $TMP_FOLDER/fake-ova/disk-1.7z $TMP_FOLDER/fake-ova/disk-1.vmdk
rm -rf $TMP_FOLDER/fake-ova/vmdk
build::common::echo_and_run tar cf $OVA_PATH -C $TMP_FOLDER/fake-ova .

rm -rf $TMP_FOLDER
