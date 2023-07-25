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
	DefaultAMICustomRoleNames      string = "projects/kubernetes-sigs/image-builder/ansible/roles/load_additional_files"
	DefaultAMIAnsibleExtraVars     string = "projects/kubernetes-sigs/image-builder/packer/config/additional_files.yaml"
	DefaultAMIManifestOutput       string = "manifest.json"

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
	codebuildCIEnvVar                     string = "CODEBUILD_CI"
	codebuildSourceDirectoryEnvVar        string = "CODEBUILD_SRC_DIR"
	releaseBranchEnvVar                   string = "RELEASE_BRANCH"
	packerAdditionalFilesConfigFileEnvVar string = "PACKER_ADDITIONAL_FILES_VAR_FILES"
	rhelUsernameEnvVar                    string = "RHSM_USERNAME"
	rhelPasswordEnvVar                    string = "RHSM_PASSWORD"
	packerTypeVarFilesEnvVar              string = "PACKER_TYPE_VAR_FILES"
)

var DefaultAMIAdditionalFiles = []File{
	{
		Source:      "projects/kubernetes-sigs/image-builder/packer/ami/additional_files/bootstrap.sh",
		Destination: "/etc/eks/",
		Owner:       "root",
		Group:       "root",
		Mode:        744,
	},
	{
		Source:      "projects/kubernetes-sigs/image-builder/packer/ami/additional_files/logging.sh",
		Destination: "/etc/eks/",
		Owner:       "root",
		Group:       "root",
		Mode:        744,
	},
}
