package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/image-builder/builder"
)

var (
	bo                        = &builder.BuildOptions{}
	vSphereConfigFile         string
	baremetalConfigFile       string
	nutanixConfigFile         string
	cloudstackConfigFile      string
	amiConfigFile             string
	additionalFilesConfigFile string
	builderType               string
	err                       error
)

var buildCmd = &cobra.Command{
	Use:   "build --os <image os> --hypervisor <target hypervisor> --release-channel <EKS-D Release channel>",
	Short: "Build EKS Anywhere Node Image",
	Long:  "This command is used to build EKS Anywhere node images corresponding to different hypervisors.",
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
	buildCmd.Flags().StringVar(&bo.OsVersion, "os-version", "", "Operating system version to use for EKS-A node image. Can be 20.04, 22.04 or 24.04 for Ubuntu, 8 or 9 for Redhat.")
	buildCmd.Flags().StringVar(&bo.Hypervisor, "hypervisor", "", "Target hypervisor for EKS-A node image")
	buildCmd.Flags().StringVar(&baremetalConfigFile, "baremetal-config", "", "Path to Baremetal Config file")
	buildCmd.Flags().StringVar(&vSphereConfigFile, "vsphere-config", "", "Path to vSphere Config file")
	buildCmd.Flags().StringVar(&nutanixConfigFile, "nutanix-config", "", "Path to Nutanix Config file")
	buildCmd.Flags().StringVar(&cloudstackConfigFile, "cloudstack-config", "", "Path to CloudStack Config file")
	buildCmd.Flags().StringVar(&amiConfigFile, "ami-config", "", "Path to AMI Config file")
	buildCmd.Flags().StringVar(&additionalFilesConfigFile, "files-config", "", "Path to Config file specifying additional files to be copied into EKS-A node image")
	buildCmd.Flags().StringVar(&bo.ReleaseChannel, "release-channel", "1-31", "EKS-D Release channel for node image. Can be 1-28, 1-29, 1-30, 1-31, 1-32, 1-33 or 1-34")
	buildCmd.Flags().BoolVar(&bo.Force, "force", false, "Force flag to clean up leftover files from previous execution")
	buildCmd.Flags().StringVar(&bo.Firmware, "firmware", "", "Desired firmware for image build. EFI is only supported for Ubuntu OVA & Raw, and Redhat 9 RAW builds.")
	buildCmd.Flags().StringVar(&bo.EKSAReleaseVersion, "eksa-release", "", "The EKS-A CLI version to build images for")
	buildCmd.Flags().StringVar(&bo.ManifestTarball, "manifest-tarball", "", "Path to Image Builder built EKS-D/A manifest tarball")
	buildCmd.Flags().IntVar(&bo.AnsibleVerbosity, "ansible-verbosity", 0, "Verbosity level for the Ansible tasks run during image building, should be in the range 0-6")
	buildCmd.Flags().BoolVar(&bo.AirGapped, "air-gapped", false, "Flag to instruct image builder to run in air-gapped mode. Requires --manifest-tarball to be set")
	buildCmd.Flags().StringVar(&builderType, "builder", builder.BuilderTypeIso, "Builder type for VSphere. Can be iso or clone")
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

	if bo.AnsibleVerbosity < 0 || bo.AnsibleVerbosity > 6 {
		log.Fatalf("Invalid Ansible verbosity level. Please provide a value in the range 0 to 6")
	}

	if bo.Os == builder.Ubuntu && bo.OsVersion == "" {
		// maintain previous default
		bo.OsVersion = "20.04"
	}

	if bo.Os == builder.RedHat && bo.OsVersion == "" {
		// maintain previous default
		bo.OsVersion = "8"
	}

	if bo.Hypervisor != builder.Nutanix && bo.Hypervisor != builder.CloudStack && bo.Hypervisor != builder.Baremetal && bo.Os == builder.RedHat && bo.OsVersion != "8" && bo.OsVersion != "9" {
		log.Fatalf("Invalid OS version for RedHat. Please choose 8 or 9")
	}

	if bo.Hypervisor != builder.VSphere && bo.BuilderType == builder.BuilderTypeClone {
		log.Fatalf("Clone builder is only supported for vSphere hypervisor")
	}

	if err = validateOSVersion(bo.Os, bo.OsVersion); err != nil {
		log.Fatal(err.Error())
	}

	if err = validateFirmware(bo.Firmware, bo.Os, bo.OsVersion, bo.Hypervisor); err != nil {
		log.Fatal(err.Error())
	}

	// Setting default to bios for everything but ubuntu/rhel 9 raw since that defaults to efi
	if bo.Firmware == "" {
		if bo.Hypervisor == builder.Baremetal && (bo.Os == builder.Ubuntu ||
			(bo.Os == builder.RedHat && builder.SliceContains(builder.SupportedRedHatEfiVersions, bo.OsVersion))) {
			bo.Firmware = builder.EFI
		} else {
			bo.Firmware = builder.BIOS
		}
	}

	// Airgapped
	if bo.AirGapped {
		if bo.ManifestTarball == "" {
			log.Fatalf("Please provide --manifest-tarball when running air-gapped builds")
		}

		if bo.Os != builder.Ubuntu {
			log.Fatalf("Only Ubuntu os is supported for air-gapped builds")
		}
		if bo.Hypervisor == builder.AMI {
			log.Fatalf("AMI hypervisor not supported for air-gapped builds")
		}
	}

	configPath := ""
	switch bo.Hypervisor {
	case builder.VSphere:
		configPath = vSphereConfigFile
		// Validate builder type for VSphere
		if builderType != builder.BuilderTypeIso && builderType != builder.BuilderTypeClone {
			return fmt.Errorf("Invalid builder type. Please choose iso or clone")
		}
		// Set the builder type in BuildOptions
		bo.BuilderType = builderType
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
			return fmt.Errorf("Error converting %s config file path to absolute path: %v", bo.Hypervisor, err)
		}
		config, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("Error reading %s config file: %v", bo.Hypervisor, err)
		}
		switch bo.Hypervisor {
		case builder.VSphere:
			if err = json.Unmarshal(config, &bo.VsphereConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				isoUrl := bo.VsphereConfig.IsoUrl
				if bo.BuilderType == builder.BuilderTypeClone {
					isoUrl = "Don't validate iso URL"
				}
				if err = validateRedhat(&bo.VsphereConfig.RhelConfig, isoUrl); err != nil {
					return err
				}
			}
			if bo.VsphereConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.VsphereConfig.IsoChecksum, bo.VsphereConfig.IsoChecksumType); err != nil {
					return err
				}
			}
			if err = validateRHSM(bo.Os, &bo.VsphereConfig.RhsmConfig); err != nil {
				return err
			}
			if bo.AirGapped {
				if err = validateAirGapped(&bo.VsphereConfig.AirGappedConfig,
					bo.VsphereConfig.ExtraRepos, bo.VsphereConfig.IsoUrl); err != nil {
					return err
				}
			}
		case builder.Baremetal:
			if err = json.Unmarshal(config, &bo.BaremetalConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				if err = validateRedhat(&bo.BaremetalConfig.RhelConfig, bo.BaremetalConfig.IsoUrl); err != nil {
					return err
				}
			}
			if bo.BaremetalConfig != nil && bo.BaremetalConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.BaremetalConfig.IsoChecksum, bo.BaremetalConfig.IsoChecksumType); err != nil {
					return err
				}
			}
			if bo.BaremetalConfig != nil && bo.BaremetalConfig.DiskSizeMb != "" {
				if _, err := strconv.Atoi(bo.BaremetalConfig.DiskSizeMb); err != nil {
					return fmt.Errorf("parsing disk_size in baremetal config: %w", err)
				}
			}
			if err = validateRHSM(bo.Os, &bo.BaremetalConfig.RhsmConfig); err != nil {
				return err
			}
			if bo.AirGapped {
				if err = validateAirGapped(&bo.BaremetalConfig.AirGappedConfig,
					bo.BaremetalConfig.ExtraRepos, bo.BaremetalConfig.IsoUrl); err != nil {
					return err
				}
			}
		case builder.Nutanix:
			if err = json.Unmarshal(config, &bo.NutanixConfig); err != nil {
				return err
			}

			if bo.Os == builder.RedHat {
				if err = validateRedhat(&bo.NutanixConfig.RhelConfig, "Don't check IsoUrl Param"); err != nil {
					return err
				}
			}

			if bo.NutanixConfig.NutanixUserName == "" || bo.NutanixConfig.NutanixPassword == "" {
				log.Fatalf("\"nutanix_username\" and \"nutanix_password\" are required fields in nutanix-config")
			}
			if bo.AirGapped {
				if err = validateAirGapped(&bo.NutanixConfig.AirGappedConfig,
					bo.NutanixConfig.ExtraRepos, bo.NutanixConfig.ImageName); err != nil {
					return err
				}
			}

			if bo.NutanixConfig.ImageSizeGb != "" {
				imageSizeGb, err := strconv.Atoi(bo.NutanixConfig.ImageSizeGb)
				if err != nil {
					return fmt.Errorf("invalid image size: %v", err)
				}

				if imageSizeGb < 0 {
					return fmt.Errorf("image size must be a positive integer")
				}
			}
			// TODO Validate other fields as well
		case builder.CloudStack:
			if err = json.Unmarshal(config, &bo.CloudstackConfig); err != nil {
				return err
			}
			if bo.Os == builder.RedHat {
				if err = validateRedhat(&bo.CloudstackConfig.RhelConfig, bo.CloudstackConfig.IsoUrl); err != nil {
					return err
				}
			}
			if bo.CloudstackConfig.IsoUrl != "" {
				if err = validateCustomIso(bo.CloudstackConfig.IsoChecksum, bo.CloudstackConfig.IsoChecksumType); err != nil {
					return err
				}
			}
			if err = validateRHSM(bo.Os, &bo.CloudstackConfig.RhsmConfig); err != nil {
				return err
			}
			if bo.AirGapped {
				if err = validateAirGapped(&bo.CloudstackConfig.AirGappedConfig,
					bo.CloudstackConfig.ExtraRepos, bo.CloudstackConfig.IsoUrl); err != nil {
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
				ManifestOutput:      builder.DefaultAMIManifestOutput,
				RootDeviceName:      builder.DefaultAMIRootDeviceName,
				VolumeSize:          builder.DefaultAMIVolumeSize,
				VolumeType:          builder.DefaultAMIVolumeType,
			}
			if err = json.Unmarshal(config, amiConfig); err != nil {
				return err
			}

			bo.AMIConfig = amiConfig
			bo.FilesConfig = &builder.AdditionalFilesConfig{
				FileVars: builder.FileVars{
					AdditionalFiles:     "true",
					AdditionalFilesList: builder.DefaultAMIAdditionalFiles,
				},
			}
		}
	}

	if additionalFilesConfigFile != "" {
		filesConfig, err := os.ReadFile(additionalFilesConfigFile)
		if err != nil {
			return fmt.Errorf("Error reading additional files config path: %v", err)
		}

		if err = json.Unmarshal(filesConfig, &bo.FilesConfig); err != nil {
			return err
		}

		bo.FilesConfig.ProcessAdditionalFiles()
		if bo.AMIConfig != nil && !builder.SameFilesProvided(bo.FilesConfig.AdditionalFilesList, builder.DefaultAMIAdditionalFiles) {
			bo.FilesConfig.AdditionalFilesList = append(bo.FilesConfig.AdditionalFilesList, builder.DefaultAMIAdditionalFiles...)
		}
	}

	return nil
}

