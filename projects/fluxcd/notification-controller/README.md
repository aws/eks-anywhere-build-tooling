## **Notification Controller**
![Version](https://img.shields.io/badge/version-v0.13.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiVDhSNCt1djNtRytnWTgvQ01BMW13b2Y1YmZPakRrSGlRWitKZ0ZLZUdaS2xxclpLOFNidnBHNjBFWjRueHpOaGRrMzV5OUhLLzhRWHgyaC85R2tET2JZPSIsIml2UGFyYW1ldGVyU3BlYyI6IlZpNGwrazFrZndNMWE4cTciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [notification-controller](https://github.com/fluxcd/notification-controller) is a Kubernetes operator specialized in handling inbound and outbound events. The controller exposes an HTTP endpoint for receiving events from other controllers. It can be configured with Kubernetes custom resources such as `Alert`, `Event`,`Provider` and `Receiver` to define how events are processed and where to dispatch them.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/fluxcd/notification-controller).
