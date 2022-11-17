package system

import (
	"testing"

	. "github.com/onsi/gomega"
)

const multiDNIs = `version = 2
[1]
dhcp4 = true
[2]
dhcp4 = true
[3]
dhcp4 = true
[eth0.static4]
addresses = ["1.2.3.4/25"]
[[eth0.route]]
to = "metadata/32"
from = "1.2.3.4"
via = "gateway"
`

func TestGenerateUserData(t *testing.T) {
	g := NewWithT(t)

	testcases := []struct {
		name   string
		data   map[string]interface{}
		output string
	}{
		{
			name: "multi dnis",
			data: map[string]interface{}{
				"dnis":              []string{"1", "2", "3"},
				"instanceIP":        "1.2.3.4",
				"defaultGateway":    "gateway",
				"metadataServiceIP": "metadata",
			},
			output: multiDNIs,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			b, err := generateNetworkTemplate(testcase.data)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(b)).To(Equal(testcase.output))
		})
	}
}
