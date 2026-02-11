## **Cluster API Provider for Nutanix**
![Version](https://img.shields.io/badge/version-v0.5.5-blue)
![Build Status]()

The [Nutanix Cloud Controller Manager](https://github.com/nutanix-cloud-native/cloud-provider-nutanix) is a the implementation of cloud-controller-manager for Nutanix AHV.


### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/nutanix-cloud-native/cloud-provider-nutanix) and decide on the new version.
   Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @jaxesn, @pokearu or @abhinavmpandey08.
2. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
3. Check the go.mod file to see if the golang version has changed when updating a version. Update the `GOLANG_VERSION` in `Makefile` to match the version upstream.
4. Update checksums and attribution using `make attribution checksums`.
5. Update the version at the top of this Readme.
6. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
