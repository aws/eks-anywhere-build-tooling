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

BASE_IMAGE="${1?Specify first argument - base image}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"

# to build git2go we need the headers and pc files for libgit which
# are included in the base image
# run buildctl to extract the files we needed into the host
$BUILD_LIB/buildkit.sh build \
  --frontend dockerfile.v0 \
  --opt platform=linux/amd64 \
  --opt build-arg:BASE_IMAGE=${BASE_IMAGE} \
  --local dockerfile=${MAKE_ROOT}/build \
  --opt filename=./Dockerfile.deps \
  --local context=. \
  --progress plain \
  --output type=local,dest=/
