#!/usr/bin/env bash
# Copyright 2020 Amazon.com Inc. or its affiliates. All Rights Reserved.
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

set -e
set -o errexit
set -o nounset
set -o pipefail
shopt -s globstar

PROJECT="$1"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
PROJECT_ROOT=$MAKE_ROOT/$PROJECT
FETCH_PROJECTS=(
    "fluxcd/flux2"
    "fluxcd/kustomize-controller"
    "kubernetes-sigs/cluster-api"
)

mkdir -p _output
touch _output/total_summary.txt

function build::attribution::generate(){
    if [ $(printf "%s\n" "${FETCH_PROJECTS[@]}" | grep -c "^$(cut -d / -f2- <<< $PROJECT)$" || true) -ne 0 ]; then
        make -C $PROJECT_ROOT create-binaries
    else
        make -C $PROJECT_ROOT binaries
    fi
    make -C $PROJECT_ROOT generate-attribution
    for summary in $PROJECT_ROOT/_output/**/summary.txt; do
        sed -i "s/+.*=/ =/g" $summary
        awk -F" =\> " '{ count[$1]+=$2} END { for (item in count) printf("%s => %d\n", item, count[item]) }' \
            $summary _output/total_summary.txt | sort > _output/total_summary.tmp && mv _output/total_summary.tmp _output/total_summary.txt
    done    
    make -C $PROJECT_ROOT clean
}


build::attribution::generate
