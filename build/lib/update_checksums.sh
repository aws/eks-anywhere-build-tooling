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

MAKE_ROOT="$1"
PROJECT_ROOT="$2"
OUTPUT_BIN_DIR="$3"

REALPATH=$(build::find::gnu_variant_on_mac realpath)

if [ ! -d ${OUTPUT_BIN_DIR} ] ;  then
    echo "${OUTPUT_BIN_DIR} not present! Run 'make binaries'"
    exit 1
fi

CHECKSUMS_FILE=$PROJECT_ROOT/CHECKSUMS

# Create associative array to store checksums
declare -A checksums

# Read existing checksums if the file exists
if [ -f "$CHECKSUMS_FILE" ]; then
    while IFS=' ' read -r checksum filepath || [ -n "$checksum" ]; do
        # Skip empty lines and malformed entries
        if [ -n "$checksum" ] && [ -n "$filepath" ]; then
            checksums["$filepath"]="$checksum"
        fi
    done < "$CHECKSUMS_FILE"
fi

# Calculate checksums for files present in output directory
for file in $(find ${OUTPUT_BIN_DIR} -type f | sort); do
    filepath=$($REALPATH --relative-base=$MAKE_ROOT $file)
    checksum=$(sha256sum $filepath | cut -d' ' -f1)
    checksums["$filepath"]="$checksum"
done

# Write all checksums to file (both updated and preserved)
rm -f $CHECKSUMS_FILE
for filepath in $(printf '%s\n' "${!checksums[@]}" | sort); do
    echo "${checksums[$filepath]}  $filepath" >> $CHECKSUMS_FILE
done

echo "*************** CHECKSUMS ***************"
cat $CHECKSUMS_FILE
echo "*****************************************"
