#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

### Configure dracut for Snow EC2
sudo mkdir -p /etc/dracut.conf.d
sudo tee /etc/dracut.conf.d/snow-ec2.conf > /dev/null <<EOF
add_drivers+=" virtio virtio_ring virtio_blk virtio_net virtio_pci ata_piix libata scsi_mod sd_mod scsi_common "
EOF

### Configure systemd network for Snow EC2
sudo mkdir -p /usr/lib/systemd/network
sudo tee /usr/lib/systemd/network/80-snow-ec2.network > /dev/null <<EOF
[Match]
Driver=virtio_net

[Link]
MTUBytes=9216

[Network]
DHCP=yes
IPv6DuplicateAddressDetection=0
LLMNR=no
DNSDefaultRoute=yes

[DHCPv4]
UseHostname=no
UseDNS=yes
UseNTP=yes
UseDomains=yes

[DHCPv6]
UseHostname=no
UseDNS=yes
UseNTP=yes
WithoutRA=solicit
EOF

### Regenerate dracut images
sudo dracut --force --verbose --regenerate-all

echo "Snow device drivers and network configuration completed"
