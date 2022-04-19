#!/usr/bin/env bash

set -x
set -o errexit
set -o nounset
set -o pipefail


BUCKET="${1?Specify second argument - s3 bucket}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${SCRIPT_ROOT}/build/lib/common.sh"

OS=linux
REPO=client
ARCH=amd64
NAME="kubernetes-sigs/kind"
RELEASE_BRANCH=$(build::eksd_releases::get_release_branch)

echo $(build::common::get_latest_eksa_asset_url $BUCKET $NAME $ARCH $RELEASE_BRANCH)