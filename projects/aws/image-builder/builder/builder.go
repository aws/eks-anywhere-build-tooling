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

var (
	codebuild = os.Getenv("CODEBUILD_CI")
)

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
			log.Fatalf("Error clonning build tooling repo")
		}
		log.Println("Cloned eks-anywhere-build-tooling repo")
	} else {
		buildToolingRepoPath = os.Getenv("CODEBUILD_SRC_DIR")
		log.Println("Using repo checked out from code commit")
	}

	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		cleanup(buildToolingRepoPath)
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
		err = ioutil.WriteFile(filepath.Join(imageBuilderProjectPath, "packer/ova/vsphere.json"), vsphereConfigData, 0644)
		if err != nil {
			log.Fatalf("Error writing vsphere config file to packer")
		}

		var buildCommand string
		switch b.Os {
		case Ubuntu:
			buildCommand = fmt.Sprintf("make -C %s local-build-ova-ubuntu-2004", imageBuilderProjectPath)
		case RedHat:
			buildCommand = fmt.Sprintf("make -C %s local-build-ova-rhel-7", imageBuilderProjectPath)
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
		if b.Os == Ubuntu {
			// Patch firmware config for tool
			upstreamPatchCommand := fmt.Sprintf("make -C %s image-builder/eks-anywhere-patched", imageBuilderProjectPath)
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
			if err := os.WriteFile(ubuntuEfiConfigPath, []byte(ubuntuPatchedEfiConfig), 0644); err != nil {
				log.Fatalf("Error writing the new ubuntu efi config file: %v", err)
			}
			log.Println("Patched upstream firmware config file")
		}
		buildCommand := fmt.Sprintf("make -C %s local-build-raw-ubuntu-2004-efi", imageBuilderProjectPath)
		err = executeMakeBuildCommand(buildCommand, commandEnvVars...)
		if err != nil {
			log.Fatalf("Error executing image-builder for raw hypervisor: %v", err)
		}

		outputImageGlob, err = filepath.Glob(filepath.Join(upstreamImageBuilderProjectPath, "output/*.gz"))
		if err != nil {
			log.Fatalf("Error getting glob for output files: %v", err)
		}

		outputArtifactPath = filepath.Join(cwd, fmt.Sprintf("%s.gz", b.Os))
	}

	// Moving artifacts from upstream directory to cwd
	log.Println("Moving artifacts from build directory to current working directory")
	err = os.Rename(outputImageGlob[0], outputArtifactPath)
	if err != nil {
		log.Fatalf("Error moving output file to current working directory")
	}

	if codebuild != "true" {
		cleanup(buildToolingRepoPath)
	}

	log.Print("Build Successful. Output artifacts located at current working directory\n")
}

func (b *BuildOptions) ValidateInputs() {
	b.Os = strings.ToLower(b.Os)
	if b.Os != Ubuntu && b.Os != RedHat {
		log.Fatalf("Invalid OS type. Please choose ubuntu or redhat")
	}

	b.Hypervisor = strings.ToLower(b.Hypervisor)
	if (b.Hypervisor != VSphere) && (b.Hypervisor != Baremetal) {
		log.Fatalf("Invalid hypervisor. Please choose vsphere or baremetal")
	}

	if b.Hypervisor == Baremetal && b.Os == RedHat {
		log.Fatalf("Redhat is not supported with baremetal hypervisor. Please choose vsphere to build Redhat")
	}

	// Validate vsphere config inputs
	if b.VsphereConfig != nil {
		// Validate Rhel username and password
		if b.Os == RedHat {
			if b.VsphereConfig.RhelUsername == "" || b.VsphereConfig.RhelPassword == "" {
				log.Fatalf("\"rhel_username\" and \"rhel_password\" are required fields in vsphere-config when os is redhat")
			}

			if b.VsphereConfig.IsoUrl == "" {
				log.Fatalf("\"iso_url\" is a required field in vsphere-config when os is redhat")
			}
		}

		// Validate iso checksum and checksum type was provided
		if b.VsphereConfig.IsoUrl != "" {
			if b.VsphereConfig.IsoChecksum == "" {
				log.Fatalf("Please provide a valid checksum for \"iso_checksum\" when providing \"iso_url\"")
			}

			if b.VsphereConfig.IsoChecksumType != "sha256" && b.VsphereConfig.IsoChecksumType != "sha512" {
				log.Fatalf("\"iso_checksum_type\" is a required field when providing iso_checksum. Checksum type can be sha256 or sha512")
			}
		}
	}
}
