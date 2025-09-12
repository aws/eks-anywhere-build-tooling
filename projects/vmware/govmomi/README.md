## **GoVMOMI**
![Version](https://img.shields.io/badge/version-v0.52.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZ1FxODROWXBIdytIZVBsNUFzODdBcngreGlZdlVwdUliRThoTGNDajBab0YzdDZ3NzVKSnBTVDBTS0lzY25sUG82MzZPMWdteE14VkZrK0F2TlppKzBjPSIsIml2UGFyYW1ldGVyU3BlYyI6IkJHNTRwbGtDV2xYRCtaZ0wiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[GoVMOMI](https://github.com/vmware/govmomi) is a Go library for interacting with VMware vSphere APIs (ESXi and/or vCenter). It primarily provides convenience functions for working with the vSphere API. It provides Go bindings to the default implementation of the VMware Managed Object Management Interface (VMOMI)

In addition to the vSphere API client, this project also includes `govc`, a CLI for vSphere. The `eks-a` tool invokes govc to perform validations on templates, fetching OVAs for building vSphere clusters, cleaning up stale VMs, etc.


### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/vmware/govmomi) and decide on new version. 
The maintainers are pretty good about calling breaking changes and other upgrade gotchas between release.  Please
review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @abhinavmpandey08 or @pokearu
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. 
ex: [0.24.0 compared to 0.27.4](https://github.com/vmware/govmomi/compare/v0.24.0...v0.27.4). Check the [gorelease config](https://github.com/vmware/govmomi/blob/master/.goreleaser.yml)
for LDFLAGS changes, these should match what is in their Makefile and the EKS-A Makefile.
1. Verify the golang version has not changed. Use the github release [action](https://github.com/vmware/govmomi/blob/master/.github/workflows/govmomi-release.yaml) as the source
of truth for the golang version upstream builds with.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
