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

org="aws"
repo="modelrocket-add-ons"
aws_region="us-west-2"

# Clone Addons
# git clone https://github.com/aws/modelrocket-add-ons.git
git clone "https://git-codecommit.${aws_region}.amazonaws.com/v1/repos/${org}.${repo}"

# cd modelrocket-add-ons/generatebundlefile
ls -la
cd aws.modelrocket-add-ons/

git checkout dont-delete/codebuild-fork

ls -la
cd generatebundlefile/


which go
go version

# Set up specific go version by using go get, additional versions apart from default can be installed by calling
# the function again with the specific parameter.
setupgo() {
    local -r version=$1
    go get golang.org/dl/go${version}
    go${version} download
    # Removing the last number as we only care about the major version of golang
    local -r majorversion=${version%.*}
    mkdir -p ${GOPATH}/go${majorversion}/bin
    ln -s ${GOPATH}/bin/go${version} ${GOPATH}/go${majorversion}/bin/go
    ln -s /root/sdk/go${version}/bin/gofmt ${GOPATH}/go${majorversion}/bin/gofmt
    go version
}

setupgo "${GOLANG117_VERSION:-1.17.5}"

go version

./vend.sh

make run

pwd=$(pwd)

echo $IMAGE_TAG
echo $DIGEST

go1.17.5 run . --input "$pwd/data/input_120.yaml"