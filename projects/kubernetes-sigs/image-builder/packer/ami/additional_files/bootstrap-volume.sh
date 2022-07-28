#!/bin/bash
source /etc/eks/logging.sh

set -uxo pipefail

SCRIPT_LOG=/var/log/eks-bootstrap.log
touch $SCRIPT_LOG

# save stdout and stderr to file descriptors 3 and 4,
# then redirect them to "$SCRIPT_LOG"
# restore stdout and stderr at bottom of script
exec 3>&1 4>&2 >>$SCRIPT_LOG 2>&1

DEVICE=vda

# stop containerd and kubelet
systemctl stop containerd
systemctl stop kubelet

# backup containerd data
mv /var/lib/containerd /var/lib/containerd_backup

# recreate dir
mkdir /var/lib/containerd

# format disk
mkfs -t ext4 /dev/$DEVICE

# mount new device
echo "/dev/$DEVICE /var/lib/containerd/     ext4    defaults        0 0" >>/etc/fstab
mount -a

# move containerd data back
mv /var/lib/containerd_backup/* /var/lib/containerd/
rm -rf /var/lib/containerd_backup

# restart containerd and kubelet
systemctl restart containerd
systemctl restart kubelet

log::info "successfully mounted new device /dev/$DEVICE"
