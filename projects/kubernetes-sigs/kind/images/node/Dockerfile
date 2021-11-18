# The base image is the unversioned base built and pushed to the local repo
# This base image was based off eks-distro-base
ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE-amd64 AS builder-amd64
FROM $BUILDER_IMAGE-arm64 AS builder-arm64

ARG TARGETARCH
FROM builder-${TARGETARCH} AS builder

ARG TARGETOS
ARG TARGETARCH

ARG IMAGE_REPO
ARG AL2_HELPER_IMAGE
ARG LOCAL_PATH_PROVISONER_IMAGE_TAG
ARG DEBIAN_BASE_IMAGE_TAG
ARG LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE
ARG PAUSE_IMAGE_TAG_OVERRIDE
ARG PAUSE_IMAGE_TAG
ARG KINDNETD_IMAGE_TAG
ARG KIND_KINDNETD_RELEASE_OVERRIDE


# remove kubeadm override script and containerd/runc
RUN rm /usr/local/bin/kubeadm && \
 rm /etc/kubeadm.config

RUN sed -i "s,image: $LOCAL_PATH_PROVISONER_IMAGE_TAG,image: $LOCAL_PATH_PROVISONER_RELEASE_OVERRIDE," /kind/manifests/default-storage.yaml 
RUN sed -i "s,$DEBIAN_BASE_IMAGE_TAG,$AL2_HELPER_IMAGE," /kind/manifests/default-storage.yaml 
RUN sed -i "s,$PAUSE_IMAGE_TAG,$PAUSE_IMAGE_TAG_OVERRIDE," /etc/containerd/config.toml
RUN sed -i "s,image: $KINDNETD_IMAGE_TAG,image: $KIND_KINDNETD_RELEASE_OVERRIDE," /kind/manifests/default-cni.yaml

# Add licenses/attribution + etcdadm + etcd tarball
COPY $TARGETOS-$TARGETARCH/LICENSES /LICENSES
COPY $TARGETOS-$TARGETARCH/etcdadm/etcdadm /usr/bin/
COPY $TARGETOS-$TARGETARCH/cache/etcdadm /var/cache/etcdadm

# AL2 doesn't setup the ip6tables alternative the same as ubuntu
# since kind's entrypoint assumes this adding the alternative
RUN update-alternatives \
	--install /usr/sbin/ip6tables ip6tables /usr/sbin/ip6tables-legacy 10 \
	--slave /usr/sbin/ip6tables-restore ip6tables-restore /usr/sbin/ip6tables-restore-legacy \
	--slave /usr/sbin/ip6tables-save ip6tables-save /usr/sbin/ip6tables-save-legacy

# During base image build using buildkit the /tmp directory's perms get changed from
# 1777 to 3777.  This probably isnt normally an issue however there is a 
# k8s conformance test that checks for 1777 specifically
# delete and recreate here as the last step
RUN rm -rf /tmp && \
    mkdir -m 1777 /tmp


FROM $BASE_IMAGE AS base-amd64
FROM $BASE_IMAGE AS base-arm64
# in the case of unstacked etcd on docker we use this image and etcd is run
# as a systemd unit, adding env var to allow it run on arm
ENV ETCD_UNSUPPORTED_ARCH=arm64


FROM base-${TARGETARCH}

COPY --from=builder / /

ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]
