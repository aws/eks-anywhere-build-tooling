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
FROM $BUILDER_IMAGE as nginx-builder

RUN set -x && \
    amazon-linux-extras enable nginx1 && \
    cp /etc/yum.repos.d/amzn2-extras.repo /newroot/etc/yum.repos.d/amzn2-extras.repo && \
    # manually "install" systemd to avoid installing the entire dep tree
    clean_install systemd true true && \
    clean_install "nginx" && \
    # remove unncessary utils and other binary deps that are not needed by nginx during runtime
    remove_package "chkconfig gawk grep info make openssl p11-kit p11-kit-trust sed shadow-utils" && \
    remove_package systemd true && \
    cleanup "nginx"

# envsubst is needed to generate the config, but it comes from the gettext package
# which pulls in a number of deps, manually installing envsubt from the rpm
RUN set -x && \
    yumdownloader --destdir=/tmp/downloads gettext && \
    cd /newroot && \
    rpm2cpio /tmp/downloads/gettext*.rpm | cpio -idv ./usr/bin/envsubst


FROM $BASE_IMAGE

COPY --from=nginx-builder /newroot /

WORKDIR /

ADD conf/index.template /usr/share/nginx/index.template
ADD conf/default.conf /etc/nginx/conf.d/
ADD conf/entrypoint.sh /entrypoint.sh

STOPSIGNAL SIGQUIT
ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
