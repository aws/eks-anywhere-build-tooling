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

# Short-circuit if script has already been sourced
[[ $(type -t build::common::loaded) == function ]] && return 0

BUILD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/" && pwd -P)"
source "${BUILD_ROOT}/eksd_releases.sh"

if [ -n "${OUTPUT_DEBUG_LOG:-}" ]; then
    set -x
fi

function build::find::gnu_variant_on_mac() {
  local -r cmd="$1"

  if [ "$(uname -s)" = "Linux" ]; then
    echo "$cmd"
    return
  fi

  local final="$cmd"
  if command -v "g$final" &> /dev/null; then
    final="g$final"
  fi

  if [[ "$final" = "$cmd" ]] && command -v "gnu$final" &> /dev/null; then
    final="gnu$final"
  fi

  if [[ "$final" = "$cmd" ]]; then
    >&2 echo " !!! Building on Mac OS X and GNU '$cmd' not found. Using the builtin version"
    >&2 echo "     *may* work, but in general you should either build on a Linux host or"
    >&2 echo "     install the gnu version via brew, usually 'brew install gnu-$cmd'"    
  fi

  echo "$final"
}

function build::common::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=$(build::find::gnu_variant_on_mac tar)
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

  build::common::echo_and_run "${TAR}" czf "${tarfile}" -C "${stagingdir}" $repository --owner=0 --group=0
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
  local -r artifactspath="${1%/}" # remove trailing slash if it exists
  local -r artifactsbucket="${2%/}" # remove trailing slash if it exists
  local -r projectpath=$3
  local -r buildidentifier=$4
  local -r githash=$5
  local -r latesttag=$6
  local -r dry_run=$7
  local -r do_not_delete=$8
  local -r create_public_acl=$9

  if [[ $artifactsbucket = /* ]]; then
    if [ "$dry_run" = "true" ]; then
      build::common::echo_and_run rsync -i "$artifactspath/" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts --recursive --dry-run
      build::common::echo_and_run rsync -i "$artifactspath/" "$artifactsbucket"/"$projectpath"/"$latesttag" --recursive --dry-run
    else
      build::common::echo_and_run rsync -i --mkpath "$artifactspath/" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts --recursive 

      if [ "$do_not_delete" = "true" ]; then
        build::common::echo_and_run rsync -i --mkpath "$artifactspath/" "$artifactsbucket"/"$projectpath"/"$latesttag" --recursive 
      else
        build::common::echo_and_run rsync -i --mkpath "$artifactspath/" "$artifactsbucket"/"$projectpath"/"$latesttag" --delete --recursive 
      fi
    fi
  else
    local public_acl=""
    if [ "$create_public_acl" = "true" ]; then
      public_acl="--acl public-read"
    fi
    
    if [ "$dry_run" = "true" ]; then
      build::common::echo_and_run aws s3 cp "$artifactspath" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts --recursive --dryrun
      build::common::echo_and_run aws s3 cp "$artifactspath" "$artifactsbucket"/"$projectpath"/"$latesttag" --recursive --dryrun
    else
      # Upload artifacts to s3 
      # 1. To proper path on s3 with buildId-githash
      # 2. Latest path to indicate the latest build, with --delete option to delete stale files in the dest path
      build::common::echo_and_run aws s3 sync "$artifactspath" "$artifactsbucket"/"$projectpath"/"$buildidentifier"-"$githash"/artifacts $public_acl --no-progress

      if [ "$do_not_delete" = "true" ]; then
        build::common::echo_and_run aws s3 sync "$artifactspath" "$artifactsbucket"/"$projectpath"/"$latesttag" $public_acl --no-progress
      else
        build::common::echo_and_run aws s3 sync "$artifactspath" "$artifactsbucket"/"$projectpath"/"$latesttag" --delete $public_acl --no-progress
      fi
    fi
  fi
}

function build::gather_licenses() {
  local -r outputdir="${1%/}" # remove trailing slash if it exists
  local -r patterns=$2
  local -r golang_version=$3
  local -r threshold=$4
  local -r cgo_enabled=$5

  # Force deps to only be pulled form vendor directories
  # this is important in a couple cases where license files
  # have to be manually created
  export GOFLAGS=-mod=vendor
  # force platform to be linux to ensure all deps are picked up
  export GOOS=linux 
  export GOARCH=amd64 
  export CGO_ENABLED=$cgo_enabled
  # Setting this variable to local so that it always uses the local
  # bundled toolchain regardless of the version specified in go.mod file
  # It was introduced in this commit from Go1.21 onwards:
  # https://github.com/golang/go/commit/83dfe5cf62234427eae04131dc6e4551fd283463
  # Upstream issue:- https://github.com/google/go-licenses/issues/244
  export GOTOOLCHAIN=local

  build::common::use_go_version "$golang_version"

  if ! command -v go-licenses &> /dev/null
  then
    echo " go-licenses not found.  If you need license or attribution file handling"
    echo " please refer to the doc in docs/development/attribution-files.md"
    exit
  fi

  mkdir -p "${outputdir}/attribution"
  # attribution file generated uses the output go-deps and go-license to gather the necessary
  # data about each dependency to generate the amazon approved attribution.txt files
  # go-deps is needed for module versions
  # go-licenses are all the dependencies found from the module(s) that were passed in via patterns
  build::common::echo_and_run go list -deps=true -json ./... | jq -s '.'  > "${outputdir}/attribution/go-deps.json"

  # go-licenses can be a bit noisy with its output and lot of it can be confusing 
  # the following messages are safe to ignore since we do not need the license url for our process
  NOISY_MESSAGES="cannot determine URL for|Error discovering license URL|unsupported package host|contains non-Go code|has empty version|\.(h|s|c)$"
 
  build::common::echo_and_run go-licenses save --confidence_threshold $threshold --force $patterns --save_path "${outputdir}/LICENSES" 2> >(grep -vE "$NOISY_MESSAGES")
  
  build::common::echo_and_run go-licenses csv --confidence_threshold $threshold $patterns 2> >(grep -vE "$NOISY_MESSAGES") > "${outputdir}/attribution/go-license.csv"  

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
  (cd $source_dir; find . -maxdepth 1 \( -name "*COPYING*" -o -name "*COPYRIGHT*" -o -name "*LICEN[C|S]E*" -o -name "*NOTICE*" \)) |
  while read file
  do
    license_dest=$destination_dir/$(dirname $file)
    mkdir -p $license_dest
    cp -r "${source_dir}/${file}" $license_dest/$(basename $file)
  done
}

function build::generate_attribution(){
  local -r project_root=$1
  local -r golang_version=$2
  local -r output_directory=${3:-"${project_root}/_output"}
  local -r attribution_file=${4:-"${project_root}/ATTRIBUTION.txt"}

  local -r root_module_name=$(cat ${output_directory}/attribution/root-module.txt)
  local -r go_path=$(build::common::get_go_path $golang_version)
  local -r golang_version_tag=$($go_path/go version | grep -o "go[0-9].* ")

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
  local -r gopathbinarypath="$GOPATH/go${version}/bin"
  if [ -d "$gorootbinarypath" ]; then
    echo $gorootbinarypath
  elif [ -d "$gopathbinarypath" ]; then
    echo $gopathbinarypath
  else
    # not in builder-base, probably running in dev environment
    # return default go installation
    local -r which_go=$(which go)
    echo "$(dirname $which_go)"
  fi
}

function build::common::use_go_version() {
  local -r version=${1:-}

  if [ -z "$version" ]; then
    return
  fi

  if (( "${version#*.}" < 16 )); then
    echo "Building with GO version $version is no longer supported!  Please update the build to use a newer version."
    exit 1
  fi

  local -r gobinarypath=$(build::common::get_go_path $version)
  echo "Adding $gobinarypath to PATH"
  # Adding to the beginning of PATH to allow for builds on specific version if it exists
  export PATH=${gobinarypath}:$PATH
  export GOCACHE=$(go env GOCACHE)/$version
  echo "$(go version)"
}

# Use a seperate build cache for each project/version to ensure there are no
# shared bits which can mess up the final checksum calculation
# this is mostly needed for create checksums locally since in the builds
# different versions of the same project are not built in the same container
function build::common::set_go_cache() {
  local -r project=$1
  local -r git_tag=$2
  export GOCACHE=$(go env GOCACHE)/$project/$git_tag
}

function build::common::re_quote() {
    local -r to_escape=$1
    sed 's/[][()\.^$\/?*+]/\\&/g' <<< "$to_escape"
}
function build::common::check_eksa_asset_url() {
  local -r s3_url_prefix="$1"
  local -r specific_uri="$2"
  local -r fallback_latest_uri="$3"
  local -r bucket_name="$4"

  if [[ "$(build::common::echo_and_run curl -I -L -s -o /dev/null -w "%{http_code}" $s3_url_prefix/$specific_uri)" == "200" ]]; then 
    echo "$s3_url_prefix/$specific_uri"
  elif [[ "$(build::common::echo_and_run curl -I -L -s -o /dev/null -w "%{http_code}" $s3_url_prefix/$fallback_latest_uri)" == "200" ]]; then 
    echo "$s3_url_prefix/$fallback_latest_uri"
  elif build::common::echo_and_run aws s3api head-object --bucket $bucket_name --key $specific_uri &> /dev/null; then
    build::common::echo_and_run aws s3 presign $bucket_name/$specific_uri
  elif build::common::echo_and_run aws s3api head-object --bucket $bucket_name --key $fallback_latest_uri &> /dev/null; then
    build::common::echo_and_run aws s3 presign $bucket_name/$fallback_latest_uri
  fi
}

function build::common::get_latest_eksa_asset_url() {
  local -r artifact_bucket="${1%/}" # remove trailing slash if it exists
  local -r project=$2
  local -r arch=${3-amd64}
  local -r s3downloadpath=${4-latest}
  local -r releasebranch=${5-}
  local -r sha=${6-false}

  s3artifactfolder=$s3downloadpath

  projectwithreleasebranch=$project
  # If not able to find git_tag w/o specifying a branch, update projectwithreleasebranch and git_tag to use a branch.
  git_tag=$(cat $BUILD_ROOT/../../projects/${projectwithreleasebranch}/GIT_TAG || echo "invalid")
  if [ "$git_tag" = "invalid" ]; then
    projectwithreleasebranch=$project/$releasebranch
    git_tag=$(cat $BUILD_ROOT/../../projects/${projectwithreleasebranch}/GIT_TAG)
  fi  

  local -r tar_file_prefix=$(MAKEOVERRIDES= MAKEFLAGS= make --no-print-directory -C $BUILD_ROOT/../../projects/${project} var-value-TAR_FILE_PREFIX)
 
  local git_tag_normalized="${git_tag//\//-}"
  local specific_uri="projects/$projectwithreleasebranch/$s3artifactfolder/$tar_file_prefix-linux-$arch-${git_tag_normalized}.tar.gz"
  local fallback_latest_uri="projects/$projectwithreleasebranch/latest/$tar_file_prefix-linux-$arch-${git_tag_normalized}.tar.gz"

  if [ "$sha" = "true" ]; then
    specific_uri+=".sha256"
    fallback_latest_uri+=".sha256"
  fi

  if [[ $artifact_bucket = /* ]]; then
    if [ -f $artifact_bucket/$specific_uri ]; then
      echo "file://$artifact_bucket/$specific_uri"
    elif [ -f $artifact_bucket/$fallback_latest_uri ]; then
      echo "file://$artifact_bucket/$fallback_latest_uri"
    fi
    return
  fi

  local -r bucket_name_without_prefix=${artifact_bucket#s3://} 
  local -r bucket_name=${bucket_name_without_prefix%%/*}
  local -r url_path=${bucket_name_without_prefix#*/}
  if [ "${url_path}" != "${bucket_name}" ]; then 
    specific_uri="${url_path%/}/${specific_uri}"
    fallback_latest_uri="${url_path%/}/${fallback_latest_uri}"
  fi

  local -r s3_url_prefix="https://$bucket_name.s3-us-west-2.amazonaws.com"
  
  local -r sleep_interval=20
  for i in {1..60}; do
    local final_url=$(build::common::check_eksa_asset_url "$s3_url_prefix" "$specific_uri" "$fallback_latest_uri" "$bucket_name")
    if [ -n "$final_url" ]; then
      echo "$final_url"
      break
    elif [ "${CODEBUILD_CI:-false}" = "false" ] || [ "$i" = "60" ]; then
      >&2 echo "******* No artifact availabe! *******"
      >&2 echo "${s3_url_prefix}/${fallback_latest_uri} does not exists!"
      >&2 echo "Please double check the value of \$ARTIFACTS_BUCKET."
      >&2 echo "${git_tag} of ${project} may not be the current latest version, verify you have the latest code from main to be sure."
      >&2 echo "*************************************"
      exit 1
    fi
    >&2 echo "Tarball does not exist!"
    >&2 echo "Waiting for tarball to be uploaded to ${specific_uri}"
    sleep $sleep_interval
  done
}

