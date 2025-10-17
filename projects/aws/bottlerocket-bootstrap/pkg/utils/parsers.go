package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubecmd "k8s.io/client-go/tools/clientcmd"
	kubecmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

func getKubeConfigRaw(path string) (kubecmdapi.Config, error) {
	// Read the kubeconfig and create config using clientcmd tool from client-go
	kubeData, err := ioutil.ReadFile(path)
	if err != nil {
		return kubecmdapi.Config{}, errors.Wrapf(err, "Error reading kubeconfig %s", path)
	}
	clientConfig, err := kubecmd.NewClientConfigFromBytes(kubeData)
	if err != nil {
		return kubecmdapi.Config{}, errors.Wrap(err, "Error generating kubeconfig from clientset")
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return kubecmdapi.Config{}, errors.Wrap(err, "Error getting rawconfig from kubeconfig")
	}
	return rawConfig, nil
}

func GetApiServerFromKubeConfig(path string) (string, error) {
	rawConfig, err := getKubeConfigRaw(path)
	if err != nil {
		return "", errors.Wrap(err, "Error getting kubeconfig parsed into raw config")
	}

	// Get the server from auth information
	var server string
	if len(rawConfig.Clusters) != 1 {
		return "", errors.Wrap(err, "More than one cluster found in control-plane init admin.conf")
	}
	fmt.Printf("\n%+v\n", rawConfig.Clusters)
	for _, clusterInfo := range rawConfig.Clusters {
		server = clusterInfo.Server
		break
	}
	return server, nil
}

func UnmarshalPodDefinition(podDef []byte) (*v1.Pod, error) {
	pod := v1.Pod{}
	err := yaml.Unmarshal(podDef, &pod)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting unmarshalling pod spec into structs")
	}
	return &pod, nil
}

// ResolveContainerPort resolves an IntOrString port to an int32 port number.
// For numeric ports (Type == Int), it returns IntVal directly.
// For named ports (Type == String), it searches the container's Ports list
// for a matching name and returns the corresponding ContainerPort.
func ResolveContainerPort(port intstr.IntOrString, container *v1.Container) (int32, error) {
	switch port.Type {
	case intstr.Int:
		return port.IntVal, nil
	case intstr.String:
		portName := port.StrVal
		for _, p := range container.Ports {
			if p.Name == portName {
				return p.ContainerPort, nil
			}
		}
		return 0, fmt.Errorf("named port %q not found in container ports", portName)
	default:
		return 0, fmt.Errorf("invalid IntOrString type: %v", port.Type)
	}
}
