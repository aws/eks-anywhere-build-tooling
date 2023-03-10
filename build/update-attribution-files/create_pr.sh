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

if [[ -z "${JOB_TYPE:-}" ]] && [[ -z "${CODEBUILD_CI:-}" ]]; then
    exit 0
fi

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"

source $SCRIPT_ROOT/../lib/common.sh

SKIP_PR="${1:-false}"

ORIGIN_ORG="eks-distro-pr-bot"
UPSTREAM_ORG="aws"
REPO="eks-anywhere-build-tooling"

MAIN_BRANCH="main"

if [[ -n "${CODEBUILD_SOURCE_VERSION:-}" ]]; then
    MAIN_BRANCH="$CODEBUILD_SOURCE_VERSION"
fi

if [[ -n "${PULL_BASE_REF:-}" ]]; then
    MAIN_BRANCH="$PULL_BASE_REF"
fi

cd ${SCRIPT_ROOT}/../../
git config --global push.default current
git config user.name "EKS Distro PR Bot"
git config user.email "aws-model-rocket-bots+eksdistroprbot@amazon.com"
git config remote.upstream.url >&- || git remote add upstream https://github.com/${UPSTREAM_ORG}/${REPO}.git
git config remote.bot.url >&- || git remote add bot https://github.com/${ORIGIN_ORG}/${REPO}.git

build::common::echo_and_run git status

# Files have already changed, stash to perform rebase
build::common::echo_and_run git stash

build::common::echo_and_run git checkout $MAIN_BRANCH

# avoid hitting github limits in presubmits
if [[ "${JOB_TYPE:-}" == "periodic" ]] || [[ -n "${CODEBUILD_CI:-}" ]]; then
    build::common::echo_and_run retry git fetch -q upstream

    # there will be conflicts before we are on the bots fork at this point
    # -Xtheirs instructs git to favor the changes from the current branch
    build::common::echo_and_run git rebase -Xtheirs upstream/$MAIN_BRANCH
fi

if [ "$(git stash list)" != "" ]; then
    build::common::echo_and_run git stash pop
    build::common::echo_and_run git status
fi

function commit_push()
{
    local -r pr_branch="$1"
    local force_push="$2"

    if [[ "$force_push" = "true" ]]; then
        force_push="-f"
    else
        force_push=""
    fi

    if build::common::echo_and_run git push $force_push -u bot $pr_branch; then
        return
    fi

    # Other jobs could be running and updating the same branch, which is expected
    # attempt to rebase and retry
    build::common::echo_and_run git fetch -q bot
    build::common::echo_and_run git branch --set-upstream-to=bot/$pr_branch $pr_branch
    
    # rebase succeeded, return 1 to retry the push
    if build::common::echo_and_run git rebase bot/$pr_branch; then
        return 1
    fi

    git status
    git diff HEAD -- $(git diff --name-only --diff-filter=U)

    # rebase failed with conflict, exit with error
    exit 1
}

function pr_create()
{
    local -r pr_title="$1"
    local -r pr_branch="$2"
    local -r pr_body="$3"

    local -r max_wait=5

    # to try and avoid burning through our rate limit, wait a random amount of time
    # check for the PR to already exist, if it doesnt, try and create it
    sleep $((RANDOM % $max_wait))

    local pr_exists=$(GH_PAGER='' gh pr list --json number -H "$pr_branch")
    if [ "$pr_exists" != "[]" ]; then
        # already exists
        echo "PR already exists."
        return
    fi

    echo "Creating PR."
    if gh pr create --title "$pr_title" --body "$pr_body" --base $MAIN_BRANCH --head $ORIGIN_ORG:$pr_branch; then
        echo "PR created."
        return
    fi

    sleep $((RANDOM % $max_wait))

    pr_exists=$(GH_PAGER='' gh pr list --json number -H "$pr_branch")
    if [ "$pr_exists" != "[]" ]; then
        # already exists
        echo "PR already exists."
        return
    fi

    # PR does not exist and creation failed, retry
    return 1
}

