package system_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/providers/snow/system"
)

var networkConfigTest = []struct {
	name       string
	data       map[string]interface{}
	wantConfig string
}{
	{
		name: "with static ip",
		data: map[string]interface{}{
			"network": &system.Network{
				DNI: "dni0",
				StaticIP: &system.StaticIP{
					Address: "address0",
					Gateway: "gateway0",
				},
			},
			"instanceIP":        "instanceip",
			"defaultGateway":    "defaultgateway",
			"metadataServiceIP": "metaserverip",
		},
		wantConfig: `version = 2
[dni0.static4]
addresses = ["address0"]
[dni0.route]
to = "default"
via = "gateway0"
[eth0.static4]
addresses = ["instanceip/25"]
[[eth0.route]]
to = "metaserverip/32"
from = "instanceip"
via = "defaultgateway"
`,
	},
	{
		name: "with dhcp",
		data: map[string]interface{}{
			"network": &system.Network{
				DNI: "dni0",
			},
			"instanceIP":        "instanceip",
			"defaultGateway":    "defaultgateway",
			"metadataServiceIP": "metaserverip",
		},
		wantConfig: `version = 2
[dni0]
dhcp4 = true
primary = true
[eth0.static4]
addresses = ["instanceip/25"]
[[eth0.route]]
to = "metaserverip/32"
from = "instanceip"
via = "defaultgateway"
`,
	},
}

func TestNetworkConfiguration(t *testing.T) {
	g := NewWithT(t)
	for _, tt := range networkConfigTest {
		t.Run(tt.name, func(t *testing.T) {
			b, err := system.GenerateNetworkTemplate(tt.data)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(b)).To(Equal(tt.wantConfig))
		})
	}
}
