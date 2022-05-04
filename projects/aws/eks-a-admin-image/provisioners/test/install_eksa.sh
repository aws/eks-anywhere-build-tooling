#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

 if type -p eksctl &>/dev/null; then
	echo "eksctl command exists. Pass."
else
	echo "eksctl command does not exist. Failing."
	exit 1
fi


if type -p eksctl-anywhere &>/dev/null; then
	echo "eksctl-anywhere binary exists. Pass."
else
	echo "eksctl-anywhere binary does not exist. Failing."
	exit 1
fi


if eksctl anywhere version &>/dev/null; then
	echo "eksctl anywhere command works. Pass."
else
	echo "eksctl anywhere command does not work. Failing."
	exit 1
fi
 