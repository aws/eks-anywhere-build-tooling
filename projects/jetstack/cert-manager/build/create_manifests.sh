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
OUTPUT_DIR="$2"
ARTIFACTS_PATH="$3"
TAG="$4"

mkdir -p $OUTPUT_DIR/assets
#TODO: find cert-manager make target to use to generate the cert-manager.yaml file
wget -q https://github.com/jetstack/cert-manager/releases/download/$TAG/cert-manager.yaml -O $OUTPUT_DIR/assets/cert-manager.yaml
cp -rf $OUTPUT_DIR/assets $ARTIFACTS_PATH
