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


# Change Go Version hack
git clone https://github.com/syndbg/goenv.git $HOME/.goenv 
echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile
echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
echo 'export PATH="$GOROOT/bin:$PATH"' >> ~/.bash_profile
echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bash_profile
source ~/.bash_profile
goenv install 1.17.2
goenv local 1.17.2
go version


# Clone Addons
# git clone https://github.com/aws/modelrocket-add-ons.git
git clone "https://git-codecommit.${aws_region}.amazonaws.com/v1/repos/${org}.${repo}"


# cd modelrocket-add-ons/generatebundlefile
ls -la
cd aws.modelrocket-add-ons/

ls -la
cd generatebundlefile/

./vend.sh

make run