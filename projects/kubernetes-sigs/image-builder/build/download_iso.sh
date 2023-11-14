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

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${SCRIPT_ROOT}/build/lib/common.sh"

ISO_FILENAME="$1"
ISO_CHECKSUM="$2"
ISO_URL="$3"

echo "$ISO_CHECKSUM  /tmp/$ISO_FILENAME" > /tmp/$ISO_FILENAME.sha256

function download_and_validate(){    
    if ! build::common::echo_and_run curl -sSL --retry 5 $ISO_URL -o /tmp/$ISO_FILENAME; then
        return 1
    fi
    
    build::common::echo_and_run sha256sum -c /tmp/$ISO_FILENAME.sha256
}

retry download_and_validate
