FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY _output/hook-embedded/$TARGETOS-$TARGETARCH/images/ /etc/embedded-images/
# the name 001 is important as that is the order in which the scripts are executed
# we need this mounting to happen before the other init.d scripts run so that
# the mount points are available to them.
COPY hook/images/hook-embedded/images-mount.sh /etc/init.d/001-images-mount.sh