func validateOSHypervisorCombinations(os, hypervisor string) error {
	if hypervisor == builder.CloudStack && os != builder.RedHat {
		return fmt.Errorf("Invalid OS type. Only redhat OS is supported for CloudStack")
	}

	if hypervisor == builder.AMI && os != builder.Ubuntu {
		return fmt.Errorf("Invalid OS type. Only ubuntu OS is supported for AMI")
	}

	return nil
}

func validateRHSM(os string, rhsmConfig *builder.RhsmConfig) error {
	if rhsmConfig.ServerHostname != "" {
		if os != builder.RedHat {
			return fmt.Errorf("RedHat Subscription Manager Config (RHSM) cannot be provided when OS is not RedHat")
		}

		if rhsmConfig.ServerReleaseVersion == "" {
			return fmt.Errorf("RHSM version required when satelite server hostname is set for RHSM")
		}

		if rhsmConfig.ActivationKey == "" || rhsmConfig.OrgId == "" {
			return fmt.Errorf("Activation key and Org ID are required to use RHSM with satellite")
		}
	}
	return nil
}

func validateRedhat(rhelConfig *builder.RhelConfig, isoUrl string) error {
	if rhelConfig.ServerHostname == "" {
		if rhelConfig.RhelUsername == "" || rhelConfig.RhelPassword == "" {
			return fmt.Errorf("\"rhel_username\" and \"rhel_password\" are required fields in config when os is redhat")
		}
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

func validateOSVersion(os, osVersion string) error {
	if os != builder.RedHat && os != builder.Ubuntu {
		return fmt.Errorf("%s is not a supported OS", os)
	}

	if os == builder.Ubuntu && !builder.SliceContains(builder.SupportedUbuntuVersions, osVersion) {
		return fmt.Errorf("%s is not a supported version of Ubuntu. Please select one of %s", osVersion, strings.Join(builder.SupportedUbuntuVersions, ","))
	}

	if os == builder.RedHat && !builder.SliceContains(builder.SupportedRedHatVersions, osVersion) {
		return fmt.Errorf("%s is not a supported version of Redhat. Please select one of %s", osVersion, strings.Join(builder.SupportedRedHatVersions, ","))
	}

	return nil
}

func validateFirmware(firmware, os, osVersion, hypervisor string) error {
	if firmware == "" {
		return nil
	}

	if !builder.SliceContains(builder.SupportedFirmwares, firmware) {
		return fmt.Errorf("%s is not a firmware. Please select one of %s", firmware, strings.Join(builder.SupportedFirmwares, ","))
	}

	if firmware == builder.EFI {
		if os == builder.Ubuntu && !builder.SliceContains([]string{builder.VSphere, builder.Baremetal}, hypervisor) {
			return fmt.Errorf("For Ubuntu, EFI firmware is only supported for OVA and Raw builds")
		}
		if os == builder.RedHat {
			if !builder.SliceContains([]string{builder.VSphere, builder.Baremetal}, hypervisor) {
				return fmt.Errorf("For RedHat, EFI firmware is only supported for OVA and Raw builds")
			}
			if !builder.SliceContains(builder.SupportedRedHatEfiVersions, osVersion) {
				return fmt.Errorf("Only RedHat version 9 supports EFI firmware")
			}
		}
		if !builder.SliceContains([]string{builder.RedHat, builder.Ubuntu}, os) {
			return fmt.Errorf("EFI firmware is only support for Ubuntu and Redhat OS")
		}
	}

	if firmware == builder.BIOS && hypervisor == builder.Baremetal {
		if os == builder.Ubuntu {
			return fmt.Errorf("Ubuntu Raw builds only support EFI firmware")
		}
		if os == builder.RedHat && builder.SliceContains(builder.SupportedRedHatEfiVersions, osVersion) {
			return fmt.Errorf("RedHat version 9 Raw builds only support EFI firmware")
		}
	}

	return nil
}

func validateAirGapped(airgappedConfig *builder.AirGappedConfig, extraRepos, isoUrl string) error {
	if airgappedConfig.EksABuildToolingRepoUrl == "" {
		return fmt.Errorf("eksa_build_tooling_repo_url must be set when using air-gapped mode")
	}
	if airgappedConfig.ImageBuilderRepoUrl == "" {
		return fmt.Errorf("image_builder_repo_url must be set when using air-gapped mode")
	}
	if extraRepos == "" {
		return fmt.Errorf("Please set extra_repos to internal os package repo when using air-gapped mode")
	}
	if airgappedConfig.PrivateServerEksDDomainUrl == "" {
		return fmt.Errorf("Please set private_artifacts_eksd_fqdn to internal artifacts server's eks-d endpoint")
	}
	if airgappedConfig.PrivateServerEksADomainUrl == "" {
		return fmt.Errorf("Please set private_artifacts_eksa_fqdn to internal artifacts server's eks-a endpoint")
	}
	if isoUrl == "" {
		return fmt.Errorf("Please provide iso_url when building in air-gapped mode")
	}
	return nil
}
