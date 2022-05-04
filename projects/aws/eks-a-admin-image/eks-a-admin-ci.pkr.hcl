packer {
  required_version = ">= 1.8.0"
  required_plugins {
    vsphere = {
      version = ">= v1.0.3"
      source  = "github.com/hashicorp/vsphere"
    }
  }
}

source "vsphere-iso" "ubuntu" {
  // vCenter Server Endpoint Settings and Credentials
  vcenter_server      = var.vsphere_endpoint
  username            = var.vsphere_username
  password            = var.vsphere_password
  insecure_connection = var.vsphere_insecure_connection

  // vSphere Settings
  datacenter = var.vsphere_datacenter
  cluster    = var.vsphere_cluster
  datastore  = var.vsphere_datastore
  folder     = var.vsphere_folder

  // Virtual Machine Settings
  guest_os_type        = var.vm_guest_os_type
  vm_name              = "${local.image_name}"
  firmware             = var.vm_firmware
  CPUs                 = var.vm_cpu_sockets
  cpu_cores            = var.vm_cpu_cores
  CPU_hot_plug         = var.vm_cpu_hot_add
  RAM                  = var.vm_mem_size
  RAM_hot_plug         = var.vm_mem_hot_add
  cdrom_type           = var.vm_cdrom_type
  disk_controller_type = var.vm_disk_controller_type
  storage {
    disk_size             = var.vm_disk_size
    disk_thin_provisioned = var.vm_disk_thin_provisioned
  }
  network_adapters {
    network      = var.vsphere_network
    network_card = var.vm_network_card
  }
  vm_version           = var.common_vm_version
  remove_cdrom         = var.common_remove_cdrom
  tools_upgrade_policy = var.common_tools_upgrade_policy
  notes                = "Version: ${var.eks-a-version}-${var.build-version}-${local.timestamp}\nBuilt on: ${local.timestamp}"

  // Removable Media Settings
  iso_url      = "${var.iso_url}"
  iso_checksum = "${var.iso_checksum_type}:${var.iso_checksum_value}"

  // Boot and Provisioning Settings
  ip_wait_timeout = var.common_ip_wait_timeout

  cd_files = var.vm_cd_files
  cd_label = var.vm_cd_label

  boot_wait    = var.vm_boot_wait
  boot_command = var.vm_boot_command

  shutdown_command = "echo '${var.build-password}' | sudo -S -E shutdown -P now"
  shutdown_timeout = var.common_shutdown_timeout

  // Communicator Settings and Credentials
  communicator = "ssh"
  ssh_username = var.build-username
  ssh_password = var.build-password
  ssh_port     = var.communicator_port
  ssh_timeout  = var.communicator_timeout

  export {
    force            = true
    output_directory = var.build-output
  }

}

build {
  name = "eks-a-admin-ci-image"
  sources = [
    "source.vsphere-iso.ubuntu",
  ]

  provisioner "shell" {
    script            = "provisioners/upgrade_linux.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    environment_vars = [
      "USER=${var.build-username}",
    ]
    // wait for reboot before starting
    pause_before      = "10s"
    script            = "provisioners/install_docker.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    environment_vars = [
      "USER=${var.build-username}",
    ]
    // wait for reboot before starting
    pause_before = "10s"
    script       = "provisioners/test/install_docker.sh"
  }

  provisioner "shell" {
    environment_vars = [
      "KUBECTL_URL=${var.kubectl-url}",
      "KIND_URL=${var.kind-url}",
      "EKSA_VERSION=${var.eks-a-version}",
      "EKSA_RELEASE_MANIFEST_URL=${var.eks-a-release-manifest-url}",
      "GO_VERSION=${var.golang-version}",
    ]

    scripts = [
      "provisioners/install_kubectl.sh",
      "provisioners/test/install_kubectl.sh",
      "provisioners/install_eksa.sh",
      "provisioners/test/install_eksa.sh",
      "provisioners/install_golang.sh",
      "provisioners/test/install_golang.sh",
      "provisioners/install_kind.sh",
      "provisioners/test/install_kind.sh",
      "provisioners/install_awscli.sh",
      "provisioners/test/install_awscli.sh",
      "provisioners/install_ssm_agent.sh",
      "provisioners/test/install_ssm_agent.sh",
    ]
  }

  provisioner "shell" {
    script = "provisioners/cleanup.sh"
  }

  post-processor "manifest" {
    output = "${var.manifest-output}"
    custom_data = {
      eks-a-version          = "${var.eks-a-version}"
      eks-a-release-manifest = "${var.eks-a-release-manifest-url}"
    }
  }

  post-processor "shell-local" {
    only   = ["vsphere-iso.ubuntu"]
    inline = ["cd ${var.build-output} && tar -cf ${local.image_name}.ova *"]
  }
}

/*
    DESCRIPTION:
    vSphere variables using the Packer Builder for VMware vSphere (vsphere-iso).
*/

//  BLOCK: variable
//  Defines the input variables.

// vSphere Credentials
variable "vsphere_endpoint" {
  type        = string
  description = "The fully qualified domain name or IP address of the vCenter Server instance. (e.g. 'sfo-w01-vc01.sfo.rainpole.io')"
  default     = null
}

variable "vsphere_username" {
  type        = string
  description = "The username to login to the vCenter Server instance. (e.g. 'svc-packer-vsphere@rainpole.io')"
  sensitive   = true
  default     = null
}

