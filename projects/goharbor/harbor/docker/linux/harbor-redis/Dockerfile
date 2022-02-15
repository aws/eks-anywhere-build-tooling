ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-base
FROM $BASE_IMAGE

RUN amazon-linux-extras enable redis6 \
    && yum install -y redis >> /dev/null \
    && yum clean all

VOLUME /var/lib/redis
WORKDIR /var/lib/redis

COPY _output/harbor/make/photon/redis/docker-healthcheck /usr/bin/
COPY --chown=redis:redis _output/harbor/make/photon/redis/redis.conf /etc/redis.conf
# COPY _output/LICENSES /LICENSES
# COPY ATTRIBUTION.txt /ATTRIBUTION.txt

RUN chmod +x /usr/bin/docker-healthcheck

HEALTHCHECK CMD ["docker-healthcheck"]
USER redis
CMD ["redis-server", "/etc/redis.conf"]