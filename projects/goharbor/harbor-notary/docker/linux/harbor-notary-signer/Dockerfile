ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 notary \
    && useradd --no-log-init -m -g 10000 -u 10000 notary

# TODO