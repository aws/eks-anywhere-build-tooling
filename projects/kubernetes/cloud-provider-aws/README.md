## **cloud-provider-aws**
![Version](https://img.shields.io/badge/version-v1.27.1-blue)

The AWS credential provider is a binary that is executed by kubelet to provide credentials for images in ECR. Refer to the [credential provider extraction Kubernetes Enhancement Proposals (KEP)](https://github.com/kubernetes/enhancements/tree/master/keps/sig-cloud-provider/2133-out-of-tree-credential-provider) for more details.

### Updating
1. Review [releases](https://github.com/kubernetes/cloud-provider-aws/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @acool or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/cloud-provider-aws/blob/master/go.mod#L3).
5. Run `make run-attribution-checksums-in-docker` in this folder.
6. Update CHECKSUMS as necessary (updated by default).
7. Update the version at the top of this Readme.
