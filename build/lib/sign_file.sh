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

FILE="${1?First argument is file to sign}"
KEY_NAME="${2?Second argument is key to sign with}"
FILE_SIGNED="signed_${FILE}"
SIGNING_ALGORITHM="ECDSA_SHA_256"

# Exclude the signature portion of a CRD yaml
alwaysexcludes='.metadata.annotations."eksa.aws.com/signature"'
has_excludes=$(<${FILE} yq -rc '.metadata.annotations."eksa.aws.com/excludes"')

# Ignore excludes for the signature field of the CRD before signing the file.
excludes=""
if [ "${has_excludes}" != "null" ]; then
    excludes=$(< ${FILE} yq -rc '.metadata.annotations."eksa.aws.com/excludes"' | base64 -d | cat <(echo ${alwaysexcludes}) - | paste -sd "," -) 
fi
fixed=$(<${FILE} yq --indentless-lists -y -S \
"del(${alwaysexcludes}$([ ! -z ${excludes} ] && echo , ${excludes})) | walk( if type == \"object\" then with_entries(select(.value != \"\" and .value != null and .value != [])) else . end)")
encoded=$(echo "${fixed}" | base64 | tr -d '\n')

# Signing the file with the KMS key ECDSA_SHA_256
signature=$(aws kms sign --key-id alias/${KEY_NAME} --message ${encoded} --message-type RAW --signing-algorithm ${SIGNING_ALGORITHM} | jq -rc '.Signature')

# Adding Signature to the bundle yaml annotation field
signature_b64=$(echo "$signature" | base64)
signed_file=$(< ${FILE} yq -y ".metadata.annotations.\"eksa.aws.com/signature\" = \"${signature_b64}\"")

# Output the signed file to a new yaml for uploading
echo "${fixed}" | openssl dgst -binary | base64
echo "${signed_file}" > ${FILE_SIGNED}
