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

EKSD_RELEASE_BRANCH="${1?Specify first argument - release branch}"
KIND_BASE_IMAGE_NAME="${2?Specify second argument - kind base tag}"
KIND_NODE_IMAGE_NAME="${3?Specify third argument - kind node image name}"
KIND_KINDNETD_IMAGE_OVERRIDE="${4?Specify the fourth argument - kindnetd image}"
IMAGE_REPO="${5?Specify fifth argument - image repo}"
IMAGE_TAG="${6?Specify sixth argument - image tag}"
ARTIFACTS_BUCKET="${7?Specify seventh argument - artifact bucket}"
PUSH="${8?Specify eighth argument - push}"
LATEST_TAG="${9?Specify ninth argument - Tag denoting build source}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

# This is used by the local-path-provisioner within the kind node
AL2_HELPER_IMAGE="public.ecr.aws/amazonlinux/amazonlinux:2"
LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE="$IMAGE_REPO/rancher/local-path-provisioner:latest"
LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE="public.ecr.aws/eks-anywhere/rancher/local-path-provisioner:$(cat $MAKE_ROOT/../../rancher/local-path-provisioner/GIT_TAG)"
KIND_KINDNETD_RELEASE_OVERRIDE="public.ecr.aws/eks-anywhere/kubernetes-sigs/kind/kindnetd:$(cat $MAKE_ROOT/GIT_TAG)"

# Preload release yaml
build::eksd_releases::load_release_yaml $EKSD_RELEASE_BRANCH

KUBE_VERSION=$(build::eksd_releases::get_eksd_component_version "kubernetes" $EKSD_RELEASE_BRANCH)
EKSD_RELEASE=$(build::eksd_releases::get_eksd_release_number $EKSD_RELEASE_BRANCH)
EKSD_KUBE_VERSION="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH-$EKSD_RELEASE"
PAUSE_IMAGE_TAG_OVERRIDE=$(build::eksd_releases::get_eksd_kubernetes_image_url "pause-image" $EKSD_RELEASE_BRANCH)
EKSD_IMAGE_REPO=$(build::eksd_releases::get_eksd_image_repo $EKSD_RELEASE_BRANCH)
EKSD_ASSET_URL=$(build::eksd_releases::get_eksd_kubernetes_asset_base_url $EKSD_RELEASE_BRANCH)/$KUBE_VERSION

# Expected versions provided by kind which are replaced in the docker build with our versions
# when updating kind check the following, they may need to be updated
# https://github.com/kubernetes-sigs/kind/blob/main/pkg/build/nodeimage/const_cni.go#L23
KINDNETD_IMAGE_TAG="docker.io/kindest/kindnetd:v20210326-1e038dc5"
# https://github.com/kubernetes-sigs/kind/blob/main/pkg/build/nodeimage/const_storage.go#L28
DEBIAN_BASE_IMAGE_TAG="k8s.gcr.io/build-image/debian-base:v2.1.0"
# https://github.com/kubernetes-sigs/kind/blob/main/pkg/build/nodeimage/const_storage.go#L28
LOCAL_PATH_PROVISONER_IMAGE_TAG="docker.io/rancher/local-path-provisioner:v0.0.14"
# https://github.com/kubernetes-sigs/kind/blob/main/images/base/files/etc/containerd/config.toml#L22
PAUSE_IMAGE_TAG="k8s.gcr.io/pause:3.5"

# Uses kind cli to build node-image, usually the fake kubernetes src
# to download the eks-d release artifacts instead of building from src
# ouput image: $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-kind
function build::kind::build_node_image(){
    export KUBE_VERSION=$KUBE_VERSION
    export EKSD_RELEASE_BRANCH=$EKSD_RELEASE_BRANCH
    export EKSD_RELEASE=$EKSD_RELEASE
    export EKSD_IMAGE_REPO=$EKSD_IMAGE_REPO
    export EKSD_ASSET_URL=$EKSD_ASSET_URL

    # base image was created using buildctl and stored as tar
    BASE_IMAGE_ID=$(docker load -i $MAKE_ROOT/_output/images/$EKSD_RELEASE_BRANCH/base.tar | sed -E 's/.*sha256:(.*)$/\1/')
    docker tag $BASE_IMAGE_ID $KIND_BASE_IMAGE_NAME:$EKSD_KUBE_VERSION

    KIND_PATH="$MAKE_ROOT/_output/bin/kind/$(uname | tr '[:upper:]' '[:lower:]')-amd64/kind"
    $KIND_PATH build node-image \
        --base-image $KIND_BASE_IMAGE_NAME:$EKSD_KUBE_VERSION  --image $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-kind \
        --kube-root $MAKE_ROOT/images/k8s.io/kubernetes
}

