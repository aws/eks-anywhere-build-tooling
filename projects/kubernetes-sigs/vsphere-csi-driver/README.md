## **Container Storage Interface (CSI) driver for vSphere**
![Version](https://img.shields.io/badge/version-v2.2.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiNFozNkl5RVJqek5vb0dYUHFzS2VaZnhUVldQYjBEalp2Wm5XSm0wV2JseXNlODhyVWdrV2NoRFhzR043L3E2NDEwOTBidHNZS3pGMTd0VDFIbCt6WVhVPSIsIml2UGFyYW1ldGVyU3BlYyI6IlBuNTdYZGduajFoa2tnTUEiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [vSphere CSI Driver project](https://github.com/kubernetes-sigs/vsphere-csi-driver), in conjunction with vSphere Cloud Provider, exposes vSphere storage and features to Kubernetes users in the form of Cloud Native Storage (CNS). The main goal of CNS is to make vSphere and vSphere storage, including vSAN, a platform to run stateful Kubernetes workloads by bringing an understanding of Kubernetes volume and pod abstractions to vSphere.

In Kubernetes, CNS provides a volume driver that has two sub-components – the CSI driver and the syncer.
* The CSI driver is responsible for volume provisioning, attaching and detaching the volume to VMs, mounting, formatting and unmounting volumes from the pod within the node VM, and so on.
* The syncer is responsible for pushing PV, PVC, and pod metadata to CNS.

You can find the latest versions of these images on ECR Public Gallery.

[CSI driver](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/vsphere-csi-driver/csi/driver) | [CSI syncer](https://gallery.ecr.aws/eks-anywhere/kubernetes-sigs/vsphere-csi-driver/csi/syncer)
