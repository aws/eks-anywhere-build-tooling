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

IMAGE_REGISTRY="${1?First argument is registry}"
HELM_DESTINATION_REPOSITORY="${2?Second argument is helm destination repository}"
OUTPUT_DIR="${3?Third argument is output directory}"
IMAGE_TAG="${4?Fourth argument is image tag}"
HELM_TAG="${5?Fifth argument is helm tag}"
PROJECT_ROOT="${6?Sixth argument is project root}"
LATEST="${7?Seventh argument is latest tag}"
HELM_USE_UPSTREAM_IMAGE="${8?Eigth argument is bool determining whether to use cached upstream images}"
HELM_IMAGE_LIST="${@:9}"

CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
DEST_DIR=${OUTPUT_DIR}/helm/${CHART_NAME}
SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
PACKAGE_DEPENDENCIES=${PACKAGE_DEPENDENCIES:=""}
FORCE_JSON_SCHEMA_FILE=${FORCE_JSON_SCHEMA_FILE:=""}

#
# Image tags
#
REQUIRES_FILE=${DEST_DIR}/requires.yaml
cat >${REQUIRES_FILE} <<!
---
kind: EksaPackageRequires
metadata:
  name: ${HELM_DESTINATION_REPOSITORY}-${HELM_TAG/v}
  namespace: eksa-packages
spec:
  images:
!
JSON_SCHEMA_FILE=$PROJECT_ROOT/helm/schema.json
SEDFILE=${OUTPUT_DIR}/helm/sedfile
export IMAGE_TAG
export HELM_TAG
export HELM_REGISTRY=$(aws ecr-public describe-registries --region us-east-1  --output text --query 'registries[*].registryUri')
envsubst <$PROJECT_ROOT/helm/sedfile.template >${SEDFILE}
# Semver requires that our version begin with a digit, so strip the v.
echo "s,version: v,version: ,g" >>${SEDFILE}
for IMAGE in ${HELM_IMAGE_LIST:-}
do
  if [ "${IMAGE}" == "${HELM_DESTINATION_REPOSITORY}" ] || [ "${IMAGE_TAG}" != "${HELM_TAG}" ]
  then
    TAG="${IMAGE_TAG}"
  else
    TAG="${LATEST}"
  fi
  IMAGE_SHASUM=$(${SCRIPT_ROOT}/image_shasum.sh ${HELM_REGISTRY} ${IMAGE} ${TAG}) ||
  IMAGE_SHASUM=$(${SCRIPT_ROOT}/image_shasum.sh ${IMAGE_REGISTRY} ${IMAGE} ${TAG})
  echo "s,{{${IMAGE}}},${IMAGE_SHASUM},g" >>${SEDFILE}
  if [ "${TAG}" == "latest" ]
  then
    USE_TAG=$(aws --region us-east-1 ecr-public describe-images --repository-name ${IMAGE} --image-ids imageTag=latest --query 'imageDetails[0].imageTags' --output yaml | grep -v latest | head -1| sed -e 's/- //') ||
    USE_TAG=$(aws ecr describe-images --repository-name ${IMAGE} --image-id imageTag=latest --query 'imageDetails[0].imageTags' --output yaml | grep -v latest | head -1| sed -e 's/- //') ||
    USE_TAG="latest"
  else
    USE_TAG=$TAG
  fi
  # If HELM_USE_UPSTREAM_IMAGE is true, we are using images from upstream.
  # Though we pull images directly from upstream for build tooling checks (i.e.
  # get images shasums), we will use cached images in the helm charts. Cached
  # images follow the convention of ${PROJECT_NAME}/${UPSTREAM_IMAGE_NAME}.
  if [ "${HELM_USE_UPSTREAM_IMAGE}" == true ]
  then
    PROJECT_NAME=$(echo "$HELM_DESTINATION_REPOSITORY" | awk -F "/" '{print $1}')
    IMAGE_REPO="${PROJECT_NAME}/${IMAGE}"
  else
    IMAGE_REPO="${IMAGE}"
  fi
  cat >>${REQUIRES_FILE} <<!
  - repository: ${IMAGE_REPO}
    tag: ${USE_TAG}
    digest: ${IMAGE_SHASUM}
!
done

if [ -n "${FORCE_JSON_SCHEMA_FILE}" ]
then 
  JSON_SCHEMA_FILE=${FORCE_JSON_SCHEMA_FILE}
fi

if [ -f ${JSON_SCHEMA_FILE} ]
then
  JSON_SCHEMA=$(cat ${JSON_SCHEMA_FILE} | gzip | base64 | tr -d '\n')
  cat >>${REQUIRES_FILE} <<!
  schema: ${JSON_SCHEMA}
!
fi

if [ -n "${PACKAGE_DEPENDENCIES}" ]; then
  echo "  dependencies:" >> ${REQUIRES_FILE}
  echo ${PACKAGE_DEPENDENCIES} | tr ',' '\n'  | while read dep; do
      echo "  - ${dep}"
  done >> ${REQUIRES_FILE}
fi
