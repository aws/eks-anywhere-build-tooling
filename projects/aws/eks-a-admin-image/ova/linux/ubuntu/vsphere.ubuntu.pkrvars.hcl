/*
    DESCRIPTION:
    Ubuntu Server 20.04 LTS  variables used by the Packer Plugin for VMware vSphere (vsphere-iso).
*/

// Virtual Machine Guest Operating System Setting
vm_guest_os_type = "ubuntu64Guest"

// Virtual Machine Hardware Settings
vm_firmware              = "efi-secure"
vm_cdrom_type            = "sata"
vm_cpu_sockets           = 4
vm_cpu_cores             = 1
vm_cpu_hot_add           = true
vm_mem_size              = 16384
vm_mem_hot_add           = true
vm_disk_size             = 51200
vm_disk_controller_type  = ["pvscsi"]
vm_disk_thin_provisioned = true
vm_network_card          = "vmxnet3"

// Removable Media Settings
iso_url           = "https://releases.ubuntu.com/20.04/ubuntu-20.04.4-live-server-amd64.iso"
iso_checksum_type  = "sha256"
iso_checksum_value = "28ccdb56450e643bad03bb7bcf7507ce3d8d90e8bf09e38f6bd9ac298a98eaad"

// Boot Settings
vm_boot_order = "disk,cdrom"

// Communicator Settings
communicator_port    = 22
communicator_timeout = "30m"

vm_cd_files = [
    "./ova/linux/ubuntu/data/meta-data",
    "./ova/linux/ubuntu/data/user-data",
  ]

vm_cd_label = "cidata"

vm_boot_command = [
    "<esc><wait>",
    "<esc><wait>",
    "linux /casper/vmlinuz --- autoinstall ds=nocloud;seedfrom=/cidata/",
    "<enter><wait>",
    "initrd /casper/initrd<enter><wait>",
    "boot<enter>",
]

// ssh creds
build-username = "eksadmin"

# passed as env var in make file
#build_password
