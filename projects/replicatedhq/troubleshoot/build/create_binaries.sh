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

TAG="$1"
BIN_PATH="$2"
OS="$3"
ARCH="$4"


VERSION_PACKAGE="github.com/replicatedhq/troubleshoot/pkg/version"
LDFLAGS="-s -w -buildid=''	-X $VERSION_PACKAGE.version=$TAG \
	-X $VERSION_PACKAGE.gitSHA=$(git rev-list -n 1 $TAG)"
BUILDTAGS="netgo containers_image_ostree_stub exclude_graphdriver_devicemapper exclude_graphdriver_btrfs containers_image_openpgp"

GOOS=$OS GOARCH=$ARCH go build -trimpath -tags "$BUILDTAGS" -installsuffix netgo -ldflags "$LDFLAGS" -o $BIN_PATH/support-bundle \
  github.com/replicatedhq/troubleshoot/cmd/troubleshoot
