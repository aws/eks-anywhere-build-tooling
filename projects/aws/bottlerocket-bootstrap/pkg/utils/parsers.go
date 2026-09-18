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

	return getCurrentContextAPIServer(rawConfig)
}

func getCurrentContextAPIServer(rawConfig kubecmdapi.Config) (string, error) {
	currentContextName := rawConfig.CurrentContext
	if currentContextName == "" {
		return "", errors.New("current context is empty in kubeconfig")
	}

	currentContext, ok := rawConfig.Contexts[currentContextName]
	if !ok || currentContext == nil {
		return "", errors.Errorf("current context %q not found in kubeconfig", currentContextName)
	}
	if currentContext.Cluster == "" {
		return "", errors.Errorf("cluster name is empty for current context %q", currentContextName)
	}

	cluster, ok := rawConfig.Clusters[currentContext.Cluster]
	if !ok || cluster == nil {
		return "", errors.Errorf("cluster %q referenced by current context %q not found in kubeconfig", currentContext.Cluster, currentContextName)
	}
	if cluster.Server == "" {
		return "", errors.Errorf("API server is empty for cluster %q", currentContext.Cluster)
	}

	return cluster.Server, nil
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
