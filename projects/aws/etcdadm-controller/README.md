## **Etcdadm Controller**
![Version](https://img.shields.io/badge/version-v1.0.25-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiUTM2ZGs4R0p2QVVLamxqeW4zWEtPZkI0SXJXcVZGbXNyM3dEZXZTOUYyYUNmdXBmRm14a3NvcTBDMjZvWWFWU2I3RkEzSFVudVhRYWNQZGFuTWdJaWNnPSIsIml2UGFyYW1ldGVyU3BlYyI6IlN1UDBjNGlNbjg0RUxNcXMiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Etcdadm controller](https://github.com/aws/etcdadm-controller) provides a mechanism for etcd cluster lifecycle management. Features include:
* Etcd cluster provisioning and upgrade
* Generating etcd certs and making them available to the KubeadmControlPlane
* Periodic healthcheck for etcd members
* Reconciliation of etcdadm clusters across all namespaces

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/aws/etcdadm-controller).
