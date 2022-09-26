packer {
  required_version = ">= 1.8.0"
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "${local.image_name}"
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
    max_attempts  = 1440
  }
}

build {
  name = "eks-a-admin-snow-image"
  sources = [
    "source.amazon-ebs.ubuntu",
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
