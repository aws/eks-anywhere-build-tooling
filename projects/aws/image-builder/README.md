## **Image Builder Tool**
![Version](https://img.shields.io/badge/version-v0.1.4-blue)

Image Builder Tool is a cli that builds EKS-A compatible Kubernetes node images. The tool is based on upstream
[image-builder](https://github.com/kubernetes-sigs/image-builder) project and uses packer to build the node images.
The tool always builds an image with the latest release artifacts and latest OS packages.

Supported Providers
* vsphere
* baremetal
* nutanix
* snow
* cloudstack

Supported OSes
* Ubuntu
* Red Hat Enterprise Linux

Supported Release Channels
* 1-27
* 1-26
* 1-25
* 1-24
* 1-23

Supported Firmwares
* bios
* efi (for Ubuntu Vsphere and Baremetal image builds)

### Building Node Images for vSphere

Vsphere is one of the supported infrastructure providers the Image builder tool can build EKS-A node images for. In order
to build a node image for vSphere, the image builder tool needs to run from an environment with network access to vcenter.

#### Pre-reqs for vSphere Image Builder
* Machine running Ubuntu 22.04 that has access to the vcenter environment and vcenter url. This machine could also run as a vm on vSphere.
* vSphere user with permissions listed below
* Minimum machine resources required
  * 50 GB disk space on the machine
  * 2 vCPUs
  * 8 GB RAM
* Machine as well as vSphere network used below must have network access to
  * vCenter endpoint
  * public.ecr.aws (to download container images from EKS-A)
  * anywhere-assets.eks.amazonaws.com (to download EKS-A binaries)
  * distro.eks.amazonaws.com (to download EKS-D binaries)
  * d2glxqk2uabbnd.cloudfront.net (for EKS-A and EKS-D ECR container images)

##### vSphere User Permissions
Inventory:
* Create new

Configuration:
* Change configuration
* Add new disk
* Add or remove device
* Change memory
* Change settings
* Set annotation

Interaction:
* Power on
* Power off
* Console interaction
* Configure CD media
* Device connection

Snapshot management:
* Create snapshot

Provisioning

* Mark as template

Resource Pool
* Assign vm to resource pool

Datastore
* Allocate space
* Browse data
* Low level file operations

Network
* Assign network to vm

#### Building a vSphere OVA Node Image
1. Install pre-requisite packages
```
sudo apt update -y
sudo apt install jq unzip make ansible -y
sudo snap install yq
```
2. Build or download the image builder tool
3. Create a content library on vSphere
```
govc library.create "<library name>"
```
4. Create vsphere-connection.json config file
```
{
  "cluster":"<vsphere cluster used for image building>",
  "convert_to_template":"false",
  "create_snapshot":"<creates a snapshot on base OVA after building if set to true>",
  "datacenter":"<vsphere datacenter used for image building>",
  "datastore":"<datastore used to store template/for image building>",
  "folder":"<folder on vsphere to create temporary vm>",
  "insecure_connection":"true",
  "linked_clone":"false",
  "network":"<vsphere network used for image building>",
  "password":"<vcenter username>",
  "resource_pool":"<resource pool used for image building vm>",
  "username":"<vcenter username>",
  "vcenter_server":"<vcenter fqdn>",
  "vsphere_library_name": "<vsphere content library name>"
}
```
5. Run the image builder tool for appropriate release channel
```
image-builder build --os ubuntu --hypervisor vsphere --vsphere-config <path to above json file> --release-channel <release channel, ex 1-23>
```

### Building Node Images for Baremetal

Baremetal is one of the supported infrastructure providers the Image builder tool can build EKS-A node images for. In order
to build a node image for baremetal, the image builder tool needs to run on baremetal machine.

#### Pre-reqs for baremetal Image Builder
* Baremetal machine running Ubuntu 22.04 with virtualization enabled
* Minimum machine resources required
    * 50 GB disk space on the machine
    * 2 vCPUs
    * 8 GB RAM
* Machine must have network access to
    * vCenter endpoint
    * public.ecr.aws (to download container images from EKS-A)
    * anywhere-assets.eks.amazonaws.com (to download EKS-A binaries)
    * distro.eks.amazonaws.com (to download EKS-D binaries)
    * d2glxqk2uabbnd.cloudfront.net (for EKS-A and EKS-D ECR container images)

#### Building a baremetal Node Image
1. Install pre-requisite packages and prep environment
```
sudo apt update -y
sudo apt install jq make qemu-kvm libvirt-daemon-system libvirt-clients virtinst cpu-checker libguestfs-tools libosinfo-bin unzip ansible -y
sudo snap install yq
sudo usermod -a -G kvm $USER
sudo chmod 666 /dev/kvm
sudo chown root:kvm /dev/kvm
echo "HostKeyAlgorithms +ssh-rsa" >> /home/$USER/.ssh/config
echo "PubkeyAcceptedKeyTypes +ssh-rsa" >> /home/$USER/.ssh/config
```
2. Build or download the image builder tool
3. Run the image builder tool for appropriate release channel
```
image-builder build --os ubuntu --hypervisor baremetal --release-channel <release channel, ex 1-23>
```

The baremetal image built from image-builder tool should be hosted and its URL should be provided to `osImageURL` under `TinkerbellDatacenterConfig`
in the cluster spec to create a cluster using the built node image.

### Additional Configuration - Proxy
The Image Builder tool also supports some additional configuration. For now this is limited to supporting a proxy. 
Users can use proxy server to route outbound requests to internet. To configure the image builder tool to use proxy, simply
export the following proxy environment variables
```
export HTTP_PROXY=<HTTP proxy URL e.g. http://proxy.corp.com:80>
export HTTPS_PROXY=<HTTPS proxy URL e.g. http://proxy.corp.com:443>
export NO_PROXY=<No proxy>
```

### Building Node Images for Nutanix AHV

Nutanix is one of the supported infrastructure providers the Image builder tool can build EKS-A node images for. In order
to build a node image for Nuntaix AHV, the image builder tool needs to run from an environment with network access to Nutanix Prism Central Endpoint.

#### Pre-reqs for Nutanix AHV Image Builder
* Machine running Ubuntu 22.04 that has access to the Nutanix Prism Central environment. This machine could also run as a vm on Nutanix AHV Prism Element.
* Nutanix Prism Central user with permissions listed below
* Minimum machine resources required
  * 50 GB disk space on the machine
  * 2 vCPUs
  * 8 GB RAM
* Machine as well as Nutanix Prism Central network used below must have network access to
  * Nutanix Prism Central endpoint
  * public.ecr.aws (to download container images from EKS-A)
  * anywhere-assets.eks.amazonaws.com (to download EKS-A binaries)
  * distro.eks.amazonaws.com (to download EKS-D binaries)
  * d2glxqk2uabbnd.cloudfront.net (for EKS-A and EKS-D ECR container images)

#### Building a Nutanix AHV Node Image
1. Install pre-requisite packages
```
sudo apt update -y
sudo apt install jq unzip make ansible -y
sudo snap install yq
```
2. Build or download the image builder tool
4. Create nutanix-connection.json config file. more details on values can be found here https://image-builder.sigs.k8s.io/capi/providers/nutanix.html
```
{
  "nutanix_cluster_name": "",
  "image_name": "",
  "source_image_name": "",
  "nutanix_endpoint": "",
  "nutanix_insecure": "",
  "nutanix_port": "9440",
  "nutanix_username": "",
  "nutanix_password": "",
  "nutanix_subnet_name": ""
}
```
5. Run the image builder tool for appropriate release channel
```
image-builder build --os ubuntu --hypervisor nutanix --nutanix-config <path to above json file> --release-channel <release channel, ex 1-23>
```