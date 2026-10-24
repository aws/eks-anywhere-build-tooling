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

ARTIFACTS_PATH="$1"
BASE_DIRECTORY="$2"
TAG="$3"

# TODO: use cert-manager make target to generate the cert-manager.yaml file when we upgrade cert-manager version (v1.7.0-alpha.0)
#
# The following commits would need to be added as patches:
# https://github.com/cert-manager/cert-manager/commit/6734e9b7469288b51848eb209597a1920e4801ea
# https://github.com/cert-manager/cert-manager/commit/32d716654a1091e99e80c10a2798cd839a705713
#
# Then would need to run the make target as follows (can specify -jN, like `-j8`, to run whatever targets that can in parallel):
# make -f make/Makefile bin/yaml/cert-manager.yaml
mkdir -p $ARTIFACTS_PATH/manifests/$TAG
cp $BASE_DIRECTORY/projects/cert-manager/cert-manager/manifests/cert-manager.yaml $ARTIFACTS_PATH/manifests/$TAG
