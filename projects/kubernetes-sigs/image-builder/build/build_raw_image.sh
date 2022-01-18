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

REPO_ROOT="${1?Specify first argument - repo root}"
PROJECT_PATH="${2?Specify second argument - Project path}"
IMAGE_BUILDER_DIR="${3?Specify third argument - Image builder directory}"
PACKER_VAR_FILES="${4?Specify fourth argument - Packer var files}"
AMI_ID="${5?Specify fifth argument - AMI ID to create instance}"
INSTANCE_TYPE="${6?Specify sixth argument - Instance type to create}"
KEY_NAME="${7?Specify seventh argument - Key name to associate with instance}"

REPO_NAME=$(basename $REPO_ROOT)
KEY_LOCATION=$REPO_ROOT/$PROJECT_PATH/$KEY_NAME.pem
IMAGE_BUILDER_MAKE_ROOT=$PROJECT_PATH/$IMAGE_BUILDER_DIR
SSH_OPTS="-i $KEY_LOCATION -o StrictHostKeyChecking=no -o ConnectTimeout=120"
CODEBUILD_CI="${CODEBUILD_CI:-false}"

terminate_instance() {
    aws ec2 terminate-instances --instance-ids $1
}
delete_key_pair() {
    aws ec2 delete-key-pair --key-name $1
}

# Delete keypair and instance when exiting script. This will prevent
# lingering instances by deleting the EC2 instance regardless of success
# or failure of build
trap 'if [ -n "$INSTANCE_ID" ]; then terminate_instance $INSTANCE_ID; fi; delete_key_pair $KEY_NAME' EXIT

# Check if key of that name already exists, else create keypair
# Query and save the key contents into a local file for
# communicating to EC2 instance via SSH 
if ! aws ec2 wait key-pair-exists --key-names=$KEY_NAME 2>/dev/null ; then
    aws ec2 create-key-pair --key-name $KEY_NAME --query "KeyMaterial" --output text > $KEY_LOCATION
fi
chmod 600 $KEY_LOCATION

# Create a single EC2 instance with provided instance type and AMI
# Query the instance ID for use in future commands
# Wait in loop until instance is running
INSTANCE_ID=$(aws ec2 run-instances --count 1 --image-id=$AMI_ID --instance-type $INSTANCE_TYPE --key-name $KEY_NAME --query "Instances[0].InstanceId" --output text)
aws ec2 wait instance-running --instance-ids $INSTANCE_ID

# Get the public DNS of the instance to use as SSH hostname
PUBLIC_DNS_NAME=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query "Reservations[].Instances[].PublicDnsName" --output text)
REMOTE_HOST=ec2-user@$PUBLIC_DNS_NAME

MAX_RETRIES=15

for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $(($i))"

    # Transfer the repo contents from the CI environment to the EC2 instance
    rsync -avzhR --progress --max-size=500K -e "ssh $SSH_OPTS" $REPO_ROOT $REMOTE_HOST:~/ && echo "Files transferred!" && break

    if [ "$i" = "$MAX_RETRIES" ]; then
        exit 1
    fi
    sleep 10
done

# If not running on Codebuild, exit gracefully
if [ "$CODEBUILD_CI" = "false" ]; then
    exit 0
fi

# Run permissions setup for KVM builds on instance
# Run make command to build raw image
ssh $SSH_OPTS $REMOTE_HOST "sudo usermod -a -G kvm ec2-user; sudo chmod 666 /dev/kvm; sudo chown root:kvm /dev/kvm; PACKER_VAR_FILES='$PACKER_VAR_FILES' PACKER_FLAGS=-force PACKER_LOG=1 make build-raw-ubuntu-2004 -C /home/ec2-user/$REPO_NAME/$IMAGE_BUILDER_MAKE_ROOT"

# Copy built raw image from the instance back into the CI build environment
scp $SSH_OPTS $REMOTE_HOST:/home/ec2-user/$REPO_NAME/$IMAGE_BUILDER_MAKE_ROOT/output/*.gz $REPO_ROOT/$IMAGE_BUILDER_MAKE_ROOT/output/
