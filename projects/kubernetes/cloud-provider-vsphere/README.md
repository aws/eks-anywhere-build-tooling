## **Kubernetes vSphere Cloud Provider**
![1.20 Version](https://img.shields.io/badge/1--20%20version-v1.20.1-blue)
![1.21 Version](https://img.shields.io/badge/1--21%20version-v1.21.3-blue)
![1.22 Version](https://img.shields.io/badge/1--22%20version-v1.22.6-blue)
![1.23 Version](https://img.shields.io/badge/1--23%20version-v1.23.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYzQ3dzRvZHVqU2MvYnVuMzB3QmRZdVd1U1RabVorWnlqTXBYUGxDSGk2NXJXUU12c3pLQ25CQUdaQmlNUE84S0JIVVZUU0ozeTJJb3J0NWxNejNSbzk4PSIsIml2UGFyYW1ldGVyU3BlYyI6IkhLNTZwQ0hiZDZVUzVRdXYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Cloud Provider vSphere](https://github.com/kubernetes/cloud-provider-vsphere) defines the vSphere-specific implementation of the Kubernetes controller-manager. The Cloud Provider Interface (CPI) allows customers to run Kubernetes clusters on vSphere infrastructure. It replaces the Kubernetes Controller Manager for only the cloud-specific control loops. The CPI integration connects to vCenter Server and maps information about the infrastructure, such as VMs, disks, and so on, back to the Kubernetes API.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kubernetes/cloud-provider-vsphere/cpi/manager).
