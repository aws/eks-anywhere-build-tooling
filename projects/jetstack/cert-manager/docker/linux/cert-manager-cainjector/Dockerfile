ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base
ARG BUILDER_IMAGE

FROM $BASE_IMAGE as base

FROM $BUILDER_IMAGE as builder

COPY --from=base /etc/group /etc/group
COPY --from=base /etc/passwd /etc/passwd
RUN yum install -y shadow-utils && \
    useradd -r -u 1000 cainjector

FROM $BASE_IMAGE

COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY _output/bin/cert-manager/linux-amd64/cert-manager-cainjector /usr/bin/cainjector
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

USER 1000
ENTRYPOINT ["/usr/bin/cainjector"]
