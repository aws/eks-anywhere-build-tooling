#!/usr/bin/env bash
# Copyright 2020 Amazon.com Inc. or its affiliates. All Rights Reserved.
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

set -o errexit
set -o nounset
set -o pipefail

GROUP_ID="$1"
USER_ID="$2"

# when running the build container on linux, the user in the container needs 
# to match the host user id and group id
# otherwise there will be perms issues due to go mods being downloaded in the
# container as root even though the host user is not

sed -i 's/^CREATE_MAIL_SPOOL=yes/CREATE_MAIL_SPOOL=no/' /etc/default/useradd
groupadd -g 100 users
groupadd --gid "$GROUP_ID" matchinguser
useradd  --no-create-home --uid "$USER_ID" --gid "$GROUP_ID" matchinguser
mkdir -p /home/matchinguser
chown -R matchinguser:matchinguser /home/matchinguser
