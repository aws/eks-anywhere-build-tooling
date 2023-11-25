## **AWS Distro for Kubernetes Metrics Server**
![1.24 Version](https://img.shields.io/badge/1--24%20version-v0.6.4-blue)
![1.25 Version](https://img.shields.io/badge/1--25%20version-v0.6.4-blue)
![1.26 Version](https://img.shields.io/badge/1--26%20version-v0.6.4-blue)
![1.27 Version](https://img.shields.io/badge/1--27%20version-v0.6.4-blue)
![1.28 Version](https://img.shields.io/badge/1--28%20version-v0.6.4-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSEFNYVlKSURxN25YRGpuWURwWmZOS05vbkl6YTdHTzNHTFJpdzdHZGJUL001ZlNqS1JhblM0QTl2VytuUzNRQ09WazJwRHVUZnp0dVRCb3dLTUVxb2w4PSIsIml2UGFyYW1ldGVyU3BlYyI6IkJIOGVvTFk2bWVVcnhUTkoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

AWS Distro for ([Metrics Server](https://github.com/kubernetes-sigs/metrics-server)) is an AWS supported version of the upstream Metrics Server and is distributed by Amazon EKS-D.

### Periodic Reviews
Review [helm chart releases](https://github.com/kubernetes-sigs/metrics-server/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updating
0. Latest images are automatically pulled in from EKS-D. GIT_TAG is generated dynamically. See EKS D releases in [ECR Gallery](https://gallery.ecr.aws/eks-distro/kubernetes-sigs/metrics-server)
1. For updating HELM_GIT_TAG, monitor [upstream releases](https://github.com/kubernetes-sigs/metrics-server/releases) and changelogs and when to bump the tag. Reach out to @jonathanmeier5 if you have any questions.
