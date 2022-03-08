ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY _output/bin/harbor/$TARGETOS-$TARGETARCH/harbor-jobservice _output/harbor-jobservice/ _output/LICENSES ATTRIBUTION.txt /

RUN yum install -y tzdata shadow-utils >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 harbor && useradd --no-log-init -r -m -g 10000 -u 10000 harbor \
    && yum erase -y shadow-utils \
    && mv /harbor-jobservice /harbor/harbor_jobservice \
    && chown -R harbor:harbor /etc/pki/tls/certs \
    && chown -R harbor:harbor /harbor/ \
    && chmod u+x /harbor/entrypoint.sh \
    && chmod u+x /harbor/install_cert.sh \
    && chmod u+x /harbor/harbor_jobservice

WORKDIR /harbor/

USER harbor

VOLUME ["/var/log/jobs/"]

HEALTHCHECK CMD curl --fail -s http://localhost:8080/api/v1/stats || curl -sk --fail --key /etc/harbor/ssl/job_service.key --cert /etc/harbor/ssl/job_service.crt https://localhost:8443/api/v1/stats || exit 1

ENTRYPOINT ["/harbor/entrypoint.sh"]