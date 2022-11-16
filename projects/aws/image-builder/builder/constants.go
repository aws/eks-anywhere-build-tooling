package builder

const (
	DefaultUbuntuAMIFilterName    string = "ubuntu/images/*ubuntu-focal-20.04-amd64-server-*"
	DefaultUbuntuAMIFilterOwners  string = "679593333241"
	DefaultAMIBuildRegion         string = "us-west-2"
	DefaultAMIBuilderInstanceType string = "t3.small"
	DefaultAMIRootDeviceName      string = "/dev/sda1"
	DefaultAMIVolumeSize          string = "25"
	DefaultAMIVolumeType          string = "gp3"
	DefaultAMICustomRoleNames     string = "/home/image-builder/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/ansible/roles/load_additional_files"
	DefaultAMIAnsibleExtraVars    string = "@/home/image-builder/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/packer/ami/ansible_extra_vars.yaml"
	DefaultAMIManifestOutput      string = "/home/image-builder/manifest.json"
)
