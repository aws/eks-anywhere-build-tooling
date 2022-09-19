package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/image-builder/builder"
)

var (
	bo                   = &builder.BuildOptions{}
	vSphereConfigFile    string
	nutanixAHVConfigFile string
	err                  error
)

var buildCmd = &cobra.Command{
	Use:   "build --os <image os> --hypervisor <target hypervisor>",
	Short: "Build EKS Anywhere Node Image",
	Long:  "This command is used to build EKS Anywhere node images",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating builder config")
		if bo.Hypervisor == builder.VSphere {
			if vSphereConfigFile == "" {
				log.Fatalf("vsphere-config is a required flag for vsphere hypervisor")
			}
			vSphereConfigFile, err = filepath.Abs(vSphereConfigFile)
			if err != nil {
				log.Fatalf("Error converting vsphere config file path to absolute path")
			}
			vSphereConfig, err := ioutil.ReadFile(vSphereConfigFile)
			if err != nil {
				log.Fatalf("Error reading vsphere config file")
			}
			err = json.Unmarshal(vSphereConfig, &bo.VsphereConfig)
			if err != nil {
				log.Fatalf("Error unmarshalling vsphere config file")
			}
		}
		bo.ValidateInputs()
		bo.BuildImage()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&bo.Os, "os", "", "Operating system to use for EKS-A node image")
	buildCmd.Flags().StringVar(&bo.Hypervisor, "hypervisor", "", "Target hypervisor EKS-A node image")
	buildCmd.Flags().StringVar(&vSphereConfigFile, "vsphere-config", "", "Path to vSphere Config file")
	buildCmd.Flags().StringVar(&nutanixAHVConfigFile, "nutanix-config", "", "Path to Nutanix AHV Config file")
	buildCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-23", "EKS-D Release channel for node image. Can be 1-20, 1-21, 1-22 or 1-23")
	buildCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
	if err := buildCmd.MarkFlagRequired("os"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	if err := buildCmd.MarkFlagRequired("hypervisor"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	if err := buildCmd.MarkFlagRequired("release-channel"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
}
