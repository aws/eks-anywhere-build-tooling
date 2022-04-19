#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

EKSA_MANIFESTS_PATH="${EKSA_MANIFESTS_PATH:-/usr/lib/eks-a/manifests}"
EKSA_ARTIFACTS_TAR_PATH="${EKSA_ARTIFACTS_TAR_PATH:-/usr/lib/eks-a/artifacts/artifacts.tar}"

if test -d $EKSA_MANIFESTS_PATH; then
	echo "Manifests dir exists. Pass."
else
	echo "Manifests dir does not exist. Failing."
	exit 1
fi

if test -f $EKSA_ARTIFACTS_TAR_PATH; then
	echo "Artifacts tarball exists. Pass"
else
	echo "Artifacts tarball does not exist. Failing."
	exit 1
fi
