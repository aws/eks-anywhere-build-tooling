ARG BASE_IMAGE 
ARG BUILDER_IMAGE
ARG APPMESH_VERSION
FROM $BUILDER_IMAGE as emissary-builder

COPY emissary/python /buildroot/python
COPY _output/bin/emissary/linux-amd64 /buildroot/bin

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin

WORKDIR /buildroot

RUN install_rpm libcap htop && \
    install_binary /usr/bin/curl /usr/bin/chgrp /usr/bin/chmod && \
    cleanup "emissary"  && \
    time pip3 install --no-deps -e python && \
    pip3 install pip-tools==6.3.1 && \
    pip3 install PyYAML==5.4.1 && \
    pip3 install orjson==3.6.6 && \
    mkdir -p /ambassador/sidecars && \
    ln -s /buildroot/python/post_update.py /ambassador/post_update.py && \
    ln -s /buildroot/python/watch_hook.py /ambassador/watch_hook.py && \
    ln -s /buildroot/python/kubewatch.py /ambassador/kubewatch.py && \
    ln -s /buildroot/bin/busyambassador /buildroot/bin/kubestatus && \
    ln -s /buildroot/bin/busyambassador /buildroot/bin/watt && \
    ln -s /buildroot/bin/busyambassador /buildroot/bin/agent && \
    ln -s /buildroot/bin/busyambassador /buildroot/bin/apiext && \
    ln -s /buildroot/bin/* /usr/local/bin/

FROM public.ecr.aws/appmesh/aws-appmesh-envoy:$APPMESH_VERSION as envoy
FROM $BASE_IMAGE

ARG TARGETARCH
ARG TARGETOS

COPY emissary/manifests/emissary/emissary-crds.yaml.in manifests/emissary/emissary-crds.yaml
COPY --from=emissary-builder /buildroot  /buildroot
COPY --from=emissary-builder /newroot /
COPY --from=envoy /usr/bin/envoy /usr/local/bin/envoy
COPY --from=envoy /THIRD-PARTY-LICENSES.txt  /ENVOY_LICENSES.txt
COPY --from=emissary-builder /ambassador /ambassador
COPY --from=emissary-builder /usr/lib/python3.9/site-packages /usr/lib/python3.9/site-packages
COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot:/buildroot/bin
ENV HOME=/tmp/ambassador

WORKDIR /buildroot

RUN cd /buildroot/python && python3 setup.py install && \
    python3 -m pip uninstall -y pip && \
    setcap cap_net_bind_service=ei /usr/local/bin/envoy && \
    setcap cap_net_bind_service=p /buildroot/bin//capabilities_wrapper && \
    chgrp -R 0 /ambassador && \
    chmod -R u+x /ambassador && \
    chmod -R g=u /ambassador /etc/passwd

WORKDIR /ambassador

ENTRYPOINT [ "bash", "/buildroot/python/entrypoint.sh" ]