function build::kind::validate_versions(){
    local -r container_id=$1
    local -r eksd_tag="$KUBE_VERSION-eks-$EKSD_RELEASE_BRANCH"
    # We expect certain versions to be in the kind base/node images since we end up replacing them
    # validate they haven't changed
    if ! docker exec -i $container_id grep "image: $KINDNETD_IMAGE_TAG"  /kind/manifests/default-cni.yaml > /dev/null 2>&1; then
        echo "Did not find expected version of kindnetd: $KINDNETD_IMAGE_TAG"
        exit 1
    fi
    if ! docker exec -i $container_id grep "$DEBIAN_BASE_IMAGE_TAG"  /kind/manifests/default-storage.yaml > /dev/null 2>&1; then
        echo "Did not find expected version of debian base: $DEBIAN_BASE_IMAGE_TAG"
        exit 1
    fi
    if ! docker exec -i $container_id grep "image: $LOCAL_PATH_PROVISONER_IMAGE_TAG"  /kind/manifests/default-storage.yaml > /dev/null 2>&1; then
        echo "Did not find expected version of local path provisoner: $LOCAL_PATH_PROVISONER_IMAGE_TAG"
        exit 1
    fi
    if ! docker exec -i $container_id grep "$PAUSE_IMAGE_TAG"  /etc/containerd/config.toml > /dev/null 2>&1; then
        echo "Did not find expected version of pause image: $PAUSE_IMAGE_TAG"
        exit 1
    fi
    if ! docker exec -i $container_id /usr/bin/kubeadm version -oshort | grep "$eksd_tag" > /dev/null 2>&1; then
        echo "Did not find expected version of kubeadm: $eksd_tag"
        exit 1
    fi
    if ! docker exec -i $container_id kubectl version --short --client=true | grep "$eksd_tag" > /dev/null 2>&1; then
        echo "Did not find expected version of kubectl: $eksd_tag"
        exit 1
    fi
    if ! docker exec -i $container_id kubelet --version | grep "$eksd_tag" > /dev/null 2>&1; then
        echo "Did not find expected version of kubelet: $eksd_tag"
        exit 1
    fi
}

# This whole process of creating a container from the image, loading images
# and commit it back is a recreation of what the kind build node-image does
# the reason being is this is the easiest way to convert a docker/oci image
# into an on disk containerd format
# output image: $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION
function build::kind::load_images(){
    CONTAINER_ID=$(docker run --entrypoint sleep -d -i $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-kind infinity)

    build::kind::validate_versions $CONTAINER_ID

    docker exec --privileged -i $CONTAINER_ID bash -c "nohup containerd > /dev/null 2>&1 & sleep 5"
    docker exec --privileged -i $CONTAINER_ID crictl images

    # remove unneeded images
    IMAGES=(
        # kind default cni        
        $KINDNETD_IMAGE_TAG 
        # kind adds this for debugging purposes + local-path-provisoner, replaced with al2
        $DEBIAN_BASE_IMAGE_TAG
        # replaced with eks-a build
        $LOCAL_PATH_PROVISONER_IMAGE_TAG
        # replaced with pause image from eks-d
        $PAUSE_IMAGE_TAG  
    )
    for image in "${IMAGES[@]}"; do
        # in case kind didnt include the image we are expected, ignore the error
        # there is a final validation to make sure all images are from public.ecr
        docker exec --privileged -i $CONTAINER_ID crictl rmi $image || true
    done

    # pull local-path-provisioner + al2 helper image
    mkdir -p $MAKE_ROOT/_output/dependencies
    IMAGES=("$AL2_HELPER_IMAGE" "$LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE" "$PAUSE_IMAGE_TAG_OVERRIDE" "$KIND_KINDNETD_IMAGE_OVERRIDE")

    declare -A release_image_overrides
    release_image_overrides["$LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE"]=$LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE
    release_image_overrides["$KIND_KINDNETD_IMAGE_OVERRIDE"]=$KIND_KINDNETD_RELEASE_OVERRIDE

    for image in "${IMAGES[@]}"; do
        if ! docker image inspect $image > /dev/null 2>&1; then
            docker pull $image
        fi
        if [[ $image =~ local-path-provisioner ]] || [[ $image =~ kindnetd ]]; then
            image_id=$(docker images $image --format "{{.ID}}")
            docker tag $image_id ${release_image_overrides[$image]}
            docker save ${release_image_overrides[$image]} -o $MAKE_ROOT/_output/dependencies/image.tar
        else
            docker save $image -o $MAKE_ROOT/_output/dependencies/image.tar
        fi
        docker exec --privileged -i $CONTAINER_ID \
            ctr --namespace=k8s.io images import --all-platforms --no-unpack - < $MAKE_ROOT/_output/dependencies/image.tar
        rm $MAKE_ROOT/_output/dependencies/image.tar
    done

    docker exec --privileged -i $CONTAINER_ID crictl images
    
    # Validate all images in the image are from public.ecr
    FINAL_IMAGES=$(docker exec -i $CONTAINER_ID crictl images -o json | jq ".images[].repoTags[]" -r) 
    mapfile -t FINAL_IMAGES <<< "$FINAL_IMAGES"
    declare -p FINAL_IMAGES
    for image in "${FINAL_IMAGES[@]}"; do
        if [[ ! $image =~ ^public\.ecr\.aws\/ ]]; then
            echo "$image is not from public.ecr.aws!"
            exit 1
        fi
    done

    docker exec --privileged -i $CONTAINER_ID pkill containerd

    NEW_IMAGE_ID=$(docker commit --change 'ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]' $CONTAINER_ID | sed -E 's/.*sha256:(.*)$/\1/')
    docker tag $NEW_IMAGE_ID $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION

    docker kill $CONTAINER_ID
    rm -rf $MAKE_ROOT/_output/dependencies
}


