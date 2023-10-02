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

#set -x
set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

RETURN=0
SED=$(build::common::gnu_variant_on_mac sed)

for GIT_TAG_FILE in projects/*/*/GIT_TAG
do
    VERSION="$(cat $GIT_TAG_FILE | $SED "s,-,--,g")"
    README="$(dirname $GIT_TAG_FILE)/README.md"
    if [ ! -f $README ]
    then
        echo "Missing file $README"
        continue
    fi
    EXPECTED_VERSION="img.shields.io/badge/version-$VERSION"
    if grep -l "$EXPECTED_VERSION" ${README} >/dev/null
    then
        continue
    fi
    if ! grep "img.shields.io/badge/version" ${README} >/dev/null
    then
        echo "Did not find version $README"
        continue
    fi
    RETURN=-1
    ACTUAL_VERSION=$(grep "img.shields.io/badge/version" ${README} | $SED -e 's,.*img.shields.io/badge/version-,,' -e 's/-blue).*$//')
    echo "Version mismatch $README expected $VERSION actual $ACTUAL_VERSION"
    $SED -i -e "s/$ACTUAL_VERSION/$VERSION/" $README
done
