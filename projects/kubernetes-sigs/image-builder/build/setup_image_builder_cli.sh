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

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

function image_build::common::download_latest_dev_image_builder_cli() {
  local -r download_location=${1-/usr/bin}
  local -r artifacts_bucket=${2-$ARTIFACTS_BUCKET}
  local -r arch=${3-amd64}

  local -r latest_tar_url=$(build::common::get_latest_eksa_asset_url $artifacts_bucket 'aws/image-builder' $arch 'latest')

  # Download the tar ball and decompress
  local -r http_code=$(curl -s -w "%{http_code}" $latest_tar_url -o image-builder.tar.gz)
  if [[ "$http_code" != "200" ]]; then
    echo "Error downloading latest image builder cli"
    exit 1
  fi
  tar -xvf image-builder.tar.gz --directory $download_location/ ./image-builder
  rm -rf image-builder.tar.gz
}
