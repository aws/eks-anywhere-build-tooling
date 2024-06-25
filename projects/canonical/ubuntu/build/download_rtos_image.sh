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

RTOS_BUCKET_NAME="${1?Specify first argument - Ubuntu RTOS image bucket name}"
RTOS_IMAGE_DATE="${2?Specify second argument - Ubuntu RTOS image build date}"
ARTIFACTS_PATH="${3?Specify third argument - artifacts path}"
RELEASE_BRANCH="${4?Specify fourth argument - release branch}"

function build::download::ubuntu::rtos::image(){
    mkdir -p $ARTIFACTS_PATH
    download_path=s3://$RTOS_BUCKET_NAME/ubuntu/jammy/$RTOS_IMAGE_DATE/ubuntu-jammy-eks-anywhere-pro-realtime-minimal-amd64-eks-anywhere-$RELEASE_BRANCH-pro-realtime.raw.gz
    aws s3 cp $download_path $ARTIFACTS_PATH/ubuntu.gz
}

build::download::ubuntu::rtos::image
