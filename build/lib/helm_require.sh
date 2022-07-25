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
set -e
set -o errexit
set -o nounset
set -o pipefail

if [ $# -lt 6 ]; then
    >&2 echo "ERROR: expected 6+ arguments, got $#"
    exit 1
fi
IMAGE_OUTPUT=${1?First argument is image output}
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
JSON_SCHEMA_FILE=helm/schema.json
SEDFILE=${OUTPUT_DIR}/helm/sedfile
export IMAGE_TAG
export HELM_REGISTRY=$(aws ecr-public describe-registries --region us-east-1  --output text --query 'registries[*].registryUri')
envsubst <helm/sedfile.template >${SEDFILE}
# Semver requires that our version begin with a digit, so strip the v.
echo "s,version: v,version: ,g" >>${SEDFILE}
for IMAGE in ${HELM_IMAGE_LIST:-}
do
  if [ "${IMAGE}" == "${HELM_DESTINATION_REPOSITORY}" ]
  then
    TAG="${IMAGE_TAG}"
  else
    TAG="${LATEST}"
  fi
  # The skopeo tool isn't programmed to read the image digest on stdin (fie!)
  # However, bash process substitution will fit the bill (yay!) See
  # https://www.gnu.org/software/bash/manual/bash.html#Process-Substitution
  # for more info.
  IMAGE_SHASUM=$(skopeo manifest-digest <(skopeo inspect --raw -n "docker-archive:$IMAGE_OUTPUT"))
  echo "s,{{${IMAGE}}},${IMAGE_SHASUM},g" >>${SEDFILE}
  if [ "${TAG}" == "latest" ]
  then
    USE_TAG=$(aws --region us-east-1 ecr-public describe-images --repository-name ${IMAGE} --image-ids imageTag=latest --query 'imageDetails[0].imageTags' --output yaml | grep -v latest | head -1| sed -e 's/- //') ||
    USE_TAG=$(aws ecr describe-images --repository-name ${IMAGE} --image-id imageTag=latest --query 'imageDetails[0].imageTags' --output yaml | grep -v latest | head -1| sed -e 's/- //') ||
    USE_TAG="latest"
  else
    USE_TAG=$TAG
  fi
  cat >>${REQUIRES_FILE} <<!
  - repository: ${IMAGE}
    tag: ${USE_TAG}
    digest: ${IMAGE_SHASUM}
!
done

if [ -f ${JSON_SCHEMA_FILE} ]
then
  JSON_SCHEMA=$(cat ${JSON_SCHEMA_FILE} | gzip | base64)
  cat >>${REQUIRES_FILE} <<!
  schema: ${JSON_SCHEMA}
!
fi
