#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

sudo apt-get install -y python python3 python3-pip
sudo snap install amazon-ssm-agent --classic
sudo systemctl stop snap.amazon-ssm-agent.amazon-ssm-agent.service

SSM_ACTIVATION_DIR=/etc/amazon/ssm

sudo mkdir -p $SSM_ACTIVATION_DIR
chmod +x /home/$USER/activate_ssm.sh
sudo mv /home/$USER/activate_ssm.sh $SSM_ACTIVATION_DIR

sudo sh -c "echo '@reboot root /etc/amazon/ssm/activate_ssm.sh' >> /etc/crontab"
