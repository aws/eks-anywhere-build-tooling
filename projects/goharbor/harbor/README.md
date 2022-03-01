## **harbor**
![Version](https://img.shields.io/badge/version-v2.4.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoickYvWTFSYk9KaTRtcUJib0VxajgrZElWNnVuaUdaK2RCU2VadmozTjVNMkVvZ3F6cnlFN05CS28zSmgwY3RLYzBNM2RIUE1lZHVhQlU5ZkxPa3NkRFBZPSIsIml2UGFyYW1ldGVyU3BlYyI6IjQ5YnVhZllRNExqayswdlEiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [harbor project](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management. Having a registry closer to the build and run environment can improve the image transfer efficiency. Harbor supports replication of images between registries, and also offers advanced security features such as user management, access control and activity auditing.

In EKS-A, harbor offers local cloud native registry service for Kubernetes clusters on vSphere infrastructure.

You can find the latest version of its images [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/harbor/).