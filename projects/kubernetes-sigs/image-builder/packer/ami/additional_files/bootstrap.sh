#!/bin/bash
source /etc/eks/logging.sh

set -uxo pipefail

SCRIPT_LOG=/var/log/eks-bootstrap.log
touch $SCRIPT_LOG

# save stdout and stderr to file descriptors 3 and 4,
# then redirect them to "$SCRIPT_LOG"
# restore stdout and stderr at bottom of script
exec 3>&1 4>&2 >>$SCRIPT_LOG 2>&1

echo "127.0.0.1   $(hostname)" >>/etc/hosts

# configure DNI
DNI=$(ip -br link | grep -Ev 'lo|ens3' | awk '{print $1}')
while [ -z "$DNI" ]
do
  # creating DNI is a separate api call, which has some delays
  log::info "DNI is not ready, retrying"
  DNI=$(ip -br link | grep -Ev 'lo|ens3' | awk '{print $1}')
done
log::info "Using DNI: $DNI"

MAC=$(ip -br link | grep -Ev 'lo|ens3' |  awk '{print $3}')
DEFAULT_GATEWAY=$(ip r | grep default | awk '{print $3}')

log::info "Using MAC: $MAC"
log::info "Using default gateway: $DEFAULT_GATEWAY"

STATIC_CONFIG_PATH=/tmp/static-ips.yaml

# only one DNI config is supported for initial launch
if [[ ! -f "$STATIC_CONFIG_PATH" ]]
then
log::info "Configuring DHCP"
cat<<EOF >/etc/netplan/config.yaml
network:
    version: 2
    renderer: networkd
    ethernets:
        $DNI:
            set-name: $DNI
            dhcp4: true
            dhcp-identifier: mac
            dhcp4-overrides:
                route-metric: 50
                send-hostname: true
                hostname: $(hostname)
            match:
                macaddress: $MAC
        ens3:
            routes:
                - to: 169.254.169.254
                  via: $DEFAULT_GATEWAY
EOF
else
log::info "Configuring static ips: $(cat $STATIC_CONFIG_PATH)"
STATIC_IP=$(grep -E address $STATIC_CONFIG_PATH | awk '{print $3}')
GATEWAY=$(grep -E gateway $STATIC_CONFIG_PATH | awk '{print $2}')
cat<<EOF >/etc/netplan/config.yaml
network:
    version: 2
    renderer: networkd
    ethernets:
        $DNI:
            set-name: $DNI
            addresses:
                - $STATIC_IP
            routes:
                - to: default
                  via: $GATEWAY
                  metric: 50
            match:
                macaddress: $MAC
        ens3:
            routes:
                - to: 169.254.169.254
                  via: $DEFAULT_GATEWAY
EOF
fi

netplan --debug apply

DNI_IP=$(ip route show | grep default | grep "$DNI" | awk '{print $9}')
while [ -z "$DNI_IP" ]
do
    # if ip is not ready in 10 s, need to retry netplan apply
    log::info "IP is not ready, retrying"
    netplan --debug apply
    sleep 10
    DNI_IP=$(ip route show | grep default | grep "$DNI" | awk '{print $9}')
done
log::info "IP leased: $DNI_IP"

log::info "network configuration finished"

# mount container data volume if exists
DEVICE=$(lsblk | grep -E vd | awk '{print $1}')
if [ -n "$DEVICE" ]; then
    log::info "found new device /dev/$DEVICE"

    # stop containerd and kubelet
    systemctl stop containerd
    systemctl stop kubelet

    MOUNT_POINT=/var/lib/containerd
    BACKUP_DIR=/var/lib/containerd_backup

    # backup containerd data
    mv "$MOUNT_POINT" "$BACKUP_DIR"

    # recreate dir
    mkdir "$MOUNT_POINT"

    # format disk
    mkfs -t ext4 /dev/"$DEVICE"

    # mount new device
    echo "/dev/$DEVICE $MOUNT_POINT     ext4    defaults        0 0" >>/etc/fstab
    mount -a

    # move containerd data back
    mv "$BACKUP_DIR"/* "$MOUNT_POINT"/
    rm -rf "$BACKUP_DIR"

    # restart containerd and kubelet
    systemctl restart containerd
    systemctl restart kubelet

    log::info "successfully mounted new device /dev/$DEVICE"
else
    log::info "no new device found, skipping volume mount process"
fi

# restore stdout and stderr
exec 1>&3 2>&4
