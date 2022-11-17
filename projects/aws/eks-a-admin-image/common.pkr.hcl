
locals {
  timestamp  = regex_replace(timestamp(), "[- TZ:]", "")
  image_name = "eks-a-admin-${replace(var.eks-a-version, "+", "-")}-${var.build-version}-${local.timestamp}"
}

variable "build-username" {
  type        = string
  description = "The username to login to the guest operating system. (e.g. 'ubuntu')"
  sensitive   = true
  default     = "ubuntu"
}

variable "build-password" {
  type        = string
  description = "The password to login to the guest operating system."
  sensitive   = true
  default     = "ubuntu"
}

variable "build-output" {
  type        = string
  description = "Path to build output dir"
  default     = "./_output"
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

variable "kind-url" {
  type    = string
  default = "https://github.com/kubernetes-sigs/kind/releases/download/v0.13.0/kind-linux-amd64"
}

variable "yq-url" {
  type    = string
  default = "https://github.com/mikefarah/yq/releases/download/v4.28.1/yq_linux_amd64"
}

variable "build-version" {
  type    = string
  default = "v0.0.0"
}

variable "golang-version" {
  type    = string
  default = "latest"
}

variable "kind-version" {
  type    = string
  default = "v0.12.0"
}

variable "manifest-output" {
  type    = string
  default = "manifest.json"
}