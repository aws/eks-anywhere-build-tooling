ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-exporter _output/harbor-exporter/ _output/LICENSES ATTRIBUTION.txt /

RUN yum install -y tzdata shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -r -m -g 10000 -u 10000 harbor \
    && yum erase -y shadow-utils \
    && mv /harbor-exporter /harbor/harbor_exporter \
    && chown -R harbor:harbor /etc/pki/tls/certs \
    && chown -R harbor:harbor /harbor/ \
    && chmod u+x /harbor/entrypoint.sh \
    && chmod u+x /harbor/install_cert.sh \
    && chmod u+x /harbor/harbor_exporter

WORKDIR /harbor
USER harbor

ENTRYPOINT ["/harbor/entrypoint.sh"]