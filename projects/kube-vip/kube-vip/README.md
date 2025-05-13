## **kube-vip**
![Version](https://img.shields.io/badge/version-v0.9.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZVJQQTc0Vk8rcHlMR0hOYnllRGNmV0NsQTNLNGFaS2hLME1MUmgwYkxpVUFoL0V0WHZzbXVCV1owQ0FUTlF6RHg1WXhWRXZLRzNwN2d2LzZGUVJvZ0pRPSIsIml2UGFyYW1ldGVyU3BlYyI6Im9jQmZMa216aHZpYmdrWDYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [kube-vip project](https://github.com/kube-vip/kube-vip) provides High-Availability and load-balancing for both the controlplane and Kubernetes services. The idea behind kube-vip is a small self-contained Highly-Available option for all environments, especially Bare-Metal, Edge (ARM/Raspberry Pi), Virtualisation, etc. kube-vip provides both a floating or virtual IP address for Kubernetes clusters as well as load-balancing the incoming traffic to various controlplane replicas. It thus simplifies the building of HA Kubernetes clusters with minimal components and configurations.

In EKS-A, kube-vip offers HA and load-balancing services for Kubernetes clusters on vSphere infrastructure.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kube-vip/kube-vip).
