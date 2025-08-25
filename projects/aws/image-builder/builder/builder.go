package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
)

var (
	eksaVersion         string
	eksaReleaseManifest string
	codebuild           = os.Getenv(codebuildCIEnvVar)
)

func (b *BuildOptions) BuildImage() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error retrieving current working directory: %v", err)
	}
	if b.AirGapped {
		var eksDArtifactsDomain, eksAArtifactsDomain string
		switch b.Hypervisor {
		case VSphere:
			eksDArtifactsDomain = b.VsphereConfig.PrivateServerEksDDomainUrl
			eksAArtifactsDomain = b.VsphereConfig.PrivateServerEksADomainUrl
		case Baremetal:
			eksDArtifactsDomain = b.BaremetalConfig.PrivateServerEksDDomainUrl
			eksAArtifactsDomain = b.BaremetalConfig.PrivateServerEksADomainUrl
		case CloudStack:
			eksDArtifactsDomain = b.CloudstackConfig.PrivateServerEksDDomainUrl
			eksAArtifactsDomain = b.CloudstackConfig.PrivateServerEksADomainUrl
		case Nutanix:
			eksDArtifactsDomain = b.NutanixConfig.PrivateServerEksDDomainUrl
			eksAArtifactsDomain = b.NutanixConfig.PrivateServerEksADomainUrl
		}
		if err = extractAndPrepManifestTarball(b.ManifestTarball, eksDArtifactsDomain, eksAArtifactsDomain); err != nil {
			log.Fatalf(err.Error())
		}
	}
	buildToolingRepoPath := getBuildToolingPath(cwd)
	_, detectedEksaVersion, err := b.prepBuildToolingRepo(buildToolingRepoPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		cleanup(buildToolingRepoPath)
		log.Fatalf("release-channel should be one of %v", supportedReleaseBranches)
	}

	imageBuilderProjectPath := filepath.Join(buildToolingRepoPath, imageBuilderProjectDirectory)
	upstreamImageBuilderProjectPath := filepath.Join(imageBuilderProjectPath, imageBuilderCAPIDirectory)
	var outputArtifactPath string
	var outputImageGlob []string
	eksAReleaseManifestUrl, err := getEksAReleasesManifestURL(b.AirGapped)
	if err != nil {
		log.Fatalf(err.Error())
	}
	commandEnvVars := []string{
		fmt.Sprintf("%s=%s", releaseBranchEnvVar, b.ReleaseChannel),
		fmt.Sprintf("%s=%s", eksAReleaseVersionEnvVar, detectedEksaVersion),
		fmt.Sprintf("%s=%s", eksAReleaseManifestURLEnvVar, eksAReleaseManifestUrl),
	}
	if b.AnsibleVerbosity != 0 {
		ansibleVerbosityArg := fmt.Sprintf("-%s", strings.Repeat("v", b.AnsibleVerbosity))
		commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", eksaAnsibleVerbosityEnvVar, ansibleVerbosityArg))
	}

	log.Printf("Initiating Image Build\n Image OS: %s\n Image OS Version: %s\n Hypervisor: %s\n Firmware: %s\n", b.Os, b.OsVersion, b.Hypervisor, b.Firmware)
	if b.FilesConfig != nil {
		additionalFilesList := filepath.Join(imageBuilderProjectPath, packerAdditionalFilesList)
		additionalFilesCustomRole := filepath.Join(imageBuilderProjectPath, ansibleAdditionalFilesCustomRole)
		b.FilesConfig.FilesAnsibleConfig.CustomRoleNames = additionalFilesCustomRole
		b.FilesConfig.FilesAnsibleConfig.AnsibleExtraVars = fmt.Sprintf("@%s", additionalFilesList)

		log.Println("Marshalling files ansible config to JSON")
		filesAnsibleConfig, err := json.Marshal(b.FilesConfig.FilesAnsibleConfig)
		if err != nil {
			log.Fatalf("Error marshalling files ansible config data: %v", err)
		}

		additionalFilesConfigFile := filepath.Join(imageBuilderProjectPath, packerAdditionalFilesConfigFile)
		log.Printf("Writing files ansible config to Packer config directory: %s", additionalFilesConfigFile)
		err = ioutil.WriteFile(additionalFilesConfigFile, filesAnsibleConfig, 0o644)
		if err != nil {
			log.Fatalf("Error writing additional files config file to Packer config directory: %v", err)
		}

		log.Println("Marshalling additional files list to YAML")
		if b.AMIConfig != nil {
			for index, file := range b.FilesConfig.AdditionalFilesList {
				for _, defaultFile := range DefaultAMIAdditionalFiles {
					if file == defaultFile {
						b.FilesConfig.AdditionalFilesList[index].Source = filepath.Join(imageBuilderProjectPath, b.FilesConfig.AdditionalFilesList[index].Source)
					}
				}
			}
		}
		additionalFileVars, err := yaml.Marshal(b.FilesConfig.FileVars)
		if err != nil {
			log.Fatalf("Error marshalling additional files list: %v", err)
		}

		log.Printf("Writing additional files list to Packer ansible directory: %s", additionalFilesList)
		err = ioutil.WriteFile(additionalFilesList, additionalFileVars, 0o644)
		if err != nil {
			log.Fatalf("Error writing additional files list to Packer ansible directory: %v", err)
		}

		commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", packerAdditionalFilesConfigFileEnvVar, additionalFilesConfigFile))
	}
	var outputImageGlobPattern string
	if b.Hypervisor == VSphere {
		if b.BuilderType == BuilderTypeClone && b.VsphereConfig.Template == "" {
			log.Fatalf("Error: When using vsphere-clone, template cannot be empty")
		}
		if b.VsphereConfig.VmxVersion != "" {
			vmxVersion, err := strconv.Atoi(b.VsphereConfig.VmxVersion)
			if err != nil {
				log.Fatalf("Error parsing vmx_version in vsphere config: %v", err)
			}
			if vmxVersion < minVmVersion {
				log.Fatalf("vmx_version cannot be less than %d, have %d in vsphere config", minVmVersion, vmxVersion)
			}
		}
		if b.AirGapped {
			airGapEnvVars, err := getAirGapCmdEnvVars(b.VsphereConfig.ImageBuilderRepoUrl, detectedEksaVersion, b.ReleaseChannel)
			if err != nil {
				log.Fatalf("Error getting air gapped env variables: %v", err)
			}
			commandEnvVars = append(commandEnvVars, airGapEnvVars...)
		}

		// Set proxy on RHSM if available
		if b.Os == RedHat && b.VsphereConfig.HttpProxy != "" {
			if err := setRhsmProxy(&b.VsphereConfig.ProxyConfig, &b.VsphereConfig.RhsmConfig); err != nil {
				log.Fatalf("Error parsing proxy host and port for RHSM: %v", err)
			}
		}

		// Read and set the vsphere connection data
		vsphereConfigData, err := json.Marshal(b.VsphereConfig)
		if err != nil {
			log.Fatalf("Error marshalling vsphere config data: %v", err)
		}
		err = ioutil.WriteFile(filepath.Join(imageBuilderProjectPath, packerVSphereConfigFile), vsphereConfigData, 0o644)
		if err != nil {
			log.Fatalf("Error writing vsphere config file to packer: %v", err)
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			if b.BuilderType == BuilderTypeClone {
				if b.Firmware == EFI {
					buildCommand = fmt.Sprintf("make -C %s local-clone-build-ova-ubuntu-%s-efi", imageBuilderProjectPath, b.OsVersion)
				} else {
					buildCommand = fmt.Sprintf("make -C %s local-clone-build-ova-ubuntu-%s", imageBuilderProjectPath, b.OsVersion)
				}
			} else {
				if b.Firmware == EFI {
					buildCommand = fmt.Sprintf("make -C %s local-build-ova-ubuntu-%s-efi", imageBuilderProjectPath, b.OsVersion)
				} else {
					buildCommand = fmt.Sprintf("make -C %s local-build-ova-ubuntu-%s", imageBuilderProjectPath, b.OsVersion)
				}
			}
		case RedHat:
			if b.BuilderType == BuilderTypeClone {
				if b.Firmware == EFI {
					buildCommand = fmt.Sprintf("make -C %s local-clone-build-ova-redhat-%s-efi", imageBuilderProjectPath, b.OsVersion)
				} else {
					buildCommand = fmt.Sprintf("make -C %s local-clone-build-ova-redhat-%s", imageBuilderProjectPath, b.OsVersion)
				}
			} else {
				if b.Firmware == EFI {
					buildCommand = fmt.Sprintf("make -C %s local-build-ova-redhat-%s-efi", imageBuilderProjectPath, b.OsVersion)
				} else {
					buildCommand = fmt.Sprintf("make -C %s local-build-ova-redhat-%s", imageBuilderProjectPath, b.OsVersion)
				}
			}
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("%s=%s", rhelUsernameEnvVar, b.VsphereConfig.RhelUsername),
				fmt.Sprintf("%s=%s", rhelPasswordEnvVar, b.VsphereConfig.RhelPassword),
				fmt.Sprintf("%s=%s", rhsmActivationKeyEnvVar, b.VsphereConfig.ActivationKey),
				fmt.Sprintf("%s=%s", rhsmOrgIDEnvVar, b.VsphereConfig.OrgId),
			)
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for vsphere hypervisor: %v", err)
		}

		outputImageGlobPattern = "output/*.ova"
		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s-%s-kube-%s.ova", b.Os, b.OsVersion, b.ReleaseChannel))

		log.Printf("Image Build Successful\n Please find the output artifact at %s\n", outputArtifactPath)
	} else if b.Hypervisor == Baremetal {
		if b.AirGapped {
			airGapEnvVars, err := getAirGapCmdEnvVars(b.BaremetalConfig.ImageBuilderRepoUrl, detectedEksaVersion, b.ReleaseChannel)
			if err != nil {
				log.Fatalf("Error getting air gapped env variables: %v", err)
			}
			commandEnvVars = append(commandEnvVars, airGapEnvVars...)
		}

		// Set proxy on RHSM if available
		if b.Os == RedHat && b.BaremetalConfig.HttpProxy != "" {
			if err := setRhsmProxy(&b.BaremetalConfig.ProxyConfig, &b.BaremetalConfig.RhsmConfig); err != nil {
				log.Fatalf("Error parsing proxy host and port for RHSM: %v", err)
			}
		}

		baremetalConfigFile := filepath.Join(imageBuilderProjectPath, packerBaremetalConfigFile)
		if b.BaremetalConfig == nil {
			b.BaremetalConfig = &BaremetalConfig{}
		}
		if b.BaremetalConfig.DiskSizeMb == "" {
			b.BaremetalConfig.DiskSizeMb = DefaultBaremetalDiskSizeMb
			log.Printf("Using default disk size: %s MB", b.BaremetalConfig.DiskSizeMb)
		} else {
			log.Printf("Using configured disk size: %s MB", b.BaremetalConfig.DiskSizeMb)
		}
		baremetalConfigData, err := json.Marshal(b.BaremetalConfig)
		if err != nil {
			log.Fatalf("Error marshalling baremetal config data: %v", err)
		}
		err = ioutil.WriteFile(baremetalConfigFile, baremetalConfigData, 0o644)
		if err != nil {
			log.Fatalf("Error writing baremetal config file to packer: %v", err)
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			buildCommand = fmt.Sprintf("make -C %s local-build-raw-ubuntu-%s", imageBuilderProjectPath, b.OsVersion)
		case RedHat:
			buildCommand = fmt.Sprintf("make -C %s local-build-raw-redhat-%s", imageBuilderProjectPath, b.OsVersion)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("%s=%s", rhelUsernameEnvVar, b.BaremetalConfig.RhelUsername),
				fmt.Sprintf("%s=%s", rhelPasswordEnvVar, b.BaremetalConfig.RhelPassword),
				fmt.Sprintf("%s=%s", rhsmActivationKeyEnvVar, b.BaremetalConfig.ActivationKey),
				fmt.Sprintf("%s=%s", rhsmOrgIDEnvVar, b.BaremetalConfig.OrgId),
			)
		}
		if b.BaremetalConfig != nil {
			commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", packerTypeVarFilesEnvVar, baremetalConfigFile))
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for raw hypervisor: %v", err)
		}

		outputImageGlobPattern = "output/*.gz"
		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s-%s-kube-%s.gz", b.Os, b.OsVersion, b.ReleaseChannel))

		log.Printf("Image Build Successful\n Please find the output artifact at %s\n", outputArtifactPath)
	} else if b.Hypervisor == Nutanix {
		// Set proxy on RHSM if available
		if b.Os == RedHat && b.NutanixConfig.HttpProxy != "" {
			if err := setRhsmProxy(&b.NutanixConfig.ProxyConfig, &b.NutanixConfig.RhsmConfig); err != nil {
				log.Fatalf("Error parsing proxy host and port for RHSM: %v", err)
			}
		}

		if b.AirGapped {
			airGapEnvVars, err := getAirGapCmdEnvVars(b.NutanixConfig.ImageBuilderRepoUrl, detectedEksaVersion, b.ReleaseChannel)
			if err != nil {
				log.Fatalf("Error getting air gapped env variables: %v", err)
			}
			commandEnvVars = append(commandEnvVars, airGapEnvVars...)
		}

		// Create config file
		nutanixConfigFile := filepath.Join(imageBuilderProjectPath, packerNutanixConfigFile)

		// Read and set the nutanix connection data
		nutanixConfigData, err := json.Marshal(b.NutanixConfig)
		if err != nil {
			log.Fatalf("Error marshalling nutanix config data: %v", err)
		}
		err = ioutil.WriteFile(nutanixConfigFile, nutanixConfigData, 0o644)
		if err != nil {
			log.Fatalf("Error writing nutanix config file to packer: %v", err)
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			buildCommand = fmt.Sprintf("make -C %s local-build-nutanix-ubuntu-%s", imageBuilderProjectPath, b.OsVersion)
		case RedHat:
			buildCommand = fmt.Sprintf("make -C %s local-build-nutanix-redhat-%s", imageBuilderProjectPath, b.OsVersion)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("%s=%s", rhelUsernameEnvVar, b.NutanixConfig.RhelUsername),
				fmt.Sprintf("%s=%s", rhelPasswordEnvVar, b.NutanixConfig.RhelPassword),
				fmt.Sprintf("%s=%s", rhelImageUrlNutanixEnvVar, b.NutanixConfig.ImageUrl),
			)
		}

		if b.NutanixConfig != nil {
			commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", packerTypeVarFilesEnvVar, nutanixConfigFile))

			if b.NutanixConfig.ImageSizeGb == "" {
				// Set default image size for Linux to 10GB as it implemented in image-builder upstream
				b.NutanixConfig.ImageSizeGb = "10"
			}

			commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", imageSizeGbNutanixEnvVar, b.NutanixConfig.ImageSizeGb))
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for nutanix hypervisor: %v", err)
		}

		log.Printf("Image Build Successful\n Please find the image uploaded under Nutanix Image Service with name %s\n", b.NutanixConfig.ImageName)
		if b.NutanixConfig.ImageExport == "true" {
			outputImageGlobPattern = "*.img"
			outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s-%s-kube-%s.img", b.Os, b.OsVersion, b.ReleaseChannel))
			log.Printf("Also please find the exported image at %s\n", outputArtifactPath)
		}
	} else if b.Hypervisor == CloudStack {
		if b.AirGapped {
			airGapEnvVars, err := getAirGapCmdEnvVars(b.CloudstackConfig.ImageBuilderRepoUrl, detectedEksaVersion, b.ReleaseChannel)
			if err != nil {
				log.Fatalf("Error getting air gapped env variables: %v", err)
			}
			commandEnvVars = append(commandEnvVars, airGapEnvVars...)
		}
		// Set proxy on RHSM if available
		if b.Os == RedHat && b.CloudstackConfig.HttpProxy != "" {
			if err := setRhsmProxy(&b.CloudstackConfig.ProxyConfig, &b.CloudstackConfig.RhsmConfig); err != nil {
				log.Fatalf("Error parsing proxy host and port for RHSM: %v", err)
			}
		}

		// Create config file
		cloudstackConfigFile := filepath.Join(imageBuilderProjectPath, packerCloudStackConfigFile)

		// Assign ansible user var for cloudstack provider
		b.CloudstackConfig.AnsibleUserVars = "provider=cloudstack"
		if b.CloudstackConfig != nil {
			cloudstackConfigData, err := json.Marshal(b.CloudstackConfig)
			if err != nil {
				log.Fatalf("Error marshalling cloudstack config data: %v", err)
			}
			err = ioutil.WriteFile(cloudstackConfigFile, cloudstackConfigData, 0o644)
			if err != nil {
				log.Fatalf("Error writing cloudstack config file to packer: %v", err)
			}
		}

		var buildCommand string
		switch b.Os {
		case RedHat:
			outputImageGlobPattern = "output/rhel-*/rhel-*"
			buildCommand = fmt.Sprintf("make -C %s local-build-cloudstack-redhat-%s", imageBuilderProjectPath, b.OsVersion)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("%s=%s", rhelUsernameEnvVar, b.CloudstackConfig.RhelUsername),
				fmt.Sprintf("%s=%s", rhelPasswordEnvVar, b.CloudstackConfig.RhelPassword),
				fmt.Sprintf("%s=%s", rhsmActivationKeyEnvVar, b.CloudstackConfig.ActivationKey),
				fmt.Sprintf("%s=%s", rhsmOrgIDEnvVar, b.CloudstackConfig.OrgId),
			)
		}
		if b.CloudstackConfig != nil {
			commandEnvVars = append(commandEnvVars, fmt.Sprintf("%s=%s", packerTypeVarFilesEnvVar, cloudstackConfigFile))
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for raw hypervisor: %v", err)
		}

		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s-%s-kube-%s.qcow2", b.Os, b.OsVersion, b.ReleaseChannel))

		log.Printf("Image Build Successful\n Please find the output artifact at %s\n", outputArtifactPath)
	} else if b.Hypervisor == AMI {
		amiConfigFile := filepath.Join(imageBuilderProjectPath, packerAMIConfigFile)

		if b.AMIConfig != nil {
			if b.AMIConfig.ManifestOutput == DefaultAMIManifestOutput {
				b.AMIConfig.ManifestOutput = filepath.Join(cwd, b.AMIConfig.ManifestOutput)
			}
			amiConfigData, err := json.Marshal(b.AMIConfig)
			if err != nil {
				log.Fatalf("Error marshalling AMI config data: %v", err)
			}

			err = ioutil.WriteFile(amiConfigFile, amiConfigData, 0o644)
			if err != nil {
				log.Fatalf("Error writing AMI config file to packer: %v", err)
			}
		}

		buildCommand := fmt.Sprintf("make -C %s local-build-ami-ubuntu-%s", imageBuilderProjectPath, b.OsVersion)
		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for AMI hypervisor: %v", err)
		}
	}

	if outputArtifactPath != "" {
		outputImageGlob, err = filepath.Glob(filepath.Join(upstreamImageBuilderProjectPath, outputImageGlobPattern))
		if err != nil {
			log.Fatalf("Error getting glob for output files: %v", err)
		}

		// Moving artifacts from upstream directory to cwd
		log.Println("Moving artifacts from build directory to current working directory")
		err = os.Rename(outputImageGlob[0], outputArtifactPath)
		if err != nil {
			log.Fatalf("Error moving output file to current working directory: %v", err)
		}
	}

	cleanup(buildToolingRepoPath)

	log.Print("Build Successful. Output artifacts located at current working directory\n")
}
