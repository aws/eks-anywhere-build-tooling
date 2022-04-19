#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source /etc/profile.d/golang.sh

if type -p go &>/dev/null; then
	echo "go command exists. Pass."
else
	echo "go command does not exist. Failing."
	exit 1
fi

if type -p git &>/dev/null; then
	echo "git command exists. Pass."
else
	echo "git command does not exist. Failing."
	exit 1
fi