## **Boots**
![Version](https://img.shields.io/badge/version-94e4b4899b383e28b6002750b14e254cfbbdd81f-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiTGRiNmxQbk5RTnZNbU41WW53bEdSTXRpVDRLaGxDRXJ1UEFnWkdlMVRGekhwdSttbXhmUWpNVFdOM200UkZZbTR3b3dTWkNXb2R1dnZDUHowQU1tU0VRPSIsIml2UGFyYW1ldGVyU3BlYyI6IjlnMlRWSTlpeXNLYmY3cmIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Boots](https://github.com/tinkerbell/boots) service handles DHCP, PXE, tftp, and iPXE for provisions in the Tinkerbell stack.

### Updating

1. Review the changelog upstream [repo](https://github.com/tinkerbell/boots) and decide on the new release tag to track.
1. Update the `GIT_TAG` file to have the new desired release tag.
1. Verify the golang version has not changed. Currently the version mentioned in the [go.mod](https://github.com/tinkerbell/boots/blob/94e4b4899b383e28b6002750b14e254cfbbdd81f/go.mod#L3) is being used to build.
1. Verify no changes have been made to the [dockerfile](https://github.com/tinkerbell/boots/blob/94e4b4899b383e28b6002750b14e254cfbbdd81f/Dockerfile) looking specifically for added runtime deps.
1. Update checksums and attribution using `make update-attribution-checksums-docker`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
1. Currently boots builds iPXE binaries that are [embedded](https://github.com/tinkerbell/boots/blob/94e4b4899b383e28b6002750b14e254cfbbdd81f/tftp/tftp.go#L14L24). These binaries are prebuilt and kept under [ipxe](https://github.com/aws/eks-anywhere-build-tooling/tree/main/projects/tinkerbell/boots/ipxe). Ensure to check for changes in the binaries when updating the release.

You should also update `BOOTS_SERVER_IMAGE` under `tinkerbell/sandbox/.env` with the new image tag once it's built.