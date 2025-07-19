## **Etcdadm Bootstrap Provider**
![Version](https://img.shields.io/badge/version-v1.0.16-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVjVpNkZxSkZLbFBTZU1RNXZHY0pnREo1VDBVKzFDTEoybVdyd0VYUGNkV0RYQjdwdEM0VGtqMkxlbTdTeDdPT1NKbDRaYWdzdFE3NlFPcWowUUMzcWdnPSIsIml2UGFyYW1ldGVyU3BlYyI6Inlyd044bVFENkpiWU1JT08iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Etcdadm bootstrap provider](https://github.com/aws/etcdadm-bootstrap-provider) converts a CAPI Machine into an etcd member. It uses [etcdadm](https://github.com/kubernetes-sigs/etcdadm) to provision etcd members. It generates a script containing etcdadm init or join commands which then gets used to initialize a CAPI Machine. 
This allows users to provision standalone etcd nodes for their Kubernetes control plane.
The two bootstrap formats are cloud-init and bottlerocket.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/aws/etcdadm-bootstrap-provider).
