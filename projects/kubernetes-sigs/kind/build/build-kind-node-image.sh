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
INTERMEDIATE_BASE_IMAGE="${2?Specify second argument - kind base tag}"
INTERMEDIATE_NODE_IMAGE="${3?Specify third argument - kind node image name}"
ARTIFACTS_BUCKET="${4?Specify fourth argument - artifact bucket}"
ARCH="${5?Specify fifth argument - Targetarch}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

# Include common constants and other vars needed when building image
source "${MAKE_ROOT}/_output/$EKSD_RELEASE_BRANCH/kind-node-image-build-args"

# Uses kind cli to build node-image, usually the fake kubernetes src
# to download the eks-d release artifacts instead of building from src
# ouput image: $INTERMEDIATE_NODE_IMAGE
function build::kind::build_node_image(){
    export KUBE_VERSION=$KUBE_VERSION
    export EKSD_RELEASE_BRANCH=$EKSD_RELEASE_BRANCH
    export EKSD_RELEASE=$EKSD_RELEASE
    export EKSD_IMAGE_REPO=$EKSD_IMAGE_REPO
    export EKSD_ASSET_URL=$EKSD_ASSET_URL
    export KUBE_ARCH=$ARCH

    KIND_PATH="$MAKE_ROOT/_output/bin/kind/$(uname | tr '[:upper:]' '[:lower:]')-$(go env GOHOSTARCH)/kind"
    $KIND_PATH build node-image $MAKE_ROOT/images/k8s.io/kubernetes \
        --base-image $INTERMEDIATE_BASE_IMAGE --image $INTERMEDIATE_NODE_IMAGE --arch $ARCH      
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
# output image: $INTERMEDIATE_NODE_IMAGE
function build::kind::load_images(){
    CONTAINER_ID=$(docker run --platform linux/$ARCH --entrypoint sleep -d -i $INTERMEDIATE_NODE_IMAGE infinity)

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
            docker pull --platform linux/$ARCH $image
        fi
        image_id=$(docker images $image --format "{{.ID}}")
        if [ ! -f $MAKE_ROOT/_output/dependencies/$image_id.tar ]; then
            if [[ $image =~ local-path-provisioner ]] || [[ $image =~ kindnetd ]]; then
                docker tag $image_id ${release_image_overrides[$image]}
                docker save ${release_image_overrides[$image]} -o $MAKE_ROOT/_output/dependencies/$image_id.tar
            else
                docker save $image -o $MAKE_ROOT/_output/dependencies/$image_id.tar
            fi
        fi
        docker exec --privileged -i $CONTAINER_ID \
            ctr --namespace=k8s.io images import --all-platforms --no-unpack - < $MAKE_ROOT/_output/dependencies/$image_id.tar
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
    docker tag $NEW_IMAGE_ID $INTERMEDIATE_NODE_IMAGE
    docker push $INTERMEDIATE_NODE_IMAGE

    docker kill $CONTAINER_ID
    rm -rf $MAKE_ROOT/_output/dependencies
}


function build::kind::download_additional_components() {
    OUTPUT_FOLDER="$MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies/linux-$ARCH"

    declare -A URLS=([$(build::eksd_releases::get_eksd_kubernetes_asset_url kubernetes-client-linux-$ARCH.tar.gz $EKSD_RELEASE_BRANCH $ARCH)]="kubernetes"
                  [$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/etcdadm' $ARCH)]="etcdadm"
                  [$(build::eksd_releases::get_eksd_component_url "cni-plugins" $EKSD_RELEASE_BRANCH $ARCH)]="cni-plugins"
                  [$(build::common::get_latest_eksa_asset_url $ARTIFACTS_BUCKET 'kubernetes-sigs/cri-tools' $ARCH)]="cri-tools")

    
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
    ETCD_HTTP_SOURCE=$(build::eksd_releases::get_eksd_component_url "etcd" $EKSD_RELEASE_BRANCH $ARCH)
    ETCD_VERSION=$(build::eksd_releases::get_eksd_component_version "etcd" $EKSD_RELEASE_BRANCH $ARCH)
    FOLDER="$OUTPUT_FOLDER/cache/etcdadm/etcd/$ETCD_VERSION"
    mkdir -p $FOLDER
    curl -sSL $ETCD_HTTP_SOURCE -o $FOLDER/etcd-$ETCD_VERSION-linux-$ARCH.tar.gz
}

if command -v docker &> /dev/null && docker info > /dev/null 2>&1 ; then
    build::kind::build_node_image
    build::kind::load_images
    build::kind::download_additional_components
fi
