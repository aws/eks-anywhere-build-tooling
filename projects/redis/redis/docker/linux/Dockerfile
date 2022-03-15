# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG BASE_IMAGE
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as redis-builder

WORKDIR /
RUN set -x && \
    amazon-linux-extras enable redis6 && \
    cp /etc/yum.repos.d/amzn2-extras.repo /newroot/etc/yum.repos.d/amzn2-extras.repo && \
    clean_install systemd true true && \
    clean_install redis && \
    remove_package "bash gawk ncurses ncurses-base info sed shadow-utils grep" && \
    remove_package systemd true && \
    cleanup "redis"
RUN set -x && \
    sed -i -e 's/^logfile .*/logfile ""/' /newroot/etc/redis/redis.conf

FROM $BASE_IMAGE
ARG IMAGE_TAG=not-set
ENV REDIS_VERSION $IMAGE_TAG

COPY --from=redis-builder /newroot /

WORKDIR /

USER redis:redis
VOLUME /data
EXPOSE 6379
ENTRYPOINT ["/usr/bin/redis-server"]
CMD ["/etc/redis/redis.conf"]
