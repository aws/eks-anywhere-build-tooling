ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install -y e2fsprogs util-linux && \
    yum clean all && \
    rm -rf /var/cache/yum

COPY _output/bin/bottlerocket-bootstrap/$TARGETOS-$TARGETARCH/bottlerocket-bootstrap-snow /bottlerocket-bootstrap-snow
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

CMD ["/bottlerocket-bootstrap-snow"]
