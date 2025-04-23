## **harbor-scanner-trivy**
![Version](https://img.shields.io/badge/version-v0.33.0--rc.2-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoieEpzUzBranRhT3NMMGdLU0lSVmh1S2RteDcyd1AwRU5LbVZFc2pnNlcvcWpaZHR4blQ3RktjbzllUmhwMmhma0pnZ2RWVEY0UEIzZ2NPc3pYQ2l1RFZvPSIsIml2UGFyYW1ldGVyU3BlYyI6IitiOTg2c2dOVW55cnVQREoiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The [Harbor Scanner Adapter for Trivy](https://github.com/goharbor/harbor-scanner-trivy) is a service that translates the Harbor scanning API into Trivy commands and allows Harbor to use Trivy for providing vulnerability reports on images stored in Harbor registry as part of its vulnerability scan feature.

### Updating

1. Update harbor-scanner-trivy tag when updating harbor tag if harbor is using a newer tag. Use the same tag that harbor uses by default. For instance [harbor v2.5.1 uses harbor-scanner-trivy v0.28.0 by default](https://github.com/goharbor/harbor/blob/v2.5.1/Makefile#L115) so when updating to harbor tag v2.5.1, update harbor-scanner-trivy tag to v0.28.0 or higher if security patching requires.
1. Review releases and changelogs in upstream [repo](https://github.com/goharbor/harbor-scanner-trivy) and decide on new version.
1. Update the `GIT_TAG` file to have the new desired version based on the upstream release tags.
1. Compare the old tag to the new, looking specifically for Makefile changes. Check the `build` target for any build flag changes, tag changes, dependencies, etc. Check that the manifest target has not changed, this is called from our Makefile.
1. Check the `go.mod` file to see if the golang version has changed when updating a version. Update the field `GOLANG_VERSION` in Makefile to match the version upstream.
1. Update checksums and attribution using make `run-attribution-checksums-in-docker`.
1. Update the version at the top of this `README`.
1. Run `make generate` to update the `UPSTREAM_PROJECTS.yaml` file.
