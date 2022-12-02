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

BUILD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/" && pwd -P)"
source "${BUILD_ROOT}/eksd_releases.sh"

USE_DOCKER="${USE_DOCKER:-false}"
USE_BUILDCTL="${USE_BUILDCTL:-false}"

function build::common::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=tar
  if which gtar &>/dev/null; then
      TAR=gtar
  elif which gnutar &>/dev/null; then
      TAR=gnutar
  fi
  if ! "${TAR}" --version | grep -q GNU; then
    echo "  !!! Cannot find GNU tar. Build on Linux or install GNU tar"
    echo "      on Mac OS X (brew install gnu-tar)."
    return 1
  fi
}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/kubernetes is what is
# being packaged.
function build::common::create_tarball() {
  build::common::ensure_tar

  local -r tarfile=$1
  local -r stagingdir=$2
  local -r repository=$3

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" $repository --owner=0 --group=0
}

# Generate shasum of tarballs. $1 is the directory of the tarballs.
function build::common::generate_shasum() {

  local -r tarpath=$1

  echo "Writing artifact hashes to shasum files..."

  if [ ! -d "$tarpath" ]; then
    echo "  Unable to find tar directory $tarpath"
    exit 1
  fi

  cd $tarpath
  for file in $(find . -name '*.tar.gz'); do
    filepath=$(basename $file)
    sha256sum "$filepath" > "$file.sha256"
    sha512sum "$filepath" > "$file.sha512"
  done
  cd -
}

function build::common::upload_artifacts() {
  local -r artifactspath=$1
  local -r artifactsbucket=$2
  local -r projectpath=$3
  local -r buildidentifier=$4
  local -r githash=$5
  local -r latesttag=$6
  local -r dry_run=$7

  if [ "$dry_run" = "true" ]; then
    aws s3 cp "$artifactspath" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts --recursive --dryrun
    aws s3 cp "$artifactspath" "$artifactsbucket"/"$projectpath"/"$latesttag" --recursive --dryrun
  else
    # Upload artifacts to s3 
    # 1. To proper path on s3 with buildId-githash
    # 2. Latest path to indicate the latest build, with --delete option to delete stale files in the dest path
    aws s3 sync "$artifactspath" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts --acl public-read
    aws s3 sync "$artifactspath" "$artifactsbucket"/"$projectpath"/"$latesttag" --delete --acl public-read
  fi
}

