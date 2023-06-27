package builder

const (
	DefaultUbuntu2004AMIFilterName  string = "ubuntu/images/*ubuntu-focal-20.04-amd64-server-*"
	DefaultUbuntu2204AMIFilterName  string = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
	DefaultUbuntuAMIFilterOwners    string = "679593333241"
	DefaultAMIBuildRegion           string = "us-west-2"
	DefaultAMIBuilderInstanceType   string = "t3.small"
	DefaultAMIRootDeviceName        string = "/dev/sda1"
	DefaultAMIVolumeSize            string = "25"
	DefaultAMIVolumeType            string = "gp3"
	DefaultAMIAnsibleExtraVars      string = "@/home/image-builder/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/packer/ami/ansible_extra_vars.yaml"
	DefaultAMIManifestOutput        string = "/home/image-builder/manifest.json"
)
