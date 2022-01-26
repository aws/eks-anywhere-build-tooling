ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && mkdir -p /etc/registry \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -m -g 10000 -u 10000 harbor

COPY ../harbor/make/photon/common/install_cert.sh /home/harbor
COPY ../harbor/make/photon/registry/entrypoint.sh /home/harbor
COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-registry /usr/bin/registry_DO_NOT_USE_GC
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R harbor:harbor /etc/pki/tls/certs \
    && chown harbor:harbor /home/harbor/entrypoint.sh && chmod u+x /home/harbor/entrypoint.sh \
    && chown harbor:harbor /home/harbor/install_cert.sh && chmod u+x /home/harbor/install_cert.sh \
    && chown harbor:harbor /usr/bin/registry_DO_NOT_USE_GC && chmod u+x /usr/bin/registry_DO_NOT_USE_GC

HEALTHCHECK CMD curl --fail -s http://localhost:5000 || curl -k --fail -s https://localhost:5443 || exit 1

USER harbor

ENTRYPOINT ["/home/harbor/entrypoint.sh"]

VOLUME ["/storage"]
