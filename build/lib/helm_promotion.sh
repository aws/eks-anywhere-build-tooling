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
set -x
set -o pipefail

export LANG=C.UTF-8

BASE_DIRECTORY=$(git rev-parse --show-toplevel)
ACCOUNT_ID="${1?Specify first argument - account id}"
IMAGE_REGISTRY="${2?Specify first argument - image registry}"
GIT_HASH="${3?Specify first argument - git hash}"

export HELM_EXPERIMENTAL_OCI=1

ECR_REPO=$(echo ${PROJECT_PATH} | sed "s/aws\///") # We sed to remove aws/ from ECR repo names
echo "ECR_REPO=${ECR_REPO}"
GIT_TAG=$(cat ${BASE_DIRECTORY}/projects/${PROJECT_PATH}/GIT_TAG)

export IMAGE_TAG=${GIT_TAG}-${GIT_HASH}
echo "IMAGE_TAG=${IMAGE_TAG}"

# Pull Helm chart from private ECR
aws ecr get-login-password --region us-west-2 | HELM_EXPERIMENTAL_OCI=1 helm registry login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com
helm pull oci://${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com/${ECR_REPO} --version ${IMAGE_TAG}-helm

# Push to public repo
aws ecr-public get-login-password --region us-east-1 | HELM_EXPERIMENTAL_OCI=1 helm registry login --username AWS --password-stdin public.ecr.aws
helm push ${ECR_REPO}-${IMAGE_TAG}-helm.tgz oci://${IMAGE_REGISTRY}

