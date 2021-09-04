#!/usr/bin/env bash
# Copyright 2020 Amazon.com Inc. or its affiliates. All Rights Reserved.
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
KIND_KINDNETD_IMAGE="${2?Specify second argument - kindnetd image}"
KIND_KINDNETD_LATEST_IMAGE="${3?Specify third argument - kindnetd latest image}"
PUSH="${4?Specify fourth argument - push}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD_LIB="${MAKE_ROOT}/../../../build/lib"
OUTPUT="dest=/tmp/kindnetd.tar"
TYPE="type=oci"

if [ $PUSH = "true" ]; then
    OUTPUT="push=true"
    TYPE="type=image"
fi

function build::kindnetd::image(){
    $BUILD_LIB/buildkit.sh build \
        --frontend dockerfile.v0 \
        --opt platform=linux/amd64 \
        --opt build-arg:BASE_IMAGE=$BASE_IMAGE \
        --local dockerfile="$MAKE_ROOT/images/kindnetd" \
        --local context=$MAKE_ROOT \
        --opt filename=Dockerfile \
        --progress plain \
        --output $TYPE,oci-mediatypes=true,\"name=${KIND_KINDNETD_IMAGE},${KIND_KINDNETD_LATEST_IMAGE}\",$OUTPUT
}

build::kindnetd::image

if command -v docker &> /dev/null && docker info > /dev/null 2>&1 && [ ! $PUSH = "true" ]; then
    # running locally and to make it easier to build the node image
    # load resulting tar into docker
    IMAGE_ID=$(docker load -i /tmp/kindnetd.tar | sed -E 's/.*sha256:(.*)$/\1/')
    docker tag $IMAGE_ID ${KIND_KINDNETD_IMAGE}
fi
