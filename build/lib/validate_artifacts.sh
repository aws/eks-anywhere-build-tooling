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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

PROJECT_ROOT="$1"
ARTIFACTS_FOLDER="$2"
GIT_TAG="$3"
FAKE_ARM_ARTIFACTS_FOR_VALIDATION="$4"
FAKE_AMD_BINARIES_FOR_VALIDATION="$5"
EXPECTED_FILES_PATH="$6"
IMAGE_OS="${7:-}"


REALPATH=$(build::find::gnu_variant_on_mac realpath)

ACTUAL_FILES=$(mktemp)
find "$ARTIFACTS_FOLDER" \
     -type f \
     -exec $REALPATH --relative-base="$ARTIFACTS_FOLDER" {} \; \
     > "$ACTUAL_FILES"

EXPECTED_FILES=$(mktemp)
build::common::echo_and_run export GIT_TAG=$GIT_TAG
build::common::echo_and_run export IMAGE_OS=$IMAGE_OS
envsubst "\$GIT_TAG:\$IMAGE_OS" \
         < "$EXPECTED_FILES_PATH" \
         > "$EXPECTED_FILES"

# Replace forward slashes with hyphens in the expected files to match actual tarball names
# When GIT_TAG has '/' inside
if [[ "$GIT_TAG" == *"/"* ]]; then
   sed -i 's|/|-|g' "$EXPECTED_FILES"
fi

if $FAKE_ARM_ARTIFACTS_FOR_VALIDATION; then
    echo "Faking arm64 artifacts"
    sed -i '/arm64/d' "$EXPECTED_FILES"
    sed -i '/aarch64/d' "$EXPECTED_FILES"
fi

if $FAKE_AMD_BINARIES_FOR_VALIDATION; then
    echo "Faking amd64 artifacts"
    sed -i '/amd64/d' "$EXPECTED_FILES"
    sed -i '/x86_64/d' "$EXPECTED_FILES"
fi

# The versions of sort found on macOS and Linux can behave
# differently. That's why we sort each file here and now, to ensure
# that whatever version is in use, it's consistent.
echo "diffing $EXPECTED_FILES $ACTUAL_FILES"
if ! diff -q <(sort "$EXPECTED_FILES") <(sort "$ACTUAL_FILES"); then
    echo "Artifacts directory does not matched expected!"
    diff -y <(sort "$EXPECTED_FILES") <(sort "$ACTUAL_FILES")
    exit 1
fi
