## **AWS Distro for Kubernetes Metrics Server**
![1.21 Version](https://img.shields.io/badge/1--21%20version-v0.6.1-blue)
![1.22 Version](https://img.shields.io/badge/1--22%20version-v0.6.1-blue)
![1.23 Version](https://img.shields.io/badge/1--23%20version-v0.6.1-blue)
![1.24 Version](https://img.shields.io/badge/1--23%20version-v0.6.1-blue)

AWS Distro for ([Metrics Server](https://github.com/kubernetes-sigs/metrics-server)) is an AWS supported version of the upstream Metrics Server and is distributed by Amazon EKS-D.

### Periodic Reviews
Review [helm chart releases](https://github.com/kubernetes-sigs/metrics-server/releases) periodically to identify new releases and decide on an update plan and an update schedule.

### Updating
0. Latest images are automatically pulled in from EKS-D. GIT_TAG is generated dynamically.
1. For updating HELM_GIT_TAG, monitor [upstream releases](https://github.com/kubernetes-sigs/metrics-server/releases) and changelogs and when to bump the tag. Reach out to @jonathanmeier5 if you have any questions.
