## **Kube RBAC Proxy**
![Version](https://img.shields.io/badge/version-v0.20.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZUxRMjRTYUl6NEhJWkI1YVh5QVB3UitEY1dCcExLTUxGR21DQ0IySUZUTEI4N3I4NnMwbnIxUW9OZ1dudm9VdTRoaHVzUHhyMjNwek9wYXY3amh3NlFVPSIsIml2UGFyYW1ldGVyU3BlYyI6ImdSc3ZLZmpxM1BMYnd0dGwiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy) is an HTTP proxy for a single upstream endpoint, that can perform RBAC authorization against the Kubernetes API using `SubjectAccessReview`. In Kubernetes clusters without NetworkPolicies, any Pod can perform requests to every other Pod in the cluster. This proxy serves to restrict requests to only those Pods that present a valid and RBAC-authorized token or client TLS certificate.

A `kube-rbac-proxy` sidecar container is injected by Cluster API Kubeadm manager when provisioning the Kubernetes controlplane.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/brancz/kube-rbac-proxy).
