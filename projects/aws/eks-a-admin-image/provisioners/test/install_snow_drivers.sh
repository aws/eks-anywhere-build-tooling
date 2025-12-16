#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Test that dracut configuration exists
if [[ -f /etc/dracut.conf.d/snow-ec2.conf ]]; then
    echo "Snow EC2 dracut config exists. Pass."
else
    echo "Snow EC2 dracut config missing. Failing."
    exit 1
fi

# Test that systemd network configuration exists
if [[ -f /usr/lib/systemd/network/80-snow-ec2.network ]]; then
    echo "Snow EC2 systemd network config exists. Pass."
else
    echo "Snow EC2 systemd network config missing. Failing."
    exit 1
fi

# Test that dracut config contains required drivers
if grep -q "virtio virtio_ring virtio_blk virtio_net virtio_pci ata_piix libata scsi_mod sd_mod scsi_common" /etc/dracut.conf.d/snow-ec2.conf; then
    echo "Required virtio drivers configured in dracut. Pass."
else
    echo "Required virtio drivers missing from dracut config. Failing."
    exit 1
fi

# Test that network config contains driver virtio_net
if grep -q "Driver=virtio_net" /usr/lib/systemd/network/80-snow-ec2.network; then
    echo "Snow Network Driver found. Pass."
else
    echo "Snow Network Driver missing. Failing."
    exit 1
fi

echo "All Snow driver tests passed."
