#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if type -p kind &>/dev/null; then
	echo "kind command exists. Pass."
else
	echo "kind command does not exist. Failing."
	exit 1
fi
