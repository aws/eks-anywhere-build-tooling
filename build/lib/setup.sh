#!/usr/bin/env bash

set -x
set -e
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

source /docker.sh 

CODEBUILD_CI="${CODEBUILD_CI:-false}"
GIT_CONFIG_SCOPE="--global"
if [[ "$CODEBUILD_CI" = "true" ]] && [[ "$CODEBUILD_BUILD_ID" =~ "aws-staging-bundle-build" ]]; then
    GIT_CONFIG_SCOPE="--system"
fi

if [ ! -d "/root/.docker" ]; then
    mkdir -p /root/.docker
fi

if [ ! -d "/root/.config/containers" ]; then
    mkdir -p /root/.config/containers
fi

cp config/docker-ecr-config.json /root/.docker/config.json
cp config/policy.json /root/.config/containers/policy.json
git config ${GIT_CONFIG_SCOPE} credential.helper '!aws codecommit credential-helper $@'
git config ${GIT_CONFIG_SCOPE} credential.UseHttpPath true

# Since the build environment is AL2, we need to use iptables in legacy mode
# as it doesn't have nftables.
update-alternatives --set iptables /usr/sbin/iptables-legacy
update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy
start::dockerd
wait::for::dockerd

build::docker::retry_pull public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v7.0.0

docker run --privileged --rm public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v7.0.0 --install aarch64,amd64
