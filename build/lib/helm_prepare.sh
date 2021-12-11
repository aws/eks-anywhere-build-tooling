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
set -o nounset
set -o pipefail

export HELM_REPOSITORY="${1?First argument is helm repoistory}"
export HELM_DIRECTORY="${2?Second argument is helm directory}"
export IMAGE_REPOSITORY="${3?Third argument is image repository}"
export IMAGE_TAG="${4?Fourth argument is image repository}"
export OUTPUT_DIR="${5?Fifth arguement is output directory}"
export CHART_NAME=$(basename ${IMAGE_REPOSITORY})

SOURCE_DIR=$(basename ${HELM_REPOSITORY})/${HELM_DIRECTORY}
DEST_DIR=${OUTPUT_DIR}/helm/${CHART_NAME}

mkdir -p ${OUTPUT_DIR}/helm/${CHART_NAME}
cp ${OUTPUT_DIR}/ATTRIBUTION.txt ${DEST_DIR}/
cp -r ${SOURCE_DIR} ${DEST_DIR}
envsubst <helm/Chart.yaml.template >${DEST_DIR}/Chart.yaml
envsubst <helm/values.yaml.template >${DEST_DIR}/values.yaml
