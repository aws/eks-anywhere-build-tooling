package builder

const (
	Ubuntu    string = "ubuntu"
	VSphere   string = "vsphere"
	Baremetal string = "baremetal"
)

type BuildOptions struct {
	Os              string
	Hypervisor      string
	VsphereConfig   string
	ReleaseChannel  string
	artifactsBucket string
	Force           bool
}
