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
    echo "❌ GIT_TAG does not match the Chart version"
    echo "GIT_TAG=${GIT_TAG}, Chart version=${CHART_VERSION}"
    exit 1
fi

SED_FILE_VERSION="$(grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" ${HELM_DIRECTORY}/../helm/sedfile.template)"
echo "Verifying version in sedfile.template matches the chart version in Chart.yaml"
if [[ ${SED_FILE_VERSION} != ${CHART_VERSION} ]]; then
    echo "❌ version in sedfile.template does not match the Chart version"
    echo "SED_FILE_VERSION=${SED_FILE_VERSION}, Chart version=${CHART_VERSION}"
    exit 1
fi

# PULL_BASE_SHA and PULL_PULL_SHA are environment variables set by the presubmit job. More info here: https://docs.prow.k8s.io/docs/jobs/
BASE_COMMIT_HASH=${PULL_BASE_SHA}
PR_COMMIT_HASH=${PULL_PULL_SHA}
PREVIOUS_CHART_VERSION=$(git show ${BASE_COMMIT_HASH}:./chart/Chart.yaml | yq '.version')

EXIT_CODE=0

if git diff ${BASE_COMMIT_HASH} ${PR_COMMIT_HASH} --quiet -- ${HELM_DIRECTORY}/templates ${HELM_DIRECTORY}/values.yaml --; then
    echo "✅ tinkerbell-chart has no changes since last release"
else
if [ "${CHART_VERSION}" = "${PREVIOUS_CHART_VERSION}" ]; then
    echo "❌ tinkerbell-chart has changed but the Chart version is the same as the last release $PREVIOUS_CHART_VERSION. Please update the chart version in Chart.yaml"
    EXIT_CODE=1
else 
    echo "✅ tinkerbell-chart has a different version since the last release ($PREVIOUS_CHART_VERSION -> $CHART_VERSION)"
fi
fi

exit $EXIT_CODE
