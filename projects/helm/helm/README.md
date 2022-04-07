## **Helm**
![Version](https://img.shields.io/badge/version-v3.8.1-blue)

[Helm](https://github.com/helm/helm) is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources.

### Updating
1. Review [releases](https://github.com/helm/helm/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @acool or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/helm/helm/blob/main/.github/workflows/build-pr.yml#L15).
4. Run `make build` in this folder.
5. Update CHECKSUMS as necessary (updated by default).
6. Update the version at the top of this Readme.
