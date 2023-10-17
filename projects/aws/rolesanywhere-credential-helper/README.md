## **rolesanywhere-credential-helper**
![Version](https://img.shields.io/badge/version-v1.0.4-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSlViVU00Vk4wbmtjdmtXTnVURVB6MUcwRktoKzduS0haUUpEeER4R3hKYkUwbUxNNldGS3F3cytyYXJmMllRQUI5b2dWQTJlanhBL1RhMERwS1lSNi9ZPSIsIml2UGFyYW1ldGVyU3BlYyI6IkdhZXIzVXk1b3JLdTFMRTAiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

The Rolesanywhere credential binary implements the [signing process](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/authentication-sign-process.html) for IAM Roles Anywhere's [CreateSession API](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/authentication-create-session.html) to get temporary credentials compatible with the `credential_process` feature available across the language SDKs.

### Updating
1. Review [releases](https://github.com/aws/rolesanywhere-credential-helper/releases) and changelogs in upstream repo and decide on new versions.Please review carefully and if there are questions about changes necessary to eks-anywhere to support the new version reach out to @junshun or @tlhowe.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile to be consistent with upstream's [go version](https://github.com/aws/rolesanywhere-credential-helper/blob/main/go.mod#L3).
4. Run `make run-attribution-checksums-in-docker` in this folder.
5. Update version at the top of this README