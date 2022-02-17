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
BUILD_TARGET="${8?Specify eighth argument - Raw build target name}"

CODEBUILD_CI="${CODEBUILD_CI:-false}"
CI="${CI:-false}"

if [ "$CODEBUILD_CI" = "true" ]; then
    KEY_NAME="$KEY_NAME-$CODEBUILD_BUILD_ID"
elif [ "$CI" = "true" ]; then
    KEY_NAME="$KEY_NAME-$PROW_JOB_ID"
fi

REPO_NAME=$(basename $REPO_ROOT)
KEY_LOCATION=$REPO_ROOT/$PROJECT_PATH/$KEY_NAME.pem
IMAGE_BUILDER_MAKE_ROOT=$PROJECT_PATH/$IMAGE_BUILDER_DIR
REMOTE_HOME_DIR="/home/ec2-user"
REMOTE_IMAGE_BUILDER_MAKE_ROOT=$REMOTE_HOME_DIR/$REPO_NAME/$IMAGE_BUILDER_MAKE_ROOT
SSH_OPTS="-i $KEY_LOCATION -o StrictHostKeyChecking=no -o ConnectTimeout=120"

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

RUN_INSTANCE_EXTRA_ARGS=""
if [ "$CODEBUILD_CI" = "true" ]; then
    # Select a random subnet from Codebuild project's private subnet list
    SUBNET_ID_LIST=$RAW_IMAGE_BUILD_SUBNET_ID
    SUBNET_COUNT=$(echo $SUBNET_ID_LIST | awk -F\, '{print NF}')
    INDEX=$((($RANDOM % $SUBNET_COUNT) + 1))
    SUBNET_ID=$(cut -d',' -f${INDEX} <<< $SUBNET_ID_LIST)

    # Define extra args to run the instance in the same subnet and use
    # the same security group as Codebuild
    RUN_INSTANCE_EXTRA_ARGS="--subnet-id $SUBNET_ID --security-group-ids $RAW_IMAGE_BUILD_SECURITY_GROUP --associate-public-ip-address"
fi

# Create a single EC2 instance with provided instance type and AMI
# Query the instance ID for use in future commands
INSTANCE_ID=$(aws ec2 run-instances --count 1 --image-id=$AMI_ID --instance-type $INSTANCE_TYPE --key-name $KEY_NAME $RUN_INSTANCE_EXTRA_ARGS --query "Instances[0].InstanceId" --output text)

# Wait in loop until instance is running
aws ec2 wait instance-running --instance-ids $INSTANCE_ID

# Get the public DNS of the instance to use as SSH hostname
PUBLIC_DNS_NAME=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query "Reservations[].Instances[].PublicDnsName" --output text)
REMOTE_HOST=ec2-user@$PUBLIC_DNS_NAME

MAX_RETRIES=15

# rsync might sometimes fail with flaky connection issues, so
# implementing retry logic will make it more robust to flakes
for i in $(seq 1 $MAX_RETRIES); do
    echo "Attempt $(($i))"

    # Transfer the repo contents from the CI environment to the EC2 instance
    rsync -avzh --progress -e "ssh $SSH_OPTS" $REPO_ROOT $REMOTE_HOST:~/ && echo "Files transferred!" && break

    if [ "$i" = "$MAX_RETRIES" ]; then
        exit 1
    fi
    sleep 10
done

# If not running on Codebuild, exit gracefully
# if [ "$CODEBUILD_CI" = "false" ]; then
#     exit 0
# fi

# Run permissions setup for KVM builds on instance
# Run make command to build raw image
ssh $SSH_OPTS $REMOTE_HOST "sudo usermod -a -G kvm ec2-user; sudo chmod 666 /dev/kvm; sudo chown root:kvm /dev/kvm; sudo wget https://redhat-iso-images.s3.amazonaws.com/8.4/rhel-8.4-x86_64-dvd.iso -P /home/ec2-user/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/image-builder/images/capi; PACKER_VAR_FILES='$PACKER_VAR_FILES' PACKER_FLAGS=-force PACKER_LOG=1 PACKER_LOG_PATH=$REMOTE_IMAGE_BUILDER_MAKE_ROOT/packer.log make build-raw-$BUILD_TARGET -C $REMOTE_IMAGE_BUILDER_MAKE_ROOT"

# Copy built raw image from the instance back into the CI build environment
mkdir -p $REPO_ROOT/$IMAGE_BUILDER_MAKE_ROOT/output
scp $SSH_OPTS $REMOTE_HOST:$REMOTE_IMAGE_BUILDER_MAKE_ROOT/output/*.gz $REPO_ROOT/$IMAGE_BUILDER_MAKE_ROOT/output/
scp $SSH_OPTS $REMOTE_HOST:$REMOTE_IMAGE_BUILDER_MAKE_ROOT/packer.log $REPO_ROOT/$IMAGE_BUILDER_MAKE_ROOT/output/
