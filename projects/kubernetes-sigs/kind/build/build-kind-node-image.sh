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
ARCH="${3?Specify third argument - Targetarch}"
BUILDER_PLATFORM_ARCH="${4?Specify fourth argument - Hostarch}"

MAKE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${MAKE_ROOT}/../../../build/lib/common.sh"

# Include common constants and other vars needed when building image
source "${MAKE_ROOT}/_output/$EKSD_RELEASE_BRANCH/kind-node-image-build-args"

INTERMEDIATE_NODE_IMAGE="kind/build/node-image:$KUBE_VERSION-$EKSD_RELEASE_BRANCH"

DEPENDENCIES_DIR=$MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies/linux-$ARCH
FILES_DIR=$DEPENDENCIES_DIR/files
ROOT_FS=$FILES_DIR/rootfs

# Uses kind cli to build node-image, usually the fake kubernetes src
# to download the eks-d release artifacts instead of building from src
# ouput image: $INTERMEDIATE_NODE_IMAGE
function build::kind::build_node_image(){

    build::common::check_for_qemu linux/$ARCH

    KIND_PATH="$MAKE_ROOT/_output/bin/kind/$(uname | tr '[:upper:]' '[:lower:]')-$BUILDER_PLATFORM_ARCH/kind"
    $KIND_PATH build node-image --type file $DEPENDENCIES_DIR/eksd/kubernetes/server.tar.gz \
        --base-image $INTERMEDIATE_BASE_IMAGE --image $INTERMEDIATE_NODE_IMAGE --arch $ARCH
}

