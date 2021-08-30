## **Cluster API Provider for vSphere**
![Version](https://img.shields.io/badge/version-v0.7.10-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYm85WnJ4aDc2ZXhhVUxOWHJuUFJwN3FlQmE2L1Q4b2ZzNG91OVpjNVNGM1ZvbVBEUUM2bkdER3N5eVNrWTBKS2VSSW9Oa051aFVWS1dzVVlTOHBBZ0NRPSIsIml2UGFyYW1ldGVyU3BlYyI6IlEwOWNtd0llNXdjUGRvQWkiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Cluster API Provider for vSphere (CAPV)](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere) is a a concrete implementation of Cluster API for vSphere, which paves the way for true vSphere hybrid deployments of Kubernetes. CAPV is designed to allow customers to use their existing vSphere infrastructure, including vCenter credentials, VMs, templates, etc. for bootstrapping and creating workload clusters.

Some of the features of Cluster API Provider vSphere include:
* Native Kubernetes manifests and API
* Manages the bootstrapping of VMs on cluster.
* Choice of Linux distribution between Ubuntu 18.04 and CentOS 7 using VM Templates based on OVA images
* Deploys Kubernetes control planes into provided clusters on vSphere.

The Cluster API Provider vSphere controller image is used in the Provider confgiration to bootstrap the vSphere Infrastructure Provider in the EKS-A CLI.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/l0g8r8j6/kubernetes-sigs/cluster-api-provider-vsphere/release/manager).