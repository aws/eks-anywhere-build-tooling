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

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
)

//go:embed config/net.toml
var netConfig string

const (
	rootfs            = "/.bottlerocket/rootfs"
	metadataServiceIP = "169.254.169.254"
)

var (
	netConfigPath = filepath.Join(rootfs, "/var/lib/netdog/net.toml")
	currentIPPath = filepath.Join(rootfs, "/var/lib/netdog/current_ip")
)

func configureDNI() error {
	if files.PathExists(netConfigPath) {
		return nil
	}

	instanceIP, err := instanceIP()
	if err != nil {
		return errors.Wrap(err, "error getting local instance ip")
	}

	dnis, err := dniList()
	if err != nil {
		return errors.Wrap(err, "error getting DNI list")
	}

	defaultGateway, err := defaultGateway()
	if err != nil {
		return errors.Wrap(err, "error getting default gateway")
	}

	val := map[string]interface{}{
		"dnis":              dnis,
		"instanceIP":        instanceIP,
		"defaultGateway":    defaultGateway,
		"metadataServiceIP": metadataServiceIP,
	}

	if err := files.WriteTemplate(netConfigPath, netConfig, val, 0o640); err != nil {
		return errors.Wrapf(err, "error writing network configuration to %s", netConfigPath)
	}

	return nil
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
		return "", errors.Errorf("requesting %s returns status code %s", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	return string(body), nil
}
