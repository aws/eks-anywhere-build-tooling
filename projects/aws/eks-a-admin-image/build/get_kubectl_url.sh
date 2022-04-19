#!/usr/bin/env bash

set -x
set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd -P)"
source "${SCRIPT_ROOT}/build/lib/common.sh"

OS=linux
REPO=client
ARCH=amd64
NAME="bin/$OS/$ARCH/kubectl"
RELEASE_BRANCH=$(build::eksd_releases::get_release_branch)

echo $(build::eksd_releases::get_eksd_kubernetes_asset_url $NAME $RELEASE_BRANCH $ARCH)