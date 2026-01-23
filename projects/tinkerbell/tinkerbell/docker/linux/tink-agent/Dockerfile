# Tink-agent image for workflow execution on worker nodes
# Minimal image with just the agent binary and CA certificates

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base

FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/tinkerbell/$TARGETOS-$TARGETARCH/tink-agent /tink-agent
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT ["/tink-agent"]
