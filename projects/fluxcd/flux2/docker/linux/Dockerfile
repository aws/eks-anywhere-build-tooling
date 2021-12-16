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
FROM $BUILDER_IMAGE as generic-builder

RUN curl https://distro.eks.amazonaws.com/kubernetes-1-21/releases/6/artifacts/kubernetes/v1.21.2/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod 755 /usr/local/bin/kubectl

FROM $BASE_IMAGE

COPY --from=generic-builder /newroot /

ARG TARGETARCH
ARG TARGETOS

COPY --from=generic-builder /usr/local/bin/kubectl /usr/local/bin/
COPY --chmod=755 _output/bin/flux2/$TARGETOS-$TARGETARCH/flux /usr/local/bin/flux
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENTRYPOINT [ "/usr/local/bin/flux" ]
