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
FROM $BASE_IMAGE
ARG IMAGE_TAG=not-set

ADD conf/index.template /usr/share/nginx/index.template
ADD conf/default.conf /etc/nginx/conf.d/
ADD conf/hello.sh /hello.sh
RUN echo $IMAGE_TAG >/IMAGE_TAG

STOPSIGNAL SIGQUIT
ENTRYPOINT ["/hello.sh"]
CMD ["nginx", "-g", "daemon off;"]
