## **Cluster API Provider for Tinkerbell**
![Version](https://img.shields.io/badge/version-v0.6.6-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZ2VEbkM5dXBVeFZhSE5IR2UvdjlNanY1RVo5S29zd0E2M1hiaCtNSEd5U3F2VUdCbkViWHVlclg5a093WVgrRUdqNnJLYUtpWjhqUWJaT0NJb3RaWWFjPSIsIml2UGFyYW1ldGVyU3BlYyI6ImtvQnZHalpsVjBCRk5jN2IiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Cluster API Provider Tinkerbell (CAPT)](https://github.com/tinkerbell/cluster-api-provider-tinkerbell) is an implementation of the Cluster API for Tinkerbell, which supports declarative infrastructure for Kubernetes cluster creation, configuration and management on Tinkerbell infrastructure. CAPT is designed to allow customers to bootstrap workload clusters using hardware managed by Tinkerbell. These clusters can be enhanced externally with remote power management and secure de-provisioning of instances using [Rufio](https://github.com/tinkerbell/rufio).

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/tinkerbell/cluster-api-provider-tinkerbell) and decide on the new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @pokearu or @abhnvp.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check if the [release](https://github.com/tinkerbell/cluster-api-provider-tinkerbell/blob/9e9c2a397288908f73a4f499ac00aaf96d15deb6/Makefile#L283)
   target has changed in the Makefile, and make the required changes in create_manifests.sh
1. Check the go.mod file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in
   Makefile to match the version upstream.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.