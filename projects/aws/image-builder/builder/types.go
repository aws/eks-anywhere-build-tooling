package builder

const (
	Ubuntu     string = "ubuntu"
	RedHat     string = "redhat"
	VSphere    string = "vsphere"
	Baremetal  string = "baremetal"
	NutanixAHV string = "nutanixahv"
)

type BuildOptions struct {
	Os               string
	Hypervisor       string
	VsphereConfig    *VsphereConfig
	NutanixAHVConfig *NutanixAHVConfig
	ReleaseChannel   string
	artifactsBucket  string
	Force            bool
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
	Username           string `json:"username"`
	Password           string `json:"password"`
	IsoUrl             string `json:"iso_url,omitempty"`
	IsoChecksum        string `json:"iso_checksum,omitempty"`
	IsoChecksumType    string `json:"iso_checksum_type,omitempty"`
	RhelUsername       string `json:"rhel_username"`
	RhelPassword       string `json:"rhel_password"`
}

type NutanixAHVConfig struct {
	ClusterName       string `json:nutanix_cluster_name`
	ImageName         string `json:image_name`
	SourceImageName   string `json:source_image_name`
	NutanixEndpoint   string `json:nutanix_endpoint`
	NutanixInsecure   string `json:nutanix_insecure`
	NutanixPort       string `json:nutanix_port`
	NutanixUserName   string `json:nutanix_username`
	NutanixPassword   string `json:nutanix_password`
	NutanixSubnetName string `json:nutanix_subnet_name`
}
