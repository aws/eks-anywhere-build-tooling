#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CUR_STATE=$(systemctl is-active docker)
if [[ $CUR_STATE == "active" ]]; then
	echo "Docker service is active. Pass"
else
	echo "Docker service is not active. State '$CUR_STATE'. Failing."
	exit 1
fi

if type -P docker &>/dev/null; then
	echo "Docker command exists. Pass."
else
	echo "Docker command does not exist. Failing."
	exit 1
fi

USER='ubuntu'
GROUP='docker'
if groups $USER | grep &>/dev/null "$GROUP"; then
	echo "The '$USER' is a member of the '$GROUP' group. Pass."
else
	echo "The '$USER' is not a member of the '$GROUP' group. Failing."
	exit 1
fi
