#!/usr/bin/env bash

# Copyright 2014 The Kubernetes Authors.
# - https://github.com/kubernetes/kubernetes/blob/master/hack/lib/version.sh
# Copyright 2020 The Kubernetes Authors.
# - https://github.com/kubernetes-sigs/cluster-api/blob/main/hack/version.sh
# Modifications Copyright Amazon.com Inc. or its affiliates. All Rights Reserved. Licensed under the Apache2 License

# Modifications:
# Adapted upsteam code from kubernetes/kubernetes and CAPI to support the CAPI use case, but also be generic for other
# projects built in this repo

set -o errexit
set -o nounset
set -o pipefail

version::get_version_vars() {
    local -r repo=$1
    GIT_RELEASE_TAG=$(git -C $repo describe --match 'v[0-9]*.[0-9]*.[0-9]**' --abbrev=0 --tags)
    GIT_RELEASE_COMMIT=$(git -C $repo rev-list -n 1  "${GIT_RELEASE_TAG}")

    GIT_COMMIT=$GIT_RELEASE_COMMIT

    # from k8s.io/hack/lib/version.sh
    # Use git describe to find the version based on tags.
    if GIT_VERSION=$GIT_RELEASE_TAG; then
        # Try to match the "git describe" output to a regex to try to extract
        # the "major" and "minor" versions and whether this is the exact tagged
        # version or whether the tree is between two tagged versions.
        if [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
            GIT_MAJOR=${BASH_REMATCH[1]}
            GIT_MINOR=${BASH_REMATCH[2]}
        fi

        # If GIT_VERSION is not a valid Semantic Version, then refuse to build.
        if ! [[ "${GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
            echo "GIT_VERSION should be a valid Semantic Version. Current value: ${GIT_VERSION}"
            echo "Please see more details here: https://semver.org"
            exit 1
        fi
    fi


}

# Prints the value that needs to be passed to the -ldflags parameter of go build
version::ldflags() {
    local -r repo=$1
    local -r package_prefix=$2
    version::get_version_vars $repo

    # from k8s.io/hack/lib/version.sh
    local -a ldflags
    function add_ldflag() {
        local key=${1}
        local val=${2}
        ldflags+=(
            "-X '${package_prefix}.${key}=${val}'"
        )
    }
    DATE=date
    if which gdate &>/dev/null; then
        DATE=gdate
    elif which gnudate &>/dev/null; then
        DATE=gnudate
    fi
    
    # buildDate is not actual buildDate to avoid it breaking reproducible checksums
    # instead it is the date of the last commit, either the upstream TAG commit or the latest patch applied
    SOURCE_DATE_EPOCH=$(git -C $repo log -1 --format=%at)
    add_ldflag "buildDate" "$(${DATE} --date=@${SOURCE_DATE_EPOCH} -u +'%Y-%m-%dT%H:%M:%SZ')"
    add_ldflag "gitCommit" "${GIT_COMMIT}"
    add_ldflag "gitTreeState" "clean"
    add_ldflag "gitMajor" "${GIT_MAJOR}"
    add_ldflag "gitMinor" "${GIT_MINOR}"
    add_ldflag "gitVersion" "${GIT_VERSION}"
    add_ldflag "gitReleaseCommit" "${GIT_RELEASE_COMMIT}"

    # The -ldflags parameter takes a single string, so join the output.
    echo "${ldflags[*]-}"
}

version::ldflags $1 $2
