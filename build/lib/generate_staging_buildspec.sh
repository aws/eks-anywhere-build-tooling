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
STAGING_BUILDSPEC_FILE="${3}"
SKIP_DEPEND_ON="${4:-false}"
EXCLUDE_VAR="${5:-EXCLUDE_FROM_STAGING_BUILDSPEC}"
BUILDSPECS_VAR="${6:-BUILDSPECS}"
FAST_FAIL="${7:-true}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

YQ_LATEST_RELEASE_URL="https://github.com/mikefarah/yq/releases/latest"
CURRENT_YQ_VERSION=$(yq -V | awk '{print $NF}')
CURRENT_YQ_MAJOR_VERSION=${CURRENT_YQ_VERSION:1:1}
LATEST_YQ_VERSION=$(curl -fIsS $YQ_LATEST_RELEASE_URL | grep "location:" | awk -F/ '{print $NF}')
LATEST_YQ_MAJOR_VERSION=${LATEST_YQ_VERSION:1:1}
if [ $CURRENT_YQ_MAJOR_VERSION -lt $LATEST_YQ_MAJOR_VERSION ]; then
    echo "Current yq major version v$CURRENT_YQ_MAJOR_VERSION.x is older than the latest (v$LATEST_YQ_MAJOR_VERSION.x)."
    echo "Please install the latest version of yq from $YQ_LATEST_RELEASE_URL"
    exit 1
fi

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

cd $MAKE_ROOT

mkdir -p $(dirname $STAGING_BUILDSPEC_FILE)
yq eval --null-input ".batch={\"fast-fail\":$FAST_FAIL,\"build-graph\":[]}" > $STAGING_BUILDSPEC_FILE # Creates an empty YAML array

