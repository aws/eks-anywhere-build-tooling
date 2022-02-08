## **harbor**
![Version](https://img.shields.io/badge/version-v0.3.7-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZVJQQTc0Vk8rcHlMR0hOYnllRGNmV0NsQTNLNGFaS2hLME1MUmgwYkxpVUFoL0V0WHZzbXVCV1owQ0FUTlF6RHg1WXhWRXZLRzNwN2d2LzZGUVJvZ0pRPSIsIml2UGFyYW1ldGVyU3BlYyI6Im9jQmZMa216aHZpYmdrWDYiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [harbor project](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management. Having a registry closer to the build and run environment can improve the image transfer efficiency. Harbor supports replication of images between registries, and also offers advanced security features such as user management, access control and activity auditing.

In EKS-A, harbor offers local cloud native registry service for Kubernetes clusters on vSphere infrastructure.

You can find the latest version of its images [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/harbor/).