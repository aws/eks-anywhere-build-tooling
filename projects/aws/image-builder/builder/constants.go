package builder

const (
	DefaultUbuntu2004AMIFilterName string = "ubuntu/images/*ubuntu-focal-20.04-amd64-server-*"
	DefaultUbuntu2204AMIFilterName string = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
	DefaultUbuntuAMIFilterOwners   string = "679593333241"
	DefaultAMIBuildRegion          string = "us-west-2"
	DefaultAMIBuilderInstanceType  string = "t3.small"
	DefaultAMIRootDeviceName       string = "/dev/sda1"
	DefaultAMIVolumeSize           string = "25"
	DefaultAMIVolumeType           string = "gp3"
	DefaultAMIManifestOutput       string = "manifest.json"

	// Paths and URLs
	buildToolingRepoUrl              string = "https://github.com/aws/eks-anywhere-build-tooling.git"
	imageBuilderProjectDirectory     string = "projects/kubernetes-sigs/image-builder"
	imageBuilderCAPIDirectory        string = "image-builder/images/capi"
	packerAdditionalFilesConfigFile  string = "packer/config/files.json"
	ansibleAdditionalFilesCustomRole string = "ansible/roles/load_additional_files"
	packerAdditionalFilesList        string = "packer/config/additional_files.yaml"
	packerVSphereConfigFile          string = "packer/ova/vsphere.json"
	packerBaremetalConfigFile        string = "packer/config/baremetal.json"
	packerUbuntuRawEFIConfigFile     string = "packer/raw/raw-ubuntu-2004-efi.json"
	OVMFCodeFile                     string = "/usr/share/edk2/ovmf/OVMF_CODE.fd"
	packerNutanixConfigFile          string = "packer/nutanix/nutanix.json"
	packerCloudStackConfigFile       string = "packer/config/cloudstack.json"
	packerAMIConfigFile              string = "packer/ami/ami.json"

	// Environment variables
	branchNameEnvVar                      string = "BRANCH_NAME"
	codebuildCIEnvVar                     string = "CODEBUILD_CI"
	codebuildSourceDirectoryEnvVar        string = "CODEBUILD_SRC_DIR"
	releaseBranchEnvVar                   string = "RELEASE_BRANCH"
	eksAReleaseVersionEnvVar              string = "EKSA_RELEASE_VERSION"
	packerAdditionalFilesConfigFileEnvVar string = "PACKER_ADDITIONAL_FILES_VAR_FILES"
	rhelUsernameEnvVar                    string = "RHSM_USERNAME"
	rhelPasswordEnvVar                    string = "RHSM_PASSWORD"
	packerTypeVarFilesEnvVar              string = "PACKER_TYPE_VAR_FILES"
)

var DefaultAMIAdditionalFiles = []File{
	{
		Source:      "packer/ami/additional_files/bootstrap.sh",
		Destination: "/etc/eks/",
		Owner:       "root",
		Group:       "root",
		Mode:        744,
	},
	{
		Source:      "packer/ami/additional_files/logging.sh",
		Destination: "/etc/eks/",
		Owner:       "root",
		Group:       "root",
		Mode:        744,
	},
}
