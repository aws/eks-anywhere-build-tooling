## **Hub**
![Version](https://img.shields.io/badge/version-352706903455cebc260fd565a38708c0e6423dc7-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQkRkY0htL2tWTlM0QmVLSS9SakxYOHBRTUxJNmczcVM4Nm1Wa0U1TFQvVkRDTHRadys0aEVIOStxc0V4aGxSQzNsdVZlaXV5R1YvaHZaOUZIZnRTTWtzPSIsIml2UGFyYW1ldGVyU3BlYyI6ImZjajIxazcybkxaZVdUR24iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Hub](https://github.com/tinkerbell/actions) is the repository that contains reusable Tinkerbell Actions. The different images are listed under [/actions](https://github.com/tinkerbell/actions/tree/main/actions).

### Updating

1. Review commits upstream [repo](https://github.com/tinkerbell/actions) and decide on release tag to track. 
1. Update the `GIT_TAG` file to have the new desired tag based on upstream.
1. Verify the golang version has not changed. Currently the version 1.15 mentioned in the [Dockerfile](https://github.com/tinkerbell/actions/blob/main/actions/cexec/v1/Dockerfile) of each action.
1. Verify no changes have been made to the Dockerfile for each action under under [actions](https://github.com/tinkerbell/actions/blob/main/actions) looking specifically for added dependencies or build 
process changes.
1. Update checksums and attribution using `make attribution checksums`.
1. Update the version at the top of this Readme.
1. Run `make generate` to update the UPSTREAM_PROJECTS.yaml file.
