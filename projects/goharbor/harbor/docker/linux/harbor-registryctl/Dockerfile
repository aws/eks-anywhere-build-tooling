ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

RUN yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -m -g 10000 -u 10000 harbor \
    && mkdir -p /etc/registry 

COPY --chown=harbor:harbor _output/make/photon/common/install_cert.sh /home/harbor
COPY --chown=harbor:harbor _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksa/distribution/distribution/registry /usr/bin/registry_DO_NOT_USE_GC
COPY --chown=harbor:harbor _output/make/photon/registryctl/start.sh /home/harbor
COPY --chown=harbor:harbor _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-registryctl /home/harbor/harbor_registryctl
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R harbor:harbor /etc/pki/tls/certs \
    && chmod u+x /home/harbor/harbor_registryctl \
    && chmod u+x /usr/bin/registry_DO_NOT_USE_GC \
    && chmod u+x /home/harbor/start.sh \
    && chmod u+x /home/harbor/install_cert.sh


HEALTHCHECK CMD curl --fail -s http://localhost:8080/api/health || curl -sk --fail --key /etc/harbor/ssl/registryctl.key --cert /etc/harbor/ssl/registryctl.crt https://localhost:8443/api/health || exit 1

VOLUME ["/var/lib/registry"]

ENTRYPOINT ["/home/harbor/start.sh"]

USER harbor
