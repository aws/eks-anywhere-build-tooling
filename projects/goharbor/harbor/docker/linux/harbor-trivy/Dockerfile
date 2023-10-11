ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS
ARG TRIVY_VERSION

COPY _output/dependencies/$TARGETOS-$TARGETARCH/eksa/aquasecurity/harbor-scanner-trivy/scanner-trivy _output/dependencies/$TARGETOS-$TARGETARCH/eksa/aquasecurity/trivy/trivy _output/harbor-trivy _output/LICENSES ATTRIBUTION.txt /

RUN yum install -y shadow-utils tar gzip >> /dev/null \
    && yum clean all \
    && groupadd -f -r -g 10000 scanner && useradd --no-log-init -m -g 10000 -u 10000 scanner \
    && yum erase -y shadow-utils \
    && mv /trivy /usr/local/bin/trivy \
    && mv /scanner-trivy /home/scanner/bin/scanner-trivy \
    && chown -R scanner:scanner /etc/pki/tls/certs /home/scanner /home/scanner/entrypoint.sh /usr/local/bin/trivy /home/scanner/bin/scanner-trivy /home/scanner/install_cert.sh \
    && chmod u+x /home/scanner/entrypoint.sh \
    && chmod u+x /usr/local/bin/trivy \
    && chmod u+x /home/scanner/bin/scanner-trivy \
    && chmod u+x /home/scanner/install_cert.sh

HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD curl --fail -s http://localhost:8080/probe/healthy || curl -k --fail -s https://localhost:8443/probe/healthy || exit 1

ENV TRIVY_VERSION=${TRIVY_VERSION}

USER scanner

ENTRYPOINT ["/home/scanner/entrypoint.sh"]
