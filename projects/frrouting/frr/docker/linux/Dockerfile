ARG BASE_IMAGE
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE AS runtime-builder

# Install the same runtime RPMs for notice generation and the final image.
# Bound RPM's file limit when its chroot has no /proc.
RUN set -eux && \
    ulimit -n 1024 && \
    install_rpm json-c c-ares readline pcre2 libcap libatomic coreutils coreutils-common sed grep && \
    install_binary /usr/bin/flock /usr/sbin/ip

COPY _output/generate_runtime_attribution.sh /source-notices/
RUN bash /source-notices/generate_runtime_attribution.sh /newroot /source-notices/RUNTIME_ATTRIBUTION.txt

FROM runtime-builder AS frr-builder

ARG BUILD_JOBS=8

RUN set -eux && \
    yum install -y gcc make autoconf automake libtool bison flex cmake pkgconf-pkg-config \
        json-c-devel readline-devel c-ares-devel pcre2-devel libcap-devel python3-devel \
        elfutils-libelf-devel diffutils libatomic texinfo perl tar gzip findutils which && \
    yum clean all

COPY _output/src/libyang /src/libyang
COPY _output/src/tini /src/tini
COPY _output/src/frr /src/frr

# Keep libyang aligned with FRR's upstream docker/alpine/libyang/APKBUILD.
RUN set -eux && \
    cd /src/libyang && \
    cmake -B build -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr \
        -DCMAKE_INSTALL_LIBDIR=lib64 -DENABLE_TESTS=OFF -DENABLE_VALGRIND_TESTS=OFF && \
    make -C build -j"$BUILD_JOBS" && \
    make -C build install && \
    make -C build install DESTDIR=/newroot

RUN set -eux && \
    cd /src/tini && \
    cmake . -DCMAKE_BUILD_TYPE=Release && \
    make tini && \
    install -D -m 0755 tini /newroot/sbin/tini

# Match the upstream Alpine package's daemon paths, users and multipath limit.
# RPKI, Lua scripting, protobuf and documentation are not needed by frr-k8s.
RUN set -eux && \
    cd /src/frr && \
    ./bootstrap.sh && \
    ./configure --prefix=/usr --sysconfdir=/etc --localstatedir=/var \
        --sbindir=/usr/lib/frr --libdir=/usr/lib64 --enable-vtysh \
        --enable-multipath=64 --enable-vty-group=frrvty --enable-user=frr \
        --enable-group=frr --enable-pcre2posix --disable-protobuf --disable-doc && \
    make -j"$BUILD_JOBS" && \
    make install DESTDIR=/newroot && \
    install -m 0755 docker/alpine/docker-start /newroot/usr/lib/frr/docker-start && \
    install -d /newroot/etc/frr /newroot/var/run/frr /newroot/var/lib/frr /newroot/var/tmp/frr && \
    install -m 0644 tools/etc/frr/daemons /newroot/etc/frr/daemons

RUN set -eux && \
    printf '%s\n' 'frr:x:101:frr' 'frrvty:x:102:frr' >> /newroot/etc/group && \
    printf '%s\n' 'frr:x:100:101:frr:/var/run/frr:/sbin/nologin' >> /newroot/etc/passwd && \
    chown -R 100:101 /newroot/etc/frr /newroot/var/run/frr /newroot/var/lib/frr /newroot/var/tmp/frr && \
    cleanup frr

ARG GIT_HASH
COPY ATTRIBUTION.txt /source-notices/

# Use the notices generated from the runtime RPMs, including the base layer.
RUN set -eux && \
    test -n "$GIT_HASH" && \
    cat /source-notices/ATTRIBUTION.txt /source-notices/RUNTIME_ATTRIBUTION.txt > /newroot/ATTRIBUTION.txt && \
    printf '\nBuild recipes and configure options: https://github.com/aws/eks-anywhere-build-tooling/tree/%s/projects/frrouting/frr\n' "$GIT_HASH" >> /newroot/ATTRIBUTION.txt && \
    for directory in licenses doc; do \
        if [ -d "/newroot/usr/share/$directory" ]; then \
            mkdir -p "/newroot/LICENSES/amazon-linux/$directory"; \
            cp -a "/newroot/usr/share/$directory/." "/newroot/LICENSES/amazon-linux/$directory/"; \
        fi; \
    done

FROM scratch AS runtime-attribution
COPY --from=runtime-builder /source-notices/RUNTIME_ATTRIBUTION.txt /RUNTIME_ATTRIBUTION.txt

FROM $BASE_IMAGE

COPY --from=frr-builder /newroot /
COPY _output/LICENSES /LICENSES

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/lib/frr/docker-start"]
