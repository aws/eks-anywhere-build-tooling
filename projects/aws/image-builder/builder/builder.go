package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var buildToolingTag string

const buildToolingRepoUrl = "https://github.com/aws/eks-anywhere-build-tooling.git"

func (b *BuildOptions) BuildImage() {
	// Clone build tooling repo
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error retrieving current working directory: %v", err)
	}
	buildToolingRepoPath := filepath.Join(cwd, "eks-anywhere-build-tooling")
	imageBuilderProjectPath := filepath.Join(buildToolingRepoPath, "projects/kubernetes-sigs/image-builder")
	upstreamImageBuilderProjectPath := filepath.Join(imageBuilderProjectPath, "image-builder/images/capi")
	var outputArtifactPath string
	var outputImageGlob []string

	if b.Force {
		// Clean up build tooling repo in cwd
		cleanup(buildToolingRepoPath)
	}
	err = cloneRepo(buildToolingRepoUrl, buildToolingRepoPath, buildToolingTag)
	if err != nil {
		log.Fatalf("Error clonning build tooling repo")
	}
	log.Println("Cloned eks-anywhere-build-tooling repo")

	log.Printf("Initiating Image Build\n Image OS: %s\n Hypervisor: %s\n", b.Os, b.Hypervisor)
	if b.Hypervisor == VSphere {
		// Read and set the vsphere connection data
		vsphereConfigData, err := os.ReadFile(b.VsphereConfig)
		if err != nil {
			log.Fatalf("Error reading the vsphere config file")
		}
		err = ioutil.WriteFile(filepath.Join(imageBuilderProjectPath, "packer/ova/vsphere.json"), vsphereConfigData, 0644)
		if err != nil {
			log.Fatalf("Error writing vsphere config file to packer")
		}
		buildCommand := fmt.Sprintf("make -C %s local-build-ova-ubuntu-2004", imageBuilderProjectPath)
		err = executeMakeBuildCommand(buildCommand, b.ReleaseChannel, b.ArtifactsBucket)
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
			if err = executeMakeBuildCommand(upstreamPatchCommand, b.ReleaseChannel, b.ArtifactsBucket); err != nil {
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
		err = executeMakeBuildCommand(buildCommand, b.ReleaseChannel, b.ArtifactsBucket)
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

	cleanup(buildToolingRepoPath)
	log.Print("Build Successful. Output artifacts located at current working directory\n")
}

func (b *BuildOptions) ValidateInputs() {
	b.Os = strings.ToLower(b.Os)
	if b.Os != Ubuntu {
		log.Fatalf("Invalid OS type. Please choose ubuntu")
	}

	b.Hypervisor = strings.ToLower(b.Hypervisor)
	if (b.Hypervisor != VSphere) && (b.Hypervisor != Baremetal) {
		log.Fatalf("Invalid hypervisor. Please choose vsphere or baremetal")
	}

	if b.Hypervisor == VSphere && b.VsphereConfig == "" {
		log.Fatalf("vsphere-config is a required flag for vsphere hypervisor")
	}

	var err error
	b.VsphereConfig, err = filepath.Abs(b.VsphereConfig)
	if err != nil {
		log.Fatalf("Error converting vsphere config path to absolute")
	}

	if b.ReleaseChannel != "1-20" && b.ReleaseChannel != "1-21" && b.ReleaseChannel != "1-22" {
		log.Fatalf("release-channel should be one of 1-20, 1-21, 1-22")
	}
}
