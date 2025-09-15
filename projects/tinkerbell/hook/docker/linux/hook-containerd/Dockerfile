ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder
# /usr/bin/service is the init system in final linuxkit produced tar
RUN mkdir -p /opt/containerd /etc/init.d && \
    ln -s /usr/bin/service /etc/init.d/020-containerd


# final linuxkit compat image
# this image is never run directly, instead the layers are extracted and combined into the vmlinuz tar
FROM scratch
ENTRYPOINT []
WORKDIR /

ARG TARGETOS
ARG TARGETARCH

COPY --from=builder /etc/init.d/ /etc/init.d/
COPY --from=builder /opt/containerd/ /opt/containerd/

COPY _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksa/containerd/containerd/$RELEASE_BRANCH/usr/local/bin /usr/bin
COPY hook/images/hook-containerd/etc /etc/
