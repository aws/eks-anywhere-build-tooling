#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

# Package cleanup and reset cloud-init
sudo apt-get -y autoremove
sudo apt-get -y clean
sudo cloud-init clean

# Remove ssh keys and hosts
sudo rm -rf /home/ubuntu/.ssh
sudo rm -rf /root/.ssh
sudo rm -f /etc/ssh/ssh_host_*

# Clean tmp
sudo rm -rf /tmp/*
sudo rm -rf /var/tmp/*

# Truncate audit logs
sudo touch /var/log/wtmp
sudo touch /var/log/lastlog

# Truncate other logs
sudo find /var/log -type f -iname '*.log' | sudo xargs truncate -s 0
sudo find /var/log -type f -name '*.gz' -exec rm {} +
sudo journalctl --rotate
sudo journalctl -m --vacuum-time=1s

# Remove bash history
cat /dev/null > ~/.bash_history && sudo rm -f /root/.bash_history && history -c
