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

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd -P)"
LATEST_RELEASE_BRANCH=''
PROD_DOMAIN="distro.eks.amazonaws.com"
DEV_DOMAIN="eks-d-postsubmit-artifacts.s3.us-west-2.amazonaws.com"
declare -A RELEASE_YAML=()

function build::eksd_releases::load_release_yaml() {
    local -r release_branch=$1
    local -r echo=${2-true}
    oldopt=$-
    set +o nounset
    set +x

    # if key exists, 1 is returned which would resolve to true
    if [ ! ${RELEASE_YAML[$release_branch]+1} ]; then
        local -r yaml_url=$(build::eksd_releases::get_release_yaml_url ${release_branch})
        RELEASE_YAML[$release_branch]=$(curl -s --retry 5 $yaml_url)
    fi
    if $echo; then
        echo "${RELEASE_YAML[$release_branch]}"
    fi
    set -$(echo $oldopt | sed 's/c//') # remove -c from options for when this is invoked with bash -c
}

function build::eksd_releases::get_release_yaml_url() {
    local -r release_branch=$1
    local -r release_number=$(yq e ".releases[] | select(.branch==\"${release_branch}\").number" ${REPO_ROOT}/EKSD_LATEST_RELEASES)
    local -r dev=$(yq e ".releases[] | select(.branch==\"${release_branch}\").dev" ${REPO_ROOT}/EKSD_LATEST_RELEASES)
    local -r yaml_path="kubernetes-${release_branch}/kubernetes-${release_branch}-eks-${release_number}.yaml"
    local yaml_url="https://$PROD_DOMAIN/$yaml_path"
    if [[ $dev == "true" ]]; then
      yaml_url="https://$DEV_DOMAIN/$yaml_path"
    fi
    echo "$yaml_url"
}

function build::eksd_releases::get_release_branch() {
    if [ -z $LATEST_RELEASE_BRANCH ]; then
        LATEST_RELEASE_BRANCH=$(yq e ".latest" ${REPO_ROOT}/EKSD_LATEST_RELEASES)
    fi
    echo $LATEST_RELEASE_BRANCH
}

function build::eksd_releases::get_eksd_release_number() {
    local -r release_branch=$1
    
    oldopt=$-
    set +x
    
    local -r yaml=$(build::eksd_releases::load_release_yaml $release_branch)

    echo "$yaml" |  yq e ".spec.number" -
    
    set -$oldopt
}

function build::eksd_releases::get_eksd_release_name() {
    local -r release_branch=$1
    
    oldopt=$-
    set +x
    
    local -r yaml=$(build::eksd_releases::load_release_yaml $release_branch)

    echo "$yaml" |  yq e ".metadata.name" -
    
    set -$oldopt
}

function build::eksd_releases::get_eksd_component_version() {
    local -r component=$1
    local -r release_branch=$2
    
    oldopt=$-
    set +x
    
    local -r yaml=$(build::eksd_releases::load_release_yaml $release_branch)

    echo "$yaml" | yq e ".status.components[] | select(.name == \"$component\") .gitTag" -
    
    set -$oldopt
}

function build::eksd_releases::get_eksd_component_asset_path() {
    local -r component=$1
    local -r release_branch=$2
    local -r path=$3
    local -r arch=$4
    local -r asset=${5-}
    
    oldopt=$-
    set +x
    
    local -r yaml=$(build::eksd_releases::load_release_yaml $release_branch)

    local query=".status.components[] | select(.name == \"$component\") .assets[]"

    if [[ $path == ".archive.uri" ]] || [[ $path == ".archive.sha256" ]]; then
        query+=" | select(.type == \"Archive\")"
    fi

    if [ -z $asset ]; then
        query+=" | select(.arch[0] == \"$arch\") $path"  
    else
        query+=" | select(.name == \"$asset\") | select(.arch[0] == \"$arch\") $path"
    fi

    echo "$yaml" | yq e "$query" -
    
    set -$(echo $oldopt | sed 's/c//') # remove -c from options for when this is invoked with bash -c
}

function build::eksd_releases::get_eksd_component_asset_url() {
    local -r component=$1
    local -r asset=$2
    local -r release_branch=${3-$(build::eksd_releases::get_release_branch)}
    local -r arch=${4-amd64}

    build::eksd_releases::get_eksd_component_asset_path $component $release_branch ".archive.uri" $arch $asset 
}

function build::eksd_releases::get_eksd_component_asset_sha() {
    local -r component=$1
    local -r asset=$2
    local -r release_branch=${3-$(build::eksd_releases::get_release_branch)}
    local -r arch=${4-amd64}

    build::eksd_releases::get_eksd_component_asset_path $component $release_branch ".archive.sha256" $arch $asset
}


function build::eksd_releases::get_eksd_kubernetes_asset_url() {
    local -r asset=$1
    local -r release_branch=${2-$(build::eksd_releases::get_release_branch)}
    local -r arch=${3-amd64}

    build::eksd_releases::get_eksd_component_asset_url "kubernetes" $asset $release_branch $arch
}

function build::eksd_releases::get_eksd_kubernetes_image_url() {
    local -r asset=$1
    local -r release_branch=${2-$(build::eksd_releases::get_release_branch)}
    local -r arch=${3-amd64}

    build::eksd_releases::get_eksd_component_asset_path "kubernetes" $release_branch ".image.uri" $arch $asset 
}

function build::eksd_releases::get_eksd_component_asset_image_tag() {
    local -r component=$1
    local -r asset=$2
    local -r release_branch=$3
    local -r arch=${4-amd64}

    build::eksd_releases::get_eksd_component_asset_path $component $release_branch ".image.uri" $arch $asset | awk -F: '{print $2}'
}

function build::eksd_releases::get_eksd_component_url() {
    local -r component=$1
    local -r release_branch=${2-$(build::eksd_releases::get_release_branch)}
    local -r arch=${3-amd64}

    build::eksd_releases::get_eksd_component_asset_path $component $release_branch ".archive.uri" $arch
}


function build::eksd_releases::get_eksd_component_sha() {
    local -r component=$1
    local -r release_branch=${2-$(build::eksd_releases::get_release_branch)}
    local -r arch=${3-amd64}

    build::eksd_releases::get_eksd_component_asset_path $component $release_branch ".archive.sha256" $arch
}

function build::eksd_releases::get_eksd_image_repo() {
    local -r release_branch=${1-$(build::eksd_releases::get_release_branch)}
    
    build::eksd_releases::get_eksd_kubernetes_image_url "kube-apiserver-image" $release_branch | sed -E 's,/kubernetes/kube-apiserver.*,,'
}

function build::eksd_releases::get_eksd_kubernetes_asset_base_url() {
    local -r release_branch=${1-$(build::eksd_releases::get_release_branch)}
    
    local -r kube_version=$(build::eksd_releases::get_eksd_component_version "kubernetes" $release_branch)
    local -r apiserver="bin/linux/amd64/kube-apiserver"
    build::eksd_releases::get_eksd_kubernetes_asset_url $apiserver $release_branch | sed -E "s,/$kube_version/$apiserver.*,,"
}

function build::eksd_releases::get_eksd_cni_asset_base_url() {
    local -r release_branch=${1-$(build::eksd_releases::get_release_branch)}
    
    local -r cni_plugins_version=$(build::eksd_releases::get_eksd_component_version "cni-plugins" $release_branch)
    build::eksd_releases::get_eksd_component_url "cni-plugins" $release_branch | sed -E "s,/$cni_plugins_version/cni-plugins.*,,"
}
