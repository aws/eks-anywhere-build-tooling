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

REPO_ROOT="$(git rev-parse --show-toplevel)"

VALIDATIONS_FAILED=0
for buildspec in "$@"; do
    echo "Validating builds count in build graph for buildspec - $buildspec"
    num_builds_in_batch=$(yq ".batch.build-graph | length" $buildspec)
    if [[ $num_builds_in_batch -gt 100 ]]; then
        echo "Maximum allowed builds in batch is 100, current number of builds: $num_builds_in_batch"
        VALIDATIONS_FAILED=1
        INVALID_BUILDSPEC="true"
    fi

    depends_on_list=($(yq "[.batch.build-graph[].depend-on[] | select(. != \"null\")] | unique | .[]" $buildspec))
    identifier_list=($(yq ".batch.build-graph[].identifier" $buildspec))

    INVALID_BUILDSPEC="false"
    echo "Validating identifier uniqueness in the buildspec - $buildspec"
    duplicate_ids=($(printf '%s\n' "${identifier_list[@]}"|awk '!($0 in seen){seen[$0];next} 1' | uniq))
    if [ "${#duplicate_ids[@]}" -gt 0 ]; then
        printf -v duplicate_id_csv '%s,' "${duplicate_ids[@]}"
        echo "Duplicate identifiers found: ${duplicate_id_csv%,}"
        VALIDATIONS_FAILED=1
        INVALID_BUILDSPEC="true"
    fi
    
    invalid_dependencies=()
    if [ "${#depends_on_list[@]}" -gt 0 ]; then
        echo "Validating identifiers in depend-on list are valid identifiers in build graph in the buildspec - $buildspec"
        invalid_dependencies=($(for dependency in ${depends_on_list[@]}; do
            [[ ${identifier_list[*]} =~ (^|[[:space:]])"$dependency"($|[[:space:]]) ]] || echo "$dependency"
        done))
    fi

    if [ "${#invalid_dependencies[@]}" -gt 0 ]; then
        printf -v invalid_deps_csv '%s,' "${invalid_dependencies[@]}"
        echo "Invalid depend-on identifiers found: ${invalid_deps_csv%,}"
        VALIDATIONS_FAILED=1
        INVALID_BUILDSPEC="true"
    fi

    if [[ "$INVALID_BUILDSPEC" == "false" ]]; then
        echo "All validations passed for the buildspec - $buildspec!"
    fi
done

exit $VALIDATIONS_FAILED
