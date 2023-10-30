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
BASE_DIRECTORY="${6:-}"
GO_MOD_CACHE="${7:-}"
BUILDER_PLATFORM_ARCH="${8:-amd64}"
REMOVE="${9:-false}"
BUILDER_BASE_TAG="${10:-latest}"
PLATFORM="${11:-}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

source "${SCRIPT_ROOT}/common.sh"

MAKE_VARS="IMAGE_REPO=$IMAGE_REPO ARTIFACTS_BUCKET=$ARTIFACTS_BUCKET"

function remove_container()
{
	docker rm -vf $CONTAINER_ID > /dev/null 2>&1
}

SKIP_RUN="false"
NAME=""
if [[ "$REMOVE" == "false" ]]; then
	NAME="--name eks-a-builder"

	if docker ps -f name=eks-a-builder | grep -w eks-a-builder; then
		SKIP_RUN="true"
		CONTAINER_ID="eks-a-builder"
	fi
else
	trap "remove_container" EXIT
fi

IMAGE="public.ecr.aws/eks-distro-build-tooling/builder-base:$BUILDER_BASE_TAG"
# since if building cgo we will specifically set the arch to something other than the host
# ensure we always explictly ask for the host platform, unless override for cgo
PLATFORM_ARG="--platform linux/$BUILDER_PLATFORM_ARCH"

if [[ -n "$PLATFORM" ]]; then
	DIGEST=$(docker buildx imagetools inspect --raw public.ecr.aws/eks-distro-build-tooling/builder-base:$BUILDER_BASE_TAG | jq -r ".manifests[] | select(.platform.architecture == \"${PLATFORM#linux/}\") | .digest")
	IMAGE="public.ecr.aws/eks-distro-build-tooling/builder-base@$DIGEST"
	PLATFORM_ARG="--platform $PLATFORM"
	MAKE_VARS+=" BINARY_PLATFORMS=$PLATFORM"
fi

DOCKER_USER_FLAG=""
NETRC_DIR="/root"
if [ "$(uname -s)" = "Linux" ] && [ -n "${USER:-}" ]; then
	# on a linux host, the uid needs to match the host user otherwise
	# all the downloaded go modules will be owned by root in the host
	USER_ID=$(id -u ${USER})
	USER_GROUP_ID=$(id -g ${USER})
	DOCKER_USER_FLAG="-u $USER_ID:$USER_GROUP_ID"
	NETRC_DIR="/home/matchinguser"
fi


if [[ "$SKIP_RUN" == "false" ]]; then
	echo "Pulling $IMAGE...."
	if ! build::docker::retry_pull $IMAGE > /dev/null 2>&1; then
		# try one more time to show the error to the user
		docker pull $IMAGE
	fi

	NETRC=""
	if [ -f $HOME/.netrc ]; then
		DOCKER_RUN_NETRC="${DOCKER_RUN_NETRC:-$HOME/.netrc}"
		NETRC="--mount type=bind,source=$DOCKER_RUN_NETRC,target=$NETRC_DIR/.netrc"
	else
		DOCKER_RUN_NETRC=""
	fi

	mkdir -p $GO_MOD_CACHE
	CONTAINER_ID=$(build::common::echo_and_run docker run -d $NAME --privileged $NETRC $PLATFORM_ARG \
		--mount type=bind,source=$BASE_DIRECTORY,target=/eks-anywhere-build-tooling \
		--mount type=bind,source=$GO_MOD_CACHE,target=/mod-cache \
		-e GOPROXY=${GOPROXY:-} -e GOMODCACHE=/mod-cache -e DOCKER_RUN_BASE_DIRECTORY=$BASE_DIRECTORY \
		--entrypoint sleep $IMAGE infinity)

	if [ -n "$DOCKER_USER_FLAG" ]; then
		build::common::echo_and_run docker exec -t $CONTAINER_ID /eks-anywhere-build-tooling/build/lib/prepare_build_container_user.sh "$USER_GROUP_ID" "$USER_GROUP_ID"
	fi

	if [[ "$REMOVE" == "false" ]]; then
		echo "****************************************************************"
		echo "A docker container with the name eks-a-builder will be launched."
		echo "It will be left running to support running consecutive runs."
		echo "Run 'make stop-docker-builder' when you are done to stop it."
		echo "****************************************************************"
	fi
fi


build::common::echo_and_run docker exec -e RELEASE_BRANCH=$RELEASE_BRANCH $DOCKER_USER_FLAG \
	-t $CONTAINER_ID \
	make --no-print-directory $TARGET -C /eks-anywhere-build-tooling/projects/$PROJECT $MAKE_VARS
