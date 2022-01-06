## **replicatedhq/troubleshoot**
![Version](https://img.shields.io/badge/version-v0.26.0-blue)

[troubleshoot](https://troubleshoot.sh/) is a kubectl plugin providing diagnostic tools for Kubernetes applications. It provides tools for collecting and analyzing cluster information including deployment statuses, cluster resources, and host logs. 

The EKS-A diagnostic bundle functionality is built on top of troubleshoot.

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/replicatedhq/troubleshoot) and decide on new version.
   The change log for the release will call out breaking changes and other gotchas between releases.  Please
   review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [0.13.2 compared to 0.23.0](https://github.com/replicatedhq/troubleshoot/compare/v0.13.10...v0.26.0). Check the `support-bundle` target for
   any build flag changes, tag changes, dependencies, etc.
- When proforming significant version upgrades, it is prudent to manually test that the new troubleshoot version 
  works with the existing EKS-A workflow. You can do this using the instructions in the section "Manually Testing Troubleshoot Version Compatibility".

1. Verify the golang version has not changed. The version specified in `go.mod` and the `Makefile` is kept up-to-date.
   - if the Go version has changed, you'll need to bump the Go version set in the [project specific `Makefile`](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/replicatedhq/troubleshoot/Makefile) to match. 
1. Update checksums and attribution by running the following command from the root of the build tooling repo:
```
make update-attribution-checksums-docker PROJECT=replicatedhq/troubleshoot
``` 
1. Update the version at the top of this Readme.
1. Run `make generate` from the root of the repo to update the UPSTREAM_PROJECTS.yaml file.


### Manually Testing Troubleshoot Version Compatibility
1. build the desired version locally (`make support-bundle`)
1. copy the support-bundle binary to your path
1. set the env var `MR_TOOLS_DISABLE=false`, so that the EKS-A commands will use local binaries
1. test the following commands against an existing EKS-A cluster:
    - run `eksctl-anywhere generate support-bundle-config`
    - run `eksctl-anywehre generate support-bundle`

#### Get and build Troubleshoot locally
```bash
git clone https://github.com/replicatedhq/troubleshoot
cd ./troubleshoot
git checkout $DESIRED_RELEASE_TAG
make support-bundle
cp ./bin/support-bundle /usr/local/bin
support-bundle version
```

#### Test Troubleshoot version with EKS-A
```bash
export MR_TOOLS_DISABLE=true
eksctl-anywhere generate support-bundle -f $MY_CLUSTER_CONFIG_FILE
eksctl-anywhere generate support-bundle-config -f $MY_CLUSTER_CONFIG_FILE
```