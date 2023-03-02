package system

import (
	"bufio"
	_ "embed"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v2"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
)

//go:embed config/net.toml
var netConfig string

const (
	rootfs            = "/.bottlerocket/rootfs"
	metadataServiceIP = "169.254.169.254"
	networkFilePath   = "/tmp/network.yaml"
)

var (
	netConfigPath = filepath.Join(rootfs, "/var/lib/netdog/net.toml")
	currentIPPath = filepath.Join(rootfs, "/var/lib/netdog/current_ip")
)

type Network struct {
	DeviceIP string     `yaml:"deviceIP"`
	DniCount int        `yaml:"dniCount"`
	Static   []StaticIP `yaml:"static,omitempty"`
}

type StaticIP struct {
	Address string `yaml:"address"`
	Gateway string `yaml:"gateway"`
	Primary bool   `yaml:"primary,omitempty"`
}

type NetworkMapping struct {
	DNI string
	*StaticIP
}

func configureDNI() error {
	if files.PathExists(netConfigPath) {
		return nil
	}

	instanceIP, err := instanceIP()
	if err != nil {
		return errors.Wrap(err, "error getting local instance ip")
	}

	defaultGateway, err := defaultGateway()
	if err != nil {
		return errors.Wrap(err, "error getting default gateway")
	}

	network, err := networkMapping()
	if err != nil {
		return errors.Wrap(err, "error generating network mapping")
	}

	data := map[string]interface{}{
		"network":           network,
		"instanceIP":        instanceIP,
		"defaultGateway":    defaultGateway,
		"metadataServiceIP": metadataServiceIP,
	}

	b, err := GenerateNetworkTemplate(data)
	if err != nil {
		return errors.Wrap(err, "error generating network template")
	}

	if err := files.Write(netConfigPath, b, 0o640); err != nil {
		return errors.Wrapf(err, "error writing network configuration to %s", netConfigPath)
	}

	return nil
}

func GenerateNetworkTemplate(data map[string]interface{}) ([]byte, error) {
	return files.ExecuteTemplate(netConfig, data)
}

// networkMapping returns a single network mapping of DNI and optional static IP.
// currently only support single primary DNI and static IP configuration.
func networkMapping() ([]NetworkMapping, error) {
	dniList, err := dniList()
	if err != nil {
		return nil, errors.Wrap(err, "error getting DNI list")
	}

	network, err := parseNetworkFile()
	if err != nil {
		return nil, errors.Wrap(err, "error parsing network file")
	}

	if len(dniList) != network.DniCount {
		return nil, errors.Wrap(err, "the number of DNI found does not match the dniCount in the network file")
	}

	if len(network.Static) > 0 && len(network.Static) < network.DniCount {
		return nil, errors.Wrap(err, "mix of using DHCP and static IP is not supported")
	}

	m := make([]NetworkMapping, 0, network.DniCount)
	for i, dni := range dniList {
		n := NetworkMapping{
			DNI: dni,
		}
		if len(network.Static) > 0 {
			n.StaticIP = &network.Static[i]
		}
		m = append(m, n)
	}

	return m, nil
}

func parseNetworkFile() (*Network, error) {
	userData, err := utils.ResolveBootstrapContainerUserData()
	if err != nil {
		return nil, errors.Wrap(err, "error resolving user data")
	}
	for _, file := range userData.WriteFiles {
		if file.Path == networkFilePath {
			network := &Network{}
			err := yaml.Unmarshal([]byte(file.Content), network)
			if err != nil {
				return nil, errors.Wrap(err, "error unmarshalling network file content")
			}
			return network, nil
		}
	}
	return nil, nil
}

func dniList() ([]string, error) {
	devices, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrap(err, "error getting the device list of links")
	}

	dnis := []string{}
	for _, device := range devices {
		name := device.Attrs().Name
		if name != "lo" && name != "eth0" {
			dnis = append(dnis, name)
		}
	}
	return dnis, nil
}

func defaultGateway() (string, error) {
	routes, err := netlink.RouteList(nil, syscall.AF_INET)
	if err != nil {
		return "", errors.Wrap(err, "error getting the route list")
	}

	for _, route := range routes {
		if route.Dst == nil {
			return route.Gw.String(), nil
		}
	}

	return "", errors.New("default gateway not found")
}

func instanceIP() (string, error) {
	url := fmt.Sprintf("http://%s/latest/meta-data/local-ipv4", metadataServiceIP)
	body, err := httpGet(url)
	if err != nil {
		return "", errors.Wrap(err, "error requesting instance ip through http")
	}

	return string(body), nil
}

func currentIP() (string, error) {
	f, err := os.Open(currentIPPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}

func deviceIP() (string, error) {
	network, err := parseNetworkFile()
	if err != nil {
		return "", errors.Wrap(err, "error parsing network file")
	}

	return network.DeviceIP, nil
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "error sending http GET request to %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("requesting %s returns status code %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	return string(body), nil
}
