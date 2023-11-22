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

# set -x
set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

RETURN=0
SED=$(build::find::gnu_variant_on_mac sed)
PROJECT=${1:-}

SUPPORTED_RELEASE_BRANCHES=($(cat $SCRIPT_ROOT/../../release/SUPPORTED_RELEASE_BRANCHES))
LATEST_RELEASE_BRANCH=${SUPPORTED_RELEASE_BRANCHES[-1]}

function check_and_update_readme() {
    local -r git_tag_file=$1
    local -r project_path=$2
    local -r release_branched=$3

    if [ ! -f $git_tag_file ]; then
        return
    fi
    VERSION="$(cat $git_tag_file | $SED "s,-,--,g")"
    README="$project_path/README.md"
    if [ ! -f $README ]
    then
        echo "Missing file $README"
        continue
    fi

    EXPECTED_VERSION="img.shields.io/badge/version-$VERSION"
    if [ "$release_branched" = "true" ]; then
        RELEASE_BRANCH=$(basename $(dirname $git_tag_file) | $SED "s,-,--,g")
        EXPECTED_VERSION=img.shields.io/badge/$RELEASE_BRANCH%20version-$VERSION
    fi
    if grep -l "$EXPECTED_VERSION" ${README} >/dev/null
    then
        echo "Actual version in README matches expected version"
        return
    fi
    VERSION_SEARCH_PATTERN="img.shields.io/badge/version"
    if [ "$release_branched" = "true" ]; then
        VERSION_SEARCH_PATTERN=img.shields.io/badge/$RELEASE_BRANCH%20version
    fi
    if ! grep $VERSION_SEARCH_PATTERN ${README} >/dev/null
    then
        echo "Did not find version  in $README"
        return
    fi

    ACTUAL_VERSION=$(grep $VERSION_SEARCH_PATTERN ${README} | $SED -e "s,.*$VERSION_SEARCH_PATTERN-,," -e 's/-blue).*$//')
    echo "Version mismatch in $README - expected $VERSION, actual $ACTUAL_VERSION"
    $SED -i -e "s/$ACTUAL_VERSION/$VERSION/" $README
}

PROJECTS="projects/*/*"
if [ -n "$PROJECT" ]; then
    PROJECTS=projects/$PROJECT
fi

for PROJECT in $PROJECTS
do
    RELEASE_BRANCHED=false
    if [ -f "$PROJECT/GIT_TAG" ] & [ -f "$PROJECT/$LATEST_RELEASE_BRANCH/GIT_TAG" ]; then
        RELEASE_BRANCHED=true
    fi
    if [ "$RELEASE_BRANCHED" = "true" ]; then
        for branch in ${SUPPORTED_RELEASE_BRANCHES[@]}; do
            check_and_update_readme $PROJECT/$branch/GIT_TAG $PROJECT true
        done
    else
        check_and_update_readme "$PROJECT/GIT_TAG" $PROJECT false
    fi
done
