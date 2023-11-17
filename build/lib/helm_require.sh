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
PACKAGE_DEPENDENCIES="${9?Ninth argument is optional project dependencies}"
FORCE_JSON_SCHEMA_FILE="${10?Tenth argument is optional schema file}"
HELM_IMAGE_LIST="${@:11}"

CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
DEST_DIR=${OUTPUT_DIR}/helm/${CHART_NAME}
SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

if ! aws sts get-caller-identity &> /dev/null; then
  echo "The AWS cli is used to find the ECR registries and repos for the current AWS account please login!"
  exit 1;
fi

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
export HELM_REGISTRY=$(aws ecr-public describe-registries --region us-east-1  --output text --query 'registries[*].registryUri' 2> /dev/null)
envsubst <$PROJECT_ROOT/helm/sedfile.template >${SEDFILE}
# Semver requires that our version begin with a digit, so strip the v.
echo "s,version: v,version: ,g" >>${SEDFILE}

# query ecr for the shasum of the image tagged with ${TAG}
function get_image_shasum() {
  local -r image=$1
  local -r tag=$2

  local image_shasum=
  if [ "${HELM_USE_UPSTREAM_IMAGE}" = true ]; then
    image_shasum=$(${SCRIPT_ROOT}/image_shasum.sh ${IMAGE_REGISTRY} ${image} ${tag})
  elif [ "${JOB_TYPE:-}" = "presubmit" ]; then
    image_shasum=${LATEST}
  fi

  if [[ -z ${image_shasum} ]] && aws --region us-east-1 ecr-public describe-repositories --repository-names ${image} &> /dev/null; then
    image_shasum=$(build::common::echo_and_run ${SCRIPT_ROOT}/image_shasum.sh ${HELM_REGISTRY} ${image} ${tag})
  fi

  if [[ -z ${image_shasum} ]] && aws ecr describe-repositories --repository-names ${image} &> /dev/null; then
    image_shasum=$(build::common::echo_and_run ${SCRIPT_ROOT}/image_shasum.sh ${IMAGE_REGISTRY} ${image} ${tag})    
  fi

  if [[ -n ${image_shasum} ]]; then
    echo ${image_shasum}
  else
    echo "${image} does not exist in ECR Public or Private"
    exit 1
  fi  
}

  # query ecr for the image by latest tag and find the first non-latest tag the image is also tagged with
function get_image_tag_not_latest() {
    local -r image=$1
    local -r tag=$2

    local use_tag=
    if [ "${JOB_TYPE:-}" = "presubmit" ]; then
      use_tag=${tag}      
    fi

    if [[ -z ${use_tag} ]] && aws --region us-east-1 ecr-public describe-repositories --repository-names ${image} &> /dev/null; then
      use_tag=$(build::common::echo_and_run aws --region us-east-1 ecr-public describe-images --repository-name ${image} --image-ids imageTag=${tag} --query 'imageDetails[0].imageTags' --output yaml 2> /dev/null | grep -v ${tag} | head -1| sed -e 's/- //')
    fi
    
    if [[ -z ${use_tag} ]] && aws ecr describe-repositories --repository-names ${image} &> /dev/null; then
      use_tag=$(build::common::echo_and_run  aws ecr describe-images --repository-name ${image} --image-id imageTag=${tag} --query 'imageDetails[0].imageTags' --output yaml 2> /dev/null | grep -v ${tag} | head -1| sed -e 's/- //')
    fi

    if [[ -n ${use_tag} ]]; then
      echo ${use_tag}
    else
      echo "${image}@${tag} does not exist in ECR Public or Private"
      exit 1
    fi  
}

for IMAGE in ${HELM_IMAGE_LIST:-}; do
  # if its the image(s) built from this project, use the image_tag
  # otherwise its an image from a different project so use latest to trigger finding the latest image
  if [ "${IMAGE}" = "${HELM_DESTINATION_REPOSITORY}" ] || [ "${IMAGE_TAG}" != "${HELM_TAG}" ]; then
    TAG="${IMAGE_TAG}"
  else
    TAG="${LATEST}"
  fi
  
  IMAGE_SHASUM=$(get_image_shasum ${IMAGE} ${TAG})

  echo "s,{{${IMAGE}}},${IMAGE_SHASUM},g" >>${SEDFILE}
  if [ "${TAG}" = "${LATEST}" ];  then
    USE_TAG=$(get_image_tag_not_latest ${IMAGE} ${LATEST})
  else
    USE_TAG=$TAG
  fi
  
  # If HELM_USE_UPSTREAM_IMAGE is true, we are using images from upstream.
  # Though we pull images directly from upstream for build tooling checks (i.e.
  # get images shasums), we will use cached images in the helm charts. Cached
  # images follow the convention of ${PROJECT_NAME}/${UPSTREAM_IMAGE_NAME}.
  if [ "${HELM_USE_UPSTREAM_IMAGE}" == true ]; then
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

if [ -n "${FORCE_JSON_SCHEMA_FILE}" ]; then 
  JSON_SCHEMA_FILE=${FORCE_JSON_SCHEMA_FILE}
fi

if [ -f ${JSON_SCHEMA_FILE} ]; then
  echo "Using schema file: ${JSON_SCHEMA_FILE}"
  JSON_SCHEMA=$(cat ${JSON_SCHEMA_FILE} | gzip -n | base64 | tr -d '\n')
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

build::common::echo_and_run cat ${SEDFILE}
build::common::echo_and_run cat ${REQUIRES_FILE}
