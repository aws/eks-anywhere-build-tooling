## **Cluster API Provider for AWS**
![Version](https://img.shields.io/badge/version-v0.6.4-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiOEtjdDI0NC9tSXE0MkJKU0hsL0RvL2pyV0pidEZzY2FnbG9YZk5IeVZvQVJwNDBZdkRRUXgra3pXeSs0dUtGbm1uSU1NRGRjbzJTeG9lcEhQVFEySEJzPSIsIml2UGFyYW1ldGVyU3BlYyI6InZUVnJxVkUvWDJQNDhmaUciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Kubernetes Cluster API Provider AWS (CAPA)](https://github.com/kubernetes-sigs/cluster-api-provider-aws) provides Kubernetes-native declarative infrastructure for AWS. It allows for true AWS hybrid deployments of Kubernetes. With CAPA, customers can re-use their existing AWS infrastructure when standing up their workload clusters.

Some of the features of Cluster API Provider AWS include:
* Native Kubernetes manifests and API
* Manages the bootstrapping of VPCs, gateways, security groups and instances.
* Choice of Linux distribution between Amazon Linux 2, CentOS 7 and Ubuntu 18.04, using pre-baked AMIs.
* Deploys Kubernetes control planes into private subnets with a separate bastion server.
* Supports control planes on EC2 instances.
* Experimental EKS support

Cluster API Provider AWS controller images are used in the Provider confgiration to bootstrap the AWS Infrastructure Provider in the EKS-A CLI.

You can find the latest versions of these images on ECR Public Gallery.

[Cluster API AWS Controller](https://gallery.ecr.aws/l0g8r8j6/kubernetes-sigs/cluster-api-provider-aws/cluster-api-aws-controller) | 
[Cluster API EKS Bootstrap Controller](https://gallery.ecr.aws/l0g8r8j6/kubernetes-sigs/cluster-api-provider-aws/eks-bootstrap-controller) | 
[Cluster API EKS Controlplane Controller](https://gallery.ecr.aws/l0g8r8j6/kubernetes-sigs/cluster-api-provider-aws/eks-control-plane-controller)

