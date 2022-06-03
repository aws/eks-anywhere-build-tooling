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

BASE_DIRECTORY="${1?Specify first argument - Base directory of build-tooling repo}"
ALL_PROJECTS="${2?Specify second argument - All projects in repo}"

STAGING_BUILDSPEC_FILE="$BASE_DIRECTORY/release/staging-build.yml"

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

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

cd $MAKE_ROOT

yq eval --null-input '.batch={"fast-fail":true,"build-graph":[]}' > $STAGING_BUILDSPEC_FILE # Creates an empty YAML array

PROJECTS=(${ALL_PROJECTS// / })
for project in "${PROJECTS[@]}"; do
    org=$(cut -d_ -f1 <<< $project)
    repo=$(cut -d_ -f2 <<< $project)

    PROJECT_PATH=$MAKE_ROOT/projects/$org/$repo

    # TODO: refactor use of release_branch to get git_tag and golang_version in makefile, we should be able to push this to common.mk and avoid needing to pass it here
    if [[ "true" == "$(make --no-print-directory -C $PROJECT_PATH var-value-EXCLUDE_FROM_STAGING_BUILDSPEC RELEASE_BRANCH=1-20)" ]]; then
        continue
    fi

    IDENTIFIER="${org//-/_}_${repo//-/_}"

    echo "Adding: $IDENTIFIER"

    DEPEND_ON=""
    PROJECT_DEPENDENCIES=$(make --no-print-directory -C $PROJECT_PATH var-value-PROJECT_DEPENDENCIES RELEASE_BRANCH=1-20)
    if [ -n "$PROJECT_DEPENDENCIES" ]; then
        DEPS=(${PROJECT_DEPENDENCIES// / })
        for dep in "${DEPS[@]}"; do
            DEP_PRODUCT="$(cut -d/ -f1 <<< $dep)"
            DEP_ORG="$(cut -d/ -f2 <<< $dep)"
            DEP_REPO="$(cut -d/ -f3 <<< $dep)"
            if [[ "$DEP_PRODUCT" == "eksd" ]]; then
                continue
            fi
            DEPEND_ON+="\"${DEP_ORG//-/_}_${DEP_REPO//-/_}\","

            if [ ! -d $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO ]; then
                echo "Non-existent project dependency: $dep!!!"
                exit 1
            fi
        done
    fi

    if [ -n "$DEPEND_ON" ]; then
        DEPEND_ON="\"depend-on\":[${DEPEND_ON%?}],"
    fi

    CLONE_URL=""
    if [[ "true" != "$(make --no-print-directory -C $PROJECT_PATH var-value-REPO_NO_CLONE RELEASE_BRANCH=1-20)" ]]; then
        REPO=$(make --no-print-directory -C $PROJECT_PATH var-value-CLONE_URL AWS_REGION=us-west-2 CODEBUILD_CI=true RELEASE_BRANCH=1-20)
        CLONE_URL=",\"CLONE_URL\":\"$REPO\""
    fi

    BUILDSPECS=$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPECS RELEASE_BRANCH=1-20)
    SPECS=(${BUILDSPECS// / })
    for buildspec in "${SPECS[@]}"; do
        if [[ "${#SPECS[@]}" != "1" ]]; then
            BUILDSPEC_NAME=$(basename $buildspec .yml)
            IDENTIFIER=${org//-/_}_${repo//-/_}_${BUILDSPEC_NAME//-/_}
        fi
        yq eval -i -P \
            ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",\"buildspec\":\"$buildspec\",$DEPEND_ON\"env\":{\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL}}}]" \
            $STAGING_BUILDSPEC_FILE 
    done
done

HEAD_COMMENT=$(cat $BASE_DIRECTORY/hack/boilerplate.yq.txt)
yq eval -i ". headComment=\"$HEAD_COMMENT\"" $STAGING_BUILDSPEC_FILE # Add a header comment with license verbiage and no-edit warning
