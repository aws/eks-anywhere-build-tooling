ARG BASE_IMAGE
ARG BUILDER_IMAGE

FROM $BUILDER_IMAGE AS frr-builder

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

# The Python minimal base supplies sh and Python. EKS-D's RPM helpers install
# the runtime packages and the libraries used by ip and flock.
# RPM 4.16 scans up to the hard file limit when its chroot has no /proc. Bound
# this build step's limit so BuildKit's very large default cannot stall scriptlets.
RUN set -eux && \
    ulimit -n 1024 && \
    install_rpm json-c c-ares readline pcre2 libcap libatomic coreutils coreutils-common sed grep && \
    install_binary /usr/bin/flock /usr/sbin/ip && \
    printf '%s\n' 'frr:x:101:frr' 'frrvty:x:102:frr' >> /newroot/etc/group && \
    printf '%s\n' 'frr:x:100:101:frr:/var/run/frr:/sbin/nologin' >> /newroot/etc/passwd && \
    chown -R 100:101 /newroot/etc/frr /newroot/var/run/frr /newroot/var/lib/frr /newroot/var/tmp/frr && \
    cleanup frr

FROM $BASE_IMAGE

COPY --from=frr-builder /newroot /
COPY ATTRIBUTION.txt /ATTRIBUTION.txt
COPY _output/LICENSES /LICENSES

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/lib/frr/docker-start"]
