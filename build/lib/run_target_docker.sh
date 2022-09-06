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
ARTIFACTS_BUCKET="${5:-}"

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"

source "${SCRIPT_ROOT}/common.sh"

echo "****************************************************************"
echo "A docker container with the name eks-a-builder will be launched."
echo "It will be left running to support running consecutive runs."
echo "Run 'make stop-docker-builder' when you are done to stop it."
echo "****************************************************************"

if ! docker ps -f name=eks-a-builder | grep -w eks-a-builder; then
	docker pull public.ecr.aws/eks-distro-build-tooling/builder-base:latest
	docker run -d --name eks-a-builder --privileged -e GOPROXY=$GOPROXY --entrypoint sleep \
		public.ecr.aws/eks-distro-build-tooling/builder-base:latest  infinity 
fi

EXTRA_INCLUDES=""
PROJECT_DEPENDENCIES=$(make --no-print-directory -C $MAKE_ROOT/projects/$PROJECT var-value-PROJECT_DEPENDENCIES RELEASE_BRANCH=$(build::eksd_releases::get_release_branch))
if [ -n "$PROJECT_DEPENDENCIES" ]; then
	DEPS=(${PROJECT_DEPENDENCIES// / })
	for dep in "${DEPS[@]}"; do
		DEP_PRODUCT="$(cut -d/ -f1 <<< $dep)"
		DEP_ORG="$(cut -d/ -f2 <<< $dep)"
		DEP_REPO="$(cut -d/ -f3 <<< $dep)"

		if [[ "$DEP_PRODUCT" == "eksd" ]]; then
			continue
		fi

		EXTRA_INCLUDES+=" --include=projects/$DEP_ORG/$DEP_REPO/***"
	done
fi

rsync -e 'docker exec -i' -t -rm --exclude='.git/***' \
	--exclude="projects/$PROJECT/_output/***" --exclude="projects/$PROJECT/$(basename $PROJECT)/***" \
	--include="projects/$PROJECT/***" --include="projects/kubernetes-sigs/image-builder/BOTTLEROCKET_RELEASES" \
	--include="release/SUPPORTED_RELEASE_BRANCHES" --include="projects/kubernetes-sigs/cri-tools/GIT_TAG" $EXTRA_INCLUDES \
	--include='*/' --exclude='projects/***' $MAKE_ROOT/ eks-a-builder:/eks-anywhere-build-tooling

# Need so git properly finds the root of the repo
CURRENT_HEAD="$(cat $MAKE_ROOT/.git/HEAD | awk '{print $2}')"
docker exec -it eks-a-builder mkdir -p /eks-anywhere-build-tooling/.git/{refs,objects} /eks-anywhere-build-tooling/.git/$(dirname $CURRENT_HEAD)
docker cp $MAKE_ROOT/.git/HEAD eks-a-builder:/eks-anywhere-build-tooling/.git
docker cp $MAKE_ROOT/.git/$CURRENT_HEAD eks-a-builder:/eks-anywhere-build-tooling/.git/$CURRENT_HEAD

docker exec -it eks-a-builder make $TARGET -C /eks-anywhere-build-tooling/projects/$PROJECT RELEASE_BRANCH=$RELEASE_BRANCH IMAGE_REPO=$IMAGE_REPO ARTIFACTS_BUCKET=$ARTIFACTS_BUCKET
