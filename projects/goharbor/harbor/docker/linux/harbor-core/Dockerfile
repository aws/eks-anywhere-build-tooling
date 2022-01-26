ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install -y tzdata shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -r -m -g 10000 -u 10000 harbor \
    && mkdir /harbor/

HEALTHCHECK CMD curl --fail -s http://localhost:8080/api/v2.0/ping || curl -k --fail -s https://localhost:8443/api/v2.0/ping || exit 1
COPY ../harbor/make/photon/common/install_cert.sh /harbor/
COPY ../harbor/make/photon/core/entrypoint.sh /harbor/
COPY ../harbor/src/core/views /harbor/views
COPY ../harbor/make/migrations /harbor/migrations
COPY ../harbor/icons /harbor/icons
COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-core /harbor/harbor_core
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R harbor:harbor /etc/pki/tls/certs \
    && chown -R harbor:harbor /harbor/ \
    && chmod u+x /harbor/entrypoint.sh \
    && chmod u+x /harbor/install_cert.sh \
    && chmod u+x /harbor/harbor_core

WORKDIR /harbor/
USER harbor
ENTRYPOINT ["/harbor/entrypoint.sh"]
COPY ../harbor/make/photon/prepare/versions /harbor/