function build::gather_licenses() {
  local -r outputdir=$1
  local -r patterns=$2
  local -r golang_version=$3

  # Force deps to only be pulled form vendor directories
  # this is important in a couple cases where license files
  # have to be manually created
  export GOFLAGS=-mod=vendor
  # force platform to be linux to ensure all deps are picked up
  export GOOS=linux 
  export GOARCH=amd64 

  # the version of go used here must be the version go-licenses was installed with
  # by default we use 1.16, but due to changes in 1.17, there are some changes that require using 1.17
  if [[ ${golang_version#1.} -ge 16 ]]; then
    build::common::use_go_version $golang_version
  else
    build::common::use_go_version 1.16
  fi

  if [ -z "${USE_HOST_GO_LICENSES:-}" ]; then    
    build::common::override_missing_tooling go-licenses
  fi  
  
  mkdir -p "${outputdir}/attribution"
  # attribution file generated uses the output go-deps and go-license to gather the necessary
  # data about each dependency to generate the amazon approved attribution.txt files
  # go-deps is needed for module versions
  # go-licenses are all the dependencies found from the module(s) that were passed in via patterns
  echo "($(pwd)) \$ go list -deps=true -json ./..."
  if ! list=$(go list -deps=true -json ./...); then
    printf "$list"
    exit 1 
  fi

  if ! echo $list | jq -s '' > "${outputdir}/attribution/go-deps.json"; then
    exit 1
  fi
  

  # go-licenses can be a bit noisy with its output and lot of it can be confusing 
  # the following messages are safe to ignore since we do not need the license url for our process
  NOISY_MESSAGES="cannot determine URL for|Error discovering license URL|unsupported package host|contains non-Go code|has empty version|vendor.*\.s$"

  echo "($(pwd)) \$ go-licenses save --force $patterns --save_path ${outputdir}/LICENSES"
  go-licenses save --force $patterns --save_path "${outputdir}/LICENSES" &>  >(grep -vE "$NOISY_MESSAGES" >&2)

  echo "($(pwd)) \$ go-licenses csv $patterns"
  go-licenses csv $patterns > "${outputdir}/attribution/go-license.csv" 2>  >(grep -vE "$NOISY_MESSAGES" >&2)

  if cat "${outputdir}/attribution/go-license.csv" | grep -q "^vendor\/golang.org\/x"; then
      echo " go-licenses created a file with a std golang package (golang.org/x/*)"
      echo " prefixed with vendor/.  This most likely will result in an error"
      echo " when generating the attribution file and is probably due to"
      echo " to a version mismatch between the current version of go "
      echo " and the version of go that was used to build go-licenses"
      exit 1
  fi

  if cat "${outputdir}/attribution/go-license.csv" | grep -e ",LGPL-" -e ",GPL-"; then
    echo " one of the dependencies is licensed as LGPL or GPL"
    echo " which is prohibited at Amazon"
    echo " please look into removing the dependency"
    exit 1
  fi

  # go-license is pretty eager to copy src for certain license types
  # when it does, it applies strange permissions to the copied files
  # which makes deleting them later awkward
  # this behavior may change in the future with the following PR
  # https://github.com/google/go-licenses/pull/28
  # We can delete these additional files because we are running go mod vendor
  # prior to this call so we know the source is the same as upstream
  # go-licenses is copying this code because it doesnt know if its be modified or not
  chmod -R 777 "${outputdir}/LICENSES"
  find "${outputdir}/LICENSES" -type f \( -name '*.yml' -o -name '*.go' -o -name '*.mod' -o -name '*.sum' -o -name '*gitignore' \) -delete

  # most of the packages show up the go-license.csv file as the module name
  # from the go.mod file, storing that away since the source dirs usually get deleted
  MODULE_NAME=$(go mod edit -json | jq -r '.Module.Path')
  if [ ! -f ${outputdir}/attribution/root-module.txt ]; then
  	echo $MODULE_NAME > ${outputdir}/attribution/root-module.txt
  fi
}

function build::non-golang::gather_licenses(){
  local -r project="$1"
  local -r git_tag="$2"
  local -r output_dir="$3"
  project_org="$(cut -d '/' -f1 <<< ${project})"
  project_name="$(cut -d '/' -f2 <<< ${project})"
  git clone https://github.com/${project_org}/${project_name}
  cd $project_name
  git checkout $git_tag
  cd ..
  build::non-golang::copy_licenses $project_name $output_dir/LICENSES/github.com/${project_org}/${project_name}
  rm -rf $project_name
}

function build::non-golang::copy_licenses(){
  local -r source_dir="$1"
  local -r destination_dir="$2"
  (cd $source_dir; find . \( -name "*COPYING*" -o -name "*COPYRIGHT*" -o -name "*LICEN[C|S]E*" -o -name "*NOTICE*" \)) |
  while read file
  do
    license_dest=$destination_dir/$(dirname $file)
    mkdir -p $license_dest
    cp -r "${source_dir}/${file}" $license_dest/$(basename $file)
  done
}

function build::generate_attribution(){
  if [ -z "${USE_HOST_GENERATE_ATTRIBUTION:-}" ]; then    
    build::common::override_missing_tooling generate-attribution
  fi

  local -r project_root=$1
  local -r golang_version=$2
  local -r output_directory=${3:-"${project_root}/_output"}
  local -r attribution_file=${4:-"${project_root}/ATTRIBUTION.txt"}

  local -r root_module_name=$(cat ${output_directory}/attribution/root-module.txt)

  local -r golang_version_tag=$(build::common::use_go_version $golang_version > /dev/null 2>&1 && go version | grep -o "go[0-9].* ")

  if cat "${output_directory}/attribution/go-license.csv" | grep -e ",LGPL-" -e ",GPL-"; then
    echo " one of the dependencies is licensed as LGPL or GPL"
    echo " which is prohibited at Amazon"
    echo " please look into removing the dependency"
    exit 1
  fi

  build::common::echo_and_run generate-attribution $root_module_name $project_root $golang_version_tag $output_directory 
  cp -f "${output_directory}/attribution/ATTRIBUTION.txt" $attribution_file
}

function build::common::get_go_path() {
  local -r version=$1

  # This is the path where the specific go binary versions reside in our builder-base image
  local -r gorootbinarypath="/go/go${version}/bin"
  # This is the path that will most likely be correct if running locally
  local -r gopathbinarypath="${GOPATH:-}/go${version}/bin"
  if [ -d "$gorootbinarypath" ]; then
    echo $gorootbinarypath
  elif [ -d "$gopathbinarypath" ]; then
    echo $gopathbinarypath
  elif command -v go &> /dev/null; then
    # not in builder-base, probably running in dev environment
    # return default go installation
    local -r which_go=$(which go)
    echo "$(dirname $which_go)"
  else
    echo ""
  fi
}

function build::common::use_go_version() {
  local -r version=$1

  if [ -z "${USE_HOST_GO:-}" ]; then    
    build::common::override_missing_tooling go $version
    return
  fi

  local -r gobinarypath=$(build::common::get_go_path $version)
  if [ -z "${gobinarypath}" ]; then
    echo "golang not available on host!"
    exit 1
  fi

  export PATH=${gobinarypath}:$PATH
  # the GOCACHE needs to be seperated, not preserved, by golang version otherwise it can leak
  # into future builds effecting checksums and builds in general
  export GOCACHE=$(go env GOCACHE)/$version
}

function build::common::find_project_root_from_pwd() {
   # find project root which may not always be in the same relative location to pwd
  local parts=()
  local path=$(pwd)
  while [[ "$(basename $path)" != "projects" ]] && [[ "$path" != "/" ]]; do
    parts+=($(basename $path))
    path=$(dirname $path)
  done

  # this shouldnt really happen in the context of our builds, but there are some cases
  # where we create tmp directories that will not follow this pattern
  if [[ "$path" = "/" ]]; then
    echo $(pwd)
    return
  fi

  local -r repo_owner=${parts[-1]}
  local -r repo=${parts[-2]}

  echo "${path}/${repo_owner}/${repo}"
}

function build::common::re_quote() {
    local -r to_escape=$1
    sed 's/[][()\.^$\/?*+]/\\&/g' <<< "$to_escape"
}

function build::common::get_latest_eksa_asset_url() {
  local -r artifact_bucket=$1
  local -r project=$2
  local -r arch=${3-amd64}
  local -r s3downloadpath=${4-latest}
  local -r gitcommitoverride=${5-false}

  s3artifactfolder=$s3downloadpath
  git_tag=$(cat $BUILD_ROOT/../../projects/${project}/GIT_TAG)
  if [ "$gitcommitoverride" = "true" ]; then
    commit_hash=$(echo $s3downloadpath | cut -d- -f2)
    git_tag=$(git show $commit_hash:projects/${project}/GIT_TAG)
    s3artifactfolder=$s3downloadpath/artifacts
  fi

  local -r tar_file_prefix=$(make --no-print-directory -C $BUILD_ROOT/../../projects/${project} var-value-TAR_FILE_PREFIX)
 
  local -r url="https://$(basename $artifact_bucket).s3-us-west-2.amazonaws.com/projects/$project/$s3artifactfolder/$tar_file_prefix-linux-$arch-${git_tag}.tar.gz"

  local -r http_code=$(curl -I -L -s -o /dev/null -w "%{http_code}" $url)
  if [[ "$http_code" == "200" ]]; then 
    echo "$url"
  else
    echo "https://$(basename $artifact_bucket).s3-us-west-2.amazonaws.com/projects/$project/latest/$tar_file_prefix-linux-$arch-${git_tag}.tar.gz"
  fi
}

function build::common::wait_for_tag() {
  local -r tag=$1
  sleep_interval=20
  for i in {1..60}; do
    echo "Checking for tag ${tag}..."
    git rev-parse --verify --quiet "${tag}" && echo "Tag ${tag} exists!" && break
    git fetch --tags > /dev/null 2>&1
    echo "Tag ${tag} does not exist!"
    echo "Waiting for tag ${tag}..."
    sleep $sleep_interval
    if [ "$i" = "60" ]; then
      exit 1
    fi
  done
}

function build::common::wait_for_tarball() {
  local -r tarball_url=$1
  sleep_interval=20
  for i in {1..60}; do
    echo "Checking for URL ${tarball_url}..."
    local -r http_code=$(curl -I -L -s -o /dev/null -w "%{http_code}" $tarball_url)
    if [[ "$http_code" == "200" ]]; then 
      echo "Tarball exists!" && break
    fi
    echo "Tarball does not exist!"
    echo "Waiting for tarball to be uploaded to ${tarball_url}"
    sleep $sleep_interval
    if [ "$i" = "60" ]; then
      exit 1
    fi
  done
}

function build::common::get_clone_url() {
  local -r org=$1
  local -r repo=$2
  local -r aws_region=$3
  local -r codebuild_ci=$4

  if [ "$codebuild_ci" = "true" ]; then
    echo "https://git-codecommit.${aws_region}.amazonaws.com/v1/repos/${org}.${repo}"
  else
    echo "https://github.com/${org}/${repo}.git"
  fi
}

function retry() {
  local n=1
  local max=120
  local delay=5
  while true; do
    "$@" && break || {
      if [[ $n -lt $max ]]; then
        ((n++))
        echo "Command failed. Attempt $n/$max:"
        sleep $delay;
      else
        fail "The command has failed after $n attempts."
      fi
    }
  done
}

function build::docker::retry_pull() {
  retry docker pull "$@"
}

function build::bottlerocket::check_release_availablilty() {
  local release_file=$1
  local release_channel=$2
  local format=$3
  retval=0
  release_version=$(yq e ".${release_channel}.${format}-release-version" $release_file)
  if [ $release_version == "null" ]; then
    retval=1
  fi
  echo $retval
}

function build::common::override_missing_tooling() {
  local -r cmd="$1"
  local -r version="${2:-}"

  # TEMP
  if [ "${USE_EXP_BUILDCTL_IN_PRESUBMIT}" = "true" ] && [ "${JOB_TYPE:-}" == "presubmit" ]; then
    rm -rf /root/sdk /go /usr/bin/generate-attribution
  fi
  
  local -r project_root="$(build::common::find_project_root_from_pwd)"
  local -r overrides_root="${project_root}/_output/.path-overrides"

  if [ "${cmd}" = "go" ] && [ ! -f $overrides_root/$cmd ]; then
    echo "************************************************************************************************************************"
    echo "By default this repo uses golang provided by the EKS-Distro golang containers instead of installed versions on the host."
    echo "If docker is available ``docker run`` will be used, otherwise ``buildctl build`` will be used."
    echo "To override this behavior, ``export USE_HOST_GO=true``."
    echo "************************************************************************************************************************"  
  fi

  mkdir -p $overrides_root
  ln -sf $BUILD_ROOT/overrides/$cmd $overrides_root
  ln -sf $BUILD_ROOT/overrides/run-base $overrides_root

  if [ -n "$version" ]; then
    echo "$version" > $overrides_root/.${cmd}version
  fi

  if [[ "$USE_BUILDCTL" = "false" ]] && command -v docker &> /dev/null && docker info > /dev/null 2>&1 ; then
    echo "Using container image for $cmd $version via docker"
    echo "true" > $overrides_root/.usedocker
  else
    echo "Using container image for $cmd $version via buildctl"
    echo "false" > $overrides_root/.usedocker
  fi

  export PATH=${overrides_root}:$PATH
}

function build::common::echo_and_run() {
  echo "($(pwd)) \$ $*"
  "$@"
}
