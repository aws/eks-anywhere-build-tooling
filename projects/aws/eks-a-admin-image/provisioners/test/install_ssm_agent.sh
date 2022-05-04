#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SSM_SERVICE_FILE=/etc/systemd/system/snap.amazon-ssm-agent.amazon-ssm-agent.service
if test -f "$SSM_SERVICE_FILE"; then
	echo "Amazon SSM Agent service is installed. Pass"
else
	echo "Amazon SSM Agent service is not installed."
	exit 1
fi