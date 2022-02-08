ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG RELEASE_BRANCH
ARG TARGETARCH
ARG TARGETOS

RUN yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && mkdir -p /etc/registry \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -m -g 10000 -u 10000 harbor

COPY --chown=harbor:harbor _output/make/photon/common/install_cert.sh /home/harbor
COPY --chown=harbor:harbor _output/make/photon/registry/entrypoint.sh /home/harbor
COPY --chown=harbor:harbor _output/$RELEASE_BRANCH/dependencies/$TARGETOS-$TARGETARCH/eksa/distribution/distribution/registry /usr/bin/registry_DO_NOT_USE_GC
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R harbor:harbor /etc/pki/tls/certs \
    && chmod u+x /home/harbor/entrypoint.sh \
    && chmod u+x /home/harbor/install_cert.sh \
    && chmod u+x /usr/bin/registry_DO_NOT_USE_GC

HEALTHCHECK CMD curl --fail -s http://localhost:5000 || curl -k --fail -s https://localhost:5443 || exit 1

USER harbor

ENTRYPOINT ["/home/harbor/entrypoint.sh"]

VOLUME ["/storage"]
