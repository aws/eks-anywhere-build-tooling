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
HELM_TAG="${3?Third argument is image tag}"
OUTPUT_DIR="${4?Fourth arguement is output directory}"

HELM_DESTINATION_OWNER=$(dirname ${HELM_DESTINATION_REPOSITORY})
CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
CHART_FILE=${OUTPUT_DIR}/helm/${CHART_NAME}-${HELM_TAG}-helm.tgz
BUILD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

DOCKER_CONFIG=${DOCKER_CONFIG:-~/.docker}
export HELM_REGISTRY_CONFIG="${DOCKER_CONFIG}/config.json"
export HELM_EXPERIMENTAL_OCI=1
TMPFILE=$(mktemp /tmp/helm-output.XXXXXX)
function cleanup() {
  if echo ${IMAGE_REGISTRY} | grep public.ecr.aws >/dev/null
  then
    echo "If authentication failed: aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws"
  else
    echo "If authentication failed: aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${IMAGE_REGISTRY}"
  fi
  rm -f "${TMPFILE}"
}
trap cleanup err
#trap "rm -f $TMPFILE" exit
helm push ${CHART_FILE} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_OWNER} | tee ${TMPFILE}
DIGEST=$(grep Digest $TMPFILE | sed -e 's/Digest: //')
echo "helm install ${CHART_NAME} oci://${IMAGE_REGISTRY}/${HELM_DESTINATION_REPOSITORY} --version ${DIGEST}"

# # Download Addons and checkout generatebundlefile
# org="aws"
# repo="modelrocket-add-ons"
# aws_region="us-west-2"
# git clone "https://git-codecommit.${aws_region}.amazonaws.com/v1/repos/${org}.${repo}"
# cd aws.modelrocket-add-ons/
# git checkout dont-delete/codebuild-fork
# cd generatebundlefile/
# ./vend.sh
# go1.17.5 run . --input "data/bundle.yaml"

aws s3 cp s3://eks-a-addons-generatebundle/generatebundlefile .
chmod +x ./generatebundlefile
# Python3 pip and yq
pip3 install yq

#  Add the new helm build to the input file
export HELM_TAG="${HELM_TAG}-helm"
export CHART_NAME=${CHART_NAME}

# Add new build to the input file
cat ${BUILD_ROOT}/addons/addons.yaml | yq -y '.addOns = [.addOns[] | select(.projects[].name == env.CHART_NAME).projects[].versions += [{"name": env.HELM_TAG}]]' > ${BUILD_ROOT}/addons/bundle.yaml
yq -y . "${BUILD_ROOT}/addons/bundle.yaml"
./generatebundlefile  --input "${BUILD_ROOT}/addons/bundle.yaml"
yq -y . "output/1.20-bundle-crd.yaml"