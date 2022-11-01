#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

EKSA_MANIFESTS_PATH="${EKSA_MANIFESTS_PATH:-/usr/lib/eks-a/manifests}" # The full path where eks-a manifests will be downloaded
EKSA_ARTIFACTS_TAR_PATH="${EKSA_ARTIFACTS_TAR_PATH:-/usr/lib/eks-a/artifacts/artifacts.tar.gz}" # The path where the tar containing all eks-a container artifacts will be stored

sudo mkdir -p $EKSA_MANIFESTS_PATH
sudo -E env "PATH=$PATH" eksctl anywhere download artifacts --retain-dir --download-dir $EKSA_MANIFESTS_PATH -v4
sudo chmod -R a+r $EKSA_MANIFESTS_PATH

ARTIFACTS_NAME_DIR=$(dirname $EKSA_ARTIFACTS_TAR_PATH)
sudo mkdir -p $ARTIFACTS_NAME_DIR
sudo -E env "PATH=$PATH" eksctl anywhere download images -o $EKSA_ARTIFACTS_TAR_PATH -v4
sudo chmod -R a+r $ARTIFACTS_NAME_DIR

# The download images command pulls down all the images in the bundle
# but after the images get saved to a tar archive, the images are still
# in the local Docker runtime, that bloats the final admin image. This can cause
# downstream resource heavy operations to fail due to lack of disk space.
# To mitigate this, we prune all the images and unused volumes to bring down
# size of the admin image.
sudo docker system prune --volumes --force
