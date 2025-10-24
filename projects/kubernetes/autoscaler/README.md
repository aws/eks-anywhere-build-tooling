## **Kubernetes Cluster Autoscaler**
![1.26 Version](https://img.shields.io/badge/1--26%20version-cluster--autoscaler--1.26.8-blue)
![1.27 Version](https://img.shields.io/badge/1--27%20version-cluster--autoscaler--1.27.8-blue)
![1.28 Version](https://img.shields.io/badge/1--28%20version-cluster--autoscaler--1.28.7-blue)
![1.29 Version](https://img.shields.io/badge/1--29%20version-cluster--autoscaler--1.29.5-blue)
![1.30 Version](https://img.shields.io/badge/1--30%20version-cluster--autoscaler--1.30.4-blue)
![1.31 Version](https://img.shields.io/badge/1--31%20version-cluster--autoscaler--1.31.2-blue)
![1.32 Version](https://img.shields.io/badge/1--32%20version-cluster--autoscaler--1.32.1-blue)
![1.33 Version](https://img.shields.io/badge/1--33%20version-cluster--autoscaler--1.33.0-blue)
![1.34 Version](https://img.shields.io/badge/1--33%20version-cluster--autoscaler--1.34.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiL0tWckptdkxsZEd1cXNiNTBncjRNVU5oekpZRlBkTDNBcFVvZkFOVHZwbTBKUm91QkR6RVN4QlhJWk42cXF3L29FMmdnTXUrVndiay8zVUQ0YjJsc21vPSIsIml2UGFyYW1ldGVyU3BlYyI6Ik1Gd2UwbmRXVWxSRTMvUHQiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Autoscaler](https://github.com/kubernetes/autoscaler) defines the cluster autoscaler.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/kubernetes/autoscaler).

### Updating
1. Review [releases](https://github.com/kubernetes/autoscaler/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @jaxsen or @jonathanmeier5.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/autoscaler/blob/master/builder/Dockerfile#L15). (specified as source of truth [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-go-version-should-be-used-to-compile-ca))
4. If adding a new version, rip out cloud providers other than clusterapi. See below for details.
5. Run `make attribution-checksums` for each release version in this folder.
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
```
git add builder_all.go
```
commit and generate a patch using `git format-patch -1 HEAD`.

Move the patch file to folder `1-xx/patches/`

To update go mod dependencies and create a patch, follow these steps:

* Run the following command:
```
 RELEASE_BRANCH=1-XX make update-go-mod-and-create-patch
```
Replace `1-XX` with the actual release branch number (e.g., `1-28`).

P.S. if the patch file `xxxx-Update-go.mod-Dependencies.patch` appears in the `1-xx/patches` folder beforehand, delete this patch file first before running the command `RELEASE_BRANCH=1-XX make update-go-mod-and-create-patch` since this command will generate a new patch `xxxx-Update-go.mod-Dependencies.patch`

This command performs the following actions:
- Invokes REMOVE_CLOUD_PROVIDER target
- Runs `go mod tidy` and updates the `go.mod` file with the latest dependencies
- Creates a patch file containing the changes
- Copies the patch file to the specified `RELEASE_BRANCH` patches directory

Rename the patch file as required.

Finally:
```
RELEASE_BRANCH=1-XX make clean
RELEASE_BRANCH=1-XX make build
```

#### Validating Helm Chart And Images

An easy way to validate your build is to install the helm chart to a kind cluster.

Install [kind](https://kind.sigs.k8s.io/) and create a cluster.

Then install the helm chart pointing at your personal registry using the command outputted when the build succeeds. For instance:
```
helm install cluster-autoscaler oci://public.ecr.aws/b9u1e4h9/cluster-autoscaler/charts/cluster-autoscaler --version 9.34.0-1.27-6444f7f1d05573c56b00d438af946ab9c36951a1 --set sourceRegistry=public.ecr.aws/a9u1e4h1 --set autoDiscovery.clusterName=foobar
```

Where `public.ecr.aws/b9u1e4h9` would be your personal registry.