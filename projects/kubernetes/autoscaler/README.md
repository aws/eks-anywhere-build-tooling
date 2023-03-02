## **Kubernetes Cluster Autoscaler**
![1.21 Version](https://img.shields.io/badge/1--21%20version-v1.21.3-blue)
![1.22 Version](https://img.shields.io/badge/1--22%20version-v1.22.3-blue)
![1.23 Version](https://img.shields.io/badge/1--23%20version-v1.23.1-blue)
![1.24 Version](https://img.shields.io/badge/1--24%20version-v1.24.0-blue)
![1.25 Version](https://img.shields.io/badge/1--25%20version-v1.25.0-blue)
![1.26 Version](https://img.shields.io/badge/1--26%20version-v1.26.1-blue)

[Autoscaler](https://github.com/kubernetes/autoscaler) defines the cluster autoscaler.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kubernetes/autoscaler).

### Updating
1. Review [releases](https://github.com/kubernetes/autoscaler/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @jaxsen or @jonathanmeier5.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/autoscaler/blob/master/builder/Dockerfile#L15). (specified as source of truth [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-go-version-should-be-used-to-compile-ca))
4. Run `make update-attribution-checksums-docker` for each release version in this folder.
5. Update CHECKSUMS as necessary (updated by default).
6. Update the versions at the top of this Readme.
7. Update the hardcoded appVersion values in sedfile.template
