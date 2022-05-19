#!/usr/bin/env bash

set -x
set -e
set -o pipefail

source /docker.sh 

CODEBUILD_CI="${CODEBUILD_CI:-false}"
GIT_CONFIG_SCOPE="--global"
if [[ "$CODEBUILD_CI" = "true" ]] && [[ "$CODEBUILD_BUILD_ID" =~ "aws-staging-bundle-build" ]]; then
    GIT_CONFIG_SCOPE="--system"
fi

if [ ! -d "/root/.docker" ]; then
    mkdir -p /root/.docker
fi

mv config/docker-ecr-config.json /root/.docker/config.json
git config ${GIT_CONFIG_SCOPE} credential.helper '!aws codecommit credential-helper $@'
git config ${GIT_CONFIG_SCOPE} credential.UseHttpPath true

start::dockerd
wait::for::dockerd

for i in {1..5}; do docker pull public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v6.1.0 && break || sleep 15; done

docker run --privileged --rm public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v6.1.0 --install aarch64