function build::common::get_latest_eksa_asset_url_sha256() {
  build::common::get_latest_eksa_asset_url $@ true
}

function build::common::wait_for_tag() {
  local -r tag=$1
  sleep_interval=20
  for i in {1..60}; do
    echo "Checking for tag/branch ${tag}..."
    
    # First try to find it as a tag
    if git rev-parse --verify --quiet "${tag}" > /dev/null 2>&1; then
      echo "Tag ${tag} exists!"
      break
    fi
    
    # If not found as a tag, try as a remote branch
    if git rev-parse --verify --quiet "origin/${tag}" > /dev/null 2>&1; then
      echo "Branch ${tag} exists!"
      break
    fi
    
    # Fetch both tags and branches
    git fetch --tags > /dev/null 2>&1
    git fetch origin > /dev/null 2>&1
    
    echo "Tag/branch ${tag} does not exist!"
    echo "Waiting for tag/branch ${tag}..."
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


function fail() {
  echo $1 >&2
  exit 1
}

function retry() {
  local n=1
  local max=120
  local delay=5
  while true; do
    "$@" && break || {
      if [[ $n -lt $max ]]; then
        ((n++))
        >&2 echo "Command failed. Attempt $n/$max:"
        sleep $delay;
      else
        fail "The command has failed after $n attempts."
      fi
    }
  done
}

# $1 - timeout value, should include unit (s/m/h/etc) ex: 10m
function retry_with_timeout() {
  TIMEOUT=$1
  shift

  local n=1
  local max=120
  local delay=5
  while true; do
    timeout $TIMEOUT "$@" && break || {
      if [[ $n -lt $max ]]; then
        ((n++))
        # multiple the numeric part of the timeout by 1.5 and suffix with the last char which is the unit
        TIMEOUT=$((${TIMEOUT:0:-1} * 3/2))${TIMEOUT: -1}
        echo "Command failed. Attempt $n/$max with timeout ${TIMEOUT}:"
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

function build::common::echo_and_run() {
  >&2 echo "($(pwd)) \$ $*"
  "$@"
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

function build::jq::update_in_place() {
  local json_file=$1
  local jq_query=$2

  cat $json_file | jq -S ''"$jq_query"'' > $json_file.tmp && mv $json_file.tmp $json_file
}

function build::common::copy_if_source_destination_different() {
  local source=$1
  local destination=$2

  STAT=$(build::find::gnu_variant_on_mac stat)
  source_inode=$($STAT -c %i $source)
  destination_inode=""
  if [ -d $destination ] && [ -e $destination/$(basename $source) ]; then
    destination_inode=$($STAT -c %i $destination/$(basename $source))
  elif [ -f $destination ] && [ -e $destination ]; then
    destination_inode=$($STAT -c %i $destination)
  fi

  if [ -n "$destination_inode" ] && [ "$source_inode" = "$destination_inode" ]; then
    echo "Source and destination are the same file"
    return
  fi

  cp -rf $source $destination
}

function build::common::is_qemu_available() {
  local -r platform="$1"

  if [ "${CODEBUILD_CI:-false}" = "true" ]; then
    # code build is running in a container so we cant rely on checking for the proc file
    # since it would be on the host directly
    return 0
  fi

  local -r normalized_platform="$(echo "${platform}" | sed 's/linux\/arm64/aarch64/g;s/linux\/amd64/x86_64/g')"

  if [ "$(uname -s)" = "Linux" ] && [ "$(uname -m)" != "$normalized_platform" ]; then
    local -r qemu_file="qemu-${normalized_platform}"
    if [ ! -f "/proc/sys/fs/binfmt_misc/$qemu_file" ] || ! grep -q "enabled" "/proc/sys/fs/binfmt_misc/$qemu_file"; then
      return 1
    fi
  fi

  return 0
}

function build::common::check_for_qemu() {
  local -r platform="$1"
  
  if build::common::is_qemu_available "$platform"; then
    return
  fi

  echo "****************************************************************"
  echo "You are trying to run, or build, a $platform based container which does not match your host architecture."
  echo "Run the following to register qemu virtualization:"
  echo "docker run --privileged --rm public.ecr.aws/eks-distro-build-tooling/binfmt-misc:qemu-v6.1.0 --install aarch64,amd64"
  echo "****************************************************************"
  exit 1
}

# Marker function to indicate script has been fully sourced
function build::common::loaded() {
  return 0
}
