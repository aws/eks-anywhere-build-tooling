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


set -x
set -o errexit
set -o pipefail

TAG="${1?Specify first argument - git version tag}"
HELM_REPO_VERSION="${2?Specify second argument - helm repo version}"
ARTIFACTS_PATH="${3?Specify third argument - artifacts path}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
HELM_BIN="${MAKE_ROOT}/_output/helm-bin"

mkdir -p $ARTIFACTS_PATH

function build::install::helm(){
  mkdir -p $HELM_BIN
  export PATH=$HELM_BIN:$PATH
  curl -s https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | HELM_INSTALL_DIR=$HELM_BIN bash
}

function build::cilium::manifests(){
  mkdir -p _output/manifests/cilium/$TAG
  helm repo add cilium https://helm.cilium.io
  helm template cilium cilium/cilium --version $HELM_REPO_VERSION --namespace kube-system -f manifests/cilium-eksa.yaml > _output/manifests/cilium/${TAG}/cilium.yaml
}

build::install::helm
build::cilium::manifests

cp -rf _output/manifests/cilium $ARTIFACTS_PATH
