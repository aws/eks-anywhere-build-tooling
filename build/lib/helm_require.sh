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

HELM_REGISTRY="${1?First argument is registry}"
HELM_DESTINATION_REPOSITORY="${2?Second argument is helm destination repository}"
OUTPUT_DIR="${3?Third argument is output directory}"
IMAGE_TAG="${4?Fourth argument is image tag}"
LATEST="${5?Fifth argument is latest tag}"
HELM_IMAGE_LIST="${@:6}"

CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
DEST_DIR=${OUTPUT_DIR}/helm/${CHART_NAME}
SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

#
# Image tags
#
REQUIRES_FILE=${DEST_DIR}/requires.yaml
cat >${REQUIRES_FILE} <<!
---
kind: EksaPackageRequires
metadata:
  name: ${HELM_DESTINATION_REPOSITORY}-${IMAGE_TAG/v}
  namespace: eksa-packages
spec:
  images:
!
REQUIRES_CONFIG_FILE=helm/requires-config.yaml
SEDFILE=${OUTPUT_DIR}/helm/sedfile
export IMAGE_TAG
export HELM_REGISTRY
envsubst <helm/sedfile.template >${SEDFILE}
echo "s,version: v,version: ,g" >>${SEDFILE}
for IMAGE in ${HELM_IMAGE_LIST:-}
do
  IMAGE_SHASUM=$(${SCRIPT_ROOT}/image_shasum.sh ${HELM_REGISTRY} ${IMAGE} ${LATEST})
  echo "s,{{${IMAGE}}},${IMAGE_SHASUM},g" >>${SEDFILE}
  cat >>${REQUIRES_FILE} <<!
  - repository: ${IMAGE}
    digest: ${IMAGE_SHASUM}
!
done
if [ -f ${REQUIRES_CONFIG_FILE} ]
then
  cat ${REQUIRES_CONFIG_FILE} >>${REQUIRES_FILE}
fi
