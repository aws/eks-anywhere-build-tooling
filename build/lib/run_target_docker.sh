#!/usr/bin/env bash
# Copyright 2020 Amazon.com Inc. or its affiliates. All Rights Reserved.
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

set -o errexit
set -o nounset
set -o pipefail

PROJECT="$1"
TARGET="$2"
IMAGE_REPO="${3:-}"
RELEASE_BRANCH="${4:-}"
ARTIFACTS_BUCKET="${5:-$ARTIFACTS_BUCKET}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${SCRIPT_ROOT}/common.sh"

echo "****************************************************************"
echo "A docker container with the name eks-a-builder will be launched."
echo "It will be left running to support running consecutive runs."
echo "Run 'make stop-docker-builder' when you are done to stop it."
echo "****************************************************************"

if ! docker ps -f name=eks-a-builder | grep -w eks-a-builder; then
	build::docker::retry_pull public.ecr.aws/eks-distro-build-tooling/builder-base:minimal-latest

	NETRC=""
	if [ -f $HOME/.netrc ]; then
		NETRC="--mount type=bind,source=$HOME/.netrc,target=/root/.netrc"
	fi

	docker run -d --name eks-a-builder --privileged $NETRC \
		--mount type=bind,source=$MAKE_ROOT,target=/eks-anywhere-build-tooling \
		-e GOPROXY=${GOPROXY:-} --entrypoint bash \
		public.ecr.aws/eks-distro-build-tooling/builder-base:minimal-latest -c "source /docker.sh && start::dockerd && sleep infinity"
fi

docker exec -it eks-a-builder make $TARGET -C /eks-anywhere-build-tooling/projects/$PROJECT RELEASE_BRANCH=$RELEASE_BRANCH IMAGE_REPO=$IMAGE_REPO ARTIFACTS_BUCKET=$ARTIFACTS_BUCKET
