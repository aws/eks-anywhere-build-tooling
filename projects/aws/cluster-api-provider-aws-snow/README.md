## **CAPAS**
![Version](https://img.shields.io/badge/version-v0.2.7-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiOTZFc2lQejUvTDFxRWtWUVFCOWRITGtHcmRPUUtWR3FVb3d2aHV6eHM3Q0Q1c1RBMzlzZFpwNHhOMlZGNjFaeFRmU0dFNVk0bWZRTnI0T0s2a3poeEdNPSIsIml2UGFyYW1ldGVyU3BlYyI6InFKOTlCeVppcmI0dEtzWWwiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

*Note: We are not currently building the Snow provider nor are we generating the manifests. For now we are checking in
the generated manifests and image tags.

### Updating
1. Update GIT_TAG file with new image tag.
1. Run `kustomize build config/default` in the CAPAS repo to generate new manifests and update in the manifests folder.
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.
