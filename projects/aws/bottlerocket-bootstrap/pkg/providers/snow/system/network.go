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
	staticIPFilePath  = "/tmp/static-ips.yaml"
)

var (
	netConfigPath = filepath.Join(rootfs, "/var/lib/netdog/net.toml")
	currentIPPath = filepath.Join(rootfs, "/var/lib/netdog/current_ip")
)

type staticIPs struct {
	static []StaticIP `yaml:"static"`
}

type StaticIP struct {
	Address string `yaml:"address"`
	Gateway string `yaml:"gateway"`
	Primary bool   `yaml:"primary,omitempty"`
}

type Network struct {
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

	network, err := dniMapping()
	if err != nil {
		return errors.Wrap(err, "error generating network mapping with dni and optional static ip")
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

// dniMapping returns a single network mapping of DNI and optional static IP.
// currently only support single primary DNI and static IP configuration.
func dniMapping() (*Network, error) {
	dniList, err := dniList()
	if err != nil {
		return nil, errors.Wrap(err, "error getting DNI list")
	}

	if len(dniList) <= 0 {
		return nil, errors.Wrap(err, "error finding any valid DNI")
	}

	staticIPList, err := parseStaticIPFile()
	if err != nil {
		return nil, errors.Wrap(err, "error parsing static ip file")
	}

	n := &Network{
		DNI: dniList[0],
	}

	if len(staticIPList) > 0 {
		n.Address = staticIPList[0].Address
		n.Gateway = staticIPList[0].Gateway
	}

	return n, nil
}

func parseStaticIPFile() ([]StaticIP, error) {
	userData, err := utils.ResolveUserData()
	if err != nil {
		return nil, errors.Wrap(err, "error resolving user data")
	}
	for _, file := range userData.WriteFiles {
		if file.Path == staticIPFilePath {
			staticIPs := &staticIPs{}
			err := yaml.Unmarshal([]byte(file.Content), staticIPs)
			if err != nil {
				return nil, errors.Wrap(err, "error unmarshalling static IP file content")
			}
			return staticIPs.static, nil
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
