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


IMAGE_BUILDER_DIR="${1?Specify first argument - image builder directory in the cloned repo}"
OS="${2?Specify second argument - base os of ova built}"

mv "${IMAGE_BUILDER_DIR}"/output/*.ova "${IMAGE_BUILDER_DIR}"/output/"${OS}".ova
sha256sum "${IMAGE_BUILDER_DIR}"/output/"${OS}".ova > "${IMAGE_BUILDER_DIR}"/output/"${OS}".ova.sha256
sha512sum "${IMAGE_BUILDER_DIR}"/output/"${OS}".ova > "${IMAGE_BUILDER_DIR}"/output/"${OS}".ova.sha512
