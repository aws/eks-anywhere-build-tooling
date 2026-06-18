## **Cluster API Provider for Nutanix**
![Version](https://img.shields.io/badge/version-v1.6.1-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiWXJwaXpMSkhOTHpUSFU0K2IrSDlnZUFGMjdIRWIvSFNRZllHVmdURTFyRHpxOXlkSmdPTVd2YXhUSDVzY0U1ajVXUDhFRkZXYVp3ZHBhQS9jd3JUTXRNPSIsIml2UGFyYW1ldGVyU3BlYyI6IkhaaGRzVUdUQzlFY29MQ0YiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Cluster API Provider for Nutanix (CAPX)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix) is a the implementation of Cluster API for Nutanix.


### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix) and decide on the new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @pokearu or @abhinavmpandey08.
2. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
3. Check the go.mod file to see if the golang version has changed when updating a version. Update the `GOLANG_VERSION` in `Makefile` to match the version upstream.
4. Compare the old tag to the new, looking specifically for Makefile changes. If `release-manifests` target has changed in the Makefile, make the required changes in `create_manifests.sh`
5. Update checksums and attribution using `make attribution checksums`.
6. Update the version at the top of this Readme.
7. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
