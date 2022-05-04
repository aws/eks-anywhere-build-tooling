#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if type -p jq &>/dev/null; then
	echo "jq command exists. Pass."
else
	echo "jq command does not exist. Failing."
	exit 1
fi

if type -p python3 &>/dev/null; then
	echo "python3 command exists. Pass."
else
	echo "python3 command does not exist. Failing."
	exit 1
fi

if type -p aws &>/dev/null; then
	echo "aws command exists. Pass."
else
	echo "aws command does not exist. Failing."
	exit 1
fi