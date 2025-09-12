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
DEFAULT_BUILDSPEC_FILE="${4}"
SKIP_DEPEND_ON="${5:-false}"
EXCLUDE_VAR="${6:-EXCLUDE_FROM_STAGING_BUILDSPEC}"
BUILDSPECS_VAR="${7:-BUILDSPECS}"
FAST_FAIL="${8:-true}"
FINAL_STAGE_BUILDSPECS_CSV="${9:-}"
NO_DEPS_FOR_FINAL_STAGE="${10:-false}"

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
yq eval --null-input ".batch={\"fast-fail\":$FAST_FAIL,\"build-graph\":[]} | . *= load(\"$DEFAULT_BUILDSPEC_FILE\")" > $STAGING_BUILDSPEC_FILE # Creates an empty YAML array

# TODO: refactor use of release_branch to get git_tag and golang_version in makefile, we should be able to push this to common.mk and avoid needing to pass it here
RELEASE_BRANCH=$(build::eksd_releases::get_release_branch)

function make_var() {
    make --no-print-directory -C $1 "$(echo "$2" | sed 's/[^ ]* */var-value-&/g')" CODEBUILD_CI=true 2>/dev/null
}

PROJECTS=(${ALL_PROJECTS// / })
ALL_PROJECT_IDS=""
for project in "${PROJECTS[@]}"; do
    org=$(cut -d_ -f1 <<< $project)
    repo=$(cut -d_ -f2- <<< $project)

    PROJECT_PATH=$MAKE_ROOT/projects/$org/$repo

    if [[ "true" == "$(make_var $PROJECT_PATH $EXCLUDE_VAR)" ]]; then
        continue
    fi

    CLONE_URL=""
    if [[ "true" != "$(make_var $PROJECT_PATH REPO_NO_CLONE)" ]]; then
        REPO=$(make_var $PROJECT_PATH CLONE_URL AWS_REGION=us-west-2)
        CLONE_URL=",\"CLONE_URL\":\"$REPO\""
    fi

    BUILDSPECS=$(make_var $PROJECT_PATH $BUILDSPECS_VAR)
    SPECS=(${BUILDSPECS// / })
    PREVIOUS_SPEC_IDENTIFIERS=""
    for (( i=0; i < ${#SPECS[@]}; i++ )); do
        IDENTIFIER="${org//-/_}_${repo//-/_}"

        buildspec=${SPECS[$i]}
        buildspec_field="\"buildspec\":\"$buildspec\","
        if [[ $(realpath $buildspec) == $(realpath $DEFAULT_BUILDSPEC_FILE) ]]; then
            buildspec_field=""
        fi

        DEPEND_ON=""
        # something other than empty string since some overrides are empty strings
        PROJECT_DEPENDENCIES="false"
        for var in "BUILDSPEC_$((( $i + 1 )))_DEPENDS_ON_OVERRIDE" "BUILDSPEC_DEPENDS_ON_OVERRIDE"; do
            BUILDSPEC_DEPENDS_ON="$(make_var $PROJECT_PATH $var)"
            HARDCODED_DEP="false"
            if [[ "none" = "$BUILDSPEC_DEPENDS_ON" ]]; then
                PROJECT_DEPENDENCIES=""
                break
            elif [[ -n "$BUILDSPEC_DEPENDS_ON" ]] && [[ $BUILDSPEC_DEPENDS_ON  == BUILDSPEC_* ]]; then
                # TODO: In the makefiles we will use `BUILDSPEC_1` as if you can control the specific buildspec
                # to depend on. This code assumes its always the previous
                HARDCODED_DEP="true"
                PROJECT_DEPENDENCIES=$PREVIOUS_SPEC_IDENTIFIERS
                break
            elif [[ -n "$BUILDSPEC_DEPENDS_ON" ]]; then
                HARDCODED_DEP="true"
                PROJECT_DEPENDENCIES=$BUILDSPEC_DEPENDS_ON
                break
            fi
        done

        PREVIOUS_SPEC_IDENTIFIERS=""

        BUILDSPEC_NAME=$(basename $buildspec .yml)
        BUILDSPEC_IDENTIFIER_OVERRIDE="$(make_var $PROJECT_PATH BUILDSPEC_$((( $i + 1 )))_IDENTIFIER_OVERRIDE)"
        if [[ -n "$BUILDSPEC_IDENTIFIER_OVERRIDE" ]]; then
            IDENTIFIER="$BUILDSPEC_IDENTIFIER_OVERRIDE"
        elif [[ "${BUILDSPEC_NAME}" != *buildspec* ]]; then
            IDENTIFIER+="_${BUILDSPEC_NAME//-/_}"
        fi

        echo "Adding: $IDENTIFIER"

        if [ "$PROJECT_DEPENDENCIES" = "false" ]; then
            PROJECT_DEPENDENCIES=$(make_var $PROJECT_PATH PROJECT_DEPENDENCIES)
        fi

        if [ -n "$PROJECT_DEPENDENCIES" ] && [[ "${buildspec}" = "buildspecs/combine-images.yml" || "$SKIP_DEPEND_ON" != "true" || "${buildspec}" = "projects/kubernetes-sigs/kind/buildspecs/node-images.yml" ]]; then
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
                DEP_RELEASE_BRANCH="$(cut -d/ -f4 <<< $dep)"

                if [ ! -d $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO ]; then
                    echo "Non-existent project dependency: $dep!!!"
                    exit 1
                fi

                DEP_IDENTIFIER=${DEP_ORG//-/_}_${DEP_REPO//-/_}

                if [ -n "${DEP_RELEASE_BRANCH}" ]; then
                    DEP_IDENTIFIER=${DEP_ORG//-/_}_${DEP_REPO//-/_}_${DEP_RELEASE_BRANCH//[-\/]/_}
                fi
                # if dep is split by binary platform, append platform/arch
                DEP_BUILDSPEC_VARS="$(make_var $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO BUILDSPEC_VARS_KEYS)"
                DEP_BUILDSPECS="$(make_var $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO BUILDSPECS)"
                DEP_HAS_RELEASE_BRANCHES="$(make_var $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO HAS_RELEASE_BRANCHES)"
                
                if [[ "$DEP_BUILDSPEC_VARS" == *"BINARY_PLATFORMS"* ]]; then
                    # Check if dependency has release branches and binaries are release branched
                    if [[ "$DEP_HAS_RELEASE_BRANCHES" == "true" ]]; then
                        # Get supported K8s versions (release branches) for the dependency
                        DEP_SUPPORTED_K8S_VERSIONS="$(make_var $MAKE_ROOT/projects/$DEP_ORG/$DEP_REPO SUPPORTED_K8S_VERSIONS)"
                        DEP_K8S_VERSIONS=(${DEP_SUPPORTED_K8S_VERSIONS// / })
                        
                        # Generate dependencies for each release branch + platform combination
                        for k8s_version in "${DEP_K8S_VERSIONS[@]}"; do
                            DEP_IDENTIFIER_WITH_BRANCH=${DEP_ORG//-/_}_${DEP_REPO//-/_}_${k8s_version//[-\/]/_}
                            DEPEND_ON+="\"${DEP_IDENTIFIER_WITH_BRANCH}_linux_amd64\",\"${DEP_IDENTIFIER_WITH_BRANCH}_linux_arm64\","
                        done
                    else
                        # Original logic for dependencies without release branches
                        DEPEND_ON+="\"${DEP_IDENTIFIER}_linux_amd64\",\"${DEP_IDENTIFIER}_linux_arm64\","
                    fi
                elif [[ $DEP_BUILDSPECS == *combine-images.yml* ]]; then
                    DEPEND_ON+="\"${DEP_IDENTIFIER}_combine_images\","
                else
                    DEPEND_ON+="\"${DEP_IDENTIFIER}\","
                fi
            done
        fi

        if [ -n "$DEPEND_ON" ]; then
            DEPEND_ON="\"depend-on\":[${DEPEND_ON%?}],"
        fi
        BUILDSPEC_VARS_KEYS=$(make_var $PROJECT_PATH BUILDSPEC_VARS_KEYS)
        if [[ -z "$BUILDSPEC_VARS_KEYS" ]]; then
            BUILDSPEC_VARS_KEYS=$(make_var $PROJECT_PATH BUILDSPEC_$((( $i + 1 )))_VARS_KEYS)
        fi

        BUILDSPEC_PLATFORM=$(make_var $PROJECT_PATH BUILDSPEC_$((( $i + 1 )))_PLATFORM)
        if [[ -z "$BUILDSPEC_PLATFORM" ]]; then
            BUILDSPEC_PLATFORM=$(make_var $PROJECT_PATH BUILDSPEC_PLATFORM)
        fi

        BUILDSPEC_COMPUTE_TYPE=$(make_var $PROJECT_PATH BUILDSPEC_$((( $i + 1 )))_COMPUTE_TYPE)
        if [[ -z "$BUILDSPEC_COMPUTE_TYPE" ]]; then
            BUILDSPEC_COMPUTE_TYPE=$(make_var $PROJECT_PATH BUILDSPEC_COMPUTE_TYPE)
        fi

        ARCH_TYPE="\"type\":\"$BUILDSPEC_PLATFORM\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","

        BUILDSPEC_VARS_VALUES=""
        if [[ "$BUILDSPECS_VAR" == "CHECKSUMS_BUILDSPECS" ]]; then
            if [[ "true" == "$(make_var $PROJECT_PATH HAS_RELEASE_BRANCHES)" ]]; then
                if [[ "false" == "$(make_var $PROJECT_PATH BINARIES_ARE_RELEASE_BRANCHED)" ]]; then
                    BUILDSPEC_VARS_KEYS=""
                else
                    # Preserve BINARY_PLATFORMS if it exists in the original BUILDSPEC_VARS_KEYS
                    if [[ "$BUILDSPEC_VARS_KEYS" == *"BINARY_PLATFORMS"* ]]; then
                        BUILDSPEC_VARS_KEYS="RELEASE_BRANCH BINARY_PLATFORMS"
                        BUILDSPEC_VARS_VALUES="SUPPORTED_K8S_VERSIONS BINARY_PLATFORMS"
                    else
                        BUILDSPEC_VARS_KEYS="RELEASE_BRANCH"
                        BUILDSPEC_VARS_VALUES="SUPPORTED_K8S_VERSIONS"
                    fi
                fi
            elif [[ "${BUILDSPEC_VARS_KEYS}" = "IMAGE_PLATFORMS" ]]; then
                BUILDSPEC_VARS_KEYS=""
            fi
        fi

        if [[ "$BUILDSPECS_VAR" == "UPGRADE_BUILDSPECS" ]]; then
            if [[ "${IDENTIFIER}" = "kubernetes_sigs_image_builder" ]]; then
                BUILDSPEC_VARS_KEYS=""
            elif [[ "true" == "$(make_var $PROJECT_PATH HAS_RELEASE_BRANCHES)" ]]; then
                if [[ "false" == "$(make_var $PROJECT_PATH BINARIES_ARE_RELEASE_BRANCHED)" ]]; then
                    BUILDSPEC_VARS_KEYS=""
                else
                    # Preserve BINARY_PLATFORMS if it exists in the original BUILDSPEC_VARS_KEYS
                    if [[ "$BUILDSPEC_VARS_KEYS" == *"BINARY_PLATFORMS"* ]]; then
                        BUILDSPEC_VARS_KEYS="RELEASE_BRANCH BINARY_PLATFORMS"
                        BUILDSPEC_VARS_VALUES="SUPPORTED_K8S_VERSIONS BINARY_PLATFORMS"
                    else
                        BUILDSPEC_VARS_KEYS="RELEASE_BRANCH"
                        BUILDSPEC_VARS_VALUES="SUPPORTED_K8S_VERSIONS"
                    fi
                fi
            fi
        fi

        if [[ -n "$BUILDSPEC_VARS_KEYS" ]]; then
            KEYS=(${BUILDSPEC_VARS_KEYS// / })

            if [[ -z "$BUILDSPEC_VARS_VALUES" ]]; then
                BUILDSPEC_VARS_VALUES=$(make_var $PROJECT_PATH BUILDSPEC_VARS_VALUES)
                if [[ -z "$BUILDSPEC_VARS_VALUES" ]]; then
                    BUILDSPEC_VARS_VALUES=$(make_var $PROJECT_PATH BUILDSPEC_$((( $i + 1 )))_VARS_VALUES)
                fi
            fi
            VARS=(${BUILDSPEC_VARS_VALUES// / })

            # Note: only support 1 or 2 vars for now since that is all we need for kind + image-builder
            if [ ${#VARS[@]} -eq 1 ]; then
                VALUES_1=$(make_var $PROJECT_PATH ${VARS[0]})
                ARR_1=(${VALUES_1// / })

                for val1 in "${ARR_1[@]}"; do
                    IDENTIFIER_WITH_VAL=${IDENTIFIER}_${val1//[-\/]/_}

                    # If building on one binary platform assume we want to run on a specific arch instance
                    ARCH_TYPE="\"type\":\"$BUILDSPEC_PLATFORM\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                    if [ "${KEYS[0]}" = "BINARY_PLATFORMS" ] || [ "${KEYS[0]}" = "IMAGE_PLATFORMS" ]; then
                        if [ "${val1}" = "linux/amd64" ]; then
                            ARCH_TYPE="\"type\":\"LINUX_CONTAINER\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                        else
                            ARCH_TYPE="\"type\":\"ARM_CONTAINER\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                        fi
                    fi

                    EXTRA_VARS=""
                    if [ "${KEYS[0]}" = "IMAGE_PLATFORMS" ]; then
                        EXTRA_VARS+=",\"BINARY_PLATFORMS\":\"${val1}\",\"IMAGE_TAG_SUFFIX\":\"-${val1#linux/}\""
                        HAS_HELM_CHART=$(make_var $PROJECT_PATH HAS_HELM_CHART)
                        if [ "${HAS_HELM_CHART}" = "true" ]; then
                            EXTRA_VARS+=",\"HAS_HELM_CHART\":\"false\""
                        fi
                    fi

                    ALL_PROJECT_IDS+="\"$IDENTIFIER_WITH_VAL\","
                    PREVIOUS_SPEC_IDENTIFIERS+="$IDENTIFIER_WITH_VAL "
                    yq eval -i -P \
                        ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER_WITH_VAL\",$buildspec_field$DEPEND_ON\"env\":{$ARCH_TYPE\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL,\"${KEYS[0]}\":\"$val1\"$EXTRA_VARS}}}]" \
                        $STAGING_BUILDSPEC_FILE

                done
            else
                VALUES_1=$(make_var $PROJECT_PATH ${VARS[0]})
                VALUES_2=$(make_var $PROJECT_PATH ${VARS[1]})

                ARR_1=(${VALUES_1// / })
                ARR_2=(${VALUES_2// / })
                for val1 in "${ARR_1[@]}"; do
                    for val2 in "${ARR_2[@]}"; do
                        IDENTIFIER_WITH_VAL=${IDENTIFIER}_${val1//[-\/]/_}_${val2//[-\/]/_}
                        # TODO: revisit this to make it more dynamic if other projects need it in the future
                        EXTRA_VARS=""
                        if [[ "$IDENTIFIER" =~ "kubernetes_sigs_image_builder_bottlerocket" ]]; then
                            EXTRA_VARS+=",\"IMAGE_OS_VERSION\":\"1\""
                        fi
                        if [ "${KEYS[1]}" = "IMAGE_PLATFORMS" ]; then
                            EXTRA_VARS+=",\"BINARY_PLATFORMS\":\"${val2}\",\"IMAGE_TAG_SUFFIX\":\"-${val2#linux/}\""
                        fi
                        # If building on one binary platform assume we want to run on a specific arch instance
                        ARCH_TYPE="\"type\":\"$BUILDSPEC_PLATFORM\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                        if [ "${KEYS[1]}" = "BINARY_PLATFORMS" ] || [ "${KEYS[1]}" = "IMAGE_PLATFORMS" ]; then
                            if [ "${val2}" = "linux/amd64" ]; then
                                ARCH_TYPE="\"type\":\"LINUX_CONTAINER\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                            else
                                ARCH_TYPE="\"type\":\"ARM_CONTAINER\",\"compute-type\":\"$BUILDSPEC_COMPUTE_TYPE\","
                            fi
                        fi
                        ALL_PROJECT_IDS+="\"$IDENTIFIER_WITH_VAL\","
                        PREVIOUS_SPEC_IDENTIFIERS+="$IDENTIFIER_WITH_VAL "
                        yq eval -i -P \
                            ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER_WITH_VAL\",$buildspec_field$DEPEND_ON\"env\":{$ARCH_TYPE\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL,\"${KEYS[0]}\":\"$val1\",\"${KEYS[1]}\":\"$val2\"$EXTRA_VARS}}}]" \
                            $STAGING_BUILDSPEC_FILE
                    done
                done
            fi
        else
            ALL_PROJECT_IDS+="\"$IDENTIFIER\","
            PREVIOUS_SPEC_IDENTIFIERS+="$IDENTIFIER "
            yq eval -i -P \
                ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",$buildspec_field$DEPEND_ON\"env\":{$ARCH_TYPE\"variables\":{\"PROJECT_PATH\": \"projects/$org/$repo\"$CLONE_URL}}}]" \
                $STAGING_BUILDSPEC_FILE
        fi
    done
done

if [ -n "${FINAL_STAGE_BUILDSPECS_CSV}" ]; then
    IFS=',' read -r -a FINAL_STAGE_BUILDSPECS <<< "$FINAL_STAGE_BUILDSPECS_CSV"
    for i in "${!FINAL_STAGE_BUILDSPECS[@]}";
        do
            IDENTIFIER="final_stage"
            if [[ "$i" -gt 0 ]]; then
                IDENTIFIER="final_stage_$((i+1))"
            fi
            ARCH_TYPE="\"type\":\"ARM_CONTAINER\",\"compute-type\":\"BUILD_GENERAL1_SMALL\""
            DEPEND_ON=",\"depend-on\":[${ALL_PROJECT_IDS%?}]"
            if [ "$NO_DEPS_FOR_FINAL_STAGE" = "true" ]; then
                DEPEND_ON=""
            fi
            yq eval -i -P \
                ".batch.build-graph += [{\"identifier\":\"$IDENTIFIER\",\"buildspec\":\"${FINAL_STAGE_BUILDSPECS[$i]}\",\"env\":{$ARCH_TYPE}$DEPEND_ON}]" \
                $STAGING_BUILDSPEC_FILE
        done
fi

HEAD_COMMENT=$(cat $BASE_DIRECTORY/hack/boilerplate.yq.txt)
yq eval -i ". headComment=\"$HEAD_COMMENT\"" $STAGING_BUILDSPEC_FILE # Add a header comment with license verbiage and no-edit warning

if [[ "${#PROJECTS[@]}" = "1" ]]; then
    # if there is only one project we do not want project_path and clone_url to be set since it will be set at the codebuild level
    yq -i 'del(.batch.build-graph.[].env.variables.PROJECT_PATH)' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.batch.build-graph.[].env.variables.CLONE_URL)' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.. | select(tag == "!!map" and length == 0))' $STAGING_BUILDSPEC_FILE
    yq -i 'del(.. | select(tag == "!!map" and length == 0))' $STAGING_BUILDSPEC_FILE
fi

$SCRIPT_ROOT/validate_release_buildspecs.sh "$STAGING_BUILDSPEC_FILE"