function pr::create()
{
    local -r pr_title="$1"
    local -r commit_message="$2"
    local -r pr_branch="$3"
    local -r pr_body="$4"
    local -r force_push="${5:-true}"

    build::common::echo_and_run git checkout -B $pr_branch

    git diff --staged
    local -r files_added=$(git diff --staged --name-only)
    if [ "$files_added" = "" ]; then
        echo "No files changed to commit"
        return 0
    fi

    build::common::echo_and_run git commit -m "$commit_message" || true

    # need to check a codebuild job "type"
    if [ "$JOB_TYPE" != "periodic" ] && [[ -z "${CODEBUILD_CI:-}" ]]; then
        return 0
    fi

    github::auth

    retry commit_push "$pr_branch" "$force_push"
    
    if [[ "${SKIP_PR}" == "false" ]]; then
        retry pr_create "$pr_title" "$pr_branch" "$pr_body"
    fi  
}

function github::auth(){
    if [[ -z "${GITHUB_TOKEN:-}" ]] && [[ ! -f /secrets/github-secrets/token ]]; then
        echo "Missing GITHUB_TOKEN or /secrets/github-secrets/token, cannot authenticate!"
        exit 1
    fi

    if [ -f /secrets/github-secrets/token ]; then
        gh auth login --with-token < /secrets/github-secrets/token
    fi

    gh auth setup-git
    gh api rate_limit
    
    gh repo set-default "${UPSTREAM_ORG}/${REPO}"
}

function pr::create::pr_body(){
    pr_body=""
    case $1 in
    attribution)
        pr_body=$(cat <<'EOF'
This PR updates the ATTRIBUTION.txt files across all dependency projects if there have been changes.

These files should only be changing due to project GIT_TAG bumps or Golang version upgrades. If changes are for any other reason, please review carefully before merging!
EOF
)
        ;;
    checksums)
        pr_body=$(cat <<EOF
This PR updates the CHECKSUMS files across all dependency projects if there have been changes.

These files should only be changing due to project GIT_TAG bumps or Golang version upgrades. If changes are for any other reason, please review carefully before merging!

These files were generated using $CODEBUILD_BUILD_IMAGE
EOF
)
        ;;
    makehelp)
        pr_body=$(cat <<'EOF'
This PR updates the Help.mk files across all dependency projects if there have been changes.
EOF
)
        ;;
    go-mod)
        pr_body=$(cat <<'EOF'
This PR updates the checked in go.mod and go.sum files across all dependency projects to support automated vulnerability scanning.
EOF
)
        ;;
    *)
        echo "Invalid argument: $1"
        exit 1
        ;;
    esac
    PROW_BUCKET_NAME=$(echo $JOB_SPEC | jq -r ".decoration_config.gcs_configuration.bucket" | awk -F// '{print $NF}')
    full_pr_body=$(printf "%s\n\n/hold\n\nBy submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice." "$pr_body")

    printf "$full_pr_body"
}

function pr::create::attribution() {
    local -r pr_title="Update ATTRIBUTION.txt files"
    local -r commit_message="[PR BOT] Update ATTRIBUTION.txt files"
    local -r pr_branch="attribution-files-update-$MAIN_BRANCH"
    local -r pr_body=$(pr::create::pr_body "attribution")

    pr::create "$pr_title" "$commit_message" "$pr_branch" "$pr_body"
}

