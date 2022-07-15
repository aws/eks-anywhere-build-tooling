ARG BASE_IMAGE

FROM public.ecr.aws/eks-distro-build-tooling/builder-base:latest as builder
ADD node-disk-manager node-disk-manager
WORKDIR node-disk-manager
RUN go build -o manager ./cmd/manager

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

# old method
# COPY _output/bin/node-disk-manager/$TARGETOS-$TARGETARCH/manager /opt/node-disk-manager/manager

# new method
COPY --from=builder /node-disk-manager/manager /opt/node-disk-manager/manager

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000

ENTRYPOINT ["/opt/node-disk-manager/manager"]