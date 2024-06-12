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

OUTPUT_DIR="${1?First arguement is output directory}"
HELM_DESTINATION_REPOSITORY="${2?Second argument is helm repository}"
HELM_CHART_FOLDER="${3?Third argument is helm chart folder}"
BUILD_HELM_DEPENDENCIES="${4?Fourth argument is whether or not to build helm dependencies}"

CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})

if [ "${HELM_CHART_FOLDER}" != "." ]; then
    CHART_NAME=${CHART_NAME}/${HELM_CHART_FOLDER}
fi


#
# Build
#
cd ${OUTPUT_DIR}/helm
if [ "${BUILD_HELM_DEPENDENCIES}" = "true" ]; then
    build::common::echo_and_run helm dependency update "${CHART_NAME}"
    build::common::echo_and_run helm dependency build "${CHART_NAME}"
fi
build::common::echo_and_run helm package "${CHART_NAME}"
