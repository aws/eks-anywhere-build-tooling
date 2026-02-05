## **cloud-provider-aws**
![1.27 Version](https://img.shields.io/badge/1--27%20version-v1.27.10-blue)
![1.28 Version](https://img.shields.io/badge/1--28%20version-v1.28.11-blue)
![1.29 Version](https://img.shields.io/badge/1--29%20version-v1.29.8-blue)
![1.30 Version](https://img.shields.io/badge/1--30%20version-v1.30.4-blue)
![1.31 Version](https://img.shields.io/badge/1--31%20version-v1.31.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiRXlTVFQzQ0dlVmNEZDZhV1lRWjlXYlFrOTNCbFA4cDJGVGNuMG9WdUVVM1BNazIzZ0hRRjVmYy9zK1NkblQ5Uk0xWmJJTlk0Um5XYTlmazg3MmxYamNZPSIsIml2UGFyYW1ldGVyU3BlYyI6ImtEM2pRV2d1QTlickRoYnUiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The AWS credential provider is a binary that is executed by kubelet to provide credentials for images in ECR. Refer to the [credential provider extraction Kubernetes Enhancement Proposals (KEP)](https://github.com/kubernetes/enhancements/tree/master/keps/sig-cloud-provider/2133-out-of-tree-credential-provider) for more details.

### Updating
1. Review [releases](https://github.com/kubernetes/cloud-provider-aws/releases) and changelogs in upstream repo and decide on new version. Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @acool or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile consistent with upstream release's [go version](https://github.com/kubernetes/cloud-provider-aws/blob/master/go.mod#L3).
5. Run `RELEASE_BRANCH=1-XX make attribution checksums` in this folder.
6. Update CHECKSUMS as necessary (updated by default).
7. Update the version at the top of this Readme.
