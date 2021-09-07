## **Etcdadm Bootstrap Provider**
![Version](https://img.shields.io/badge/version-v0.1.0‒beta‒4.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVjVpNkZxSkZLbFBTZU1RNXZHY0pnREo1VDBVKzFDTEoybVdyd0VYUGNkV0RYQjdwdEM0VGtqMkxlbTdTeDdPT1NKbDRaYWdzdFE3NlFPcWowUUMzcWdnPSIsIml2UGFyYW1ldGVyU3BlYyI6Inlyd044bVFENkpiWU1JT08iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Etcdadm bootstrap provider](https://github.com/mrajashree/etcdadm-bootstrap-provider) provides a bootstrap provider for creating an etcd cluster using etcdadm. This allows users to provision standalone etcd nodes for their Kubernetes control plane.
Features include:
* Unstacked etcd topology to support standalone etcd node creation
* Custom schema and cloudconfig schema for bootstrapping etcd clusters
* Cluster API providers support (core, AWS, Docker, vSphere)
* [Bottlerocket](https://github.com/bottlerocket-os/bottlerocket) bootstrap support

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/mrajashree/etcdadm-bootstrap-provider).
