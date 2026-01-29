#!/usr/bin/env sh
set -euo pipefail

# Check reboot marker first - skip if already processed
INSTANCE_REBOOTED=/.bottlerocket/bootstrap-containers/current/rebooted
if [[ -f "$INSTANCE_REBOOTED" ]]; then
    echo "Configuration already processed"
    exit 0
fi

# Build net.toml configuration as a string
NET_CONFIG="version = 3\n"

# Get list of network interfaces (excluding loopback)
INTERFACE_LIST=$(ip -o link | grep -v 'lo' | awk '{print $2}' | sed 's/://g')

# Configure each interface
INDEX=0
for INTERFACE in $INTERFACE_LIST; do
    NET_CONFIG+="[$INTERFACE]\n"
    NET_CONFIG+="dhcp4 = true\n"
    
    # Set first interface as primary
    if [ $INDEX -eq 0 ]; then
        NET_CONFIG+="primary = true\n"
    fi
    
    INDEX=$((INDEX + 1))
done

# Encode the configuration to base64
NET_CONFIG_BASE64=$(echo -e "$NET_CONFIG" | base64 -w 0)

# Apply network configuration
apiclient network configure "base64:$NET_CONFIG_BASE64"


# Current system - handle reboot ourselves
touch "$INSTANCE_REBOOTED"
apiclient reboot