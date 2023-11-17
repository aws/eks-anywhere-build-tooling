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

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

REGISTRY="${1?First argument is registry}"
REPOSITORY="${2?Second argument is repository}"
IMAGE_TAG="${3?Third argument is image tag}"

TMPFILE=$(mktemp)
trap "rm -f $TMPFILE" exit
TARGET=${REGISTRY}/${REPOSITORY}:${IMAGE_TAG}

>&2 echo -n "Checking for the existence of ${TARGET}..."
if skopeo inspect -n --raw docker://${TARGET} >${TMPFILE} 2>/dev/null; then
  >&2 echo "Found!"
  skopeo manifest-digest ${TMPFILE}
else
  >&2 echo "Not Found!"
fi

