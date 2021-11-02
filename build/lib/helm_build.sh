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

IMAGE_REPOSITORY="${1?First argument is image repository}"
IMAGE_TAG="${2?Second argument is image tag}"
IMAGE_DESCRIPTION="${3?Third argument is image description}"
IMAGE_REGISTRY="${4:-}"

cd helm
cat >${IMAGE_REPOSITORY}/Chart.yaml <<!
apiVersion: v2
name: ${IMAGE_REPOSITORY}
description: ${IMAGE_DESCRIPTION}
type: application
version: ${IMAGE_TAG}-helm
appVersion: "${IMAGE_TAG}-helm"
!
trap "rm -f ${IMAGE_REPOSITORY}-${IMAGE_TAG}-helm.tgz ${IMAGE_REPOSITORY}/Chart.yaml" err exit
helm package ${IMAGE_REPOSITORY}

if [ -n "${IMAGE_REGISTRY}" ]
then
  export HELM_EXPERIMENTAL_OCI=1
  export DOCKER_CONFIG=~/.docker
  export HELM_REGISTRY_CONFIG="${DOCKER_CONFIG}/config.json"
  helm push ${IMAGE_REPOSITORY}-${IMAGE_TAG}-helm.tgz oci://${IMAGE_REGISTRY} ||
   (echo "If authentication failed: aws ecr get-login-password --region ${AWS_REGION} | helm registry login --username AWS --password-stdin ${IMAGE_REGISTRY}" &&
   false)
fi
