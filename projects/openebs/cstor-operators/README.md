## **OpenEBS**

Official website: https://openebs.io/
Upstream repository: https://github.com/openebs

OpenEBS is a Kubernetes native Container Attached Storage solution that makes it possible for Stateful applications to easily access Dynamic Local PVs or Replicated PVs. By using the CAS pattern users report lower costs, easier management, and more control of their teams.

[Upstream setup examples](https://openebs.io/docs/user-guides/quickstart)

### Updating

1. Review [releases notes](https://openebs.io/docs/introduction/releases)
    * Any changes to the upstream configuration needs a thorough review + testing
    * Deprecation or removal of any protocol must be considered breaking 
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Verify the golang version has not changed. 
1. Verify no changes have been made to the dockerfiles