variable "vsphere_password" {
  type        = string
  description = "The password for the login to the vCenter Server instance."
  sensitive   = true
  default     = null
}

variable "vsphere_insecure_connection" {
  type        = bool
  description = "Do not validate vCenter Server TLS certificate."
  default     = true
}

// vSphere Settings
variable "vsphere_datacenter" {
  type        = string
  description = "The name of the target vSphere datacenter. (e.g. 'sfo-w01-dc01')"
  default     = null
}

variable "vsphere_cluster" {
  type        = string
  description = "The name of the target vSphere cluster. (e.g. 'sfo-w01-cl01')"
  default     = null
}

variable "vsphere_datastore" {
  type        = string
  description = "The name of the target vSphere datastore. (e.g. 'sfo-w01-cl01-vsan01')"
  default     = null
}

variable "vsphere_network" {
  type        = string
  description = "The name of the target vSphere network segment. (e.g. 'sfo-w01-dhcp')"
  default     = null
}

variable "vsphere_folder" {
  type        = string
  description = "The name of the target vSphere cluster. (e.g. 'sfo-w01-fd-templates')"
  default     = null
}

// Virtual Machine Settings
variable "vm_guest_os_type" {
  type        = string
  description = "The guest operating system type, also know as guestid. (e.g. 'ubuntu64Guest')"
  default     = "ubuntu64Guest"
}

variable "vm_firmware" {
  type        = string
  description = "The virtual machine firmware. (e.g. 'efi-secure'. 'efi', or 'bios')"
  default     = "efi-secure"
}

variable "vm_cdrom_type" {
  type        = string
  description = "The virtual machine CD-ROM type. (e.g. 'sata', or 'ide')"
  default     = "sata"
}

variable "vm_cpu_sockets" {
  type        = number
  description = "The number of virtual CPUs sockets. (e.g. '2')"
  default     = 4
}

variable "vm_cpu_cores" {
  type        = number
  description = "The number of virtual CPUs cores per socket. (e.g. '1')"
  default     = 1
}

variable "vm_cpu_hot_add" {
  type        = bool
  description = "Enable hot add CPU."
  default     = true
}

variable "vm_mem_size" {
  type        = number
  description = "The size for the virtual memory in MB. (e.g. '2048')"
  default     = 2048
}

variable "vm_mem_hot_add" {
  type        = bool
  description = "Enable hot add memory."
  default     = true
}

variable "vm_disk_size" {
  type        = number
  description = "The size for the virtual disk in MB. (e.g. '40960')"
  default     = 40960
}

variable "vm_disk_controller_type" {
  type        = list(string)
  description = "The virtual disk controller types in sequence. (e.g. 'pvscsi')"
  default     = ["pvscsi"]
}

variable "vm_disk_thin_provisioned" {
  type        = bool
  description = "Thin provision the virtual disk."
  default     = true
}

variable "vm_network_card" {
  type        = string
  description = "The virtual network card type. (e.g. 'vmxnet3' or 'e1000e')"
  default     = "vmxnet3"
}

variable "vm_boot_command" {
  type        = list(string)
  description = "The boot command to run on the VM"
  default     = [""]
}

variable "vm_cd_label" {
  type        = string
  description = "The CD label for the boot command to reference"
  default     = "cidata"
}

variable "vm_cd_files" {
  type        = list(string)
  description = "Files to be mounted via cdrom for boot process to reference"
  default     = ["./ova/linux/ubuntu/data/meta-data", "./ova/linux/ubuntu/data/user-data"]
}

variable "common_vm_version" {
  type        = number
  description = "The vSphere virtual hardware version. (e.g. '19')"
  default     = "19"
}

variable "common_tools_upgrade_policy" {
  type        = bool
  description = "Upgrade VMware Tools on reboot."
  default     = true
}

variable "common_remove_cdrom" {
  type        = bool
  description = "Remove the virtual CD-ROM(s)."
  default     = true
}

// Removable Media Settings
variable "iso_url" {
  type        = string
  description = "The url name of the ISO image used by the vendor. (e.g. 'http://ubuntu-<version>-live-server-amd64.iso')"
  default     = null
}

variable "iso_checksum_type" {
  type        = string
  description = "The checksum algorithm used by the vendor. (e.g. 'sha256')"
  default     = "sha256"
}

variable "iso_checksum_value" {
  type        = string
  description = "The checksum value provided by the vendor."
  default     = null
}

// Boot Settings
variable "vm_boot_order" {
  type        = string
  description = "The boot order for virtual machines devices. (e.g. 'disk,cdrom')"
  default     = "disk,cdrom"
}

variable "vm_boot_wait" {
  type        = string
  description = "The time to wait before boot."
  default     = "3s"
}

variable "common_ip_wait_timeout" {
  type        = string
  description = "Time to wait for guest operating system IP address response."
  default     = "15m"
}

variable "common_shutdown_timeout" {
  type        = string
  description = "Time to wait for guest operating system shutdown."
  default     = "5m"
}

// Communicator Settings and Credentials
variable "communicator_port" {
  type        = string
  description = "The port for the communicator protocol."
  default     = 22
}

variable "communicator_timeout" {
  type        = string
  description = "The timeout for the communicator protocol."
  default     = "30m"
}
