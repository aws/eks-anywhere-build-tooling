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
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

SED=$(build::common::gnu_variant_on_mac sed)

HELM_DESTINATION_REPOSITORY="${1?First argument is helm destination repository}"
OUTPUT_DIR="${2?Second arguement is output directory}"

CHART_NAME=$(basename ${HELM_DESTINATION_REPOSITORY})
DEST_DIR=${OUTPUT_DIR}/helm/${CHART_NAME}

#
# Search and replace
#
TEMPLATE_DIR=helm/templates
SEDFILE=${OUTPUT_DIR}/helm/sedfile
for file in Chart.yaml values.yaml
do
  $SED -f ${SEDFILE} -i ${DEST_DIR}/${file}
done

if [ -d ${OUTPUT_DIR}/helm/${CHART_NAME}/crds ]; then
  for file in crds/*.yaml 
  do
    $SED -f ${SEDFILE} -i ${DEST_DIR}/${file}
  done
fi