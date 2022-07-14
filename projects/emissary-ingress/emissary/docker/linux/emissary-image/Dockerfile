ARG BASE_IMAGE

FROM public.ecr.aws/eks-distro-build-tooling/builder-base:latest as emissary-builder

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin

RUN yum install -y sudo python3-pip curl wget make tar openssl-devel bzip2-devel libffi-devel

WORKDIR /opt

RUN wget https://www.python.org/ftp/python/3.9.7/Python-3.9.7.tgz && \
    tar xzf Python-3.9.7.tgz 

WORKDIR /opt/Python-3.9.7

# Install target Python runtime for wheel compliation
RUN ./configure --enable-optimizations && \
    make altinstall && \
    rm -rf /opt/Python-3.9.7* && \
    rm -rf /usr/bin/python3* && \
    ln -sf /usr/local/bin/python3.9 /usr/bin/python && \
    ln -sf /usr/local/bin/python3.9 /usr/bin/python3

# Install Rust and Cargo for Orjson package
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y && \
    rm -rf /root/.rustup

WORKDIR /buildroot

ADD emissary/python python

RUN time pip3 install --no-deps -e python

RUN pip3 install pip-tools==6.3.1
RUN pip3 install PyYAML==5.4.1
RUN pip3 install orjson==3.6.6

# Build Go binaries
ADD emissary/api api
ADD emissary/cmd cmd
ADD emissary/pkg pkg
ADD emissary/vendor vendor
ADD emissary/go.mod go.mod
ADD emissary/go.sum go.sum

RUN mkdir -p /buildroot/bin && \
	time go build -o /buildroot/bin/ ./cmd/...
## Remove base Python3 version before compliation to force on 3.9.7

RUN rm -rf vendor/

RUN mkdir -p /ambassador/sidecars && \
    ln -s /buildroot/python/post_update.py /ambassador/post_update.py && \
    ln -s /buildroot/python/watch_hook.py /ambassador/watch_hook.py && \
    ln -s /buildroot/python/kubewatch.py /ambassador/kubewatch.py

# Not sure where this should go, found in the upstream build
ADD emissary/manifests/emissary/emissary-crds.yaml.in manifests/emissary/emissary-crds.yaml

FROM public.ecr.aws/appmesh/aws-appmesh-envoy:v1.22.2.0-prod as envoy

FROM $BASE_IMAGE

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot:/buildroot/bin

WORKDIR /opt

RUN yum update -y && \
    yum install -y python3-pip docker bash wget gzip make tar libcap htop gcc openssl-devel bzip2-devel libffi-devel

RUN wget https://www.python.org/ftp/python/3.9.7/Python-3.9.7.tgz && \
    tar xzf Python-3.9.7.tgz 

WORKDIR /opt/Python-3.9.7

# Install python 3.9.7 and remove all other python versions
RUN ./configure --enable-optimizations && \
    make altinstall && \
    rm -rf /opt/Python-3.9.7* && \
    rm -rf /usr/bin/python3* \

WORKDIR /buildroot

RUN yum remove -y make wget tar

RUN ln -sf /usr/local/bin/python3.9 /usr/bin/python && \
    ln -sf /usr/local/bin/python3.9 /usr/bin/python3

# Copy GO binaries
COPY --from=emissary-builder /buildroot  /buildroot
RUN ln -s /buildroot/bin/* /usr/local/bin/

# Copy Python code and wheels
COPY --from=emissary-builder /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages

RUN cd /buildroot/python && python3 setup.py install

RUN chmod u+s $(which docker)

# uninstall un-needed software
RUN python -m pip uninstall -y pip && \
    find /var \( -iname '*sqlite*' \) -exec rm -rf -- '{}' \+ && \
    find /usr \( -iname '*gdbm*' -o -iname '*smtplib*' -o -iname '*piptools*' -o -iname '*pip_tools*' -o -iname '*ensurepip*' \) -exec rm -rf -- '{}' \+ && \
    rm -rf /usr/bin/dockerd /usr/bin/docker /usr/bin/containerd /usr/bin/ctr /usr/local/lib/python3.9/test /usr/bin/python2* /usr/bin/containerd-shim*

# Setup Envoy
COPY --from=envoy /usr/bin/envoy /usr/local/bin/envoy
RUN setcap cap_net_bind_service=ei /usr/local/bin/envoy && \
    setcap cap_net_bind_service=p /buildroot/bin//capabilities_wrapper

########
COPY --from=emissary-builder /ambassador /ambassador

RUN chgrp -R 0 /ambassador && \
    chmod -R u+x /ambassador && \
    chmod -R g=u /ambassador /etc/passwd

COPY _output/LICENSES /LICENSES
COPY ATTRIBUTION.txt /ATTRIBUTION.txt

WORKDIR /ambassador

ENV HOME=/tmp/ambassador

ENTRYPOINT [ "bash", "/buildroot/python/entrypoint.sh" ]
