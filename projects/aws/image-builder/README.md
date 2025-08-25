## **Image Builder Tool**
![Version](https://img.shields.io/badge/version-v0.6.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiRHQ0UnNzTElaQyt5eDI5OG9XYUhYQW85WXE5RzI3Sjd5YWFwK2d2aHBVb2R4dS8xek5aeUcrVHJFN05JR2JnbWx2aGRURlAxdDZrNFQwMFRaMzY4MWU0PSIsIml2UGFyYW1ldGVyU3BlYyI6InIxUHNId1RQcCs3SzlFWWQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

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
* 1-34
* 1-33
* 1-32
* 1-31
* 1-30
* 1-29
* 1-28
* 1-27

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
3. Create vsphere-connection.json config file
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
  "vmx_version":"<hardware version of virtual machine>",
  "template": "<template used by clone builder>"
}
```
4. Run the image builder tool for appropriate release channel
```
image-builder build --os ubuntu --hypervisor vsphere --vsphere-config <path to above json file> --release-channel <release channel, ex 1-23> --builder <vsphere builder type, can be iso or clone>
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

### Air Gapped Image Building
Image builder only supports building Ubuntu in an airgapped mode for now.

1. Create the config json file for respective provider and make sure to include the fields required for airgapped building.
   An example of baremetal config json for ubuntu airgapped builds are below
   ```
   {
      "eksa_build_tooling_repo_url": "https://internal-repos/eks-anywhere-build-tooling.git",
      "image_builder_repo_url": "https://internal-repos/image-builder.git",
      "private_artifacts_eksd_fqdn": "http://artifactory:8081/artifactory",
      "private_artifacts_eksa_fqdn": "http://artifactory:8081/artifactory/EKS-A",
      "extra_repos": "/home/airgapped/sources.list",
      "disable_public_repos": "true",
      "iso_url": "http://artifactory:8081/artifactory/EKS-A/ISO/ubuntu-20.04.1-legacy-server-amd64.iso",
      "iso_checksum": "f11bda2f2caed8f420802b59f382c25160b114ccc665dbac9c5046e7fceaced2",
      "iso_checksum_type": "sha256"
    }
   ```
2. Install pre-requisites required for image builder in the environment or admin machine.
   1. Packer version 1.9.4
   2. Ansible version 2.15.3
   3. Packer provisioner goss version 3.1.4
   4. jq
   5. yq
   6. unzip 
   7. make 
   8. python3-pip
3. From an environment with internet access run the following command to generate the manifest tarball
    ```
   image-builder download manifests
    ```
   This will download a eks-a-manifests.tar in the current working directory. This tarball is required for airgapped building.
4. Replicate all the required EKS-D and EKS-A artifacts to the internal artifacts server like artifactory.
   Required artifacts are as follows
   EKS-D amd64 artifacts for specific release branch
   1. kube-apiserver.tar
   2. kube-scheduler.tar
   3. kube-proxy.tar
   4. kube-controller-manager.tar
   5. etcd.tar
   6. coredns.tar
   7. pause.tar
   8. kubectl
   9. kubeadm
   10. kubelet
   11. etcd-linux-amd64-v<version>.tar.gz
   12. cni-plugins-linux-amd64-v<version>.tar.gz
   
   
   EKS-A amd64 artifacts for specific release bundle version
   1. containerd-linux-amd64.tar.gz
   2. etcdadm-linux-amd64.tar.gz
   3. cri-tools-amd64.tar.gz

   In addition to these EKS-D and EKS-A artifacts please ensure the base ubuntu iso is also hosted internally.
   
5. Run image builder in airgapped mode
   ```
   image-builder build --os ubuntu --hypervisor baremetal --release-channel 1-27 --air-gapped --baremetal-config baremetal.json --manifest-tarball eks-a-manifests.tar
   ```