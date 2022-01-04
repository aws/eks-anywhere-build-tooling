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

COPY redis /redis
RUN set -x && yum install -y gcc make
WORKDIR /redis
RUN case "$(uname -p)" in \
    x86_64 | amd64) JEMALLOC_CONFIG_FLAGS="--with-lg-page=12 --with-lg-hugepage=21" ;; \
    *) JEMALLOC_CONFIG_FLAGS="--with-lg-page=16 --with-lg-hugepage=21" ;; \
    esac; \
    make JEMALLOC_CONFIG_FLAGS=${JEMALLOC_CONFIG_FLAGS} -j "$(nproc)" all; \
    make install;
WORKDIR /
RUN mkdir -p /data

FROM $BASE_IMAGE
ARG IMAGE_TAG=not-set
ENV REDIS_VERSION $IMAGE_TAG

COPY --from=redis-builder /redis/src/redis-server /redis-server
COPY --from=redis-builder /data /

WORKDIR /
VOLUME /data
WORKDIR /data

EXPOSE 6379
ENTRYPOINT ["/redis-server"]
