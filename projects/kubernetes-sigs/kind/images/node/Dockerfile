# The base image is the unversioned base built and pushed to the local repo
# This base image was based off eks-distro-base
ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-kind
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE as builder
FROM $BASE_IMAGE AS base-amd64
FROM $BASE_IMAGE AS base-arm64
# in the case of unstacked etcd on docker we use this image and etcd is run
# as a systemd unit, adding env var to allow it run on arm
# an env conf file is also added when copying files below
ENV ETCD_UNSUPPORTED_ARCH=arm64

ARG TARGETARCH
FROM base-${TARGETARCH} as node

ARG TARGETOS
ARG TARGETARCH

ARG PAUSE_IMAGE_TAG_OVERRIDE
ARG PAUSE_IMAGE_TAG

RUN set -x && \
	# remove kubeadm override script and containerd/runc
	rm /usr/local/bin/kubeadm && \
	rm /etc/kubeadm.config && \
	# Update containerd config to have eks-d pause image tag
	sed -i "s,$PAUSE_IMAGE_TAG,$PAUSE_IMAGE_TAG_OVERRIDE," /etc/containerd/config.toml && \
	# During base image build using buildkit the /tmp directory's perms get changed from
	# 1777 to 3777.  This probably isnt normally an issue however there is a 
	# k8s conformance test that checks for 1777 specifically
	# delete and recreate here as the last step
	rm -rf /tmp && \
	mkdir -m 1777 /tmp

# Copy kind manifests, containerd blobs, kube binaries from result of kind build node-image process
# Also includes licenses/attribution + etcdadm + etcd tarball
COPY $TARGETOS-$TARGETARCH/files/rootfs /

# This image is based off a minimal image, the /etc/os-release is not the same as the standard
# AL os-release. Since kubelet reads the os-release to render the `get nodes` data it would be better
# to for it to show the real AL info since even though this is a minimal image, its effectively AL
COPY --from=builder /etc/os-release /etc/os-release

# Non-reproducible containerd db adding as last layer to utilize as much
# cache as possible
COPY $TARGETOS-$TARGETARCH/files/io.containerd.metadata.v1.bolt /var/lib/containerd/io.containerd.metadata.v1.bolt

ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]
