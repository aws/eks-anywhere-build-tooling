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

ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-nginx
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

ARG IMAGE_TAG=not-set

ADD conf/index.template /newroot/usr/share/nginx/index.template
ADD conf/default.conf /newroot/etc/nginx/conf.d/
ADD conf/hello.sh /newroot/hello.sh

RUN set -x && \
    echo $IMAGE_TAG >/newroot/IMAGE_TAG && \
    mkdir -p /newroot/usr/share/nginx/txt/ && \
    chown -R 999:999 /newroot/usr/share/nginx /newroot/var/log/nginx && \
    chmod 777 /newroot/run && \
    ln -sf /dev/stdout /newroot/var/log/nginx/access.log && \
    ln -sf /dev/stderr /newroot/var/log/nginx/error.log

FROM $BASE_IMAGE

COPY --from=builder /newroot /

# Ideally we should run as nginx, but the chart which is pulled from upstream
# is configured to bind to port 80
# if/when we get that changed we can remove this
USER root

STOPSIGNAL SIGQUIT
ENTRYPOINT ["/hello.sh"]
CMD ["nginx", "-g", "daemon off;"]
