## **Helm**
![Version](https://img.shields.io/badge/version-v4.1.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoieVZ2Vm4zalcvTTRlVHk3ODJMLy80a2hqaGw1eUNEMlBEQktYOGxLdkZYQmxMK2tWUTMyUHlxZDVIK2lYak9qM25OZm9IYTFkUGlXZ3dCOEhRb0dHMzBjPSIsIml2UGFyYW1ldGVyU3BlYyI6Im9EemRhdkg1Tll6d1lSaVciLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Helm](https://github.com/helm/helm) is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources.

### Updating
1. Review [releases](https://github.com/helm/helm/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @acool or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/helm/helm/blob/main/.github/workflows/build-pr.yml#L15).
4. Ensure correct patch has been used. The base patch to be used can be found [here](https://github.com/helm/helm/pull/10408). 
5. Run `make attribution checksums` in this folder.
6. Update CHECKSUMS as necessary (updated by default).
7. Update the version at the top of this Readme.
