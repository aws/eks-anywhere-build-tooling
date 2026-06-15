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
ENV UV_PYTHON_INSTALL_DIR=/opt/python

# Copy python code, binaries, pyproject.toml, LICENSE, and CRDs
COPY emissary/python /newroot/buildroot/python
COPY emissary/pyproject.toml emissary/LICENSE /newroot/buildroot/
COPY _output/bin/emissary/${TARGETOS}-${TARGETARCH} /newroot/buildroot
COPY emissary/pkg/api/getambassador.io/crds.yaml /newroot/buildroot/
COPY --from=envoy /usr/bin/envoy /THIRD-PARTY-LICENSES.txt /

WORKDIR /newroot/buildroot

RUN set -x && \
    pip3 install --no-cache-dir uv && \
    uv python install 3.13.5 && \
    mkdir -p /newroot/ambassador/sidecars /newroot/usr/local/bin /ambassador && \
    # Create venv with Python 3.13.5 and build/install wheels
    uv venv /newroot/opt/ambassador/emissary --python 3.13.5 && \
    . /newroot/opt/ambassador/emissary/bin/activate && \
    cd python/ambassador && uv build --wheel --out-dir /tmp/wheels && \
    cd ../ambassador-diag && uv build --wheel --out-dir /tmp/wheels && \
    cd /tmp && \
    uv pip install /tmp/wheels/*.whl && \
    deactivate && \
    # Fix venv paths: replace /newroot prefix in shebangs and pyvenv.cfg
    find /newroot/opt/ambassador/emissary/bin -type f -exec sed -i 's|/newroot/opt|/opt|g' {} + && \
    sed -i 's|/newroot/opt|/opt|g' /newroot/opt/ambassador/emissary/pyvenv.cfg && \
    # Create diagd entry point (uv pip install doesn't generate console_scripts from wheels)
    printf '#!/opt/ambassador/emissary/bin/python3\nfrom ambassador_diag.diagd import main\nmain()\n' > /newroot/opt/ambassador/emissary/bin/diagd && \
    chmod +x /newroot/opt/ambassador/emissary/bin/diagd && \
    # Copy Python runtime into newroot for final image
    cp -r /opt/python /newroot/opt/python && \
    # Set up binaries
    mv /envoy /newroot/usr/local/bin/envoy && \
    mv /THIRD-PARTY-LICENSES.txt /newroot/ENVOY_LICENSES.txt && \
    mkdir -p /newroot/buildroot/manifests/emissary && \
    mv /newroot/buildroot/crds.yaml /newroot/buildroot/manifests/emissary/emissary-crds.yaml && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/kubestatus && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/watt && \
    ln -sf /buildroot/busyambassador /newroot/buildroot/agent && \
    ln -sf /opt/ambassador/emissary/bin/diagd /newroot/usr/local/bin/diagd && \
    ln -sf /newroot/buildroot/* /usr/local/bin/

FROM $BASE_IMAGE

ENV PATH=/opt/ambassador/emissary/bin:/opt/python/cpython-3.13.5-linux-${TARGETARCH}-gnu/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/buildroot
ENV HOME=/tmp/ambassador

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

ENTRYPOINT [ "bash", "/buildroot/python/entrypoint.sh" ]
