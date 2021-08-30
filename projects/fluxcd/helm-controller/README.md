## **Helm Controller**
![Version](https://img.shields.io/badge/version-v0.10.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiS045T05yUXhCRzNPeXZwczkwcjgrbm8wOWJmSXZ6dll3eHBlVTV3bERUSlhadlRyOGE1Q1AzeWpEQTlvN2RISG9MNnMrMGRmOG1FZ2N2d0Nxc0l0b2UwPSIsIml2UGFyYW1ldGVyU3BlYyI6IlpJMTJ1cUxhdzc4bWlqNFUiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [helm-controller](https://github.com/fluxcd/helm-controller) is a Kubernetes operator that allows users to declaratively manage Helm chart releases. The desired state of a Helm release is described through a Kubernetes Custom Resource named HelmRelease. Based on the creation, mutation or removal of a HelmRelease resource in the cluster, Helm actions are performed by the operator.

Some of the features of the Helm controller are:

* Watches for `HelmRelease` objects and generates `HelmChart` objects
* Supports `HelmChart` artifacts produced from `HelmRepository`, `GitRepository` and `Bucket` sources
* Fetches artifacts produced by source-controller from `HelmChart` objects
* Watches `HelmChart` objects for revision changes (including semver ranges for charts from `HelmRepository` sources)
* Performs automated Helm actions, including Helm tests, rollbacks and uninstalls
* Offers extensive configuration options for automated remediation (rollback, uninstall, retry) on failed Helm install, upgrade or test actions
* Runs Helm install/upgrade in a specific order, taking into account the depends-on relationship defined in a set of `HelmRelease` objects
* Reports Helm release statuses
* Built-in Kustomize compatible Helm post renderer, providing support for strategic merge, JSON 6902 and images patches

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/l0g8r8j6/fluxcd/helm-controller).
