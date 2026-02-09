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

set -x
set -o errexit
set -o pipefail

HELM_DIRECTORY="$1"
GIT_TAG="$2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Running helm lint"
helm lint ${HELM_DIRECTORY}

CHART_VERSION="$(yq eval '.version' ${PROJECT_DIR}/chart/Chart.yaml)"
echo "Verifying sedfile.template version matches the chart version in Chart.yaml"
SED_FILE_VERSION="$(grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" ${PROJECT_DIR}/helm/sedfile.template)"
if [[ ${SED_FILE_VERSION} != ${CHART_VERSION} ]]; then
    echo "❌ version in sedfile.template does not match the Chart version"
    echo "SED_FILE_VERSION=${SED_FILE_VERSION}, Chart version=${CHART_VERSION}"
    exit 1
fi

# Verify all CRDs have required annotations
echo "Verifying CRDs have required annotations..."

# CRDs that should NOT have clusterctl labels (to prevent move during cluster migration)
SKIP_CLUSTERCTL_LABELS="tinkerbell.org_workflows.yaml bmc.tinkerbell.org_jobs.yaml bmc.tinkerbell.org_tasks.yaml"

for crd_file in "${HELM_DIRECTORY}"/templates/*.yaml; do
    filename=$(basename "$crd_file")
    
    # Check helm.sh/resource-policy annotation
    if ! yq -e '.metadata.annotations["helm.sh/resource-policy"] == "keep"' "$crd_file" > /dev/null 2>&1; then
        echo "❌ ${filename} missing helm.sh/resource-policy: keep annotation"
        exit 1
    fi
    
    # Check clusterctl labels (skip for transient CRDs)
    if [[ ! " ${SKIP_CLUSTERCTL_LABELS} " =~ " ${filename} " ]]; then
        if ! yq -e '.metadata.labels["clusterctl.cluster.x-k8s.io"] == ""' "$crd_file" > /dev/null 2>&1; then
            echo "❌ ${filename} missing clusterctl.cluster.x-k8s.io label"
            exit 1
        fi
    fi
    
    echo "✅ ${filename} has required annotations/labels"
done

echo "✅ All verifications passed"
