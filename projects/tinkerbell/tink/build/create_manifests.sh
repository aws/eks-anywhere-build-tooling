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
set -o pipefail

REPO="$1"
ARTIFACTS_PATH="$2"
IMAGE_REPO="$3"
TINK_SERVER_IMAGE_COMPONENT="$4"
TINK_CONTROLLER_IMAGE_COMPONENT="$5"
IMAGE_TAG="$6"
GOLANG_VERSION="$7"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

build::common::use_go_version ${GOLANG_VERSION}

cd $REPO

npm install prettier

PATH=$(pwd)/node_modules/.bin/:$PATH make release-manifests \
  TINK_SERVER_IMAGE=${IMAGE_REPO}/${TINK_SERVER_IMAGE_COMPONENT} \
  TINK_CONTROLLER_IMAGE=${IMAGE_REPO}/${TINK_CONTROLLER_IMAGE_COMPONENT} \
  TINK_SERVER_TAG=${IMAGE_TAG} \
  TINK_CONTROLLER_TAG=${IMAGE_TAG}

cp out/release/tink.yaml $ARTIFACTS_PATH
