## **CloudMonkey**
![Version](https://img.shields.io/badge/version-6.5.0-blue)
![Build Status](https://camo.githubusercontent.com/e659888f1b9d9c9c3046644ef48da811f5b75b28fdff97596401ed5f39e866d5/68747470733a2f2f7472617669732d63692e636f6d2f6170616368652f636c6f7564737461636b2d636c6f75646d6f6e6b65792e7376673f6272616e63683d6d61696e)

`cloudmonkey` is a command line interface (CLI) for
[Apache CloudStack](http://cloudstack.apache.org).
CloudMonkey can be use both as an interactive shell and as a command line tool
which simplifies Apache CloudStack configuration and management.

The modern cloudmonkey is a re-written and simplified port in Go and can be used
with Apache CloudStack 4.9 and above. The legacy cloudmonkey written in Python
can be used with Apache CloudStack 4.0-incubating and above.

For documentation, kindly see the [wiki](https://github.com/apache/cloudstack-cloudmonkey/wiki).

The `eks-a` tool invokes cloudmonkey to perform validations on templates, fetching OVAs for building Cloudstack clusters, cleaning up stale VMs, etc.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/apache/cloudstack-cloudmonkey) and decide on new version.
   CloudMonkey maintainers are pretty good about calling breaking changes and other upgrade gotchas between release.  Please
   review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @vignesh-goutham
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check the `build` target for
   any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called
   from our Makefile.
1. Verify the golang version has not changed. The version specified in `go.mod` seems to be kept up to date.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
