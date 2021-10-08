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
set -o nounset
set -o pipefail


BOTTLEROCKET_DOWNLOAD_PATH="${1?Specify first argument - Download path for Bottlerocket-related files}"
CARGO_HOME="${2?Specify second argument - Root directory for Cargo installation}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
RUSTUP_HOME=$MAKE_ROOT/_output/rustup
BOTTLEROCKET_ROOT_JSON_URL="https://cache.bottlerocket.aws/root.json"

mkdir -p $BOTTLEROCKET_DOWNLOAD_PATH
mkdir -p $CARGO_HOME
mkdir -p $RUSTUP_HOME
export PATH=$CARGO_HOME/bin:$PATH

# This configuration supports local installations and checksum validations
# of root.json file
export BOTTLEROCKET_ROOT_JSON_PATH=$BOTTLEROCKET_DOWNLOAD_PATH/root.json
envsubst '$BOTTLEROCKET_ROOT_JSON_PATH' \
    < $MAKE_ROOT/bottlerocket-root-json-checksum \
    > $BOTTLEROCKET_DOWNLOAD_PATH/bottlerocket-root-json-checksum
curl $BOTTLEROCKET_ROOT_JSON_URL -o $BOTTLEROCKET_ROOT_JSON_PATH
sha512sum -c $BOTTLEROCKET_DOWNLOAD_PATH/bottlerocket-root-json-checksum

# On AL2, the Cargo build system requires the openssl-devel package
# for installing OpenSSL libraries and the pkgconfig utility to 
# locate these headers/libs
if [ "$(uname)" = "Linux" ]; then
    yum install -y openssl-devel pkgconfig
fi

# This code installs the Rust toolchain manager called rustup along
# with other Rust binaries such as rustc, rustfmt. It also installs Cargo,
# the Rust package manager which is then used to install Tuftool.
curl https://sh.rustup.rs -sSf | CARGO_HOME=$CARGO_HOME RUSTUP_HOME=$RUSTUP_HOME sh -s -- -y
$CARGO_HOME/bin/rustup default stable
CARGO_NET_GIT_FETCH_WITH_CLI=true $CARGO_HOME/bin/cargo install --force --root $CARGO_HOME tuftool
