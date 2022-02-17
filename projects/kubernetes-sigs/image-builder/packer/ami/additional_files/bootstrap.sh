#!/bin/bash

VIP="$1"

DNI=$(ip -br link | egrep -v 'lo|ens3|docker0' | awk '{print $1}')

while [ -z $DNI ]
do
  echo "DNI is not ready, retry after 5 s"
  sleep 5
  DNI=$(ip -br link | egrep -v 'lo|ens3|docker0' | awk '{print $1}')
done

MAC=$(ip -br link | egrep -v 'lo|ens3|docker0' |  awk '{print $3}')

DEFAULT_GATEWAY=$(ip r | grep default | awk '{print $3}')

cat<<EOF >/etc/netplan/config.yaml
network:
    version: 2
    renderer: networkd
    ethernets:
        $DNI:
            set-name: $DNI
            dhcp4: true
            dhcp4-overrides:
                route-metric: 50
            match:
                macaddress: $MAC
        ens3:
            routes:
                - to: 169.254.169.254
                  via: $DEFAULT_GATEWAY
EOF

netplan apply

# DNI needs several seconds to lease an ip via DHCP
MY_IP=$(ip route show | grep default | grep $DNI | awk '{print $9}')
while [ -z $MY_IP ]
do
    echo "ip is not ready, retry after 5 s"
    sleep 5
    MY_IP=$(ip route show | grep default | grep $DNI | awk '{print $9}')
done

# Using instance id as a unique hostname before we implement hostname in capas
INSTANCE_ID=$(curl 169.254.169.254/latest/meta-data/instance-id | sed -r 's/[.]+/-/g')

hostnamectl set-hostname $INSTANCE_ID
echo "preserve_hostname: true" > /etc/cloud/cloud.cfg.d/99_hostname.cfg

cat <<EOF > /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1   $INSTANCE_ID
EOF

cat <<EOF >/etc/sysctl.d/kubernetes.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF

sysctl --system

swapoff -a

# if vip is not provided, it's a worker node, we don't need kube-vip manifest
# if `kubeadm init` command doesn't exist in the user-data, it's not the first control plane node, we should generate the kube-vip manifest after the `kubeadm join` command finishes
if [ ! -z $VIP ] && grep -q "kubeadm init --config /run/kubeadm/kubeadm.yaml" /var/lib/cloud/instance/user-data.txt ; then
  /etc/eks/generate-kube-vip-manifest.sh $VIP
fi
