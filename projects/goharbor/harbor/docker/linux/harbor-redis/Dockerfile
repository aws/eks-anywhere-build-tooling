ARG BASE_IMAGE
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as redis-builder

WORKDIR /

RUN set -x \
    && amazon-linux-extras enable redis6 \
    && cp /etc/yum.repos.d/amzn2-extras.repo /newroot/etc/yum.repos.d/amzn2-extras.repo \
    && clean_install systemd true true \
    && clean_install redis \
    && remove_package "bash gawk ncurses ncurses-base info sed shadow-utils grep" \
    && remove_package systemd true \
    && cleanup "redis"

FROM $BASE_IMAGE

COPY --from=redis-builder /newroot/ /
COPY --chown=redis:redis _output/harbor-redis _output/LICENSES ATTRIBUTION.txt /

VOLUME /var/lib/redis
WORKDIR /var/lib/redis

USER redis
CMD ["/usr/bin/redis-server", "/etc/redis.conf"]