ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/files/notification-controller /
COPY _output/bin/notification-controller/$TARGETOS-$TARGETARCH/notification-controller /usr/local/bin/notification-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER controller

ENTRYPOINT ["notification-controller"]