function build::kind::validate_versions(){
    local -r container_id=$1
    local -r eksd_tag="$KUBE_VERSION-eks-" # there is a commit hash at the end
    # We expect certain versions to be in the kind base/node images since we end up replacing them
    # validate they haven't changed
    if ! docker exec -i $container_id grep "image: $KINDNETD_IMAGE_TAG"  /kind/manifests/default-cni.yaml > /dev/null 2>&1; then
        echo "Did not find expected version of kindnetd: $KINDNETD_IMAGE_TAG"
        exit 1
    fi
    if ! docker exec -i $container_id grep "$LOCAL_PATH_HELPER_IMAGE_TAG"  /kind/manifests/default-storage.yaml > /dev/null 2>&1; then
        echo "Did not find expected version of debian base: $LOCAL_PATH_HELPER_IMAGE_TAG"
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
    if ! docker exec -i $container_id kubectl version -o json --client=true | jq -r ".clientVersion.gitVersion"  | grep "$eksd_tag" > /dev/null 2>&1; then
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

    docker cp $ROOT_FS/usr/local/bin/crictl $CONTAINER_ID:/usr/local/bin
    docker exec --privileged -i $CONTAINER_ID bash -c "nohup containerd > /dev/null 2>&1 &"
    # Wait for containerd socket to be ready (similar to upstream kind's WaitForReady)
    # This is needed because containerd may take longer to initialize under QEMU emulation
    docker exec --privileged -i $CONTAINER_ID bash -c '
for i in $(seq 0 10); do
  if [ -S /run/containerd/containerd.sock ]; then
    ctr info > /dev/null 2>&1 && exit 0
  fi
  sleep "$i"
done
echo "Timed out waiting for containerd socket"
exit 1
'
    docker exec --privileged -i $CONTAINER_ID crictl images

    # pull local-path-provisioner + al2 helper image
    IMAGES_FOLDER=$MAKE_ROOT/_output/$EKSD_RELEASE_BRANCH/dependencies/images
    mkdir -p $IMAGES_FOLDER
    IMAGES=("$AL2_HELPER_IMAGE" "$LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE" "$PAUSE_IMAGE_TAG_OVERRIDE" "$KIND_KINDNETD_IMAGE_OVERRIDE" "$ETCD_IMAGE_TAG" "$COREDNS_IMAGE_TAG")

    declare -A release_image_overrides
    release_image_overrides["$LOCAL_PATH_PROVISONER_IMAGE_TAG_OVERRIDE"]=$LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE
    release_image_overrides["$KIND_KINDNETD_IMAGE_OVERRIDE"]=$KIND_KINDNETD_RELEASE_OVERRIDE

    for image in "${IMAGES[@]}"; do
        # docker pull when passing a platform will pull the image for the given platform
        # and that new image id will "take over" the tag which we can use to trigger the save
        build::docker::retry_pull --platform linux/$ARCH $image
        image_id=$(docker images $image --format "{{.ID}}")

        if [ ! -f $IMAGES_FOLDER/$image_id.tar ]; then
            if [[ $image =~ local-path-provisioner ]] || [[ $image =~ kindnetd ]]; then
                docker tag $image_id ${release_image_overrides[$image]}
                docker save ${release_image_overrides[$image]} -o $IMAGES_FOLDER/$image_id.tar
            else
                docker save $image -o $IMAGES_FOLDER/$image_id.tar
            fi
        fi
        docker exec --privileged -i $CONTAINER_ID \
            ctr --namespace=k8s.io images import --all-platforms --no-unpack - < $IMAGES_FOLDER/$image_id.tar
    done

    docker exec --privileged -i $CONTAINER_ID crictl images
    
    # Validate all images in the image are from public.ecr
    EXPECTED_FINAL_IMAGES=("amazonlinux/amazonlinux" "kind/kindnetd" "rancher/local-path-provisioner" \
        "kubernetes/kube-apiserver" "kubernetes/kube-controller-manager" "kubernetes/kube-proxy" \
        "kubernetes/kube-scheduler" "kubernetes/pause" "coredns/coredns" "etcd-io/etcd")
    declare -a FOUND_EXPECTED_IMAGES
    FINAL_IMAGES=$(docker exec -i $CONTAINER_ID crictl images -o json | jq ".images[].repoTags[]" -r) 
    mapfile -t FINAL_IMAGES <<< "$FINAL_IMAGES"
    declare -p FINAL_IMAGES
    for image in "${FINAL_IMAGES[@]}"; do
        if [[ ! $image =~ ^public\.ecr\.aws\/ ]]; then
            echo "$image is not from public.ecr.aws!"
            exit 1
        fi
        EXPECTED_IMAGE=false
        for expected in "${EXPECTED_FINAL_IMAGES[@]}"; do
            if [[ $image =~ $expected ]]; then
                FOUND_EXPECTED_IMAGES+=( "$expected}" )
                EXPECTED_IMAGE=true
                break
            fi
        done
        if ! $EXPECTED_IMAGE; then
            echo "$image is not expected to be included in final image!"
            exit 1
        fi
        if [[ $ARCH != $(docker exec --privileged -i $CONTAINER_ID \
                crictl inspecti -o go-template --template={{.info.imageSpec.architecture}} $image) ]]; then
            echo "saved image: $image is not the correct arch: $ARCH!"
            exit 1
        fi
    done

    if [[ "${#FOUND_EXPECTED_IMAGES[@]}" != "${#EXPECTED_FINAL_IMAGES[@]}" ]]; then
        echo "${EXPECTED_FINAL_IMAGES[*]} are expected to be included in the final image but only ${FOUND_EXPECTED_IMAGES[*]} exist!"
        exit 1
    fi

    if [[ $(docker exec --privileged -i $CONTAINER_ID ctr -n k8s.io snapshots ls | wc -l) -gt 1 ]]; then
        echo "Snapshots exists, all images should have been loaded but not unpacked!"
        exit 1
    fi

    docker exec --privileged -i $CONTAINER_ID pkill containerd

    docker commit --change 'ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]' $CONTAINER_ID | sed -E 's/.*sha256:(.*)$/\1/'

    docker kill $CONTAINER_ID

    # Copy files created by the node build process and ctr import out to be used in next build
    mkdir -p $ROOT_FS/var/lib/containerd $ROOT_FS/etc/containerd
    docker cp $CONTAINER_ID:/etc/containerd/config.toml $ROOT_FS/etc/containerd
    docker cp $CONTAINER_ID:/kind $ROOT_FS
    docker cp $CONTAINER_ID:/var/lib/containerd/io.containerd.content.v1.content $ROOT_FS/var/lib/containerd

    mkdir -p $ROOT_FS/usr/bin
    for binary in kubeadm kubelet kubectl; do
        docker cp $CONTAINER_ID:/usr/bin/$binary $ROOT_FS/usr/bin
    done
    
    # meta.db is the databse for containerd and its not reproducible since dates are included
    # moving out of file directly to support reusing buildkit cache as much as possible
    docker cp $CONTAINER_ID:/var/lib/containerd/io.containerd.metadata.v1.bolt $FILES_DIR

    # update kind default manifests to use overridden images
	sed -i "s,image: $LOCAL_PATH_PROVISONER_IMAGE_TAG,image: $LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE," $ROOT_FS/kind/manifests/default-storage.yaml
	sed -i "s,$LOCAL_PATH_HELPER_IMAGE_TAG,$AL2_HELPER_IMAGE," $ROOT_FS/kind/manifests/default-storage.yaml
	sed -i "s,image: $KINDNETD_IMAGE_TAG,image: $KIND_KINDNETD_RELEASE_OVERRIDE," $ROOT_FS/kind/manifests/default-cni.yaml
	# Update containerd config to have eks-d pause image tag
	sed -i "s,$PAUSE_IMAGE_TAG,$PAUSE_IMAGE_TAG_OVERRIDE," $ROOT_FS/etc/containerd/config.toml
}

if command -v docker &> /dev/null && docker info > /dev/null 2>&1 ; then
    build::kind::build_node_image
    build::kind::load_images
fi
