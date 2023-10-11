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

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

SED=$(build::find::gnu_variant_on_mac sed)

IMAGE_REGISTRY="${1?First argument is image registry}"
HELM_REPO_URL="${2?Second arguement is update repo url}"
HELM_PULL_NAME="${3?Third argument is name of helm chart}"
REPO="${4?Fourth argument is helm source repository}"
HELM_DIRECTORY="${5?Fifth argument is helm directory name}"
CHART_VERSION="${6?Sixth argument is the version of helm chart we need to pull}"
COPY_CRDS="${7?Seventh argument is whether we need add crds to the helm chart}"

function cleanup() {
  if [ -d "$HELM_PULL_NAME" ]; then
    rm -rf "${HELM_PULL_NAME}"
  fi
}

function cleanup-chart() {
  if [ -d "${REPO}/${HELM_DIRECTORY}" ]; then
    rm -rf "${REPO}/${HELM_DIRECTORY}"
  fi
}

function copy-crds() {
  if [ -d "${COPY_CRDS}" ]; then
    rm -rf  ${REPO}/${HELM_DIRECTORY}/crds/*
    mkdir ${REPO}/${HELM_DIRECTORY}/crds || true
    cp ${COPY_CRDS}/* ${REPO}/${HELM_DIRECTORY}/crds
  fi
}

trap cleanup err

helm repo add ${IMAGE_REGISTRY} ${HELM_REPO_URL}
helm repo update
helm pull ${IMAGE_REGISTRY}/${HELM_PULL_NAME} --version ${CHART_VERSION} --untar

cleanup-chart

mv ${HELM_PULL_NAME} ${REPO}/${HELM_DIRECTORY}

copy-crds

cleanup
