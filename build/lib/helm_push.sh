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

IMAGE_REGISTRY="${1?First argument is image registry}"
HELM_DESTINATION_REPOSITORY="${2?Second argument is helm repository}"
IMAGE_TAG="${3?Third argument is image tag}"
OUTPUT_DIR="${4?Fourth arguement is output directory}"

HELM_DESTINATION_OWNER=$(dirname ${HELM_DESTINATION_REPOSITORY})
CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
CHART_FILE=${OUTPUT_DIR}/helm/${CHART_NAME}-${IMAGE_TAG}-helm.tgz

export HELM_EXPERIMENTAL_OCI=1
export DOCKER_CONFIG=~/.docker
export HELM_REGISTRY_CONFIG="${DOCKER_CONFIG}/config.json"
if echo ${IMAGE_REGISTRY} | grep public.ecr.aws >/dev/null
then
  echo "If authentication fails: aws ecr-public get-login-password --region us-east-1 | helm registry login --username AWS --password-stdin public.ecr.aws"
else
  echo "If authentication fails: aws ecr get-login-password --region ${AWS_REGION} | helm registry login --username AWS --password-stdin ${IMAGE_REGISTRY}"
fi
TMPFILE=$(mktemp /tmp/helm-output.XXXXXX)
helm push ${CHART_FILE} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER} | tee ${TMPFILE}
DIGEST=$(grep Digest $TMPFILE | sed -e 's/Digest: //')
#rm -f $TMPFILE
echo "helm install oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_REPOSITORY} --version ${DIGEST} --generate-name"
