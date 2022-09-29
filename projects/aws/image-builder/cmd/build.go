package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/image-builder/builder"
)

var (
	bo                   = &builder.BuildOptions{}
	vSphereConfigFile    string
	baremetalConfigFile  string
	nutanixAHVConfigFile string
	err                  error
)

var buildCmd = &cobra.Command{
	Use:   "build --os <image os> --hypervisor <target hypervisor>",
	Short: "Build EKS Anywhere Node Image",
	Long:  "This command is used to build EKS Anywhere node images",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating builder config")
		err = ValidateInputs(bo)
		if err != nil {
			log.Fatalf(err.Error())
		}
		bo.BuildImage()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVar(&bo.Os, "os", "", "Operating system to use for EKS-A node image")
	buildCmd.Flags().StringVar(&bo.Hypervisor, "hypervisor", "", "Target hypervisor EKS-A node image")
	buildCmd.Flags().StringVar(&baremetalConfigFile, "baremetal-config", "", "Path to Baremetal Config file")
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

func ValidateInputs(bo *builder.BuildOptions) error {
	if bo.Os != builder.Ubuntu && bo.Os != builder.RedHat {
		log.Fatalf("Invalid OS type. Please choose ubuntu or redhat")
	}

	if (bo.Hypervisor != builder.VSphere) && (bo.Hypervisor != builder.Baremetal) {
		log.Fatalf("Invalid hypervisor. Please choose vsphere or baremetal")
	}

	configPath := ""
	switch bo.Hypervisor {
	case builder.VSphere:
		configPath = vSphereConfigFile
	case builder.Baremetal:
		configPath = baremetalConfigFile
	}
	bo.Os = strings.ToLower(bo.Os)
	bo.Hypervisor = strings.ToLower(bo.Hypervisor)

	if configPath == "" {
		if bo.Hypervisor == builder.VSphere || (bo.Hypervisor == builder.Baremetal && bo.Os == builder.RedHat) {
			return fmt.Errorf("%s-config is a required flag for %s hypervisor or when os is redhat", bo.Hypervisor, bo.Hypervisor)
		}
	} else {
		configPath, err = filepath.Abs(configPath)
		if err != nil {
			return fmt.Errorf("Error converting %s config file path to absolute path", bo.Hypervisor)
		}
		config, err := ioutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Error reading %s config file", bo.Hypervisor)
		}
		switch bo.Hypervisor {
		case builder.VSphere:
			if err = json.Unmarshal(config, &bo.VsphereConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				if err = validateRedhat(bo.VsphereConfig.RhelUsername, bo.VsphereConfig.RhelPassword, bo.VsphereConfig.IsoUrl); err != nil {
					return err
				}
			}
			if bo.VsphereConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.VsphereConfig.IsoChecksum, bo.VsphereConfig.IsoChecksumType); err != nil {
					return err
				}
			}
		case builder.Baremetal:
			if err = json.Unmarshal(config, &bo.BaremetalConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				if err = validateRedhat(bo.BaremetalConfig.RhelUsername, bo.BaremetalConfig.RhelPassword, bo.BaremetalConfig.IsoUrl); err != nil {
					return err
				}
			}
			if bo.BaremetalConfig != nil && bo.BaremetalConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.BaremetalConfig.IsoChecksum, bo.BaremetalConfig.IsoChecksumType); err != nil {
					return err
				}
			}
		case builder.NutanixAHV:
			if err = json.Unmarshal(config, &bo.NutanixAHVConfig); err != nil {
				return err
			}
			if bo.NutanixAHVConfig.NutanixUserName == "" || bo.NutanixAHVConfig.NutanixPassword == "" {
				log.Fatalf("\"nutanix_username\" and \"nutanix_password\" are required fields in nutanix-config")
			}
			// TODO Validate other fields as well
		}
	}

	return nil
}

func validateRedhat(rhelUsername, rhelPassword, isoUrl string) error {
	if rhelUsername == "" || rhelPassword == "" {
		return fmt.Errorf("\"rhel_username\" and \"rhel_password\" are required fields in config when os is redhat")
	}
	if isoUrl == "" {
		return fmt.Errorf("\"iso_url\" is a required field in config when os is redhat")
	}
	return nil
}

func validateCustomIso(isoChecksum, isoChecksumType string) error {
	if isoChecksum == "" {
		return fmt.Errorf("Please provide a valid checksum for \"iso_checksum\" when providing \"iso_url\"")
	}
	if isoChecksumType != "sha256" && isoChecksumType != "sha512" {
		return fmt.Errorf("\"iso_checksum_type\" is a required field when providing iso_checksum. Checksum type can be sha256 or sha512")
	}
	return nil
}
