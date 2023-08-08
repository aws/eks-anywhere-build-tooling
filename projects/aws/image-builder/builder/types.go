package builder

const (
	Ubuntu     string = "ubuntu"
	RedHat     string = "redhat"
	VSphere    string = "vsphere"
	Baremetal  string = "baremetal"
	Nutanix    string = "nutanix"
	CloudStack string = "cloudstack"
	AMI        string = "ami"
	EFI        string = "efi"
	BIOS       string = "bios"
)

var SupportedHypervisors = []string{
	VSphere,
	Baremetal,
	Nutanix,
	CloudStack,
	AMI,
}

var SupportedUbuntuVersions = []string{
	"20.04",
	"22.04",
}

var SupportedRedHatVersions = []string{
	"8",
}

var SupportedFirmwares = []string{
	BIOS,
	EFI,
}

type BuildOptions struct {
	Os                 string
	OsVersion          string
	Hypervisor         string
	VsphereConfig      *VsphereConfig
	BaremetalConfig    *BaremetalConfig
	NutanixConfig      *NutanixConfig
	CloudstackConfig   *CloudstackConfig
	AMIConfig          *AMIConfig
	FilesConfig        *AdditionalFilesConfig
	ReleaseChannel     string
	Force              bool
	Firmware           string
	EKSAReleaseVersion string
}

type VsphereConfig struct {
	Cluster            string `json:"cluster"`
	ConvertToTemplate  string `json:"convert_to_template"`
	CreateSnapshot     string `json:"create_snapshot"`
	Datacenter         string `json:"datacenter"`
	Datastore          string `json:"datastore"`
	Folder             string `json:"folder"`
	InsecureConnection string `json:"insecure_connection"`
	LinkedClone        string `json:"linked_clone"`
	Network            string `json:"network"`
	ResourcePool       string `json:"resource_pool"`
	Template           string `json:"template"`
	VcenterServer      string `json:"vcenter_server"`
	VsphereLibraryName string `json:"vsphere_library_name"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	AdditionalFiles    []File `json:"files"`
	IsoConfig
	RhelConfig
	ExtraPackagesConfig
}

type BaremetalConfig struct {
	AdditionalFiles []File `json:"files"`
	IsoConfig
	RhelConfig
	ExtraPackagesConfig
}

type CloudstackConfig struct {
	AnsibleUserVars string `json:"ansible_user_vars"`
	IsoConfig
	RhelConfig
	ExtraPackagesConfig
}

type IsoConfig struct {
	IsoUrl          string `json:"iso_url,omitempty"`
	IsoChecksum     string `json:"iso_checksum,omitempty"`
	IsoChecksumType string `json:"iso_checksum_type,omitempty"`
}

type RhelConfig struct {
	RhelUsername string `json:"rhel_username"`
	RhelPassword string `json:"rhel_password"`
}

type NutanixConfig struct {
	ClusterName       string `json:"nutanix_cluster_name"`
	ImageName         string `json:"image_name"`
	SourceImageName   string `json:"source_image_name"`
	NutanixEndpoint   string `json:"nutanix_endpoint"`
	NutanixInsecure   string `json:"nutanix_insecure"`
	NutanixPort       string `json:"nutanix_port"`
	NutanixUserName   string `json:"nutanix_username"`
	NutanixPassword   string `json:"nutanix_password"`
	NutanixSubnetName string `json:"nutanix_subnet_name"`
	ExtraPackagesConfig
}

type AMIConfig struct {
	AMIFilterName       string `json:"ami_filter_name"`
	AMIFilterOwners     string `json:"ami_filter_owners"`
	AMIRegions          string `json:"ami_regions"`
	AWSRegion           string `json:"aws_region"`
	BuilderInstanceType string `json:"builder_instance_type"`
	ManifestOutput      string `json:"manifest_output"`
	RootDeviceName      string `json:"root_device_name"`
	SubnetID            string `json:"subnet_id"`
	VolumeSize          string `json:"volume_size"`
	VolumeType          string `json:"volume_type"`

	ExtraPackagesConfig
}

type ExtraPackagesConfig struct {
	ExtraDebs  string `json:"extra_debs,omitempty"`
	ExtraRepos string `json:"extra_repos,omitempty"`
	ExtraRpms  string `json:"extra_rpms,omitempty"`
}
