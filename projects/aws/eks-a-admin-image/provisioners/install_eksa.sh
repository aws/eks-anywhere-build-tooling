#!/bin/bash

echo $EKSA_VERSION

set -x
set -o errexit
set -o nounset
set -o pipefail



EKSCTL_VERSION="${EKSCTL_VERSION:-latest}" # The eksctl version to install (example v0.0.0)
EKSA_RELEASE_MANIFEST_URL="${EKSA_RELEASE_MANIFEST_URL:-https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml}" # The url pointing to the eks-a releases manifest
EKSA_VERSION="${EKSA_VERSION:-latest}" # The eks-a version to install (example v0.0.0)

# Install depencies
sudo apt-get install -y tar
sudo snap install yq

# Install eksctl
mkdir eksctl_tmp
if [[ "$EKSCTL_VERSION" == "latest" ]]; then
	EKSCTL_URL="https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_linux_amd64.tar.gz"
else
	EKSCTL_URL="https://github.com/weaveworks/eksctl/releases/download/${EKSCTL_VERSION}/eksctl_linux_amd64.tar.gz"
fi
curl -L --silent $EKSCTL_URL | tar xz -C ./eksctl_tmp

sudo mv eksctl_tmp/eksctl /usr/local/bin/
rm -rf eksctl_tmp

# Install eksctl-anywhere
mkdir eksa_tmp
if [[ "$EKSA_VERSION" == "latest" || "$EKSA_VERSION" == "v0.0.0" ]]; then
	EKSA_VERSION=$(curl -L --silent $EKSA_RELEASE_MANIFEST_URL |  yq e '.spec.latestVersion' -)
fi
EKSA_TAR_URL=$(curl -L --silent $EKSA_RELEASE_MANIFEST_URL | yq e '.spec.releases[] | select(.version == "'$EKSA_VERSION'") | .eksABinary.linux.uri')
curl -L --silent $EKSA_TAR_URL | tar xz -C ./eksa_tmp

sudo mv eksa_tmp/eksctl-anywhere /usr/local/bin/
rm -rf eksa_tmpq
