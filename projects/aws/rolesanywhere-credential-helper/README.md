## **rolesanywhere-credential-helper**
![Version](https://img.shields.io/badge/version-v1.0.4-blue)

The Rolesanywhere credential binary implements the [signing process](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/authentication-sign-process.html) for IAM Roles Anywhere's [CreateSession API](https://docs.aws.amazon.com/rolesanywhere/latest/userguide/authentication-create-session.html) to get temporary credentials compatible with the `credential_process` feature available across the language SDKs.

### Updating
1. Review [releases](https://github.com/aws/rolesanywhere-credential-helper/releases) and changelogs in upstream repo and decide on new versions.
2. Update GIT_TAG file based on the upstream release tags.
3. Update GOLANG_VERSION in Makefile to be consistent with upstream's [go version](https://github.com/aws/rolesanywhere-credential-helper/blob/main/go.mod#L3).
4. Run `make update-attribution-checksums-docker` in this folder.
5. Update version at the top of this README