PROJECTS=(${ALL_PROJECTS// / })
for project in "${PROJECTS[@]}"; do
    org=$(cut -d_ -f1 <<< $project)
    repo=$(cut -d_ -f2- <<< $project)

    PROJECT_PATH=$MAKE_ROOT/projects/$org/$repo

    # TODO: refactor use of release_branch to get git_tag and golang_version in makefile, we should be able to push this to common.mk and avoid needing to pass it here
    RELEASE_BRANCH=$(build::eksd_releases::get_release_branch)
    if [[ "true" == "$(make --no-print-directory -C $PROJECT_PATH var-value-$EXCLUDE_VAR RELEASE_BRANCH=$RELEASE_BRANCH)" ]]; then
        continue
    fi

    CLONE_URL=""
    if [[ "true" != "$(make --no-print-directory -C $PROJECT_PATH var-value-REPO_NO_CLONE RELEASE_BRANCH=$RELEASE_BRANCH)" ]]; then
        REPO=$(make --no-print-directory -C $PROJECT_PATH var-value-CLONE_URL AWS_REGION=us-west-2 CODEBUILD_CI=true RELEASE_BRANCH=$RELEASE_BRANCH)
        CLONE_URL=",\"CLONE_URL\":\"$REPO\""
    fi

    BUILDSPECS=$(make --no-print-directory -C $PROJECT_PATH var-value-$BUILDSPECS_VAR RELEASE_BRANCH=$RELEASE_BRANCH)
    SPECS=(${BUILDSPECS// / })
    for (( i=0; i < ${#SPECS[@]}; i++ )); do
        IDENTIFIER="${org//-/_}_${repo//-/_}"

        buildspec=${SPECS[$i]}

        DEPEND_ON=""
        # something other than empty string since some overrides are empty strings
        PROJECT_DEPENDENCIES="false"
        for var in "BUILDSPEC_DEPENDS_ON_OVERRIDE" "BUILDSPEC_$((( $i + 1 )))_DEPENDS_ON_OVERRIDE"; do
            BUILDSPEC_DEPENDS_ON="$(make --no-print-directory -C $PROJECT_PATH var-value-$var RELEASE_BRANCH=$RELEASE_BRANCH 2>/dev/null)"
            HARDCODED_DEP="false"
            if [[ "none" = "$BUILDSPEC_DEPENDS_ON" ]]; then
                PROJECT_DEPENDENCIES=""
            elif [[ -n "$BUILDSPEC_DEPENDS_ON" ]]; then
                HARDCODED_DEP="true"
                PROJECT_DEPENDENCIES=$BUILDSPEC_DEPENDS_ON
            fi
        done

        BUILDSPEC_IDENTIFIER_OVERRIDE="$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPEC_$((( $i + 1 )))_IDENTIFIER_OVERRIDE RELEASE_BRANCH=$RELEASE_BRANCH 2>/dev/null)"
        if [[ -n "$BUILDSPEC_IDENTIFIER_OVERRIDE" ]]; then
            IDENTIFIER="$BUILDSPEC_IDENTIFIER_OVERRIDE"
        fi   
       
        echo "Adding: $IDENTIFIER"

        if [ "$PROJECT_DEPENDENCIES" = "false" ]; then
            PROJECT_DEPENDENCIES=$(make --no-print-directory -C $PROJECT_PATH var-value-PROJECT_DEPENDENCIES RELEASE_BRANCH=$RELEASE_BRANCH)
        fi

        if [ -n "$PROJECT_DEPENDENCIES" ] && [ "$SKIP_DEPEND_ON" != "true" ]; then
            DEPS=(${PROJECT_DEPENDENCIES// / })
            for dep in "${DEPS[@]}"; do
                if [ "$HARDCODED_DEP" = "true" ]; then
                    DEPEND_ON+="\"$dep\","
                    continue
                fi

                DEP_PRODUCT="$(cut -d/ -f1 <<< $dep)"
                if [[ "$DEP_PRODUCT" == "eksd" ]]; then
                    continue
                fi
                DEP_ORG="$(cut -d/ -f2 <<< $dep)"
                DEP_REPO="$(cut -d/ -f3 <<< $dep)"
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
        BUILDSPEC_VARS_KEYS=$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPEC_VARS_KEYS RELEASE_BRANCH=$RELEASE_BRANCH)
        if [[ -z "$BUILDSPEC_VARS_KEYS" ]]; then
            BUILDSPEC_VARS_KEYS=$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPEC_$((( $i + 1 )))_VARS_KEYS RELEASE_BRANCH=$RELEASE_BRANCH 2>/dev/null)
        fi

        if [[ -n "$BUILDSPEC_VARS_KEYS" ]]; then
            KEYS=(${BUILDSPEC_VARS_KEYS// / })

            BUILDSPEC_VARS_VALUES=$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPEC_VARS_VALUES RELEASE_BRANCH=$RELEASE_BRANCH)
            if [[ -z "$BUILDSPEC_VARS_VALUES" ]]; then
                BUILDSPEC_VARS_VALUES=$(make --no-print-directory -C $PROJECT_PATH var-value-BUILDSPEC_$((( $i + 1 )))_VARS_VALUES RELEASE_BRANCH=$RELEASE_BRANCH 2>/dev/null)
            fi
            VARS=(${BUILDSPEC_VARS_VALUES// / })
            
            # Note: only support 1 or 2 vars for now since that is all we need for kind + image-builder 
            if [ ${#VARS[@]} -eq 1 ]; then
                VALUES_1=$(make --no-print-directory -C $PROJECT_PATH var-value-${VARS[0]} RELEASE_BRANCH=$RELEASE_BRANCH)
                ARR_1=(${VALUES_1// / })
                
                for val1 in "${ARR_1[@]}"; do                
                    BUILDSPEC_NAME=$(basename $buildspec .yml)
                    IDENTIFIER=${org//-/_}_${repo//-/_}_${val1//[-\/]/_}
                    
                    # If building on one binary platform assume we want to run on a specific arch instance
                    ARCH_TYPE=""
                    if [ "${KEYS[0]}" = "BINARY_PLATFORMS" ]; then
                        if [ "${val1}" = "linux/amd64" ]; then
                            ARCH_TYPE="\"type\":\"LINUX_CONTAINER\","
                        else
                            ARCH_TYPE="\"type\":\"ARM_CONTAINER\",\"compute-type\":\"BUILD_GENERAL1_LARGE\","
                        fi
                    fi

                    yq eval -i -P \
                        ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",\"buildspec\":\"$buildspec\",$DEPEND_ON\"env\":{$ARCH_TYPE\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL,\"${KEYS[0]}\":\"$val1\"}}}]" \
                        $STAGING_BUILDSPEC_FILE
                done
            else
                VALUES_1=$(make --no-print-directory -C $PROJECT_PATH var-value-${VARS[0]} RELEASE_BRANCH=$RELEASE_BRANCH)
                VALUES_2=$(make --no-print-directory -C $PROJECT_PATH var-value-${VARS[1]} RELEASE_BRANCH=$RELEASE_BRANCH)

                ARR_1=(${VALUES_1// / })
                ARR_2=(${VALUES_2// / })
                for val1 in "${ARR_1[@]}"; do
                    for val2 in "${ARR_2[@]}"; do
                        BUILDSPEC_NAME=$(basename $buildspec .yml)
                        IDENTIFIER=${org//-/_}_${repo//-/_}_${val1//-/_}_${val2//-/_}_${BUILDSPEC_NAME//-/_}
                        yq eval -i -P \
                            ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",\"buildspec\":\"$buildspec\",$DEPEND_ON\"env\":{\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL,\"${KEYS[0]}\":\"$val1\",\"${KEYS[1]}\":\"$val2\"}}}]" \
                            $STAGING_BUILDSPEC_FILE 
                    done
                done
            fi
        else
            yq eval -i -P \
                ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",\"buildspec\":\"$buildspec\",$DEPEND_ON\"env\":{\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL}}}]" \
                $STAGING_BUILDSPEC_FILE 
        fi
        
    done
done

HEAD_COMMENT=$(cat $BASE_DIRECTORY/hack/boilerplate.yq.txt)
yq eval -i ". headComment=\"$HEAD_COMMENT\"" $STAGING_BUILDSPEC_FILE # Add a header comment with license verbiage and no-edit warning

if [[ "${#PROJECTS[@]}" = "1" ]]; then
    # if there is only one project we do not want project_path and clone_url to be set since it will be set at the codebuild level
    yq -i 'del(.batch.build-graph.[].env.variables.PROJECT_PATH)' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.batch.build-graph.[].env.variables.CLONE_URL)' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.. | select(tag == "!!map" and length == 0))' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.. | select(tag == "!!map" and length == 0))' $STAGING_BUILDSPEC_FILE
fi
