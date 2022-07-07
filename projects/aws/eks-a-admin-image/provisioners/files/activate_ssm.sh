#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

getprops_from_ovfxml() {
/usr/bin/python - <<EOS
from xml.dom.minidom import parseString
ovfEnv = open("$1", "r").read()
dom = parseString(ovfEnv)
section = dom.getElementsByTagName("PropertySection")[0]
for property in section.getElementsByTagName("Property"):
  key = property.getAttribute("oe:key").replace('.','_')
  value = property.getAttribute("oe:value")
  print "{0}='{1}'".format(key,value)
dom.unlink()
EOS
}

[ -x /usr/bin/vmtoolsd ] || {
  echo "ERROR: VMware Tools are not installed. Exiting ..."
  exit 1
}

/usr/bin/vmtoolsd --cmd='info-get guestinfo.ovfEnv' >/tmp/ovf.xml 2>/dev/null

[ -s /tmp/ovf.xml ] || {
  echo "ERROR: Cannot get OVF parameters through VMware Tools. Exiting ..."
  exit 1
}

eval `getprops_from_ovfxml /tmp/ovf.xml`

SSM_ACTIVATION_DIR=/etc/amazon/ssm

mkdir -p $SSM_ACTIVATION_DIR

SSM_ACTIVATION_FILE=$SSM_ACTIVATION_DIR/.activated

if test -f "$SSM_ACTIVATION_FILE"; then
	echo "SSM is already activated."
else
	echo "Activating SSM Agent..."
    /snap/amazon-ssm-agent/current/amazon-ssm-agent -register -code "$ssm_activation_code" -id "$ssm_activation_id" -region "$ssm_activation_region"
    systemctl enable snap.amazon-ssm-agent.amazon-ssm-agent.service
    systemctl start snap.amazon-ssm-agent.amazon-ssm-agent.service
    touch $SSM_ACTIVATION_FILE
fi
