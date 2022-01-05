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

REPO_OWNER="${1?First argument is repository owner}"
REPO="${2?Second argument is repository}"
OUTPUT_DIR="${3?Third argument is output directory}"
PACKAGE_FILTER="${4?Fourth argument is package filter}"
REPO_SUBPATH="${5?Fifth argument is repository subpath}"
GO_LICENSES="${6?Sixth argument is go license flag}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

if [ "${GO_LICENSES}" == "false" ]
then
  build::non-golang::copy_licenses ${REPO}/${REPO_SUBPATH} $OUTPUT_DIR/LICENSES/github.com/${REPO_OWNER}/${REPO}
else
  cd $REPO/$REPO_SUBPATH
  build::gather_licenses $OUTPUT_DIR "$PACKAGE_FILTER"
fi
