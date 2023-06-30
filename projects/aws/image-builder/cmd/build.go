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
	nutanixConfigFile    string
	cloudstackConfigFile string
	amiConfigFile        string
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
	buildCmd.Flags().StringVar(&bo.OsVersion, "os-version", "", "Operating system version to use for EKS-A node image. Can be 20.04 or 22.04 for Ubuntu or 8 for Redhat. ")
	buildCmd.Flags().StringVar(&bo.Hypervisor, "hypervisor", "", "Target hypervisor EKS-A node image")
	buildCmd.Flags().StringVar(&baremetalConfigFile, "baremetal-config", "", "Path to Baremetal Config file")
	buildCmd.Flags().StringVar(&vSphereConfigFile, "vsphere-config", "", "Path to vSphere Config file")
	buildCmd.Flags().StringVar(&nutanixConfigFile, "nutanix-config", "", "Path to Nutanix Config file")
	buildCmd.Flags().StringVar(&cloudstackConfigFile, "cloudstack-config", "", "Path to CloudStack Config file")
	buildCmd.Flags().StringVar(&amiConfigFile, "ami-config", "", "Path to AMI Config file")
	buildCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-27", "EKS-D Release channel for node image. Can be 1-23, 1-24, 1-25, 1-26 or 1-27")
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

	if err = validateSupportedHypervisors(bo.Hypervisor); err != nil {
		log.Fatal(err.Error())
	}

	if err = validateOSHypervisorCombinations(bo.Os, bo.Hypervisor); err != nil {
		log.Fatal(err.Error())
	}

	if bo.Os == builder.Ubuntu && bo.OsVersion == "" {
		// maintain previous default
		bo.OsVersion = "20.04"
	}

	if bo.Os == builder.RedHat && bo.OsVersion == "" {
		// maintain previous default
		bo.OsVersion = "8"
	}

	if err = validateOSVersion(bo.Os, bo.OsVersion); err != nil {
		log.Fatal(err.Error())
	}

	configPath := ""
	switch bo.Hypervisor {
	case builder.VSphere:
		configPath = vSphereConfigFile
	case builder.Baremetal:
		configPath = baremetalConfigFile
	case builder.Nutanix:
		configPath = nutanixConfigFile
	case builder.CloudStack:
		configPath = cloudstackConfigFile
	case builder.AMI:
		configPath = amiConfigFile
	}
	bo.Os = strings.ToLower(bo.Os)
	bo.Hypervisor = strings.ToLower(bo.Hypervisor)

	if bo.OsVersion != "" {
		// From this point forward use 2004 instead of 20.04 for Ubuntu versions to upstream image-builder
		bo.OsVersion = strings.ReplaceAll(bo.OsVersion, ".", "")
	}

	if configPath == "" {
		if bo.Hypervisor == builder.VSphere ||
			(bo.Hypervisor == builder.Baremetal && bo.Os == builder.RedHat) ||
			(bo.Hypervisor == builder.Nutanix) ||
			(bo.Hypervisor == builder.CloudStack) {
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
		case builder.Nutanix:
			if err = json.Unmarshal(config, &bo.NutanixConfig); err != nil {
				return err
			}

			if bo.NutanixConfig.NutanixUserName == "" || bo.NutanixConfig.NutanixPassword == "" {
				log.Fatalf("\"nutanix_username\" and \"nutanix_password\" are required fields in nutanix-config")
			}
			// TODO Validate other fields as well
		case builder.CloudStack:
			if err = json.Unmarshal(config, &bo.CloudstackConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				if err = validateRedhat(bo.CloudstackConfig.RhelUsername, bo.CloudstackConfig.RhelPassword, bo.CloudstackConfig.IsoUrl); err != nil {
					return err
				}
			}
			if bo.CloudstackConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.CloudstackConfig.IsoChecksum, bo.CloudstackConfig.IsoChecksumType); err != nil {
					return err
				}
			}
		case builder.AMI:
			// Default configuration for AMI builds
			amiFilter := builder.DefaultUbuntu2004AMIFilterName
			if bo.OsVersion == "2204" {
				amiFilter = builder.DefaultUbuntu2204AMIFilterName
			}

			amiConfig := &builder.AMIConfig{
				AMIFilterName:       amiFilter,
				AMIFilterOwners:     builder.DefaultUbuntuAMIFilterOwners,
				AMIRegions:          builder.DefaultAMIBuildRegion,
				AWSRegion:           builder.DefaultAMIBuildRegion,
				BuilderInstanceType: builder.DefaultAMIBuilderInstanceType,
				CustomRole:          "true",
				AnsibleExtraVars:    builder.DefaultAMIAnsibleExtraVars,
				ManifestOutput:      builder.DefaultAMIManifestOutput,
				RootDeviceName:      builder.DefaultAMIRootDeviceName,
				VolumeSize:          builder.DefaultAMIVolumeSize,
				VolumeType:          builder.DefaultAMIVolumeType,
			}
			if err = json.Unmarshal(config, amiConfig); err != nil {
				return err
			}
			if amiConfig.CustomRole == "true" {
				if (amiConfig.CustomRoleNameList == nil && amiConfig.CustomRoleNames == "") ||
					(amiConfig.CustomRoleNameList != nil && amiConfig.CustomRoleNames != "") {
					log.Fatalf("Exactly one of \"custom_role_name_list\" or \"custom_role_names\" must be provided")
				}

				if amiConfig.CustomRoleNameList != nil {
					amiConfig.CustomRoleNames = strings.Join(amiConfig.CustomRoleNameList, " ")
					amiConfig.CustomRoleNameList = nil
				}
			}

			bo.AMIConfig = amiConfig
		}
	}

	return nil
}

func validateOSHypervisorCombinations(os, hypervisor string) error {
	if hypervisor == builder.CloudStack && os != builder.RedHat {
		return fmt.Errorf("Invalid OS type. Only redhat OS is supported for CloudStack")
	}

	if hypervisor == builder.Nutanix && os != builder.Ubuntu {
		return fmt.Errorf("Invalid OS type. Only ubuntu OS is supported for Nutanix")
	}

	if hypervisor == builder.AMI && os != builder.Ubuntu {
		return fmt.Errorf("Invalid OS type. Only ubuntu OS is supported for AMI")
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

func validateSupportedHypervisors(hypervisor string) error {
	if builder.SliceContains(builder.SupportedHypervisors, hypervisor) {
		return nil
	}
	return fmt.Errorf("%s is not supported yet. Please select one of %s", hypervisor, strings.Join(builder.SupportedHypervisors, ","))
}

func validateOSVersion(os string, osVersion string) error {
	if os != builder.RedHat && os != builder.Ubuntu {
		return fmt.Errorf("%s is not a supported OS.", os)
	}

	if os == builder.Ubuntu && !builder.SliceContains(builder.SupportedUbuntuVersions,osVersion) {
		return fmt.Errorf("%s is not a supported version of Ubuntu. Please select one of %s", osVersion, strings.Join(builder.SupportedUbuntuVersions, ","))
	}

	if os == builder.RedHat && !builder.SliceContains(builder.SupportedRedHatVersions,osVersion) {
		return fmt.Errorf("%s is not a supported version of Redhat. Please select one of %s", osVersion, strings.Join(builder.SupportedRedHatVersions, ","))
	}

	return nil
}
