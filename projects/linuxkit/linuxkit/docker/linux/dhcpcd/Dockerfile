ARG BASE_IMAGE # https://gallery.ecr.aws/eks-distro-build-tooling/eks-distro-minimal-base-glibc

ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

# Bash and the following utils are needed
# by the various dhcpcd hook scripts
# and the dhcp.sh provided by hook
RUN set -x && \
    install_rpm bash && \
    install_binary \
        /usr/bin/cat \
        /usr/bin/chmod \
        /sbin/ifconfig \    
        /usr/bin/mkdir \
        /usr/bin/sed \
        /usr/bin/sleep \        
        /usr/bin/rm  && \
    cleanup "deps"

FROM $BASE_IMAGE
ARG TARGETARCH
ARG TARGETOS

ARG MOBY_CONFIG
LABEL "org.mobyproject.config"="${MOBY_CONFIG}"

COPY --from=builder /newroot /
COPY _output/bin/linuxkit/$TARGETOS-$TARGETARCH/dhcpcd /
# put the config file in / because we bind mount /etc to handle resolv.conf symlinks
COPY linuxkit/pkg/dhcpcd/dhcpcd.conf /
COPY linuxkit/pkg/dhcpcd/usr/ /usr/

CMD ["/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf"]
