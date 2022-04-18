#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

EKSA_MANIFESTS_PATH="${EKSA_MANIFESTS_PATH:-/usr/lib/eks-a/manifests}" # The full path where eks-a manifests will be downloaded
EKSA_ARTIFACTS_TAR_PATH="${EKSA_ARTIFACTS_TAR_PATH:-/usr/lib/eks-a/artifacts/artifacts.tar}" # The path where the tar containing all eks-a container artifacts will be stored

sudo mkdir -p $EKSA_MANIFESTS_PATH
sudo eksctl anywhere download artifacts --retain-dir --download-dir $EKSA_MANIFESTS_PATH -v4
sudo chmod -R a+r $EKSA_MANIFESTS_PATH

ARTIFACTS_NAME_DIR=$(dirname $EKSA_ARTIFACTS_TAR_PATH)
sudo mkdir -p $ARTIFACTS_NAME_DIR
sudo eksctl anywhere download images -o $EKSA_ARTIFACTS_TAR_PATH -v4
sudo chmod -R a+r $ARTIFACTS_NAME_DIR
