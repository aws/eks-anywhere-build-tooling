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

# This script generates the tinkerbell-crds helm chart from the mono-repo CRDs.
# It copies CRDs from crd/bases/ and adds required annotations/labels.

set -o errexit
set -o pipefail
set -o nounset

REPO_ROOT="$1"
OUTPUT_DIR="$2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Generating tinkerbell-crds helm chart..."
echo "  Source: ${REPO_ROOT}/crd/bases/"
echo "  Output: ${OUTPUT_DIR}"

# Clean and create output directory
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/templates"

# Copy chart metadata files from project
cp "${PROJECT_DIR}/chart/Chart.yaml" "${OUTPUT_DIR}/"
cp "${PROJECT_DIR}/chart/values.yaml" "${OUTPUT_DIR}/"
cp "${PROJECT_DIR}/chart/.helmignore" "${OUTPUT_DIR}/"

# Transient CRDs that should NOT have clusterctl labels (to prevent move during cluster migration).
# The move label causes clusterctl to discover and move these resources during cluster move.
# When moved, they arrive with empty status (clusterctl only moves spec) causing controllers
# to re-process them
SKIP_CLUSTERCTL_LABELS="tinkerbell.org_workflows.yaml bmc.tinkerbell.org_jobs.yaml bmc.tinkerbell.org_tasks.yaml"

# Copy all CRDs from mono-repo
for crd_file in "${REPO_ROOT}"/crd/bases/*.yaml; do
    filename=$(basename "$crd_file")
    echo "  Processing: ${filename}"
    
    # Copy CRD to templates directory
    cp "$crd_file" "${OUTPUT_DIR}/templates/${filename}"
    
    # Add helm.sh/resource-policy: keep annotation using yq
    # This prevents Helm from deleting CRDs on uninstall
    yq -i '.metadata.annotations["helm.sh/resource-policy"] = "keep"' "${OUTPUT_DIR}/templates/${filename}"
    
    # Add clusterctl labels (skip for transient CRDs)
    if [[ ! " ${SKIP_CLUSTERCTL_LABELS} " =~ " ${filename} " ]]; then
        yq -i '.metadata.labels["clusterctl.cluster.x-k8s.io"] = ""' "${OUTPUT_DIR}/templates/${filename}"
        yq -i '.metadata.labels["clusterctl.cluster.x-k8s.io/move"] = ""' "${OUTPUT_DIR}/templates/${filename}"
    fi
done

echo "Generated tinkerbell-crds chart with $(ls -1 "${OUTPUT_DIR}/templates/" | wc -l | tr -d ' ') CRDs"
