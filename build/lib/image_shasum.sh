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

REGISTRY="${1?First argument is registry}"
REPOSITORY="${2?Second argument is repository}"
IMAGE_TAG="${3?Third argument is image tag}"

JSONPATH="imageDetails[?(@.imageTags[?(@ == '"${IMAGE_TAG}"')])].imageDigest"
if [[ "${REGISTRY}" == *"public.ecr.aws"* ]]
then
  aws ecr-public describe-images --region us-east-1 --repository-name "${REPOSITORY}" --output text  --query "${JSONPATH}"
else
  aws ecr describe-images --repository-name "${REPOSITORY}" --output text  --query "${JSONPATH}"
fi
