package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	buildToolingRepoUrl = "https://github.com/aws/eks-anywhere-build-tooling.git"
)

var codebuild = os.Getenv("CODEBUILD_CI")

func (b *BuildOptions) BuildImage() {
	// Clone build tooling repo
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error retrieving current working directory: %v", err)
	}
	buildToolingRepoPath := getBuildToolingPath(cwd)

	if b.Force && codebuild != "true" {
		// Clean up build tooling repo in cwd
		cleanup(buildToolingRepoPath)
	}

	if codebuild != "true" {
		err = cloneRepo(buildToolingRepoUrl, buildToolingRepoPath)
		if err != nil {
			log.Fatalf("Error cloning build tooling repo")
		}
		log.Println("Cloned eks-anywhere-build-tooling repo")

		gitCommitFromBundle, err := getGitCommitFromBundle(buildToolingRepoPath)
		if err != nil {
			log.Fatalf("Error getting git commit from bundle: %v", err)
		}

		err = checkoutRepo(buildToolingRepoPath, gitCommitFromBundle)
		if err != nil {
			log.Fatalf("Error checking out build tooling repo at commit %s", gitCommitFromBundle)
		}
		log.Printf("Checked out eks-anywhere-build-tooling repo at commit %s\n", gitCommitFromBundle)
	} else {
		buildToolingRepoPath = os.Getenv("CODEBUILD_SRC_DIR")
		log.Println("Using repo checked out from code commit")
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		if codebuild != "true" {
			cleanup(buildToolingRepoPath)
		}
		log.Fatalf("release-channel should be one of %v", supportedReleaseBranches)
	}

	imageBuilderProjectPath := filepath.Join(buildToolingRepoPath, "projects/kubernetes-sigs/image-builder")
	upstreamImageBuilderProjectPath := filepath.Join(imageBuilderProjectPath, "image-builder/images/capi")
	var outputArtifactPath string
	var outputImageGlob []string
	commandEnvVars := []string{fmt.Sprintf("RELEASE_BRANCH=%s", b.ReleaseChannel)}

	log.Printf("Initiating Image Build\n Image OS: %s\n Hypervisor: %s\n", b.Os, b.Hypervisor)
	if b.Hypervisor == VSphere {
		// Read and set the vsphere connection data
		vsphereConfigData, err := json.Marshal(b.VsphereConfig)
		if err != nil {
			log.Fatalf("Error marshalling vsphere config data")
		}
		err = ioutil.WriteFile(filepath.Join(imageBuilderProjectPath, "packer/ova/vsphere.json"), vsphereConfigData, 0o644)
		if err != nil {
			log.Fatalf("Error writing vsphere config file to packer")
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			buildCommand = fmt.Sprintf("make -C %s local-build-ova-ubuntu-2004", imageBuilderProjectPath)
		case RedHat:
			buildCommand = fmt.Sprintf("make -C %s local-build-ova-rhel-8", imageBuilderProjectPath)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("RHSM_USERNAME=%s", b.VsphereConfig.RhelUsername),
				fmt.Sprintf("RHSM_PASSWORD=%s", b.VsphereConfig.RhelPassword),
			)
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for vsphere hypervisor: %v", err)
		}

		// Move the output ova to cwd
		outputImageGlob, err = filepath.Glob(filepath.Join(upstreamImageBuilderProjectPath, "output/*.ova"))
		if err != nil {
			log.Fatalf("Error getting glob for output files: %v", err)
		}
		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s.ova", b.Os))

		log.Printf("Image Build Successful\n Please find the output artifact at %s\n", outputArtifactPath)
	} else if b.Hypervisor == Baremetal {
		baremetalConfigFile := filepath.Join(imageBuilderProjectPath, "packer/config/baremetal.json")
		if b.BaremetalConfig != nil {
			baremetalConfigData, err := json.Marshal(b.BaremetalConfig)
			if err != nil {
				log.Fatalf("Error marshalling baremetal config data")
			}
			err = ioutil.WriteFile(baremetalConfigFile, baremetalConfigData, 0o644)
			if err != nil {
				log.Fatalf("Error writing baremetal config file to packer")
			}
		}

		if b.Os == Ubuntu {
			// Patch firmware config for tool
			upstreamPatchCommand := fmt.Sprintf("make -C %s patch-repo", imageBuilderProjectPath)
			if err = executeMakeBuildCommand(upstreamPatchCommand, commandEnvVars...); err != nil {
				log.Fatalf("Error executing upstream patch command")
			}

			ubuntuEfiConfigPath := filepath.Join(upstreamImageBuilderProjectPath, "packer/raw/raw-ubuntu-2004-efi.json")
			ubuntuEfiConfigFileData, err := os.ReadFile(ubuntuEfiConfigPath)
			if err != nil {
				log.Fatalf("Error reading ubuntu efi config file: %v", err)
			}
			ubuntuEfiConfigFileString := string(ubuntuEfiConfigFileData)
			// This comes from our patch for AL2 on CodeBuild Image Builder
			ubuntuPatchedEfiConfig := strings.ReplaceAll(ubuntuEfiConfigFileString, "/usr/share/edk2/ovmf/OVMF_CODE.fd", "OVMF.fd")
			if err := os.Remove(ubuntuEfiConfigPath); err != nil {
				log.Fatalf("Error removing the old ubuntu efi config file: %v", err)
			}
			if err := os.WriteFile(ubuntuEfiConfigPath, []byte(ubuntuPatchedEfiConfig), 0o644); err != nil {
				log.Fatalf("Error writing the new ubuntu efi config file: %v", err)
			}
			log.Println("Patched upstream firmware config file")
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			buildCommand = fmt.Sprintf("make -C %s local-build-raw-ubuntu-2004-efi", imageBuilderProjectPath)
		case RedHat:
			buildCommand = fmt.Sprintf("make -C %s local-build-raw-rhel-8", imageBuilderProjectPath)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("RHSM_USERNAME=%s", b.BaremetalConfig.RhelUsername),
				fmt.Sprintf("RHSM_PASSWORD=%s", b.BaremetalConfig.RhelPassword),
			)
		}
		if b.BaremetalConfig != nil {
			commandEnvVars = append(commandEnvVars, fmt.Sprintf("PACKER_TYPE_VAR_FILES=%s", baremetalConfigFile))
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for raw hypervisor: %v", err)
		}

		outputImageGlob, err = filepath.Glob(filepath.Join(upstreamImageBuilderProjectPath, "output/*.gz"))
		if err != nil {
			log.Fatalf("Error getting glob for output files: %v", err)
		}

		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s.gz", b.Os))
	} else if b.Hypervisor == Nutanix {
		// Patch firmware config for tool
		upstreamPatchCommand := fmt.Sprintf("make -C %s patch-repo", imageBuilderProjectPath)
		if err = executeMakeBuildCommand(upstreamPatchCommand, commandEnvVars...); err != nil {
			log.Fatalf("Error executing upstream patch command")
		}

		// Read and set the nutanix connection data
		nutanixConfigData, err := json.Marshal(b.NutanixConfig)
		if err != nil {
			log.Fatalf("Error marshalling nutanix config data")
		}
		err = ioutil.WriteFile(filepath.Join(upstreamImageBuilderProjectPath, "packer/nutanix/nutanix.json"), nutanixConfigData, 0o644)
		if err != nil {
			log.Fatalf("Error writing nutanix config file to packer: %v", err)
		}

		buildCommand := fmt.Sprintf("make -C %s local-build-nutanix-ubuntu-2004", imageBuilderProjectPath)
		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for nutanix hypervisor: %v", err)
		}

		log.Printf("Image Build Successful\n Please find the image uploaded under Nutanix Image Service with name %s\n", b.NutanixConfig.ImageName)
	} else if b.Hypervisor == CloudStack {
		// Create config file
		cloudstackConfigFile := filepath.Join(imageBuilderProjectPath, "packer/config/cloudstack.json")

		// Assign ansible user var for cloudstack provider
		b.CloudstackConfig.AnsibleUserVars = "provider=cloudstack"
		if b.CloudstackConfig != nil {
			cloudstackConfigData, err := json.Marshal(b.CloudstackConfig)
			if err != nil {
				log.Fatalf("Error marshalling cloudstack config data")
			}
			err = ioutil.WriteFile(cloudstackConfigFile, cloudstackConfigData, 0o644)
			if err != nil {
				log.Fatalf("Error writing cloudstack config file to packer")
			}
		}

		var buildCommand string
		var outputImageGlobPattern string
		switch b.Os {
		case RedHat:
			outputImageGlobPattern = "output/rhel-*/rhel-*"
			buildCommand = fmt.Sprintf("make -C %s local-build-qemu-rhel-8", imageBuilderProjectPath)
			commandEnvVars = append(commandEnvVars,
				fmt.Sprintf("RHSM_USERNAME=%s", b.CloudstackConfig.RhelUsername),
				fmt.Sprintf("RHSM_PASSWORD=%s", b.CloudstackConfig.RhelPassword),
			)
		}
		if b.CloudstackConfig != nil {
			commandEnvVars = append(commandEnvVars, fmt.Sprintf("PACKER_TYPE_VAR_FILES=%s", cloudstackConfigFile))
		}

		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for raw hypervisor: %v", err)
		}

		outputImageGlob, err = filepath.Glob(filepath.Join(upstreamImageBuilderProjectPath, outputImageGlobPattern))
		if err != nil {
			log.Fatalf("Error getting glob for output files: %v", err)
		}

		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s.qcow2", b.Os))
	} else if b.Hypervisor == AMI {
		amiConfigFile := filepath.Join(imageBuilderProjectPath, "packer/ami/ami.json")

		upstreamPatchCommand := fmt.Sprintf("make -C %s patch-repo", imageBuilderProjectPath)
		if err = executeMakeBuildCommand(upstreamPatchCommand, commandEnvVars...); err != nil {
			log.Fatalf("Error executing upstream patch command")
		}

		if b.AMIConfig != nil {
			amiConfigData, err := json.Marshal(b.AMIConfig)
			if err != nil {
				log.Fatalf("Error marshalling AMI config data")
			}

			err = ioutil.WriteFile(amiConfigFile, amiConfigData, 0o644)
			if err != nil {
				log.Fatalf("Error writing AMI config file to packer")
			}
		}

		buildCommand := fmt.Sprintf("make -C %s local-build-ami-ubuntu-2004", imageBuilderProjectPath)
		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for AMI hypervisor: %v", err)
		}
	}

	if outputArtifactPath != "" {
		// Moving artifacts from upstream directory to cwd
		log.Println("Moving artifacts from build directory to current working directory")
		err = os.Rename(outputImageGlob[0], outputArtifactPath)
		if err != nil {
			log.Fatalf("Error moving output file to current working directory")
		}
	}

	if codebuild != "true" {
		cleanup(buildToolingRepoPath)
	}

	log.Print("Build Successful. Output artifacts located at current working directory\n")
}
