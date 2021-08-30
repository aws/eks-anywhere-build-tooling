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

INLINE_VARS_FILE="${1?Specify first argument - inline vars file for goss}"

TEST_VM_NAME=vendor-validation-testing
SSH_KEY=$(cat ${HOME}/.ssh/id_rsa.pub)
USERNAME=tester

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

GOSS_DIR=$MAKE_ROOT/image-builder/images/capi/packer/goss
REMOTE_GOSS_DIR=/home/$USERNAME

# Clone the template to a vm
govc vm.clone -on=false -vm $TEMPLATE -net=sddc-cgw-network-1 $TEST_VM_NAME

# Set the meta data to the clonned vm
govc vm.change -vm $TEST_VM_NAME \
	-e guestinfo.metadata="$(MACADDRESS=$(govc device.info -vm $TEST_VM_NAME -json ethernet-0 | jq ".Devices[0].MacAddress") TEST_VM_NAME=$TEST_VM_NAME envsubst < $MAKE_ROOT/validate/metadata.yaml | base64)" \
	-e guestinfo.metadata.encoding="base64"

# Set the user data to the clonned vm to add users
govc vm.change -vm $TEST_VM_NAME \
	-e guestinfo.userdata="$(SSH_KEY=$SSH_KEY USERNAME=$USERNAME envsubst < $MAKE_ROOT/validate/cloudconfig.yaml | base64)" \
	-e guestinfo.userdata.encoding="base64"

govc vm.power -on $TEST_VM_NAME

# Wait for the vm to get an ip
echo "Waiting for vm"
govc vm.info -waitip $TEST_VM_NAME

REMOTE_IP=$(govc vm.ip $TEST_VM_NAME)
echo "Remote vm ip address is" $REMOTE_IP

# Run all the pre-reqs commands
declare -a cmds=("sudo apt install git -y"
				 "sudo curl -L https://github.com/aelsabbahy/goss/releases/latest/download/goss-linux-amd64 -o /usr/local/bin/goss"
				 "sudo chmod +rx /usr/local/bin/goss"
)

for cmd in "${cmds[@]}"
do
	ssh -o StrictHostKeyChecking=no $USERNAME@$REMOTE_IP "$cmd"
done

# Copy over all the goss scripts
scp -o StrictHostKeyChecking=no -r $GOSS_DIR/* $USERNAME@$REMOTE_IP:$REMOTE_GOSS_DIR

# Run the goss scripts
INLINE_VARS=$(cat $INLINE_VARS_FILE)
ssh -o StrictHostKeyChecking=no $USERNAME@$REMOTE_IP "sudo goss -g goss.yaml --vars goss-vars.yaml --vars-inline '$INLINE_VARS' validate --retry-timeout 0s --sleep 1s -f json -o pretty"

if [ $? -ne 0 ]; then
	echo "Goss validation failed"
	echo "The vm is left on, please access with $USERNAME@$REMOTE_IP"
	exit 1
else
	echo "Goss validation passed"
fi

# Deleting and removing the vm
govc vm.destroy $TEST_VM_NAME
exit 0