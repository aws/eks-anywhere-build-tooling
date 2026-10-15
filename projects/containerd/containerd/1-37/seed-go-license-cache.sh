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

set -euo pipefail

project_root="${1:?Specify containerd project root}"
golang_version="${2:?Specify Go version}"
base_directory="$(cd "$project_root/../../.." && pwd -P)"

if [[ "$(uname -s)" != Linux || ! -f /buildkit.sh ]]; then
    echo "Run Go license generation in builder-base"
    exit 1
fi

source "$base_directory/build/lib/common.sh"
cd "$project_root/containerd"
export GOTOOLCHAIN=local
build::common::use_go_version "$golang_version"
compiler_version="$(go env GOVERSION)"
compiler_root="$(go env GOROOT)"
test "$(head -n 1 "$compiler_root/VERSION")" = "$compiler_version"
test -s "$compiler_root/LICENSE"
sha256sum "$compiler_root/LICENSE"

node - "$compiler_root/LICENSE" "$compiler_version" <<'NODE'
const fs = require('fs');
const os = require('os');
const path = require('path');
const [licenseFile, version] = process.argv.slice(2);
const cacheDirectory = path.join(os.homedir(), '.generate-attribution');
const cacheFile = path.join(cacheDirectory, '.cache');
const license = fs.readFileSync(licenseFile, 'utf8');
if (!license.startsWith('Copyright 2009 The Go Authors.') ||
    !license.includes('Redistribution and use in source and binary forms')) {
    throw new Error('Unexpected compiler license contents');
}
fs.mkdirSync(cacheDirectory, { recursive: true });
const cache = fs.existsSync(cacheFile)
    ? JSON.parse(fs.readFileSync(cacheFile, 'utf8'))
    : {};
const url = `https://github.com/golang/go/blob/${version}/LICENSE?raw=true`;
cache[url] = license;
const temporaryFile = `${cacheFile}.containerd-${process.pid}`;
fs.writeFileSync(temporaryFile, JSON.stringify(cache));
fs.renameSync(temporaryFile, cacheFile);
console.log(`Go ${version} license cached from its compiler distribution`);
NODE
