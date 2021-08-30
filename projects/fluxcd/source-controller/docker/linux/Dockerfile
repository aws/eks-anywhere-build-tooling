ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-git
ARG BUILDER_IMAGE

FROM $BASE_IMAGE as base

FROM $BUILDER_IMAGE as builder

COPY --from=base /etc/group /etc/group
COPY --from=base /etc/passwd /etc/passwd
RUN yum install -y shadow-utils && \
    groupadd -r controller && \
    useradd -r -g controller controller


FROM $BASE_IMAGE

COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

COPY _output/bin/source-controller/linux-amd64/source-controller /usr/local/bin/source-controller
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER controller

ENTRYPOINT ["source-controller"]
