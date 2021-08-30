#!/usr/bin/env bash

set -x
set -e
set -o pipefail

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
