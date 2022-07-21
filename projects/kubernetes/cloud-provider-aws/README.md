## **cloud-provider-aws**
![Version](https://img.shields.io/badge/version-v1.22.3-blue)

The [AWS cloud provider](https://github.com/kubernetes/cloud-provider-aws) provides the interface between a Kubernetes cluster and AWS service APIs. This project allows a Kubernetes cluster to provision, monitor and remove AWS resources necessary for operation of the cluster.

### Updating
1. Review [releases](https://github.com/kubernetes/cloud-provider-aws/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @acool or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/cloud-provider-aws/blob/master/go.mod#L3).
5. Run `make update-attribution-checksums-docker` in this folder.
6. Update CHECKSUMS as necessary (updated by default).
7. Update the version at the top of this Readme.
