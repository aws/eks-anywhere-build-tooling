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

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

source "${SCRIPT_ROOT}/eksd_releases.sh"

FAILED=0
for release_branch in $(cat $SCRIPT_ROOT/../../release/SUPPORTED_RELEASE_BRANCHES); do
    YAML_URL=$(build::eksd_releases::get_release_yaml_url $release_branch)
    HTTP_CODE=$(curl -I -L -s -o /dev/null -w "%{http_code}" $YAML_URL)
    if [[ "$HTTP_CODE" != "200" ]]; then
        echo "EKS-D manifest does not exist: $YAML_URL"
        FAILED=1
        continue
    fi   

    build::eksd_releases::load_release_yaml $release_branch > /dev/null
    COMPONENT_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $release_branch)
    COMPONENT_VERSION_FROM_YAML=$(yq e ".releases[] | select(.branch==\"${release_branch}\").kubeVersion" $SCRIPT_ROOT/../../EKSD_LATEST_RELEASES)

    if [[ "$COMPONENT_VERSION" != "$COMPONENT_VERSION_FROM_YAML" ]]; then
        echo "kubeVersion: $COMPONENT_VERSION_FROM_YAML does not match version from EKS-D release manifest: $COMPONENT_VERSION"
        FAILED=1
    fi
done

if [[ "$FAILED" == "0" ]]; then
    echo "All EKS-D versions validated!"
fi
exit $FAILED
