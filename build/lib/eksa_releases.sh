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

# Key for bundles manifest is environment/latest
# example key = "dev/release-1.20"
declare -A BUNDLE_MANIFEST=()

function build::eksa_releases::load_bundle_manifest() {
  local -r dev_release=${1-false}
  local -r latest=${2-latest}
  local -r echo=${3-true}
  oldopt="$(set +o)"
  set +o nounset
  set +x

  local -r bundle_manifest_key=$(build::eksa_releases::get_bundle_manifest_key $dev_release $latest)
  if [ ! ${BUNDLE_MANIFEST[$bundle_manifest_key]+1} ]; then
    local -r release_manifest_url=$(build::eksa_releases::get_eksa_release_manifest_url $dev_release $latest)
    local -r release_manifest=$(curl -s --retry 5 $release_manifest_url)

    # The EKSA_RELEASE_VERSION variable is set only when this script is run from the image-builder CLI.
    # When running the image-builder CLI in dev, the EKSA_RELEASE_VERSION will be set to a dev version
    # such as v0.0.0-dev, but without the build metadata. This incomplete version is not available in the
    # dev EKS-A releases manifest and so the yq search will fail. Hence if are running in dev, we append
    # a wildcard build metadata to the EKSA_RELEASE_VERSION var that will make it pass the yq select check.
    EKSA_RELEASE_VERSION="${EKSA_RELEASE_VERSION:-}"
    local -r eksa_release_version=${EKSA_RELEASE_VERSION:-$(echo "$release_manifest" | yq e ".spec.latestVersion" -)}
    if [ $dev_release = true ] && [ -n "$EKSA_RELEASE_VERSION" ]; then
      eksa_release_version="$eksa_release_version+build.*"
    fi
    local -r bundle_manifest_url=$(echo "$release_manifest" | yq e ".spec.releases[] | select(.version == \"$eksa_release_version\") .bundleManifestUrl" -)
    BUNDLE_MANIFEST[$bundle_manifest_key]=$(curl -s --retry 5 "$bundle_manifest_url" | yq)
  fi
  if $echo; then
    echo "${BUNDLE_MANIFEST[$bundle_manifest_key]}"
  fi
  eval "$oldopt"
}

function build::eksa_releases::get_eksa_component_asset_url() {
  local -r component=$1
  local -r asset=$2
  local -r release_branch=$3
  local -r dev_release=${4-false}
  local -r latest=${5-main}

  build::eksa_releases::get_eksa_component_asset_path $release_branch ".$component.$asset.uri" $dev_release $latest
}

function build::eksa_releases::get_eksa_component_asset_artifact_checksum() {
  local -r component=$1
  local -r asset=$2
  local -r type=$3
  local -r release_branch=$4
  local -r dev_release=${5-false}
  local -r latest=${6-main}

  if [[ $type != "sha256" ]] && [[ $type != "sha512" ]]; then
    echo "Invalid shasum type. Allowed types are sha256 and sha512"
  fi

  build::eksa_releases::get_eksa_component_asset_path $release_branch ".$component.$asset.$type" $dev_release $latest
}

function build::eksa_releases::get_eksa_component_asset_path() {
  local -r release_branch=$1
  local -r path=$2
  local -r dev_release=${3-false}
  local -r latest=${4-main}

  oldopt="$(set +o)"
  set +x

  # Get latest bundle manifest url
  local bundle_manifest=$(build::eksa_releases::load_bundle_manifest $dev_release $latest)
  local kube_version=$(echo $release_branch | sed 's/\-/\./g')

  local query=".spec.versionsBundles[] | select(.kubeVersion == \"$kube_version\") $path"

  echo "$bundle_manifest" | yq e "$query" -

  eval "$oldopt"
}

function build::eksa_releases::get_eksa_release_manifest_url() {
  local -r dev_release=${1-false}
  local -r latest=${2-latest}

  if [[ $dev_release == false ]]; then
    echo "https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml"
  elif [[ $latest == "latest" ]]; then
    echo "https://dev-release-assets.eks-anywhere.model-rocket.aws.dev/eks-a-release.yaml"
  else
    echo "https://dev-release-assets.eks-anywhere.model-rocket.aws.dev/$latest/eks-a-release.yaml"
  fi
}

function build::eksa_releases::get_bundle_manifest_key() {
  local -r dev_release=${1-false}
  local -r latest=${2-latest}

  local environment="prod"
  if [ $dev_release == true ]; then
    environment="dev"
  fi

  echo "$environment/$latest"
}
