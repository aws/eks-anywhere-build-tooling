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
export LC_ALL=C

runtime_root="${1:?Usage: $0 RUNTIME_ROOT OUTPUT_FILE}"
output_file="${2:?Usage: $0 RUNTIME_ROOT OUTPUT_FILE}"
work_dir=$(mktemp -d)
trap 'rm -rf "$work_dir"' EXIT

rpm --root "$runtime_root" -qa \
    --qf '%{NAME}\t%{VERSION}-%{RELEASE}\t%{LICENSE}\t%{SOURCERPM}\n' \
    > "$work_dir/installed"
awk -F '\t' 'toupper($3) ~ /GPL|SISSL|SLEEPYCAT/' "$work_dir/installed" \
    | sort > "$work_dir/copyleft"
if [ ! -s "$work_dir/copyleft" ]; then
    echo "No copyleft runtime RPMs found in $runtime_root" >&2
    exit 1
fi

mapfile -t sources < <(cut -f 4 "$work_dir/copyleft" | sort -u)
source_specs=()
for source in "${sources[@]}"; do
    if [[ "$source" != *.src.rpm ]]; then
        echo "Missing source RPM metadata: $source" >&2
        exit 1
    fi
    source_specs+=("${source%.src.rpm}")
done

# Query the installed source NVR, not the newest version of a package name.
if command -v dnf > /dev/null; then
    dnf -q repoquery --repo=amazonlinux-source --arch=src --location \
        "${source_specs[@]}" > "$work_dir/query"
else
    yumdownloader -q --source --urls --enablerepo=amzn2-core-source \
        "${source_specs[@]}" > "$work_dir/query"
fi

# AL2023 repository URLs contain ../ segments. Keep the permanent blobstore URL.
touch "$work_dir/urls"
while IFS= read -r url; do
    [[ "$url" == https://cdn.amazonlinux.com/* ]] || continue
    url="https://cdn.amazonlinux.com$(realpath -ms -- "/${url#https://cdn.amazonlinux.com/}")"
    if [[ ! "$url" =~ ^https://cdn\.amazonlinux\.com/(al2023/)?blobstore/[[:xdigit:]]{64}/[^/]+\.src\.rpm$ ]]; then
        echo "Unexpected source RPM URL: $url" >&2
        exit 1
    fi
    printf '%s\n' "$url" >> "$work_dir/urls"
done < "$work_dir/query"
sort -u "$work_dir/urls" -o "$work_dir/urls"

cat > "$work_dir/attribution" <<'EOF'
Amazon Linux runtime source notices

These links identify the exact source RPMs for the copyleft runtime packages.
RPM license and documentation files are included under /LICENSES/amazon-linux.
EOF

for source in "${sources[@]}"; do
    mapfile -t urls < <(awk -F / -v source="$source" '$NF == source' "$work_dir/urls")
    # In particular, repoquery can return success without finding a version.
    if [ "${#urls[@]}" -ne 1 ]; then
        echo "Expected one exact-version source URL for $source; found ${#urls[@]}" >&2
        exit 1
    fi
    packages=$(awk -F '\t' -v source="$source" '$4 == source {print $1 " " $2}' "$work_dir/copyleft" \
        | sort -u | paste -sd ',' - | sed 's/,/, /g')
    licenses=$(awk -F '\t' -v source="$source" '$4 == source {print $3}' "$work_dir/copyleft" \
        | sort -u | paste -sd ';' - | sed 's/;/; /g')
    printf '\nPackages: %s\nLicense metadata: %s\nSource: %s\n' \
        "$packages" "$licenses" "${urls[0]}" >> "$work_dir/attribution"
done

# Do not overwrite a previous notice file when any lookup or validation fails.
install -m 0644 "$work_dir/attribution" "$output_file"
