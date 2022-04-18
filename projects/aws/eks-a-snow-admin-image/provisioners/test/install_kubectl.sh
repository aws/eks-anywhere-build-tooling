#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if type -p kubectl &>/dev/null; then
	echo "kubectl command exists. Pass."
else
	echo "kubectl command does not exist. Failing."
	exit 1
fi
