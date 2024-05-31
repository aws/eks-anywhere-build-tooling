package system_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/aws/eks-anywhere-build-tooling/bottlerocket-bootstrap/pkg/providers/snow/system"
)

var networkConfigTest = []struct {
	name       string
	data       map[string]interface{}
	wantConfig string
}{
	{
		name: "with static ip",
		data: map[string]interface{}{
			"network": []system.NetworkMapping{
				{
					DNI: "dni0",
					StaticIP: &system.StaticIP{
						Address: "address0",
						Gateway: "gateway0",
						Primary: true,
					},
				},
				{
					DNI: "dni1",
					StaticIP: &system.StaticIP{
						Address: "address1",
						Gateway: "gateway1",
						Primary: false,
					},
				},
				{
					DNI: "dni2",
					StaticIP: &system.StaticIP{
						Address: "address2",
						Gateway: "gateway2",
					},
				},
			},
			"instanceIP":        "instanceip",
			"defaultGateway":    "defaultgateway",
			"metadataServiceIP": "metaserverip",
		},
		wantConfig: `version = 2
[dni0]
primary = true
[dni0.static4]
addresses = ["address0"]
[[dni0.route]]
to = "default"
via = "gateway0"
[dni1.static4]
addresses = ["address1"]
[[dni1.route]]
to = "default"
via = "gateway1"
[dni2.static4]
addresses = ["address2"]
[[dni2.route]]
to = "default"
via = "gateway2"
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
			"network": []system.NetworkMapping{
				{
					DNI: "dni0",
				},
				{
					DNI: "dni1",
				},
				{
					DNI: "dni2",
				},
			},
			"instanceIP":        "instanceip",
			"defaultGateway":    "defaultgateway",
			"metadataServiceIP": "metaserverip",
		},
		wantConfig: `version = 2
[dni0]
dhcp4 = true
primary = true
[dni1]
dhcp4 = true
[dni2]
dhcp4 = true
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
