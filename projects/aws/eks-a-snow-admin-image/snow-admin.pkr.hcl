packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "eks-a-version" {
  type    = string
  default = "latest"
}

variable "eks-a-release-manifest-url" {
  type    = string
  default = "https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml"
}

variable "kubectl-url" {
  type    = string
  default = "https://distro.eks.amazonaws.com/kubernetes-1-22/releases/4/artifacts/kubernetes/v1.22.6/bin/linux/amd64/kubectl"
}

variable "build-version" {
  type    = string
  default = "X"
}

variable "manifest-output" {
  type    = string
  default = "manifets.json"
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "snow-eks-a-admin-${var.eks-a-version}-${var.build-version}-${local.timestamp}"
  instance_type = "t3.xlarge"
  // with t3.micro 48 min
  // with t3.large 45 min
  // with t3.xlarge 40 min
  // with t3.2xlarge 36 min
  region = "us-west-2"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    // canonical account
    owners = ["099720109477"]
  }
  ssh_username = "ubuntu"

  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    delete_on_termination = true
    volume_size           = 50
    volume_type           = "gp3"
    iops                  = 3000
    throughput            = 125
  }

  // Increasing the default polling to avoid timeouts
  // Generating the AMI from the instance takes a while because
  // the AMI is very big in size
  aws_polling {
    delay_seconds = 5
    max_attempts  = 800
  }
}

build {
  name = "snow-eks-a-admin-ami"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    script            = "provisioners/upgrade_linux.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    // wait for reboot before starting
    pause_before      = "10s"
    script            = "provisioners/install_docker.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    // wait for reboot before starting
    pause_before = "10s"
    script       = "provisioners/test/install_docker.sh"
  }

  provisioner "shell" {
    environment_vars = [
      "KUBECTL_URL=${var.kubectl-url}",
      "EKSA_VERSION=${var.eks-a-version}",
      "EKSA_RELEASE_MANIFEST_URL=${var.eks-a-release-manifest-url}",
    ]
    scripts = [
      "provisioners/install_kubectl.sh",
      "provisioners/test/install_kubectl.sh",
      "provisioners/install_eksa.sh",
      "provisioners/test/install_eksa.sh",
      "provisioners/download_eksa_artifacts.sh",
      "provisioners/test/download_eksa_artifacts.sh",
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
}