function build::kind::download_additional_components() {
    OUTPUT_FOLDER="$MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies"

    declare -A URLS=([$(build::eksd_releases::get_eksd_kubernetes_asset_url kubernetes-client-linux-amd64.tar.gz)]="kubernetes"
                  [$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/etcdadm')]="etcdadm"
                  [$(build::eksd_releases::get_eksd_component_url "cni-plugins" $EKSD_RELEASE_BRANCH)]="cni-plugins"
                  [$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/cri-tools')]="cri-tools")

    
    mkdir -p $OUTPUT_FOLDER/LICENSES
    for URL in "${!URLS[@]}"
    do
        FOLDER=(${URLS[$URL]})

        mkdir -p $OUTPUT_FOLDER/$FOLDER
        TARBALL="$OUTPUT_FOLDER/tmp.tar.gz"
        curl -sSL $URL -o ${TARBALL}
        tar xzf ${TARBALL} -C $OUTPUT_FOLDER/$FOLDER

        FOLDER_PATH=$FOLDER
        if [ $FOLDER == 'kubernetes' ]; then
            FOLDER_PATH="kubernetes/kubernetes"
        fi
        cp -rf $OUTPUT_FOLDER/$FOLDER_PATH/LICENSES "$OUTPUT_FOLDER/LICENSES/$(echo $FOLDER | tr a-z A-Z  | tr -d '-'  )_LICENSES"
        cp $OUTPUT_FOLDER/$FOLDER_PATH/ATTRIBUTION.txt "$OUTPUT_FOLDER/LICENSES/$(echo $FOLDER | tr a-z A-Z  | tr -d '-'  )_ATTRIBUTION.txt"
    done
    
    # Download etcd tarball to be placed in etcdadm cache directory to avoid downloading at runtime
    ETCD_HTTP_SOURCE=$(build::eksd_releases::get_eksd_component_url "etcd" $EKSD_RELEASE_BRANCH)
    ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH)
    FOLDER="$OUTPUT_FOLDER/cache/etcdadm/etcd/$ETCD_VERSION"
    mkdir -p $FOLDER
    curl -sSL $ETCD_HTTP_SOURCE -o $FOLDER/etcd-$ETCD_VERSION-linux-amd64.tar.gz
}

function build::kind::build_final_node_image(){
    # squash image since we removed images
    docker build \
        -f $MAKE_ROOT/images/node/Dockerfile.squash \
        -t $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-$IMAGE_TAG \
        --build-arg BASE_IMAGE=$KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION \
        --build-arg IMAGE_REPO=$IMAGE_REPO \
        --build-arg AL2_HELPER_IMAGE=$AL2_HELPER_IMAGE \
        --build-arg DEBIAN_BASE_IMAGE_TAG=$DEBIAN_BASE_IMAGE_TAG \
        --build-arg LOCAL_PATH_PROVISONER_IMAGE_TAG=$LOCAL_PATH_PROVISONER_IMAGE_TAG \
        --build-arg LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE=$LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE \
        --build-arg PAUSE_IMAGE_TAG_OVERRIDE=$PAUSE_IMAGE_TAG_OVERRIDE \
        --build-arg PAUSE_IMAGE_TAG=$PAUSE_IMAGE_TAG \
        --build-arg KIND_KINDNETD_IMAGE_OVERRIDE=$KIND_KINDNETD_RELEASE_OVERRIDE \
        --build-arg KINDNETD_IMAGE_TAG=$KINDNETD_IMAGE_TAG \
        $MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies

    if [ "$PUSH" == "true" ] ; then
        docker push $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-$IMAGE_TAG
        docker tag $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-$IMAGE_TAG $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-$LATEST_TAG
        docker push $KIND_NODE_IMAGE_NAME:$EKSD_KUBE_VERSION-$LATEST_TAG
    fi
}

if command -v docker &> /dev/null && docker info > /dev/null 2>&1 ; then
    build::kind::build_node_image
    build::kind::load_images
    build::kind::download_additional_components
    build::kind::build_final_node_image
fi
