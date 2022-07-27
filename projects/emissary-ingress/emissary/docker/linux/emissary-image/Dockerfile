ARG BASE_IMAGE 
ARG BUILDER_IMAGE
ARG APPMESH_VERSION
ARG TARGETARCH
ARG TARGETOS
FROM public.ecr.aws/appmesh/aws-appmesh-envoy:$APPMESH_VERSION as envoy

FROM $BUILDER_IMAGE as emissary-builder

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin:/newroot/buildroot

COPY emissary/python _output/bin/emissary/linux-amd64 emissary/manifests/emissary/emissary-crds.yaml.in /newroot/buildroot
COPY --from=envoy /usr/bin/envoy /THIRD-PARTY-LICENSES.txt /

WORKDIR /newroot/buildroot

RUN install_rpm libcap htop && \
    install_binary /usr/bin/curl /usr/bin/chgrp /usr/bin/chmod && \
    time pip3 install --no-deps -e . && \
    pip3 install pip-tools==6.3.1 && \
    pip3 install PyYAML==5.4.1 && \
    pip3 install orjson==3.6.6 && \
    mkdir -p /newroot/ambassador/sidecars /newroot/usr/lib/python3.9/site-packages /newroot/usr/local/bin /ambassador /buildroot && \
    mv /envoy /newroot/usr/local/bin/envoy && \
    mv /THIRD-PARTY-LICENSES.txt /newroot/ENVOY_LICENSES.txt && \
    cp -r /usr/lib/python3.9/site-packages /newroot/usr/lib/python3.9 && \
    mkdir -p /newroot/buildroot/manifests/emissary && \
    mv emissary-crds.yaml.in /newroot/buildroot/manifests/emissary/emissary-crds.yaml && \
    ln -sf /buildroot/post_update.py /newroot/ambassador/post_update.py && \
    ln -sf /buildroot/watch_hook.py /newroot/ambassador/watch_hook.py && \
    ln -sf /buildroot/kubewatch.py /newroot/ambassador/kubewatch.py && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/kubestatus && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/watt && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/agent && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/apiext && \
    ln -sf /newroot/buildroot/* /usr/local/bin/ && \
    cleanup "emissary"

FROM $BASE_IMAGE

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot HOME=/tmp/ambassador

COPY --from=emissary-builder /newroot /

WORKDIR /buildroot

RUN cd /buildroot && python3 setup.py install && \
    python3 -m pip uninstall -y pip && \
    setcap cap_net_bind_service=ei /usr/local/bin/envoy && \
    setcap cap_net_bind_service=p /buildroot/capabilities_wrapper && \
    chgrp -R 0 /ambassador && \
    chmod -R u+x /ambassador && \
    chmod -R g=u /ambassador /etc/passwd

WORKDIR /ambassador

COPY ATTRIBUTION.txt _output/LICENSES /

ENTRYPOINT [ "bash", "/buildroot/entrypoint.sh" ]
