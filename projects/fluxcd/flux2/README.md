## **Flux**
![Version](https://img.shields.io/badge/version-v0.13.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiYzRDM0E2d3BGeHZNenB4aVdRY0RqMkhoMUZBdjVHdjZsTSsrVEdhVEw1Sy9DREIwRUlwSEx4MFpoUVBiK2grUnhyT2JodmNVWUVaemFGR2JTOWhkWC9VPSIsIml2UGFyYW1ldGVyU3BlYyI6Im1VckJkV25QbHdyc0hRbmgiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Flux](https://github.com/fluxcd/flux2) is a tool for keeping Kubernetes clusters in sync with sources of configuration (like Git repositories), and automating updates to configuration when new code is deployed.

Flux v2 is constructed with the `GitOps` Toolkit, a set of composable APIs and specialized tools for building Continuous Delivery on top of Kubernetes. In version 2, Flux supports multi-tenancy and support for syncing an arbitrary number of Git repositories, among other long-requested features. It is built from the ground up to use Kubernetes' API extension system, and to integrate with Prometheus and other core components of the Kubernetes ecosystem.

EKS-A customers can store the configurations for all their clusters under version control, and leave the heavy-lifting of syncing and reconciliation to Flux.
