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


set -e
set -o pipefail
set -x

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
REPO="eks-anywhere-build-tooling"

TAG_FILES=(
    "EKS_DISTRO_BASE_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_CSI_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_DOCKER_CLIENT_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_GIT_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_GLIBC_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_IPTABLES_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_NONROOT_TAG_FILE"
    "EKS_DISTRO_MINIMAL_BASE_TAG_FILE"
)

count=0
for FILE in "${TAG_FILES[@]}"; do
    count=$(expr $count + 1)
    OLD_TAG="$(cat ${SCRIPT_ROOT}/../../${FILE})"
    LATEST_TAG="$(curl https://raw.githubusercontent.com/aws/eks-distro-build-tooling/main/${FILE})"
    if [ "$OLD_TAG" = "$LATEST_TAG" ]; then
        echo "Tag file is up to date!"
        continue
    else
        echo "Tag file is out of date! Updating"
        sed -i "s,.*,${LATEST_TAG}," ./${FILE}
        git add ./${FILE}
    fi
done

if [ $count -eq 0 ]; then
    echo "No files to update!"
    exit 0
fi

if [ -z "$REPO_OWNER" ]; then
    echo "No org information was provided, please set and export REPO_OWNER environment variable. \
      This is used to raise a pull request against your org after updating tags in the respective files."
    exit 1
fi
if [ "$REPO_OWNER" = "aws" ]; then
    ORIGIN_ORG="eks-distro-pr-bot"
    UPSTREAM_ORG="aws"
else
    ORIGIN_ORG=$REPO_OWNER
    UPSTREAM_ORG=$REPO_OWNER
fi

COMMIT_MESSAGE="[PR BOT] Update EKS Distro base image tag(s)"
PR_TITLE="Update EKS Distro base image tag in Tag file(s)"
PR_BODY=$(cat ${SCRIPT_ROOT}/eks_distro_base_pr_body)
PR_BRANCH="image-tag-update"

git config --global push.default current
git remote add origin git@github.com:${ORIGIN_ORG}/${REPO}.git
git remote add upstream git@github.com:${UPSTREAM_ORG}/${REPO}.git
git checkout -b $PR_BRANCH

git commit -m "$COMMIT_MESSAGE"
ssh-agent bash -c 'ssh-add /secrets/ssh-secrets/ssh-privatekey; ssh -o StrictHostKeyChecking=no git@github.com; git fetch upstream; git rebase upstream/main; git push -u origin $PR_BRANCH -f'

gh auth login --with-token < /secrets/github-secrets/token

PR_EXISTS=$(gh pr list | grep -c "${PR_BRANCH}" || true)
if [ $PR_EXISTS -eq 0 ]; then
  gh pr create --title "$PR_TITLE" --body "$PR_BODY"
fi
