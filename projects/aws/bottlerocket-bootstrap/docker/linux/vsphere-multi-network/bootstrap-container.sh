#!/usr/bin/env sh
set -euo pipefail

# Check reboot marker first - skip if already processed
INSTANCE_REBOOTED=/.bottlerocket/bootstrap-containers/current/rebooted
if [[ -f "$INSTANCE_REBOOTED" ]]; then
    echo "Configuration already processed"
    exit 0
fi

# Place network configuration - dynamically detect network interfaces
NET_CONFIG_PATH=/.bottlerocket/rootfs/var/lib/netdog/net.toml

# Start with version header
echo "version = 2" > "$NET_CONFIG_PATH"

# Get list of network interfaces (excluding loopback)
INTERFACE_LIST=$(ip -o link | grep -v 'lo' | awk '{print $2}' | sed 's/://g')

# Configure each interface
INDEX=0
for INTERFACE in $INTERFACE_LIST; do
    echo "[$INTERFACE]" >> "$NET_CONFIG_PATH"
    echo "dhcp4 = true" >> "$NET_CONFIG_PATH"
    
    # Set first interface as primary
    if [ $INDEX -eq 0 ]; then
        echo "primary = true" >> "$NET_CONFIG_PATH"
    fi
    
    INDEX=$((INDEX + 1))
done

# Forward-compatible guard logic
if [ -d "/.bottlerocket/rootfs/.bottlerocket" ]; then
    echo "New system detected - OS will handle migration"
    exit 0
fi

# Current system - handle reboot ourselves
touch "$INSTANCE_REBOOTED"
apiclient reboot