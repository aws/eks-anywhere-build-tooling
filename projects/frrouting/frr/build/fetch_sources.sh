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

set -euo pipefail

frr_tag="${1:?FRR tag is required}"
libyang_tag="${2:?libyang tag is required}"
tini_tag="${3:?tini tag is required}"
output_dir="${4:?output directory is required}"

mkdir -p "$output_dir/sources" "$output_dir/src"

fetch() {
    local owner="$1" repo="$2" tag="$3"
    local archive="$output_dir/sources/$repo-$tag.tar.gz"
    if [[ ! -f "$archive" ]]; then
        curl --fail --location --retry 5 \
            "https://codeload.github.com/$owner/$repo/tar.gz/refs/tags/$tag" \
            --output "$archive.tmp"
        mv "$archive.tmp" "$archive"
    fi
}

fetch FRRouting frr "$frr_tag"
fetch CESNET libyang "$libyang_tag"
fetch krallin tini "$tini_tag"

# Verify every archive before extracting or compiling any source.
sha256sum -c CHECKSUMS

extract() {
    local repo="$1" tag="$2"
    local staging_dir
    staging_dir="$(mktemp -d "$output_dir/src/.$repo.XXXXXX")"
    tar -xzf "$output_dir/sources/$repo-$tag.tar.gz" \
        --strip-components=1 -C "$staging_dir"
    rm -rf "$output_dir/src/$repo"
    mv "$staging_dir" "$output_dir/src/$repo"
}

extract frr "$frr_tag"
extract libyang "$libyang_tag"
extract tini "$tini_tag"
