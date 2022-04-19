#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo snap install amazon-ssm-agent --classic
sudo systemctl stop snap.amazon-ssm-agent.amazon-ssm-agent.service
