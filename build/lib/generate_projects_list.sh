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
set -o pipefail

BASE_DIRECTORY="${1?Specify first argument - Base directory of build-tooling repo}"
UPSTREAM_PROJECTS_FILE="$BASE_DIRECTORY/UPSTREAM_PROJECTS.yaml"

YQ_LATEST_RELEASE_URL="https://github.com/mikefarah/yq/releases/latest"
CURRENT_YQ_VERSION=$(yq -V | awk '{print $NF}')
CURRENT_YQ_MAJOR_VERSION=${CURRENT_YQ_VERSION:1:1}
LATEST_YQ_VERSION=$(curl --retry 5 -fIsS $YQ_LATEST_RELEASE_URL | grep "location:" | awk -F/ '{print $NF}')
LATEST_YQ_MAJOR_VERSION=${LATEST_YQ_VERSION:1:1}
if [ $CURRENT_YQ_MAJOR_VERSION -lt $LATEST_YQ_MAJOR_VERSION ]; then
    echo "Current yq major version v$CURRENT_YQ_MAJOR_VERSION.x is older than the latest (v$LATEST_YQ_MAJOR_VERSION.x)."
    echo "Please install the latest version of yq from $YQ_LATEST_RELEASE_URL"
    exit 1
fi

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

cd $MAKE_ROOT

declare -i org_count=0 # Counter variable for each org

yq eval --null-input ".projects = []" > $UPSTREAM_PROJECTS_FILE # Creates an empty YAML array

# Iterating over orgs under projects folder
for org_path in projects/*; do
    repos=() # Empty array for repos in the org
    org=$(cut -d/ -f2 <<< $org_path)
    declare -i repo_count=0
    for repo_path in projects/$org/*; do
        repo="$(cut -d/ -f3 <<< $repo_path)"
        if [ "$org" = "aws" ] & [[ $repo =~ "eks-anywhere" ]]; then # Ignore self-referential repos
            continue
        fi
        if [ "$org" = "kubernetes-sigs" ] & [[ "$repo" =~ "metrics-server" ]]; then # Ignore helm builds backed by eks-d images to reduce toil
            continue
        fi
        if curl --retry 5 -fIsS https://github.com/$org/$repo &> /dev/null; then # Check if org/repo combination is a Github repo
            repos+=("$repo")
        fi
    done
    if [ ${#repos[@]} -gt 0 ]; then
        yq eval -i -P ".projects += [{\"org\": \"$org\", \"repos\": []}]" $UPSTREAM_PROJECTS_FILE # Add an entry for this org in the projects array
        for repo in "${repos[@]}"; do
            yq eval -i -P ".projects[$org_count].repos += [{\"name\": \"$repo\", \"versions\": []}]" $UPSTREAM_PROJECTS_FILE # Add each repo to the repos array
            git_tags=$(find projects/$org/$repo -type f -name "GIT_TAG" | sort)
            for file in $git_tags; do
                tag=$(cat $file)
                golang_version="N/A"
                golang_version_file=$(dirname $file)/GOLANG_VERSION
                if [ -f $golang_version_file ]; then
                    golang_version="$(cat $golang_version_file)"
                fi
                if [[ $tag =~ ^[0-9a-f]{7,40}$ ]]; then
                    yq eval -i -P ".projects[$org_count].repos[$repo_count].versions += [{\"commit\": \"$tag\", \"go_version\": \"$golang_version\"}]" $UPSTREAM_PROJECTS_FILE
                else
                    yq eval -i -P ".projects[$org_count].repos[$repo_count].versions += [{\"tag\": \"$tag\", \"go_version\": \"$golang_version\"}]" $UPSTREAM_PROJECTS_FILE
                fi
            done
            repo_count+=1
        done
        org_count+=1
    fi
done

HEAD_COMMENT=$(cat $BASE_DIRECTORY/hack/boilerplate.yq.txt)
yq eval -i ". headComment=\"$HEAD_COMMENT\"" $UPSTREAM_PROJECTS_FILE # Add a header comment with license verbiage and no-edit warning
