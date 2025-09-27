## **replicatedhq/troubleshoot**
![Version](https://img.shields.io/badge/version-v0.122.0-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiWlJsdnRmNnRYUjhUV20xaHJYTng2WXVlVXFBbHZPQnpnblh2bzFLYk1VUHAra2VpWFRFNWpMY0ovTC9PWnBBN2JEcDBXcjRSeVoxd3pyWWxQVzQzZFY4PSIsIml2UGFyYW1ldGVyU3BlYyI6IjZxRUdIK2N6TVZNNUdqL0oiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[troubleshoot](https://troubleshoot.sh/) is a kubectl plugin providing diagnostic tools for Kubernetes applications. It provides tools for collecting and analyzing cluster information including deployment statuses, cluster resources, and host logs. 

The EKS-A diagnostic bundle functionality is built on top of troubleshoot.

> Pro tip: `make help` to see the various targets and actions available

### Updating

1. Review releases and changelogs in upstream [repo](https://github.com/replicatedhq/troubleshoot) and decide on new version.
   The change log for the release will call out breaking changes and other gotchas between releases.  Please
   review carefully and if there are questions about changes necessary to eks-anywhere to support the new version
   and/or automatically update between eks-anywhere version reach out to @danbudris
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes.
   ex: [0.36.0 compared to 0.37.1](https://github.com/replicatedhq/troubleshoot/compare/v0.36.0...v0.57.1). Check the `support-bundle` target for
   any build flag changes, tag changes, dependencies, etc.
- When performing significant version upgrades, it is prudent to manually test that the new troubleshoot version 
  works with the existing EKS-A workflow. You can do this using the instructions in the section "Manually Testing Troubleshoot Version Compatibility".

1. Verify the golang version has not changed. The version specified in `go.mod` and the `Makefile` is kept up-to-date.
   - if the Go version has changed, you'll need to bump the Go version set in the [project specific `Makefile`](https://github.com/aws/eks-anywhere-build-tooling/blob/main/projects/replicatedhq/troubleshoot/Makefile) to match. 
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.


### Testing Troubleshoot Version Compatibility
1. build the desired version locally from the Build Tooling repo.
   - The make target `make binaries` will generate linux and darwin binaries for Troubleshoot, outputting them to `./_output/bin/troubleshoot`
   - For more options, see [the Development documentation](https://github.com/aws/eks-anywhere-build-tooling/blob/main/docs/development/building-locally.md)
1. add the `support-bundle` binary for your platform to your path
1. set the env var `MR_TOOLS_DISABLE=true`, so that the EKS-A commands will use local binaries
1. test the following commands against an existing EKS-A cluster:
```bash
export MR_TOOLS_DISABLE=true
eksctl-anywhere generate support-bundle -f $MY_CLUSTER_CONFIG_FILE
eksctl-anywhere generate support-bundle-config -f $MY_CLUSTER_CONFIG_FILE
```
