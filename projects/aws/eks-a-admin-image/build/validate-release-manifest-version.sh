#!/usr/bin/env bash

set -x
set -o errexit
set -o nounset
set -o pipefail

EKSA_RELEASE_MANIFEST_S3_URL="${1?Specify first argument - Canonical S3 URL for the EKS-A releases manifest}"
EKSA_RELEASE_MANIFEST_CDN_URL="${2?Specify second argument - CloudFront URL for the EKS-A releases manifest}"

n=0
max_attempts=120
delay=30
while true; do
    LATEST_EKSA_VERSION_S3=$(curl $EKSA_RELEASE_MANIFEST_S3_URL | yq e '.spec.releases[-1].version' -)
    LATEST_EKSA_VERSION_CLOUDFRONT=$(curl $EKSA_RELEASE_MANIFEST_CDN_URL | yq e '.spec.releases[-1].version' -)
    [[ "$LATEST_EKSA_VERSION_S3" == "$LATEST_EKSA_VERSION_CLOUDFRONT" ]] && break || {
        if [[ $n -lt $max_attempts ]]; then
            ((n++))
            echo "Older version of release manifest is cached in CloudFront, waiting for cache to be updated - Attempt $n/$max"
            sleep $delay;
        else
            echo "Timeout occured after $n attempts waiting for CloudFront cache to be updated"
            exit 1
        fi
    }
done
