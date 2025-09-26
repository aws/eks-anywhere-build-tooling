packer {
  required_version = ">= 1.14.1"
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "amazonlinux2" {
  ami_name      = "${local.image_name}"
  instance_type = "t3.xlarge"
  // with t3.micro 48 min
  // with t3.large 45 min
  // with t3.xlarge 40 min
  // with t3.2xlarge 36 min
  region = "us-west-2"
  source_ami_filter {
    filters = {
      name                = "amzn2-ami-kernel-5.*-hvm-2*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
      architecture        = "x86_64"
    }
    most_recent = true
    owners      = ["amazon"]
  }
  ssh_username = "ec2-user"

  launch_block_device_mappings {
    device_name           = "/dev/xvda"
    delete_on_termination = true
    volume_size           = 70
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

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
  }
}

build {
  name = "eks-a-admin-snow-image"
  sources = [
    "source.amazon-ebs.amazonlinux2",
  ]

  provisioner "shell" {
    script            = "provisioners/al2/upgrade_linux.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    environment_vars = [
      "YQ_URL=${var.yq-url}",
    ]
    script = "provisioners/al2/install_deps.sh"
  }

  provisioner "shell" {
    environment_vars = [
      "USER=ec2-user",
    ]
    // wait for reboot before starting
    pause_before      = "10s"
    script            = "provisioners/setup_docker.sh"
    expect_disconnect = true
  }

  provisioner "shell" {
    environment_vars = [
      "USER=ec2-user",
    ]
    // wait for reboot before starting
    pause_before = "10s"
    script       = "provisioners/test/install_docker.sh"
  }

  provisioner "shell" {
    environment_vars = [
      "KUBECTL_URL=${var.kubectl-url}",
      "YQ_URL=${var.yq-url}",
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
    script = "provisioners/al2/cleanup.sh"
  }

  post-processor "manifest" {
    output = "${var.manifest-output}"
    custom_data = {
      eks-a-version          = "${var.eks-a-version}"
      eks-a-release-manifest = "${var.eks-a-release-manifest-url}"
    }
  }
}
