# syntax=docker/dockerfile:experimental

# To around the lack of a overlayfs in presubmit jobs since they run on fargate
# we have to keep the layers to a minimum.
# All the source files and librarys are gatered into once folder on the base and then
# copied in the builder
ARG BASE_IMAGE
ARG BUILDER_IMAGE
FROM ${BASE_IMAGE} as base

FROM scratch as deps

COPY --from=base /usr/include/git2 ./usr/include/git2
COPY --from=base /usr/include/git2.h ./usr/include
COPY --from=base /usr/lib64/pkgconfig/*.pc ./usr/lib64/pkgconfig/
COPY --from=base /usr/lib64/libgit2* ./usr/lib64
COPY --from=base /usr/lib64/libssh2* ./usr/lib64


FROM ${BUILDER_IMAGE} as builder

ARG TARGETARCH
ARG TARGETOS

ARG GOPROXY
ENV GOPROXY=$GOPROXY

RUN --mount=type=bind,source=/,target=/files set -x && \
    cp -rf /files/* / && \
    cd /eks-anywhere-build-tooling && \
    make -C projects/fluxcd/source-controller simple-create-binaries BINARY_PLATFORMS=$TARGETOS/$TARGETARCH && \
    readelf -n projects/fluxcd/source-controller/_output/bin/source-controller/linux-amd64/source-controller && \
    if [ $TARGETARCH = "amd64" ] ; then make -C projects/fluxcd/source-controller generate-attribution ; fi

FROM scratch

COPY --from=builder /eks-anywhere-build-tooling/projects/fluxcd/source-controller/_output/ .
