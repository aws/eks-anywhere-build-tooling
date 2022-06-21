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

PROJECT_ROOT="$1"
ARTIFACTS_FOLDER="$2"
GIT_TAG="$3"
FAKE_ARM_ARTIFACTS_FOR_VALIDATION="$4"
IMAGE_FORMAT="${5:-}"
IMAGE_OS="${6:-}"

EXPECTED_FILES_PATH=$PROJECT_ROOT/expected_artifacts
if [ -n "$IMAGE_FORMAT" ]; then
  if [ "$IMAGE_OS" = "bottlerocket" ]; then
    EXPECTED_FILES_PATH=${PROJECT_ROOT}/expected-artifacts/expected_artifacts_${IMAGE_FORMAT}_bottlerocket
  else
    EXPECTED_FILES_PATH=${PROJECT_ROOT}/expected-artifacts/expected_artifacts_${IMAGE_FORMAT}
  fi
fi

ACTUAL_FILES=$(mktemp)
find "$ARTIFACTS_FOLDER" \
     -type f \
     -exec realpath --relative-base="$ARTIFACTS_FOLDER" {} \; \
     > "$ACTUAL_FILES"

EXPECTED_FILES=$(mktemp)
export GIT_TAG=$GIT_TAG
export IMAGE_OS=$IMAGE_OS
envsubst "\$GIT_TAG:\$IMAGE_OS" \
         < "$EXPECTED_FILES_PATH" \
         > "$EXPECTED_FILES"

if $FAKE_ARM_ARTIFACTS_FOR_VALIDATION; then
    sed -i '/arm64/d' "$EXPECTED_FILES"
fi

# The versions of sort found on macOS and Linux can behave
# differently. That's why we sort each file here and now, to ensure
# that whatever version is in use, it's consistent.
if ! diff -q <(sort "$EXPECTED_FILES") <(sort "$ACTUAL_FILES"); then
    echo "Artifacts directory does not matched expected!"
    diff -y <(sort "$EXPECTED_FILES") <(sort "$ACTUAL_FILES")
    exit 1
fi
