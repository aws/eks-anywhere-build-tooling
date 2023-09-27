## **Kubernetes Cluster Autoscaler**
![1.21 Version](https://img.shields.io/badge/1--21%20version-v1.21.3-blue)
![1.22 Version](https://img.shields.io/badge/1--22%20version-v1.22.3-blue)
![1.23 Version](https://img.shields.io/badge/1--23%20version-v1.23.1-blue)
![1.24 Version](https://img.shields.io/badge/1--24%20version-v1.24.0-blue)
![1.25 Version](https://img.shields.io/badge/1--25%20version-v1.25.0-blue)
![1.26 Version](https://img.shields.io/badge/1--26%20version-v1.26.1-blue)
![1.28 Version](https://img.shields.io/badge/1--26%20version-v1.28.0-blue)

[Autoscaler](https://github.com/kubernetes/autoscaler) defines the cluster autoscaler.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kubernetes/autoscaler).

### Updating
1. Review [releases](https://github.com/kubernetes/autoscaler/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @jaxsen or @jonathanmeier5.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/autoscaler/blob/master/builder/Dockerfile#L15). (specified as source of truth [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-go-version-should-be-used-to-compile-ca))
4. If adding a new version, rip out cloud providers other than clusterapi. See below for details.
5. Run `make run-attribution-checksums-in-docker` for each release version in this folder.
6. Update CHECKSUMS as necessary (updated by default).
7. Update the versions at the top of this README.
8. Update the hardcoded appVersion values in sedfile.template


#### Removing Cloud Providers
We strip out all cloud providers except for clusterapi to reduce our CVE and maintenance surface area.

Setup a 1-XX directory like for other versions. `make build` will fail on `autoscaler/cluster-autoscaler/eks-anywhere-go-mod-download` because the `REMOVE_CLOUD_PROVIDERS` target removes dependencies used in upstream repo's `cluster-autoscaler/cloudprovider/builders` directory.

To get the build working, clean out other provider references. You will generate a patch like 1-26/patches/0001-Remove-Cloud-Provider-Builders-Except-Cluster-API.patch
```
cd autoscaler/cluster-autoscaler/cloudprovider/builders
 ls . | grep -v -e _all.go -e clusterapi.go -e _builder.go | xargs rm
git add .
```

Then clean out references to the other providers in:
```
builder_all.go
builder_clusterapi.go
cloud_provider_builder.go
```

commit and generate a patch using `git format-patch -1 HEAD`.

Then go into the cluster-autoscaler directory and tidy up and generate patch for go.mod and go.sum.
```
cd ../..
go mod tidy
git add go.mod go.sum
```

Commit and generate a patch for these changes.

Finally:
```
RELEASE_BRANCH=1-XX make clean
RELEASE_BRANCH=1-XX make build
```
