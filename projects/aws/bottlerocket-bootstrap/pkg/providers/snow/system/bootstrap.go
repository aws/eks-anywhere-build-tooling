package system

import (
	"fmt"
	"os"
)

func Bootstrap() {
	if err := configureKubernetesSettings(); err != nil {
		fmt.Printf("Error configuring Kubernetes settings: %v\n", err)
		os.Exit(1)
	}

	if err := configureDNI(); err != nil {
		fmt.Printf("Error configuring snow DNI: %v\n", err)
		os.Exit(1)
	}

	if err := mountContainerdVolume(); err != nil {
		fmt.Printf("Error configuring container volume: %v\n", err)
		os.Exit(1)
	}

	if err := rebootInstanceIfNeeded(); err != nil {
		fmt.Printf("Error rebooting instance: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Snow bootstrap tasks succeeded")
	os.Exit(0)
}
