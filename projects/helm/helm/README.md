## **Helm**
![Version](https://img.shields.io/badge/version-v3.8.1-blue)

[Helm](https://github.com/helm/helm) is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources.

### Updating
1. Update GIT_TAG file with new Helm version.
2. Update GOLANG_VERSION in Makefile consistent with upstream release
3. Run `make build` in this folder.
4. Update CHECKSUMS as necessary
5. Update the version at the top of this Readme.
