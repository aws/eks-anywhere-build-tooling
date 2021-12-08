#!/usr/bin/env bash
# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
set -o errexit
set -o nounset
set -o pipefail

# the upstream version of this script builds the image
# download them from eks-distro public.ecr
# and pulls the binaries from the eks-d release bucket

# since kind is calling this directly, cannot require args
# instead the following env vars are required
if [[ -z "${KUBE_VERSION}" ]]; then
    echo "KUBE_VERSION env var not set"
    exit 1
fi

if [[ -z "${EKSD_RELEASE_BRANCH}" ]]; then
    echo "EKSD_RELEASE_BRANCH env var not set"
    exit 1
fi

if [[ -z "${EKSD_RELEASE}" ]]; then
    echo "EKSD_RELEASE env var not set"
    exit 1
fi

if [[ -z "${EKSD_IMAGE_REPO}" ]]; then
    echo "EKSD_IMAGE_REPO env var not set"
    exit 1
fi
if [[ -z "${EKSD_ASSET_URL}" ]]; then
    echo "EKSD_ASSET_URL env var not set"
    exit 1
fi

if [[ -z "${KUBE_ARCH}" ]]; then
    echo "KUBE_ARCH env var not set"
    exit 1
fi

SOURCE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"

# symlink to project root _output to make sure files are properly cleaned up
# when running make clean
PROJECT_ROOT="$(cd "$SOURCE_ROOT/../../.." && pwd -P)"
rm -f $SOURCE_ROOT/_output
ln -s $PROJECT_ROOT/_output/$EKSD_RELEASE_BRANCH $SOURCE_ROOT/_output

# Download binaries
mkdir -p $SOURCE_ROOT/_output/dockerized/bin/linux/$KUBE_ARCH 

for binary in "kubeadm" "kubelet" "kubectl"; do
    FILE="$SOURCE_ROOT/_output/dockerized/bin/linux/$KUBE_ARCH/$binary "
    if [ ! -f $FILE ]; then
        curl $EKSD_ASSET_URL/bin/linux/$KUBE_ARCH/$binary -o $FILE
    fi
done

# Download container images
EKSD_TAG="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"

mkdir -p $SOURCE_ROOT/_output/release-images/$KUBE_ARCH/   

for container in "kube-apiserver" "kube-controller-manager" "kube-scheduler" "kube-proxy"; do
    IMAGE_TAG="$EKSD_IMAGE_REPO/kubernetes/$container:$EKSD_TAG"
    FILE="$SOURCE_ROOT/_output/release-images/$KUBE_ARCH/$container.tar"
    if [ ! -f $FILE ]; then
        docker pull --platform linux/$KUBE_ARCH $IMAGE_TAG
        docker save $IMAGE_TAG -o $FILE
    fi
done
