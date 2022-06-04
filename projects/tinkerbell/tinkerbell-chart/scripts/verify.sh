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

HELM_DIRECTORY="$1"
GIT_TAG="$2"

echo "Running helm lint"
helm lint ${HELM_DIRECTORY}

CHART_VERSION="$(yq eval '.version' ${HELM_DIRECTORY}/Chart.yaml)"
echo "Verifying GIT_TAG matches the chart version in Chart.yaml"
if [[ ${GIT_TAG} != ${CHART_VERSION} ]]; then
    echo "GIT_TAG does not match the Chart version"
    echo "GIT_TAG=${GIT_TAG}, Chart version=${CHART_VERSION}"
    exit 1
fi
