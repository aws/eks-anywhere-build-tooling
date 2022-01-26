ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum update -y && yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -m -g 10000 -u 10000 harbor \
    && mkdir -p /etc/registry 

COPY ../harbor/make/photon/common/install_cert.sh /home/harbor
COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-registry /usr/bin/registry_DO_NOT_USE_GC
COPY ../harbor/make/photon/registryctl/start.sh /home/harbor
COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-registryctl /home/harbor/harbor_registryctl
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R harbor:harbor /etc/pki/tls/certs \
    && chown harbor:harbor /home/harbor/harbor_registryctl && chmod u+x /home/harbor/harbor_registryctl \
    && chown harbor:harbor /usr/bin/registry_DO_NOT_USE_GC && chmod u+x /usr/bin/registry_DO_NOT_USE_GC \
    && chown harbor:harbor /home/harbor/start.sh && chmod u+x /home/harbor/start.sh \
    && chown harbor:harbor /home/harbor/install_cert.sh && chmod u+x /home/harbor/install_cert.sh


HEALTHCHECK CMD curl --fail -s http://localhost:8080/api/health || curl -sk --fail --key /etc/harbor/ssl/registryctl.key --cert /etc/harbor/ssl/registryctl.crt https://localhost:8443/api/health || exit 1

VOLUME ["/var/lib/registry"]

ENTRYPOINT ["/home/harbor/start.sh"]

USER harbor
