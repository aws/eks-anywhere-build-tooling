#!/usr/bin/env bash

set -o errexit
set -o pipefail

BASE_DIRECTORY="${1?Specify first argument - Base directory of build-tooling repo}"
UPSTREAM_PROJECTS_FILE=$BASE_DIRECTORY/UPSTREAM_PROJECTS.yaml

YQ_LATEST_RELEASE_URL="https://github.com/mikefarah/yq/releases/latest"
CURRENT_YQ_VERSION=$(yq -V | awk '{print $NF}')
CURRENT_YQ_MAJOR_VERSION=${CURRENT_YQ_VERSION:0:1}
LATEST_YQ_VERSION=$(curl -fIsS $YQ_LATEST_RELEASE_URL | grep "location:" | awk -F/ '{print $NF}')
LATEST_YQ_MAJOR_VERSION=${LATEST_YQ_VERSION:1:1}
if [ $CURRENT_YQ_MAJOR_VERSION -lt $LATEST_YQ_MAJOR_VERSION ]; then
    echo "Current yq major version v$CURRENT_YQ_MAJOR_VERSION.x is older than the latest (v$LATEST_YQ_MAJOR_VERSION.x)."
    echo "Please install the latest version of yq from $YQ_LATEST_RELEASE_URL"
    exit 1
fi

declare -i org_count=0 # Counter variable for each org

yq eval --null-input ".projects = []" > $UPSTREAM_PROJECTS_FILE # Creates an empty YAML array

# Iterating over orgs under projects folder
for org_path in projects/*; do
    repos=() # Empty array for repos in the org
    org=$(cut -d/ -f2 <<< $org_path)
    for repo_path in projects/$org/*; do
        repo="$(cut -d/ -f3 <<< $repo_path)"
        if [ "$org" = "aws" ] & [[ $repo =~ "eks-anywhere" ]]; then # Ignore self-referential repos
            continue
        fi
        if curl -fIsS https://github.com/$org/$repo &> /dev/null; then # Check if org/repo combination is a Github repo
            repos+=("$repo")
        fi
    done
    if [ ${#repos[@]} -gt 0 ]; then
        yq eval -i -P ".projects += [{\"org\": \"$org\", \"repos\": []}]" $UPSTREAM_PROJECTS_FILE # Add an entry for this org in the projects array
        for repo in "${repos[@]}"; do
            yq eval -i -P ".projects[$org_count].repos += [{\"name\": \"$repo\"}]" $UPSTREAM_PROJECTS_FILE # Add each repo to the repos array
        done
        org_count+=1
    fi
done
HEAD_COMMENT=$(cat $BASE_DIRECTORY/hack/boilerplate.yq.txt)
yq eval -i ". headComment=\"$HEAD_COMMENT\"" $UPSTREAM_PROJECTS_FILE # Add a header comment with license verbiage and no-edit warning
yq eval $UPSTREAM_PROJECTS_FILE # Print generated YAML

echo "Contents written to $UPSTREAM_PROJECTS_FILE"
