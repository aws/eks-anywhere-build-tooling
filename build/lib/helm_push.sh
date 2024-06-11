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

SED=$(build::find::gnu_variant_on_mac sed)

IMAGE_REGISTRY="${1?First argument is image registry}"
HELM_DESTINATION_REPOSITORY="${2?Second argument is helm repository}"
HELM_CHART_FOLDER="${3?Third arguemet is helm chart folder}"
HELM_TAG="${4?Fourth argument is helm tag}"
GIT_TAG="${5?Fifth arguement is the Git Tag}"
OUTPUT_DIR="${6?Sixth arguement is output directory}"
LATEST_TAG="${7?Seventh arguement is latest tag}"
SEMVER_GIT_TAG="${GIT_TAG#[^0-9:main]}"

SEMVER="${HELM_TAG#[^0-9]}" # remove any leading non-digits
SEMVER_REGEX='^([0-9]+\.){0,2}(\*|[0-9]+)$'
if [[ ! $SEMVER_GIT_TAG =~ $SEMVER_REGEX ]]; then
  # if not a valid semver, fallback to helm tag semver
  SEMVER_GIT_TAG=$SEMVER
fi

HELM_DESTINATION_OWNER=$(dirname ${HELM_DESTINATION_REPOSITORY})
CHART_NAME=$(basename ${HELM_CHART_FOLDER})
CHART_FILE="${OUTPUT_DIR}/helm/${CHART_NAME}-${SEMVER}.tgz"

DOCKER_CONFIG=${DOCKER_CONFIG:-~/.docker}
export HELM_REGISTRY_CONFIG="${DOCKER_CONFIG}/config.json"
export HELM_EXPERIMENTAL_OCI=1
TMPFILE=$(mktemp /tmp/helm-output.XXXXXX)
function cleanup() {
  if grep -q "blobs/uploads/\": EOF" $TMPFILE || grep -q "blobs/uploads.*404 Not Found" $TMPFILE; then
    echo "******************************************************"
    echo "Ensure container registry and repository exists!!"
    echo "Try running make create-ecr-repos to create ecr repositories in your aws account."
    echo "******************************************************"
  else
    cat $TMPFILE
    if [[ "${IMAGE_REGISTRY}" == *"public.ecr.aws"* ]]; then
      echo "If authentication failed: aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws"
    else
      echo "If authentication failed: aws ecr get-login-password --region \${AWS_REGION} | docker login --username AWS --password-stdin ${IMAGE_REGISTRY}"
    fi
  fi
  rm -f "${TMPFILE}"
}

trap cleanup err
trap "rm -f $TMPFILE" exit
echo "($(pwd)) \$ helm push ${CHART_FILE} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER}"
helm push ${CHART_FILE} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER} 2>&1 | tee ${TMPFILE}
DIGEST=$(grep Digest $TMPFILE | $SED -e 's/Digest: //')

# Adds a 2nd tag to the helm chart for the bundle-release jobs.
build::common::echo_and_run skopeo copy docker://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER}/${CHART_NAME}@${DIGEST} docker://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER}/${CHART_NAME}:${SEMVER_GIT_TAG}-${LATEST_TAG}-helm

{
    set +x
    echo
    echo
    echo "helm install ${CHART_NAME} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER}/${CHART_NAME} --version ${SEMVER} --set sourceRegistry=${IMAGE_REGISTRY}"
}
