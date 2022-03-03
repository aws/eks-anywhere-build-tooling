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


set -xe
set -o errexit
set -o nounset

BUILDER_ROOT="${1?Specify first argument - builder root}"
BASE_IMAGE="${2?Specify second argument - base rhel 8 image}"
IMAGE_BUILDER_DIR="${3?Specify third argument - Image builder directory}"
PACKER_VAR_FILES="${4?Specify fourth argument - Packer var files}"

# Validate env vars
if [ -z $BASE_IMAGE ]; then
  echo "BASE_IMAGE env variable not set. Please set and re-try"
  exit 1
fi
if [ -z $RHSM_USER ]; then
  echo "RHSM_USER env variable not set. Please set and re-try"
  exit 1
fi
if [ -z $RHSM_PASS ]; then
  echo "RHSM_PASS env variable not set. Please set and re-try"
  exit 1
fi

RHEL_QEMU_CONFIG_FILE=$IMAGE_BUILDER_DIR/packer/qemu/qemu-rhel-8.json

echo "Generating base image sha256sum"
BASE_IMAGE_CHECKSUM=$(sha256sum $BASE_IMAGE| awk '{print $1}')

echo "Base Image Checksum - $BASE_IMAGE_CHECKSUM"

cat <<< $(jq '.iso_url="$(BUILDER_ROOT)/$(BASE_IMAGE)"|.iso_checksum="$(BASE_IMAGE_CHECKSUM)"' $RHEL_QEMU_CONFIG_FILE) > $RHEL_QEMU_CONFIG_FILE

PACKER_FLAGS="-force" PACKER_LOG=1 PACKER_VAR_FILES="$PACKER_VAR_FILES" make -C $IMAGE_BUILDER_DIR build-qemu-rhel-8

echo "The output image is located at $IMAGE_BUILDER_DIR/output"