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

IMAGE_REPOSITORY="${1?First argument is image repository}"
OUTPUT_DIR="${2?Second arguement is output directory}"
CHART_NAME=$(basename ${IMAGE_REPOSITORY})
DESTINATION=${OUTPUT_DIR}/helm/${CHART_NAME}

mkdir -p ${DESTINATION}
SEDFILE=${OUTPUT_DIR}/helm/sedfile
envsubst <helm/sedfile.template >${SEDFILE}
TEMPLATE_DIR=helm/templates
cat helm/files.txt | while read SOURCE_FILE DESTINATION_FILE
do
  TMPFILE=/tmp/$(basename ${SOURCE_FILE})
  cp ${SOURCE_FILE} ${TMPFILE}
  sed -f ${SEDFILE} ${TMPFILE} >${DESTINATION}/${DESTINATION_FILE}
  rm -f ${TMPFILE}
done
