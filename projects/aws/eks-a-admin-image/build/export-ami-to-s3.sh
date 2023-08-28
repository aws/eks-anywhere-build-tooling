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

terminate_instance() {
    aws ec2 terminate-instances --instance-ids $1
}
trap 'if [ -n "$INSTANCE_ID" ]; then terminate_instance $INSTANCE_ID; fi;' EXIT

AMI_MANIFEST_OUTPUT="${1?Specify first argument - AMI manifest output file from packer}"
IMAGE_FORMAT="${2?Specify second argument - Format for exported image \(vmdk\|raw\|vhd\)}"
S3_DST_EXPORT_PATH="${3?Specify third argument - Destination S3 path}"
REPLICAS="${4?Specify fourth argument - Comma separated list of s3 destinations for exported image copies}"

ARTIFACT_ID=$(cat $AMI_MANIFEST_OUTPUT | jq -r '.builds[0].artifact_id')
ARTIFACT_ID_SPLIT=(${ARTIFACT_ID//:/ })
AMI_ID=${ARTIFACT_ID_SPLIT[1]}

REPLICAS_SPLIT=(${REPLICAS//,/ })

DST_BUCKET_PATH=${S3_DST_EXPORT_PATH#"s3://"}
DST_BUCKET_NAME=${DST_BUCKET_PATH%%/*}
DST_PATH=${DST_BUCKET_PATH#"$DST_BUCKET_NAME/"}

EXPORTED_IMAGE_PREFIX="$DST_PATH"

EXPORT_RESPONSE=$(aws ec2 export-image --disk-image-format $IMAGE_FORMAT --s3-export-location S3Bucket=$DST_BUCKET_NAME,S3Prefix=$EXPORTED_IMAGE_PREFIX --image-id $AMI_ID)
echo $EXPORT_RESPONSE

EXPORT_TASK_ID=$(echo $EXPORT_RESPONSE | jq -r '.ExportImageTaskId')
EXPORTED_IMAGE_KEY="${EXPORTED_IMAGE_PREFIX}${EXPORT_TASK_ID}.${IMAGE_FORMAT}"
EXPORTED_IMAGE_LOCATION="s3://${DST_BUCKET_NAME}/${EXPORTED_IMAGE_KEY}"
EXPORTED_IMAGE_URL="https://${DST_BUCKET_NAME}.s3.amazonaws.com/${EXPORTED_IMAGE_KEY}"

FINAL_STATUSES=(completed deleted)
STATUS=$(echo $EXPORT_RESPONSE | jq -r '.Status')
STATUS_MESSAGE=$(echo $EXPORT_RESPONSE | jq -r '.StatusMessage')
PROGRESS=$(echo $EXPORT_RESPONSE | jq -r '.Progress')

until [[ "${FINAL_STATUSES[*]}" =~ "${STATUS}" ]]; do
  echo "Image import is $STATUS: $STATUS_MESSAGE $PROGRESS%"
  sleep 30

  DESCRIBE_RESPONSE=$(aws ec2 describe-export-image-tasks --export-image-task-ids $EXPORT_TASK_ID)
  STATUS=$(echo $DESCRIBE_RESPONSE | jq -r '.ExportImageTasks[0].Status')
  STATUS_MESSAGE=$(echo $DESCRIBE_RESPONSE | jq -r '.ExportImageTasks[0].StatusMessage')
  PROGRESS=$(echo $DESCRIBE_RESPONSE | jq -r '.ExportImageTasks[0].Progress')
done

if [[ "$STATUS" != "completed" ]]; then
    echo "Image import failed: $STATUS - $STATUS_MESSAGE"
    exit 1
fi

echo "Image import for ami $AMI_ID succeeded"

aws s3api put-object-acl --bucket $DST_BUCKET_NAME --key $EXPORTED_IMAGE_KEY --acl public-read
EXPORTED_IMAGE_SHA256=$(curl -s $EXPORTED_IMAGE_URL | sha256sum | cut -d" " -f1)
for dst in "${REPLICAS_SPLIT[@]}"
do
  echo "Copying exported image to $dst"
  aws s3 cp --no-progress --acl public-read $EXPORTED_IMAGE_LOCATION $dst
  echo -n "$EXPORTED_IMAGE_SHA256  $(basename $dst)" > $(basename $dst).sha256
  aws s3 cp --no-progress --acl public-read $(basename $dst).sha256 $dst.sha256
done

echo "Launching EC2 instance from AMI $AMI_ID for Amazon Inspector scan"
MAX_RETRIES=20
for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $(($i)) of instance launch"

    # Create a single EC2 instance with provided instance type and AMI
    # Query the instance ID for use in future commands
    INSTANCE_ID=$(aws ec2 run-instances --count 1 --image-id=$AMI_ID --instance-type $SNOW_ADMIN_IMAGE_INSTANCE_TYPE  --metadata-options "HttpEndpoint=enabled,HttpTokens=required,HttpPutResponseHopLimit=2" --associate-public-ip-address --iam-instance-profile Name=$SNOW_ADMIN_IMAGE_IAM_INSTANCE_PROFILE --query "Instances[0].InstanceId" --output text) && break

    if [ "$i" = "$MAX_RETRIES" ]; then
        echo "Failed to launch EC2 instance after $i retries"
        exit 1
    fi
    sleep 30
done

# Wait in loop until instance is running
aws ec2 wait instance-running --instance-ids $INSTANCE_ID

# Amazon Inspector requires that the instance be managed by AWS Systems Manager, so we wait until
# the condition is satisfied
MAX_RETRIES=40
for i in $(seq 1 $MAX_RETRIES); do
  echo "Attempt $(($i)) of checking if instance is managed by AWS Systems Manager"
  MANAGED_INSTANCES=($(aws ssm describe-instance-information --filters Key=ResourceType,Values=EC2Instance --query "InstanceInformationList[].InstanceId" --output text))
  if [[ "${MANAGED_INSTANCES[@]}" =~ "INSTANCE_ID" ]]; then
    echo "EC2 instance $INSTANCE_ID successfully registered as a managed instance"
    break
  fi

  if [ "$i" = "$MAX_RETRIES" ]; then
    echo "EC2 instance $INSTANCE_ID failed to register as a managed instance"
    exit 1
  fi

  sleep 15
fi

# Generate the findings report for the EC2 instance using Amazon Inspector V2
# and export it to an S3 bucket
INSPECTOR_FINDINGS_REPORT_ID=$(aws inspector2 create-findings-report --filter-criteria '{"resourceId": [{"comparison": "EQUALS", "value": "'"$INSTANCE_ID"'"}], "findingStatus": [{"comparison": "EQUALS", "value": "ACTIVE"}]}' --report-format CSV --s3-destination '{"bucketName": "'"$SNOW_ADMIN_IMAGE_INSPECTOR_BUCKET"'", "kmsKeyArn": "'"$SNOW_ADMIN_IMAGE_INSPECTOR_KMS_KEY_ARN"'"}' --query "reportId" --output text) 

# Make findings report public 
aws s3api put-object-acl --bucket $SNOW_ADMIN_IMAGE_INSPECTOR_BUCKET --key $INSPECTOR_FINDINGS_REPORT_ID --acl public-read
