ARG BASE_IMAGE 
ARG BUILDER_IMAGE
ARG COMPILER_IMAGE
ARG APPMESH_VERSION
FROM public.ecr.aws/appmesh/aws-appmesh-envoy:$APPMESH_VERSION as envoy

FROM $BUILDER_IMAGE as yum-installer

RUN install_rpm libcap htop && \
    install_binary /usr/bin/curl /usr/bin/chgrp /usr/bin/chmod && \
    cleanup "yum-installer"


FROM $COMPILER_IMAGE as emissary-builder

ARG TARGETARCH
ARG TARGETOS

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin:/newroot/buildroot

COPY emissary/python _output/bin/emissary/${TARGETOS}-${TARGETARCH} emissary/manifests/emissary/emissary-crds.yaml.in /newroot/buildroot
COPY --from=envoy /usr/bin/envoy /THIRD-PARTY-LICENSES.txt /

WORKDIR /newroot/buildroot

RUN set -x && \
    mkdir -p /newroot/ambassador/sidecars /newroot/usr/lib/python3.11/site-packages /newroot/usr/local/bin /ambassador /buildroot && \
    time pip3 install --user --no-cache-dir --no-deps -e . && \
    pip3 install --user --no-cache-dir pip-tools==6.12.1 PyYAML==6.0 orjson==3.8.10 && \ 
    python3 setup.py install --user && \
    mv /envoy /newroot/usr/local/bin/envoy && \
    mv /THIRD-PARTY-LICENSES.txt /newroot/ENVOY_LICENSES.txt && \
    cp -r /root/.local/lib/python3.11/site-packages /newroot/usr/lib/python3.11 && \
    cp -r /root/.local/bin /newroot/usr/bin && \
    mkdir -p /newroot/buildroot/manifests/emissary && \
    mv emissary-crds.yaml.in /newroot/buildroot/manifests/emissary/emissary-crds.yaml && \
    ln -sf /buildroot/post_update.py /newroot/ambassador/post_update.py && \
    ln -sf /buildroot/watch_hook.py /newroot/ambassador/watch_hook.py && \
    ln -sf /buildroot/kubewatch.py /newroot/ambassador/kubewatch.py && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/kubestatus && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/watt && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/agent && \
    ln -sf /newroot/buildroot/* /usr/local/bin/

FROM $BASE_IMAGE

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot HOME=/tmp/ambassador

COPY --from=yum-installer /newroot /
COPY --from=emissary-builder /newroot /

WORKDIR /buildroot

RUN set -x && \
    setcap cap_net_bind_service=ei /usr/local/bin/envoy && \
    setcap cap_net_bind_service=p /buildroot/capabilities_wrapper && \
    chgrp -R 0 /ambassador && \
    chmod -R u+x /ambassador && \
    chmod -R g=u /ambassador /etc/passwd

WORKDIR /ambassador

COPY ATTRIBUTION.txt _output/LICENSES /

ENTRYPOINT [ "bash", "/buildroot/entrypoint.sh" ]
