## **Cluster API**
![Version](https://img.shields.io/badge/version-v0.3.23-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQVZ3TDBZZVVXZUZiVmtqLzVoOVcrV2FaMmxRRzJXRmJCRlZtQkNodXdWZ0FrNm0zQ3l5UzNqTkdsQXgwdzc0bTBZc1RIcjBhMUVFbEhIK3d2VDVPek1rPSIsIml2UGFyYW1ldGVyU3BlYyI6IkVuOGJxNXBPZEtDek81Q3giLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Cluster API](https://github.com/kubernetes-sigs/cluster-api) is a Kubernetes sub-project focused on providing declarative APIs and tooling to simplify provisioning, upgrading, and operating multiple Kubernetes clusters. It uses Kubernetes-style APIs and patterns to automate cluster lifecycle management for platform operators. The supporting infrastructure, like virtual machines, networks, load balancers, and VPCs, as well as the Kubernetes cluster configuration are all defined in the same way that application developers operate deploying and managing their workloads. This enables consistent and repeatable cluster deployments across a wide variety of infrastructure environments. Cluster API can be extended to support any infrastructure provider (AWS, Azure, vSphere, etc.) or bootstrap provider (kubeadm is default) as required by the customer.

The `eks-a` CLI uses Cluster API to generate configurations for various infrasurcture providers, and uses it to create and manage multiple workload clusters.

You can find the latest versions of these images on ECR Public Gallery.

[Cluster API Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/cluster-api-controller) | 
[Cluster API Kubeadm Bootstrap Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-bootstrap-controller) | 
[Cluster API Kubeadm Controlplane Controller](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/cluster-api/kubeadm-control-plane-controller)
