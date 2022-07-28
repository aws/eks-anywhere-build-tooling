#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo bash -c 'cat <<'EOF' >> /etc/docker/daemon.json
{
  "log-driver": "journald",
  "log-level": "debug"
}
EOF'

sudo systemctl restart docker --no-block