function pr::create::checksums() {
    local -r pr_title="Update CHECKSUMS files"
    local pr_branch="checksums-files-update-$MAIN_BRANCH"
    local -r pr_body=$(pr::create::pr_body "checksums")

    if [ -n "${CODEBUILD_RESOLVED_SOURCE_VERSION:-}" ]; then
        pr_branch+="-$CODEBUILD_RESOLVED_SOURCE_VERSION"
    fi

    # This file is being added as it may have been updated by the last lines of ./build/lib/update_go_versions.sh,
    # which replaces the go version in this script with the go version(s) in the builder base if they are newer 
    # when running `make update-checksum-files`
    git add ./build/lib/install_go_versions.sh

    # stash checksums and attribution and help.mk files
    git stash --keep-index
    
    pr::create "$pr_title" "[PR BOT] Update install_go_versions.sh" "$pr_branch" "$pr_body" "false"

    if [ "$(git stash list)" != "" ]; then
        build::common::echo_and_run git stash pop
        build::common::echo_and_run git status
    fi

    # Add checksum files
    for FILE in $(find . -type f -name CHECKSUMS); do
        pr::file:add $FILE
    done

    local -r changed_files=$(git diff --staged --name-only | grep "^projects" | cut -d '/' -f2- | uniq | tr  '\n' ' ')
    local -r commit_message="[PR BOT] Update $changed_files"
 
    pr::create "$pr_title" "$commit_message" "$pr_branch" "$pr_body" "false"

    # in codebuild we run this seperately after all other jobs has finished
    # in that case there will be no files changed so no commit, but we still want to open the pr
    if [[ "${SKIP_PR}" == "false" ]]; then
        github::auth
        retry pr_create "$pr_title" "$pr_branch" "$pr_body"
    fi
}

function pr::create::help() {
    local -r pr_title="Update Makefile generated help"
    local -r commit_message="[PR BOT] Update Help.mk files"
    local -r pr_branch="help-makefiles-update-$MAIN_BRANCH"
    local -r pr_body=$(pr::create::pr_body "makehelp")

    pr::create "$pr_title" "$commit_message" "$pr_branch" "$pr_body"
}

function pr::create::go-mod() {
    local -r pr_title="Update go.mod files"
    local -r commit_message="[PR BOT] Update go.mod files"
    local -r pr_branch="go-mod-update-$MAIN_BRANCH"
    local -r pr_body=$(pr::create::pr_body "go-mod")

    pr::create "$pr_title" "$commit_message" "$pr_branch" "$pr_body"
}

function pr::file:add() {
    local -r file="$1"

    if git check-ignore -q $FILE; then
        return
    fi

    local -r diff="$(git diff --ignore-blank-lines --ignore-all-space $FILE)"
    if [[ -z $diff ]]; then
        return
    fi

    git add $file
}

echo "Adding Checksum files"

pr::create::checksums

git checkout $MAIN_BRANCH

if [ "$(git stash list)" != "" ]; then
    build::common::echo_and_run git stash pop
    build::common::echo_and_run git status
fi

echo "Adding ATTRIBUTION files"
# Add attribution files
for FILE in $(find . -type f \( -name "*ATTRIBUTION.txt" ! -path "*/_output/*" \)); do
    pr::file:add $FILE
done

# stash help.mk files
git stash --keep-index

pr::create::attribution

git checkout $MAIN_BRANCH

if [ "$(git stash list)" != "" ]; then
    build::common::echo_and_run git stash pop
    build::common::echo_and_run git status
fi

echo "Adding Help.mk files"

# Add help.mk/Makefile files
for FILE in $(find . -type f \( -name Help.mk -o -name Makefile \)); do
    pr::file:add $FILE
done

# Not handling go.sum/go.mod yet in eks-a
# stash go.sum files
#git stash --keep-index

pr::create::help

#git checkout $MAIN_BRANCH

# if [ "$(git stash list)" != "" ]; then
#     build::common::echo_and_run git stash pop
#     build::common::echo_and_run git status
# fi

#echo "Adding go.mod and go.sum files"
# Add go.mod files
# for FILE in $(find . -type f \( -name go.sum -o -name go.mod \)); do    
#    git check-ignore -q $FILE || git add $FILE
# done

# pr::create::go-mod
