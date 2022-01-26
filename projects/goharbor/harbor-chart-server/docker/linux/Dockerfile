ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

RUN yum install -y shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 chart \
    && useradd --no-log-init -m -g 10000 -u 10000 chart

COPY _output/bin/chartmuseum/$TARGETOS-$TARGETARCH/harbor-chart-server /home/chart/chartm
COPY ../harbor/make/photon/chartserver/docker-entrypoint.sh /home/chart/
COPY ../harbor/make/photon/common/install_cert.sh /home/chart/
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chown -R chart:chart /etc/pki/tls/certs \
    && chown -R chart:chart /home/chart \
    && chmod u+x /home/chart/chartm \
    && chmod u+x /home/chart/docker-entrypoint.sh \
    && chmod u+x /home/chart/install_cert.sh

USER chart

WORKDIR /home/chart

ENTRYPOINT ["./docker-entrypoint.sh"]

VOLUME ["/chart_storage"]

HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD curl -sS http://localhost:9999/health || curl -k -sS https://localhost:9443/health || exit 1