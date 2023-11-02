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

# args: <trace|log|docker> <true|false> <-c|-eu -o pipefail -c>

ACTION="$1"
TRACE="$2"
if [ "$TRACE" = "true" ]; then
    >&2 echo "Shell trace: $@"

    if [ -n "${LOGGING_TARGET:-}" ]; then
        >&2 echo "LOGGING_TARGET set to: ${LOGGING_TARGET}"
    fi
    
    if [ -n "${RUN_IN_DOCKER_ARGS:-}" ]; then
        >&2 echo "RUN_IN_DOCKER_ARGS set to: ${RUN_IN_DOCKER_ARGS}"
    fi
fi

shift
shift

# remove action and shellflags up to the -c
for var; do
    shift
    [ "$var" = '-c' ] && break;
done

if [ -z "${LOGGING_TARGET:-}" ] || [ "$ACTION" = "trace" ]; then
    eval "$@"
    exit $?
fi

# in case of recursive make calls, unset the TARGET env var to avoid
# logging in the child call when not desired
TARGET=$LOGGING_TARGET
unset LOGGING_TARGET

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
source "${SCRIPT_ROOT}/common.sh"

DATE=$(build::find::gnu_variant_on_mac date)
DATE_NANO=$(if [ "$(uname -s)" = "Linux" ] || [ "$DATE" = "gdate" ]; then echo %3N; fi)

START_TIME=$($DATE +%s.$DATE_NANO)
if [ "$ACTION" = "docker" ]; then
    TARGET="run-in-docker/$TARGET"
fi

echo -e "\n------------------- $($DATE +"%Y-%m-%dT%H:%M:%S.$DATE_NANO%z") $([ -n "${DOCKER_RUN_BASE_DIRECTORY:-}" ] && echo "(In Docker) ")Starting target=$TARGET -------------------"
if [ "$ACTION" = "docker" ]; then
    echo "($(pwd)) \$ $RUN_IN_DOCKER_ARGS"
    eval $SCRIPT_ROOT/run_target_docker.sh $RUN_IN_DOCKER_ARGS
else
    echo "($(pwd)) \$ $@"
    eval "$@"
fi
echo -e "------------------- $($DATE +"%Y-%m-%dT%H:%M:%S.$DATE_NANO%z") $([ -n "${DOCKER_RUN_BASE_DIRECTORY:-}" ] && echo "(In Docker) ")Finished target=$TARGET duration=$(echo $($DATE +%s.$DATE_NANO) - $START_TIME | bc) seconds -------------------\n"
