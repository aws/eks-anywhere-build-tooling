#!/usr/bin/env bash

set -x
set -o errexit
set -o nounset
set -o pipefail

LATEST_EKSA_RELEASE_NUMBER=$(curl $EKSA_RELEASE_MANIFEST_URL | yq e '.spec.releases[] | select(.version == '\"$EKSA_VERSION\"') | .number')
PROD_RELEASE_ADMIN_AMI_DESTINATION="s3://$EKSA_PRODUCTION_ARTIFACTS_BUCKET/releases/eks-a/$LATEST_EKSA_RELEASE_NUMBER/artifacts/eks-a-admin-ami/$EKSA_VERSION/eks-anywhere-admin-ami-$EKSA_VERSION-eks-a-$LATEST_EKSA_RELEASE_NUMBER.raw"

echo $PROD_RELEASE_ADMIN_AMI_DESTINATION
