package kubeadm

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/files"
	"github.com/eks-anywhere-build-tooling/aws/bottlerocket-bootstrap/pkg/utils"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

const (
	apiServerManifestPath = "/.bottlerocket/rootfs/etc/kubernetes/manifests/kube-apiserver"
	ebsInitMarker         = "/tmp/run-ebs-init"
)

func getBootstrapToken() (string, error) {
	token, err := exec.Command(kubeadmBinary, "token", "list", "-o", "jsonpath=\"{.token}\"").Output()
	if err != nil {
		return "", errors.Wrap(err, "Error getting token")
	}

	// Remove leading and ending double quotes
	// tokens do not contain quotes in the middle
	replacedToken := strings.ReplaceAll(string(token), "\"", "")
	return replacedToken, nil
}

func getEncodedCA() (string, error) {
	caFileData, err := ioutil.ReadFile("/etc/kubernetes/pki/ca.crt")
	if err != nil {
		return "", errors.Wrap(err, "Error reading the ca data")
	}
	return base64.StdEncoding.EncodeToString(caFileData), nil
}

func setHostName(filepath string) error {
	fmt.Println("Setting hostname in config files")
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "Error getting hostname")
	}
	fileData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "Error reading kubeadm file")
	}
	fileDataStr := string(fileData)
	fileDataStr = strings.ReplaceAll(fileDataStr, "{{ ds.meta_data.hostname }}", hostname)

	// Write the file back
	if err := files.Write(filepath, []byte(fileDataStr), 0o640); err != nil {
		return errors.Wrap(err, "Error writing file")
	}
	fmt.Println("Wrote config file back to kubeadm")
	fmt.Println(fileDataStr)
	return nil
}

func getDNS(path string) (string, error) {
	dns, err := exec.Command(kubectl, "--kubeconfig", path, "get", "svc", "kube-dns", "-n", "kube-system", "-o", "jsonpath='{..clusterIP}'").Output()
	if err != nil {
		return "", errors.Wrap(err, "Error getting api server")
	}

	// Remove leading and ending double quotes
	// dns ip doesnt have quotes in the middle
	replacedDns := strings.ReplaceAll(string(dns), "'", "")
	return replacedDns, nil
}

// TODO: ioutil.ReadFile on the kubelet config fails saying no file found,
// but awk and cat command pass. Really weird error. Must try and debug this out
// Possibly we kill the kubeadm command prematurely because we just stat the file
// while the file is being written?
func getDNSFromJoinConfig(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	fmt.Println(string(data))
	dns, err := exec.Command("bash", "-c", "awk '/clusterDNS:/ { getline; print $2 }' "+path).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "Error getting api server")
	}

	return strings.TrimSuffix(string(dns), "\n"), nil
}

func getBootstrapFromJoinConfig(path string) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", errors.Wrap(err, "Error reading kubeadm join config file")
	}
	joinConfig := strings.TrimPrefix(string(data), "---")

	kubeadmJoinData, err := unmarshalIntoMap([]byte(joinConfig))
	if err != nil {
		return "", "", errors.Wrap(err, "failed unmarshalling yaml kubeadm join config to interfaces")
	}
	discovery := kubeadmJoinData["discovery"].(map[string]interface{})
	bootstrapToken := discovery["bootstrapToken"].(map[string]interface{})
	serverEndpoint := bootstrapToken["apiServerEndpoint"].(string)
	token := bootstrapToken["token"].(string)

	return "https://" + serverEndpoint, token, nil
}

func getLocalApiServerReadinessEndpoint() (string, error) {
	data, err := ioutil.ReadFile(apiServerManifestPath)
	if err != nil {
		return "", errors.Wrap(err, "Error reading ApiServer manifest file")
	}

	podDef, err := utils.UnmarshalPodDefinition(data)
	if err != nil {
		return "", errors.Wrap(err, "Error parsing pod def from manifest")
	}

	for _, container := range podDef.Spec.Containers {
		// Validate if readiness probe exists on the definition
		if container.ReadinessProbe != nil {
			readinessProbeHandler := container.ReadinessProbe.HTTPGet
			url := fmt.Sprintf("%s://%s:%d%s", readinessProbeHandler.Scheme, readinessProbeHandler.Host, readinessProbeHandler.Port.IntVal, readinessProbeHandler.Path)
			return url, nil
		}
	}
	return "", errors.New("Cannot find readiness probe exists on pod definition")
}

func isClusterWithExternalEtcd(kubeconfigPath string) (bool, error) {
	clusterConfiguration, err := getClusterConfigurationFromCluster(kubeconfigPath)
	if err != nil {
		return false, err
	}

	return isExternalEtcd(clusterConfiguration)
}

func getClusterConfigurationFromCluster(kubeconfigPath string) ([]byte, error) {
	clusterConfiguration, err := exec.Command(kubectl, "--kubeconfig", kubeconfigPath, "-n", "kube-system", "get", "cm", "kubeadm-config", "-o", "jsonpath='{..data.ClusterConfiguration}'").Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed getting kubeadm-config configMap")
	}

	return clusterConfiguration, nil
}

func isExternalEtcd(clusterConfiguration []byte) (bool, error) {
	config, err := unmarshalIntoMap(clusterConfiguration)
	if err != nil {
		return false, errors.Wrap(err, "error unmarshalling ClusterConfiguration")
	}

	etcd := config["etcd"].(map[string]interface{})

	return etcd["external"] != nil, nil
}

func unmarshalIntoMap(content []byte) (map[string]interface{}, error) {
	var parsedMap map[string]interface{}
	content = []byte(strings.Trim(string(content), "'"))
	err := yaml.Unmarshal(content, &parsedMap)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling into map of empty interface")
	}

	return parsedMap, err
}

type EbsInitControl struct {
	Cancel  context.CancelFunc
	OkChan  chan bool
	Timeout <-chan time.Time
}

// startEbsInit starts the ebs-init goroutine in the background
// if marker file is present. Currenlty ebs-init is best-effort
// and if it fails it won't prevent the instance from bootstrapping.
// A 2 minutes timeout is added so that instance bootstrap time is
// not inflated.
func startEbsInit() *EbsInitControl {
	if _, err := os.Stat(ebsInitMarker); err == nil {
		okChan := make(chan bool)
		ctx, cancel := context.WithCancel(context.Background())
		fmt.Printf("Starting ebs-init \n")
		ebsInitControl := &EbsInitControl{
			Timeout: time.After(2 * time.Minute),
			Cancel:  cancel,
			OkChan:  okChan,
		}

		go func(ctx context.Context, okChan chan bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					readFiles()
					okChan <- true
					return
				}
			}

		}(ctx, okChan)

		return ebsInitControl
	} else {
		fmt.Printf("Skipping ebs-init \n")
		return nil
	}
}

// readFiles walks the filesystem tree and reads all regular files
// skipping certain folders that contain special files.
func readFiles() {
	root := "/.bottlerocket/rootfs/"
	filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, err error) error {

		if dirEntry.IsDir() {
			if dirEntry.Name() == "proc" || dirEntry.Name() == "sys" || dirEntry.Name() == "run" {
				return filepath.SkipDir
			}
		} else {
			entryInfo, _ := dirEntry.Info()
			if entryInfo.Mode().IsRegular() {
				ioutil.ReadFile(path)
			}
		}
		return nil
	})
	fmt.Printf("All files read \n")
}
